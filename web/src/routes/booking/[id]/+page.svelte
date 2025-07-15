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
  import { DateTime } from 'luxon';
  import type { PageProps } from './$types';

  let { data }: PageProps = $props();
  let err: string = $state('');
  let isPaid: boolean = $state(data.booking.isPaid);
  let loading: boolean = $state(false);

  async function checkin() {
    loading = true;
    const response = await fetch(`/api/reservation`, {
      method: 'PUT',
      body: JSON.stringify({
        id: data.booking.id,
      }),
    });

    if (!response.ok) {
      err = response.statusText;
      return;
    }

    await response.json();

    isPaid = true;
  }
</script>

{#if err}
  <article class="message is-danger">
    <div class="message-header">Error</div>
    <div class="message-body">{err}</div>
  </article>
{/if}

<div class="content">
  <p>Thanks for your booking</p>
  <p>Ref: {data.booking.id}</p>
  <dl>
    <dt>Check-in Date:</dt>
    <dd>
      {DateTime.fromISO(data.booking.checkIn).toLocaleString(
        DateTime.DATE_MED_WITH_WEEKDAY,
      )}
    </dd>

    <dt>Check-out Date:</dt>
    <dd>
      {DateTime.fromISO(data.booking.checkOut).toLocaleString(
        DateTime.DATE_MED_WITH_WEEKDAY,
      )}
    </dd>

    <dt>Total cost:</dt>
    <dd>&pound;{data.booking.totalCostPence / 100}</dd>
  </dl>

  <p>
    {#if isPaid}
      You're checked into Room &infin;. Happy walking
    {:else}
      <button
        class="button is-fullwidth is-primary"
        class:is-loading={loading}
        onclick={() => checkin()}
      >
        Check in now
      </button>
    {/if}
  </p>
</div>
