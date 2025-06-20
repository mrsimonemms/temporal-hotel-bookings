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
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func BootHotel(ctx workflow.Context, data BookHotelWorkflowInput) (*BookHotelWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Running workflow")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour,
	})

	var a *activities

	var createBookingResult *BookResult
	if err := workflow.ExecuteActivity(ctx, a.CreateBooking, BookInput{}).Get(ctx, &createBookingResult); err != nil {
		logger.Error("Error creating a booking", "error", err)
		return nil, fmt.Errorf("error creating a booking: %w", err)
	}

	return nil, nil
}
