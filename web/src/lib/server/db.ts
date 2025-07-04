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

import { Low } from 'lowdb';
import { JSONFilePreset } from 'lowdb/node';

export type Booking = {
  id: string;
  temporalId: string;
  paymentDate: Date;
  totalCostPence: number;
  payOnCheckIn: boolean;
  prePaymentDate: boolean;
  isPaid: boolean;
  checkIn: Date;
  checkOut: Date;
};

export type Data = {
  bookings: Booking[];
};

// Handle singleton
let db: Low<Data>;

export async function ensureDB(): Promise<Low<Data>> {
  if (!db) {
    const defaultData: Data = { bookings: [] };

    db = await JSONFilePreset('db.json', defaultData);
    await db.write();
  }

  return db;
}
