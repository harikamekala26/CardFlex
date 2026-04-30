## 📋 Project Name

# CardFlex - Multi-Tenant Credit Card Account Management Portal

## 📋 Project Description

**CardFlex** is a scalable, multi-tenant white-label SaaS platform that enables multiple credit card issuing companies to deliver a fully branded online account management experience through a single configurable application. The platform dynamically adapts feature availability based on tenant-specific configurations, branding and credit card usage. It eliminates duplicated development efforts and allows rapid rollout of new features across all clients. The architecture aims to support feature toggling per issuer, reduce operational costs, and significantly accelerate time-to-market for credit card companies.

### Problem Statement

Credit card issuing companies need to provide their customers with online account management portals. Currently, each company must build and maintain separate applications, leading to:

- High development and maintenance costs
- Duplicated effort across the industry
- Inconsistent user experiences
- Difficulty implementing new features across multiple platforms

### Solutions

CardFlex provides a single, multi-tenant platform where:

- Multiple credit card companies share the same codebase
- Each company gets a fully branded, customized experience
- URL parameters control tenant-specific configurations
- Features can be enabled/disabled per company
- Reduces development costs and accelerates time-to-market

### Current Functionality

- Tenant-aware branding and routing using the `company` query parameter.
- Registration and login flows scoped to the selected tenant.
- JWT-protected dashboard, profile, and payment endpoints.
- Dashboard account summary and transaction history loaded from the backend database.
- Tenant feature flags returned by the dashboard API.
- Payment UI hidden when a tenant has payments disabled.
- Payment submission with validation, balance updates, and transaction creation.
- Profile retrieval for the authenticated tenant user.
- Default account creation for new registered users.
- Login-time account backfill for existing users missing an account.

### Example Usage

```
https://cardflex.com?company=chase-bank
https://cardflex.com?company=capital-one
https://cardflex.com?company=citibank
```

## Run Locally

### Prerequisites

- Go 1.22 or newer
- Node.js and npm
- Chrome or Chromium for Angular unit tests and Cypress
- SQLite support, preconfigured through the Go SQLite driver

### Backend Environment

Create `backend/.env` from the example file:

```bash
cd backend
cp .env.example .env
```

Expected local values:

```dotenv
PORT=8080
DB_DRIVER=sqlite
DB_DSN=cardflex.db
JWT_SECRET=replace-with-a-local-development-secret
```

Environment variables:

- `PORT`: backend HTTP port. The frontend expects `8080` by default.
- `DB_DRIVER`: use `sqlite` for local development.
- `DB_DSN`: SQLite database path. `cardflex.db` creates/uses the local repo database.
- `JWT_SECRET`: required signing secret for login tokens.

### Start Backend

```bash
cd backend
go run .
```

Backend runs on `http://localhost:8080`.

Health check:

```bash
curl http://localhost:8080/ping
```

### Start Frontend

Install dependencies and start Angular:

```bash
cd frontend
npm install
npm start
```

Frontend runs on `http://localhost:4200` and calls `http://localhost:8080` from `frontend/src/environments/environment.ts`.

### Example Tenant URLs

```text
http://localhost:4200/?company=chase-bank
http://localhost:4200/login?company=chase-bank
http://localhost:4200/dashboard?company=chase-bank
http://localhost:4200/profile?company=chase-bank
http://localhost:4200/?company=capital-one
http://localhost:4200/?company=wells-fargo
```

### Demo Users

Seeded demo users are available after the backend migrations run.

Password for all demo users:

```text
secret123
```

Example demo accounts:

```text
demo+chase-bank@cardflex.local
demo+wells-fargo@cardflex.local
demo+capital-one@cardflex.local
```

Feature flag examples:

- `chase-bank`: payments enabled, profile enabled
- `wells-fargo`: payments disabled, profile enabled
- `capital-one`: payments enabled, profile disabled

## Tests

### Backend Tests

```bash
cd backend
GOCACHE=/private/tmp/cardflex-go-build-cache go test -v ./...
```

### Frontend Unit Tests

```bash
cd frontend
npm test
```

### Frontend Production Build

```bash
cd frontend
npm run build
```

### Cypress End-to-End Tests

Start both app servers first:

```bash
cd backend
go run .
```

In another terminal:

```bash
cd frontend
npm start
```

Then run Cypress:

```bash
cd frontend
npm run cypress:run
```

### All Local Checks

```bash
cd backend
GOCACHE=/private/tmp/cardflex-go-build-cache go test -v ./...

cd ../frontend
npm test
npm run build
npm run cypress:run
```

## Sprint Documentation

- `Sprint1.md`
- `Sprint2.md`
- `Sprint3.md`
- `Sprint4.md`

## Backend API Reference

### Tenant Resolution

- `POST /register` resolves the tenant from `tenantId` or `companyCode` in the JSON body.
- `POST /register` also supports the legacy `?company=<company-code>` query parameter as a fallback.
- `POST /login`, `GET /dashboard`, `GET /profile`, and `POST /payment` require `?company=<company-code>` in the URL.

### Authentication

- `POST /register` does not require a JWT.
- `POST /login` does not require a JWT.
- `GET /dashboard`, `GET /profile`, and `POST /payment` require both `?company=<company-code>` and an `Authorization: Bearer <jwt>` header.
- Successful login returns a JWT containing `userId` and `tenantId` claims.
- A JWT can only access data for the tenant it was issued for.

### `POST /register`

Creates a new user under a specific tenant.
Registration also creates a default account for the new user so the dashboard can load immediately.

Request body:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "tenantId": "acme"
}
```

Request fields:

- `name`: required string, minimum 1 non-space character after trimming
- `email`: required valid email address
- `password`: required string, minimum 6 characters
- `tenantId`: optional string tenant/company code
- `companyCode`: optional string alternative to `tenantId`

Success response `200 OK`:

```json
{
  "message": "user registered",
  "dummy": {
    "name": "John Doe",
    "email": "john@example.com",
    "userId": 3,
    "tenantId": 1,
    "companyCode": "acme"
  }
}
```

Common errors:

- `400 Bad Request`: invalid JSON, invalid email, short password, or missing tenant identifier
- `404 Not Found`: tenant not found
- `409 Conflict`: email already exists for this tenant
- `500 Internal Server Error`: password hashing or persistence failure

### `POST /login?company=<company-code>`

Authenticates a user for a tenant and returns a JWT.
If an existing user is missing an account, login initializes a default account before returning the token.

Example request:

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

Success response `200 OK`:

```json
{
  "message": "user logged in",
  "token": "<jwt>"
}
```

Common errors:

- `400 Bad Request`: invalid email or short password
- `401 Unauthorized`: invalid credentials
- `404 Not Found`: tenant not found
- `500 Internal Server Error`: account initialization or token generation failure

### `GET /dashboard?company=<company-code>`

Returns tenant-branded dashboard data for the authenticated user.

Required header:

```text
Authorization: Bearer <jwt>
```

Success response `200 OK`:

```json
{
  "tenant": {
    "name": "Acme Card",
    "companyCode": "acme",
    "themeColor": "#0B6E4F"
  },
  "features": {
    "paymentsEnabled": true,
    "profileEnabled": true
  },
  "card": {
    "maskedCardNumber": "**** **** **** 4821",
    "creditLimit": 12000,
    "availableBalance": 8250,
    "currency": "USD"
  },
  "transactions": [
    {
      "date": "2026-02-14",
      "merchant": "Grocery Mart",
      "amount": -82.41,
      "status": "Posted"
    }
  ]
}
```

Common errors:

- `400 Bad Request`: missing `company` query parameter
- `401 Unauthorized`: missing or invalid JWT
- `403 Forbidden`: token tenant mismatch
- `404 Not Found`: tenant or account not found

### `GET /profile?company=<company-code>`

Returns the authenticated user's profile for the active tenant.

Required header:

```text
Authorization: Bearer <jwt>
```

Success response `200 OK`:

```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

Common errors:

- `400 Bad Request`: missing `company` query parameter
- `401 Unauthorized`: missing or invalid JWT
- `403 Forbidden`: token tenant mismatch
- `404 Not Found`: tenant or user not found
- `500 Internal Server Error`: database is unavailable or the user cannot be loaded

### `POST /payment?company=<company-code>`

Records a payment for the authenticated user's tenant account.
The tenant must have `payments_enabled` set to `true`; disabled payment tenants receive `403 Forbidden`.

Required header:

```text
Authorization: Bearer <jwt>
```

Request body:

```json
{
  "amount": 250
}
```

Request fields:

- `amount`: required positive number that cannot exceed the account available balance

Success response `200 OK`:

```json
{
  "message": "payment recorded successfully",
  "updatedBalance": 4750,
  "transactionId": 12,
  "amount": 250,
  "timestamp": "2026-04-29T14:30:00Z"
}
```

Common errors:

- `400 Bad Request`: missing `company`, invalid JSON body, non-positive amount, or amount exceeds available balance
- `401 Unauthorized`: missing or invalid JWT
- `403 Forbidden`: token tenant mismatch or tenant payments are disabled
- `404 Not Found`: tenant or account not found
- `500 Internal Server Error`: database update, transaction creation, or commit failure

### Team: Resolvers

### Members

- Pooja Aslesha Kunchepu (Back-end)
- Ashrita Yanala (Back-end)
- Harika Mekala (Front-end)
- Venkata Ratna Chandu Gembali (Front-end)
