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
  import Navbar from '$lib/components/navbar.svelte';

  let { children } = $props();
</script>

<svelte:head>
  <title>Hilbert's Hotel</title>
</svelte:head>

<div class="container">
  <Navbar />
</div>

<section class="section">
  <div class="container">
    <h1 class="title">Welcome to Hilbert's Hotel</h1>
    <h2 class="subtitle">The world's best hotel with infinite rooms</h2>

    <div class="columns">
      <div class="column is-half">
        {@render children()}
      </div>

      <div class="column is-half">
        <figure class="image is-3-by-2">
          <img
            src="/img/hotel-room.jpg"
            alt="Standard room at Hilbert's Hotel"
          />
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
  </div>
</section>
