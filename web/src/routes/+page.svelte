<!--
  ~ Copyright 2025 Simon Emms <simon@simonemms.com>
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~     http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
-->

<script lang="ts">
  import { goto } from '$app/navigation';
  import { DateTime } from 'luxon';

  interface IRate {
    id: string;
    name: string;
    description: string;
    discount: number;
    totalInPence: number;
    payOnCheckIn?: boolean;
    prePaymentDate?: Date;
  }

  function calculateMinCheckoutDate(date: string): string {
    if (date) {
      return DateTime.fromISO(date).plus({ day: 1 }).toISODate() ?? '';
    }
    return '';
  }

  function calculateRates(checkinStr: string, checkoutStr: string): IRate[] {
    const rates: IRate[] = [];
    if (checkinStr === '' || checkoutStr === '') {
      return rates;
    }

    const checkin = DateTime.fromISO(checkinStr).set({
      hour: 15,
    });
    const checkout = DateTime.fromISO(checkoutStr).set({
      hour: 15,
    });
    const today = DateTime.fromISO(DateTime.now().toISODate()).set({
      hour: 15,
    });

    const { days: stayInDays } = checkout.diff(checkin, ['days']).toObject();
    const { days: daysTillCheckin } = checkin.diff(today, ['day']).toObject();

    if (!stayInDays || !daysTillCheckin) {
      return rates;
    }

    rates.push(
      ...[
        {
          id: 'flex',
          name: 'Flex',
          description: 'Pay on arrival.',
          totalInPence: dailyPriceInPence * stayInDays,
          discount: 0,
          payOnCheckIn: true,
        },
        {
          id: 'standard',
          name: 'Standard',
          description: 'Pay 2 days before arrival.',
          totalInPence: dailyPriceInPence * 0.9 * stayInDays, // Offer 10% discount
          discount: 10,
          prePaymentDate: checkin.plus({ days: -2 }).toJSDate(),
        },
        {
          id: 'nonflex',
          name: 'Non-Flex',
          description: 'Pay now.',
          totalInPence: dailyPriceInPence * 0.8 * stayInDays, // Offer 20% discount
          discount: 20,
          prePaymentDate: DateTime.now()
            .set({ hour: 15, minute: 0, second: 0, millisecond: 0 })
            .toJSDate(),
        },
      ],
    );

    return rates.filter((rate) => {
      if (rate.payOnCheckIn || !rate.prePaymentDate) {
        return true;
      }

      // Disable if payment date is in the past, or less than 2 days until check-in
      const target = DateTime.fromJSDate(rate.prePaymentDate).set({
        hour: 15,
      });

      const { days = 0 } = target.diff(today, ['days']).toObject();

      if (days < 0) {
        // Payment day in past - not available
        return false;
      }

      if (daysTillCheckin < 2 && rate.id !== 'flex') {
        // Payment date is less than two days away - only available if flex
        return false;
      }

      return true;
    });
  }

  async function submit(rate: IRate) {
    loading = true;
    const data = {
      hotelId: '12345',
      totalCostInPence: rate.totalInPence,
      checkInDate: DateTime.fromISO(checkin).toJSDate(),
      checkOutDate: DateTime.fromISO(checkout).toJSDate(),
      payOnCheckIn: rate.payOnCheckIn ?? false,
      prePaymentDate: rate.prePaymentDate,
      cardDetails: {
        number: '5555555555554444',
        expiryMonth: 12,
        expiryYear: new Date().getFullYear() + 3,
        securityCode: 123,
      },
    };

    const response = await fetch('/api/reservation', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      console.log(response);
      loading = false;
      err = response.statusText;
      return;
    }

    const res = await response.json();

    goto(`/booking/${res.bookingId}`);
  }

  const dailyPriceInPence = 15000;

  // Minimum booking date is tomorrow
  const minDate = DateTime.now().plus({ day: 1 });
  let minCheckoutDate = $state(minDate.plus({ day: 1 }).toISODate());

  let err: string = $state('');
  let checkin: string = $state('');
  let checkout: string = $state('');
  let rates: IRate[] = $state([]);
  let loading: boolean = $state(false);
</script>

<h1 class="title">Welcome to Hilbert's Hotel</h1>
<h2 class="subtitle">The world's best hotel with infinite rooms</h2>

<div class="columns">
  <div class="column is-half">
    {#if err}
      <article class="message is-danger">
        <div class="message-header">Error</div>
        <div class="message-body">{err}</div>
      </article>
    {/if}

    <div class="field">
      <label class="label" for="checkin">Check-in date</label>
      <div class="control">
        <input
          class="input"
          type="date"
          id="checkin"
          name="checkin"
          min={minDate.toISODate()}
          bind:value={checkin}
          required
          onchange={() => {
            minCheckoutDate = calculateMinCheckoutDate(checkin);
            rates = calculateRates(checkin, checkout);
          }}
        />
      </div>
    </div>

    <div class="field">
      <label class="label" for="checkout">Check-out date</label>
      <div class="control">
        <input
          class="input"
          type="date"
          id="checkout"
          name="checkout"
          min={minCheckoutDate}
          bind:value={checkout}
          required
          onchange={() => {
            rates = calculateRates(checkin, checkout);
          }}
        />
      </div>
    </div>

    {#if rates.length > 0}
      <div class="card">
        <div class="card-content">
          <p class="is-size-4 has-text-weight-semibold pb-4">Rates</p>
          {#each rates as rate (rate.id)}
            <div class="columns is-vcentered">
              <div class="column">
                <p class="is-size-5 has-text-weight-semibold">
                  {rate.name}
                </p>
                <p>{rate.description}</p>
              </div>
              <div class="column is-3 has-text-right">
                <p class="has-text-weight-semibold">
                  &pound;{rate.totalInPence / 100}
                </p>
                {#if rate.discount > 0}
                  <p class="has-text-danger">
                    -{rate.discount}%
                  </p>
                {/if}
              </div>
              <div class="column is-narrow">
                <div class="field">
                  <div class="control">
                    <button
                      type="submit"
                      onclick={() => submit(rate)}
                      class="button is-link"
                      class:is-loading={loading}
                    >
                      Book
                    </button>
                  </div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>

  <div class="column is-half">
    <figure class="image is-3-by-2">
      <img src="/img/hotel-room.jpg" alt="Standard room at Hilbert's Hotel" />
    </figure>
    <p class="is-size-7">
      Photo by <a
        href="https://unsplash.com/@3dottawa?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash"
        target="_blank">Point3D Commercial Imaging Ltd.</a
      >
      on
      <a
        href="https://unsplash.com/photos/white-bed-linen-on-bed-oxeCZrodz78?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash"
        target="_blank">Unsplash</a
      >
    </p>
  </div>
</div>
