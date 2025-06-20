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

import "time"

type BookHotelWorkflowInput struct {
	HotelID          string    // Nominal ID for hotel
	TotalCostInPence int32     // Cost of the hotel in pence
	CheckInDate      time.Time // Date of check-in
	CheckOutDate     time.Time // Date of check-out
	PayOnCheckIn     bool      // If paying on check-in, no pre-payment available

	// If not paying on check-in, pre-payment can be may any time between date
	// of booking and day before check-in date
	PrePaymentDate time.Time

	CardDetails CardDetails // Card details are required for all bookings
}

type BookHotelWorkflowResult struct {
	BookingID   string
	HotelID     string
	PaymentDate time.Time
}

type CardDetails struct {
	Number       string
	ExpiryMonth  int
	ExpiryYear   int
	SecurityCode int
}

type PayHotelInput struct {
	PaymentDate      time.Time
	CardDetails      CardDetails
	BookingID        string
	TotalCostInPence int32
}

type PayHotelResult struct{}

type ReserveHotelInput struct {
	HotelID      string
	CheckInDate  time.Time
	CheckOutDate time.Time
}

type ReserveHotelResult struct {
	BookingID string
}
