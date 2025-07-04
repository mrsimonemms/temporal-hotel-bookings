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

import { ensureDB } from '$lib/server/db';
import { ensureConnection } from '$lib/server/temporal';
import { json } from '@sveltejs/kit';
import { nanoid } from 'nanoid';
import type { RequestHandler } from './$types';

interface Booking {
  bookingId: string;
  hotelId: string;
  paymentDate: Date;
}

export const POST: RequestHandler = async ({ request }) => {
  const temporal = await ensureConnection();
  const db = await ensureDB();

  const data = await request.json();

  const workflowId = `book-${nanoid()}`;

  const handle = await temporal.workflow.start<
    (data: unknown) => Promise<Booking>
  >('BookHotel', {
    taskQueue: 'hotel-bookings',
    args: [data],
    workflowId,
  });

  const booking = await handle.result();

  db.data.bookings.push({
    id: booking.bookingId,
    temporalId: workflowId,
    paymentDate: new Date(booking.paymentDate).toString(),
    totalCostPence: data.totalCostInPence,
    payOnCheckIn: data.payOnCheckIn,
    prePaymentDate: data.prePaymentDate,
    isPaid: !data.payOnCheckIn,
    checkIn: data.checkInDate,
    checkOut: data.checkOutDate,
  });
  await db.write();

  return json(booking);
};
