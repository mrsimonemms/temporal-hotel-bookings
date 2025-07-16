# temporal-hotel-bookings

<!-- markdownlint-disable-next-line MD013 MD034 -->
[![Go Report Card](https://goreportcard.com/badge/github.com/mrsimonemms/temporal-hotel-bookings)](https://goreportcard.com/report/github.com/mrsimonemms/temporal-hotel-bookings)

Book a hotel with different payment terms

<!-- toc -->

* [Overview](#overview)
* [Requirements](#requirements)
* [Getting started](#getting-started)
* [Contributing](#contributing)
  * [Open in a container](#open-in-a-container)
  * [Commit style](#commit-style)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

## Overview

This demonstrates how Temporal can help ensure long-lived processes are robust
and reliable.

In this example, [Hilbert's Hotel](https://en.wikipedia.org/wiki/Hilbert%27s_paradox_of_the_Grand_Hotel)
wants to allow their customers to book rooms. As there are an infinite number of
rooms, we don't need to check for availability, so we only need to code in the
reservation and payment system. The cost of a room is £150 per night, but if they
pay in advance they can get a discount - 10% if paid two days before check-in or
20% if paid now.

To achieve this, we have a couple of Temporal workflows:

* `BookHotel`: this makes the booking with the hotel, taking booking, customer
  and payment details. This returns the booking ID and the date when payment will
  been taken, depending upon the level of discount the customer has chosen. This
  then sends a confirmation message to the customer and then triggers the payment
  as an orphaned child workflow.
* `PayHotel`: this checks the guest in and makes the payment. If the customer
  has opted for 20% discount, this will be run immediately. If they have selected
  0%/10% discount, this workflow sleeps until the desired time is reached when
  it will run. This also listens for a `check-in` signal so that the hotelier
  can check the customer in early, should they arrive before the timer expires.
  A confirmation is also sent to the customer.

By design, all actions in the workflows are unreliable, with a one-in-three chance
of failing. This is done to demonstrate Temporal's retry configuration and how
it's able to make a reliable service from unreliable components.

## Requirements

This requires:
* [NodeJS](https://nodejs.org)
* [Go](https://go.dev)
* [Temporal CLI](https://docs.temporal.io/cli)

## Getting started

1. Start your Temporal server (only required if not using VSCode):

   ```sh
   temporal server start-dev
   ```

1. Start the backend server:

   ```sh
   air
   ```

1. Start the web application:

   ```sh
   cd web
   npm ci
   npm run dev
   ```

1. You can also run a lot of bookings in one go (optional):

   ```sh
   go run ./client 1000
   ```

Now you can open the [Temporal UI](http://localhost:8233) and the [web app](http://localhost:5173)

## Contributing

### Open in a container

* [Open in a container](https://code.visualstudio.com/docs/devcontainers/containers)

### Commit style

All commits must be done in the [Conventional Commit](https://www.conventionalcommits.org)
format.

```git
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```
