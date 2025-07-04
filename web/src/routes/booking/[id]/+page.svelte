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
</script>

<div class="content">
  <p>Thanks for your booking</p>
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
    {#if data.booking.isPaid}
      Paid with thanks
    {:else}
      <button class="button is-fullwidth is-primary">Check in now</button>
    {/if}
  </p>
</div>
