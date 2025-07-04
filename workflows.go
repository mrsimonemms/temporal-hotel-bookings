/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package temporalhotelbookings

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

func BookHotel(ctx workflow.Context, data *BookHotelWorkflowInput) (*BookHotelWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Running workflow")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
	})

	var a *activities

	// First job is to reserve the hotel. This reserves the booking with the hotel, but does not pay for it yet.
	var reserveHotelResult *ReserveHotelResult
	if err := workflow.ExecuteActivity(ctx, a.ReserveHotel, ReserveHotelInput{
		HotelID:      data.HotelID,
		CheckInDate:  data.CheckInDate,
		CheckOutDate: data.CheckOutDate,
	}).Get(ctx, &reserveHotelResult); err != nil {
		logger.Error("Error reserving hotel", "error", err)
		return nil, fmt.Errorf("error reserving hotel: %w", err)
	}

	paymentDate := data.CheckInDate
	if !data.PayOnCheckIn {
		paymentDate = data.PrePaymentDate
	}

	// The payment is a child workflow that is decoupled from the booking workflow.
	// This allows the booking workflow to create the payment workflow, but it
	// doesn't know anything about the execution as this may happen any point
	// between now and check-in date.
	paymentCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
		// Need to disconnect the child workflow else it terminates when the parent finishes
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		WorkflowID:        fmt.Sprintf("%s_payment", workflow.GetInfo(ctx).WorkflowExecution.ID),
	})

	paymentWorkflow := workflow.ExecuteChildWorkflow(paymentCtx, PayHotel, &PayHotelInput{
		BookingID:        reserveHotelResult.BookingID,
		CardDetails:      data.CardDetails,
		TotalCostInPence: data.TotalCostInPence,
		PaymentDate:      paymentDate,
	})

	// Ensure that the workflow is executed, but we don't want to wait for the
	// result as this might be some time.
	if err := paymentWorkflow.GetChildWorkflowExecution().Get(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to start payment workflow: %w", err)
	}

	return &BookHotelWorkflowResult{
		BookingID:   reserveHotelResult.BookingID,
		HotelID:     data.HotelID,
		PaymentDate: paymentDate,
	}, nil
}

func PayHotel(ctx workflow.Context, data *PayHotelInput) (*PayHotelResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Paying for hotel")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
	})

	timerCtx, cancelHandler := workflow.WithCancel(ctx)

	// Declare a new selector
	selector := workflow.NewSelector(ctx)

	// Configure a signal to cancel the timer
	channel := workflow.GetSignalChannel(ctx, "check-in")
	selector.AddReceive(channel, func(c workflow.ReceiveChannel, more bool) {
		// Run the receiver to clear the signal - we don't receive any data
		c.Receive(ctx, nil)

		logger.Info("Check in received - cancelling timer")
		cancelHandler()
	})

	// Activate the selectors
	selector.Select(ctx)

	// Calculate the delay until payment is taken
	t := time.Now().UTC()
	delay := max(data.PaymentDate.UTC().Sub(t), 0)
	logger.Info("Delaying workflow", "time", data.PaymentDate, "delay", delay)

	// Activate the sleep
	if err := workflow.Sleep(timerCtx, delay); err != nil {
		// Allow sleep to be cancelled
		if !errors.Is(err, workflow.ErrCanceled) {
			logger.Error("Error sleeping payment")
			return nil, fmt.Errorf("error sleeping payment: %w", err)
		}
	}

	var a *activities

	var result *PayHotelResult
	if err := workflow.ExecuteActivity(ctx, a.PayHotel, data).Get(ctx, &result); err != nil {
		logger.Error("Error paying hotel", "error", err)
		return nil, fmt.Errorf("error paying hotel: %w", err)
	}

	return result, nil
}
