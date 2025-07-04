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
	"context"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
)

type activities struct{}

// Pay the hotel
func (a *activities) PayHotel(ctx context.Context, data *PayHotelInput) (*PayHotelResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Paying hotel")

	// Simulate a time delay
	time.Sleep(time.Second * 5)

	if err := SimulateFailure(ctx); err != nil {
		return nil, err
	}

	logger.Info("Hotel successfully paid")
	return &PayHotelResult{
		TransactionID: uuid.NewString(),
	}, nil
}

// Reserve the hotel booking
func (a *activities) ReserveHotel(ctx context.Context, data *ReserveHotelInput) (*ReserveHotelResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Reserve hotel", "hotelId", data.HotelID)

	// Simulate a time delay
	time.Sleep(time.Second)

	if err := SimulateFailure(ctx); err != nil {
		return nil, err
	}

	logger.Info("Hotel successfully reserved")
	return &ReserveHotelResult{
		//nolint:gosec
		BookingID: strconv.Itoa(rand.Int()),
	}, nil
}

func NewActivities() (*activities, error) {
	return &activities{}, nil
}
