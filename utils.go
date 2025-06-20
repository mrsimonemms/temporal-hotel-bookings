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
	"fmt"
	"math/rand/v2"

	"go.temporal.io/sdk/activity"
)

// Pseudo-randomise failure - this is obviously not going to be in a real-world
// version, but it exists to demonstrate that APIs are a black box and we have
// no control over the failures
func SimulateFailure(ctx context.Context) error {
	logger := activity.GetLogger(ctx)

	//nolint:gosec
	r := rand.IntN(3)
	logger.Debug("Simulating failure", "randomisedInt", r)

	if r == 1 {
		return fmt.Errorf("simulate failure")
	}
	return nil
}
