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
	HotelID          string    `json:"hotelId"`          // Nominal ID for hotel
	TotalCostInPence int32     `json:"totalCostInPence"` // Cost of the hotel in pence
	CheckInDate      time.Time `json:"checkInDate"`      // Date of check-in
	CheckOutDate     time.Time `json:"checkOutDate"`     // Date of check-out
	PayOnCheckIn     bool      `json:"payOnCheckIn"`     // If paying on check-in, no pre-payment available

	// If not paying on check-in, pre-payment can be may any time between date
	// of booking and day before check-in date
	PrePaymentDate time.Time `json:"prePaymentDate"`

	CardDetails CardDetails `json:"cardDetails"` // Card details are required for all bookings
}

type BookHotelWorkflowResult struct {
	BookingID   string    `json:"bookingId"`
	HotelID     string    `json:"hotelId"`
	PaymentDate time.Time `json:"paymentDate"`
}

type CardDetails struct {
	Number       string `json:"number"`
	ExpiryMonth  int    `json:"expiryMonth"`
	ExpiryYear   int    `json:"expiryYear"`
	SecurityCode int    `json:"securityCode"`
}

type PayHotelInput struct {
	PaymentDate      time.Time   `json:"paymentDate"`
	CardDetails      CardDetails `json:"cardDetails"`
	BookingID        string      `json:"bookingId"`
	TotalCostInPence int32       `json:"totalCostPence"`
}

type PayHotelResult struct {
	TransactionID string `json:"transactionId"`
}

type ReserveHotelInput struct {
	HotelID      string    `json:"hotelId"`
	CheckInDate  time.Time `json:"checkInDate"`
	CheckOutDate time.Time `json:"checkOutDate"`
}

type ReserveHotelResult struct {
	BookingID string `json:"bookingId"`
}

type ConfirmationType string

const (
	ConfirmationTypeBooking = "booking"
	ConfirmationTypePayment = "payment"
)

type SendConfirmationInput struct {
	BookingID string           `json:"bookingId"`
	Type      ConfirmationType `json:"type"`
}

type SendConfirmationOutput struct {
	TransactionID string `json:"transactionId"`
}
