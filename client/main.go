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

package main

import (
	"context"
	"log"
	"time"

	temporalhotelbookings "github.com/mrsimonemms/temporal-hotel-bookings"
	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: "hotel-bookings",
	}

	now := time.Now().UTC()

	ctx := context.Background()

	we, err := c.ExecuteWorkflow(
		ctx,
		workflowOptions,
		temporalhotelbookings.BookHotel,
		&temporalhotelbookings.BookHotelWorkflowInput{
			HotelID:          "12345", // Arbitrary hotel ID
			TotalCostInPence: 18999,   // Total cost £189.99
			CheckInDate: func() time.Time {
				// Set check-in to one week in the future at 15:00
				target := now.Add(time.Hour * 24 * 7)
				return time.Date(target.Year(), target.Month(), target.Day(), 15, 0, 0, 0, time.UTC)
			}(),
			CheckOutDate: func() time.Time {
				// Set check-in to 8 days in the future at 11:00
				target := now.Add(time.Hour * 24 * 8)
				return time.Date(target.Year(), target.Month(), target.Day(), 11, 0, 0, 0, time.UTC)
			}(),
			PayOnCheckIn:   false,
			PrePaymentDate: time.Now().Add(time.Minute), // This rate is only available with immediate payment

			// This is a demo. Do **NOT** use your real card details.
			CardDetails: temporalhotelbookings.CardDetails{
				Number:       "5555555555554444",
				ExpiryMonth:  1,
				ExpiryYear:   now.Year() + 3,
				SecurityCode: 123,
			},
		},
	)
	if err != nil {
		//nolint:gocritic
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result temporalhotelbookings.BookHotelWorkflowResult
	err = we.Get(ctx, &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Printf("Workflow result: %+v", result)
	log.Printf("Payment will be taken on %s", result.PaymentDate)
}
