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

// This is a way of running lots of bookings and check-ins concurrently.
// This defaults to 100 concurrent bookings, but can be configured by passing a
// number through as the first command line argument.
//
// This creates a booking, waits for a few seconds then checks the customer in

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	temporalhotelbookings "github.com/mrsimonemms/temporal-hotel-bookings"
	"go.temporal.io/sdk/client"
)

type future struct {
	ctx      context.Context
	workflow client.WorkflowRun
}

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

	// Default to 100 runs, but receive number from first input arg
	count := 100
	if len(os.Args) >= 2 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			//nolint:gocritic
			log.Fatalln("Error converting number", err)
		}
		count = i
	}

	futures := make([]future, 0)

	for range count {
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
				PrePaymentDate: time.Now().Add(time.Hour),

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
			log.Fatalln("Unable to execute workflow", err)
		}
		futures = append(futures, future{
			ctx:      ctx,
			workflow: we,
		})

		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	}

	for _, f := range futures {
		// Synchronously wait for the workflow completion.
		var result temporalhotelbookings.BookHotelWorkflowResult
		err = f.workflow.Get(f.ctx, &result)
		if err != nil {
			log.Fatalln("Unable get workflow result", err)
		}
		log.Printf("Workflow result: %+v", result)
		log.Printf("Payment will be taken on %s", result.PaymentDate)
	}

	fmt.Println("Sleeping before check-in")

	time.Sleep(time.Second * 15)

	for _, f := range futures {
		workflowID := fmt.Sprintf("%s_payment", f.workflow.GetID())
		err = c.SignalWorkflow(f.ctx, workflowID, "", "check-in", nil)
		if err != nil {
			log.Fatalf("Unable to signal workflow: %v", err)
		}
	}

	fmt.Println("Everyone checked-in - have a lovely stay")
}
