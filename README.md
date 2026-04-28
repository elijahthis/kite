# Kite by Grey: A Multi-currency Wallet

Kite is a secure prototype of a multi-currency wallet system that allows users store money in a variety of different currenies (currently EUR, GBP, KES, NGN, USD), as well as convert between these currencies and receive payouts. A user's transaction history is visible on their dashboard.

## Key Features

- **Auth:** Simple. Sign up with username, email and password.
- **Balances:** A user starts off with zero in every currency.
- **Deposits:** user can deposit money into account.
- **Conversion:** user can convert money between currencies.
- **Payouts:** user can withdraw money.
- **Transaction History:** User can view their full transaction history.

## This project was built with:

- Go
- React + TS

### 1. Prerequisites

- Docker
- Docker Compose

### 2. Running the system

```bash
git clone https://github.com/elijahthis/kite.git
cd kite
make up
```

You should be able to interact with the

- frontend via http://localhost:5173
- backend via http://localhost:8080

### 3. Testing

```bash
make test
```

### 3. E-R Diagram

https://mermaid.ai/d/e1bb123e-27f2-449f-bb3e-03035a9f2e63
<img width="8179" height="6257" alt="kite-erd" src="https://github.com/user-attachments/assets/bf17ad52-29c7-4ee1-8444-88961029be8c" />

### 4. Loom Video:

https://www.loom.com/share/5c794b8b55ca4cd4936cd599bcbfccaa

## Core Architectural Decisions

The backend is structured to ensure separation of business logic from infra, DB and HTTP layers. This makes it easier to test the core logic of the system without http/db overhead, as well as makes it such that the system an easily be extended in the future to allow for other means of interfacing with it (e.g. gRPC, CLI, etc.)

- `domain/`: Contains business logic, structs, and repository interfaces.
- `application/`: Contains use-cases/services (e.g., `ConversionService`, `PayoutService`).
- `interfaces/`: Handles HTTP routing, middleware, and JSON serialization.
- `infra/`: Contains Postgres implementations, FX API integrations, and cryptography.

## Auth Choice

I chose a token-based authentication mechanism using JWTs, but sent as session-like cookies.

Instead of sending the JWT in the response body to be stored in the frontend's localStorage (which is highly vulnerable to XSS), the backend directly sets the JWT in a strictly `HttpOnly`, `Secure`, and `SameSite=Strict` cookie. This ccombines the stateless scalability of token-based auth with the security of traditional server-side sessions.

### The Double-Entry Ledger System

Kite does not store balances as mutable rows in a database. Instead, balances are calculated dynamically from an immutable, append-only double-entry ledger.

- `transactions` table: Represents the business event (Deposit, Payout, Conversion).
- `ledger_entries` table: Represents the immutable ledger entries.
  The TransactionBuilder is built to guarantee that the sum of Debits perfectly equals the sum of Credits before any DB insertion occurs.
  `BIGINT` was used at the DB layer to represent money (cents, kobo). THe only floating points in the system is reserved for FX multiplication, and the results are converted to integer before the DB gets involved

### Other Features

- **idempotency of transactions:** done at the DB layer. Every transaction has a unique `reference`
- **Concurrency & Race Conditions:** In FX, this was handled using quote_id. Even if 500 concurrent threads attempt to execute the same quote, the database will accept exactly 1 and reject the other 499
- **Async bank flow:** a goroutine was used to move the transaction between statuses, from `PENDING` to `SUCCESS`
- **ACID:** DB transactions are wrapped in `AtomicUnit` to ensure they're correctly rolled back upon failure.
- **Observability**: structired logs using the `zerolog` package, and Threaded Request IDs to go with every log

### Other Features

- **Session Management:** The backend issues a strictly `HttpOnly`, `Secure`, `SameSite=Strict` cookie. The frontend does not handle raw JWTs, because raw JWT tokens stored in `localStorage` pose an XSS risk.
- **Tanstack Query for API integration**
- **Vite & Nginx Proxy**

### Testing Strategy

Most notably:

- **Ledger Integrity:** Proves the TransactionBuilder actively rejects unbalanced debits/credits or negative values.
- **Idempotency Locks:** Proves that identical references yield a domain-specific duplicate error without mutating the ledger.
- **Concurrency Simulation:** Spins up 500 concurrent goroutines attempting to execute the same FX quote simultaneously, verifying that only 1 succeeds.
- **Reversal Verification:** Proves that simulated failed payouts automatically append an offsetting double-entry reversal to the ledger.

I added a few more unit tests and an integration test for reversal verification, and expect to test more extensively in a prod environment

## Future Work

At production scale, there are a few features I'd implement to ensure things don't fall apart:

- Queue for transactions (possibly Kafka), to ensure no event ever gets dropped, even in high-throughput environments or server downtime
- implement a Dead Letter Queues (DLQ) to ensure failed transactions can be retried
- Automated Reconciliation: A nightly or weekly (depending on regulation) CRON job to sum the entire ledger_entries table to ensure the sum of all currencies equal 0. Currently I'm only balancing transactions, with the expectation that if every transcction is mathematically balanced, the ledger (an aggregation of these transactions) must also be perfectly balanced.
  A periodic CRON job would improve consistency. If it fails, an SRE alert should triggered.
- Rate Limiting
- On SRE alerts, I'd also implement a metrics monitoring system (most likely Prometheus, as the system is in Go), as well as Grafana charts and automated alets for when metrics surpass or fall below certain thresholds.
- i might also move quotes from a postgres table to a Redis cache with a ttl
- cookies are currently active for 30 mins for ease of testability during development. in prod, i'd also consider bring in that down to 5 mins, and implementing refresh tokens and a token revocation list for logouts.
- AML limits, as mentioned in the assessment docs
- password length and complexity enforcement, as well as OTP
- currently the 1% spread on fx conversions doesn't go to any account. i'd consider setting up an FX Revenue Account and route the money there.

## What breaks first at 1M users?

- most likely the goroutine that handles payout. Too many goroutines spawned, we'd run out of memory quickly.
  Fix: Move to a worker pool consuming from a messaging queue (Kafka/SQS) with a Dead Letter Queue (DLQ) for retries.
- DB as well: We'll need to horizontally scale Postgres.
- FX Cache: We'd need to get rid of the in-memory conversion cache.
- Lack of periodic automatic ledger-wide balancing to spot data corruption quickly.
