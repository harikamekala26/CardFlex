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

### Example Usage

```
https://cardflex.com?company=chase-bank
https://cardflex.com?company=capital-one
https://cardflex.com?company=citibank
```

## Run Locally

### Prerequisites

- Node.js and npm
- Go (1.22+ recommended)
- SQL database (SQLite is preconfigured)

### 1) Start Backend

```bash
cd backend
cp .env.example .env
# Fill DB_DRIVER, DB_DSN, and JWT_SECRET in .env
go run .
```

Backend runs on `http://localhost:8080` (health: `GET /ping`).

### 2) Start Frontend

```bash
cd frontend
npm install
npm start
```

Frontend runs on `http://localhost:4200`.

### 3) Open Tenant URLs

```text
http://localhost:4200/?company=chase-bank
http://localhost:4200/?company=capital-one
http://localhost:4200/?company=wells-fargo
```

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
