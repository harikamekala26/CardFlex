# Sprint 2 Report
Project: CardFlex - Multi-Tenant Credit Card Account Management Portal  
Team: Resolvers

## Sprint 2 Overview

Sprint 2 focused on completing the integration between the Angular frontend and the Go backend, adding automated tests for both layers, and documenting the backend API for submission.

All major Sprint 2 requirements were completed:

- frontend and backend integration
- frontend unit tests
- Cypress end-to-end test
- backend unit tests
- backend API documentation in this file

## Work Completed In Sprint 2

### Entire Team

- Integrated the Angular frontend with the Go backend.
- Extended Sprint 1 foundation work into a working authenticated application flow.
- Prepared the application for demonstration with successful frontend and backend test execution.

### Frontend Work Completed

- Implemented tenant-aware registration flow connected to the backend.
- Implemented tenant-aware login flow connected to the backend.
- Implemented authenticated dashboard data loading from the backend.
- Preserved tenant context in the URL using the `company` query parameter.
- Added Angular unit tests for services and components.
- Added a Cypress end-to-end test for the frontend flow.

### Backend Work Completed

- Implemented tenant-aware registration API.
- Implemented tenant-aware login API.
- Implemented JWT authentication middleware for protected routes.
- Implemented a protected dashboard API that returns tenant-specific data.
- Added backend unit tests for controller logic and authentication middleware.
- Documented the backend API.

## Progress From Sprint 1

Sprint 1 established the project foundation, including:

- Angular project setup
- frontend layout and tenant detection
- Go server setup
- database connection
- tenant and user models

Sprint 2 built on that foundation by adding:

- full frontend-backend integration
- registration and login functionality
- JWT-based protection for private backend routes
- dashboard API consumption from the frontend
- automated frontend and backend tests

## Frontend Tests

### Cypress Test

File:

- `frontend/cypress/e2e/auth-flow.cy.ts`

What it covers:

- tenant-aware navigation on the home page
- registration form submission
- login flow and dashboard navigation

Execution result:

- `3 passing`

### Frontend Unit Tests

Files:

- `frontend/src/app/services/auth.service.spec.ts`
- `frontend/src/app/services/tenant.service.spec.ts`
- `frontend/src/app/components/login/login.component.spec.ts`
- `frontend/src/app/components/register/register.component.spec.ts`

What they cover:

- API calls for registration and login
- tenant resolution behavior
- session storage behavior
- login form validation
- registration form validation
- frontend error handling for backend failures

Execution result:

- `TOTAL: 16 SUCCESS`

## Backend Unit Tests

### Controller Tests

File:

- `backend/controllers/auth_controller_test.go`

What it covers:

- successful registration for a valid tenant
- registration using company query parameter fallback
- rejection of duplicate email registration
- rejection when tenant identifier is missing
- successful login response with JWT generation
- rejection of invalid login inputs

Verbose test results:

- `TestRegisterReturnsDummyResponseForTenant` - passed
- `TestRegisterSupportsCompanyQueryFallback` - passed
- `TestRegisterRejectsDuplicateEmailForTenant` - passed
- `TestRegisterRequiresTenantIdentifier` - passed
- `TestLoginReturnsDummyResponseForTenant` - passed
- `TestLoginValidationRejectsInvalidInputs` - passed

File:

- `backend/controllers/dashboard_controller_test.go`

What it covers:

- successful retrieval of tenant-scoped dashboard data for an authenticated user

Verbose test results:

- `TestGetDashboardReturnsTenantScopedData` - passed

### Middleware Tests

File:

- `backend/middleware/auth_middleware_test.go`

What it covers:

- valid JWT allows access to protected route
- expired JWT is rejected
- wrong-tenant token is rejected
- missing authorization header is rejected

Verbose test results:

- `TestJWTAuthAllowsProtectedRouteWithValidToken` - passed
- `TestJWTAuthRejectsExpiredToken` - passed
- `TestJWTAuthRejectsTokenForDifferentTenant` - passed
- `TestJWTAuthRejectsMissingAuthorizationHeader` - passed

Execution result:

- backend controller and middleware tests passed successfully with `go test ./...`

## Backend API Documentation

### Base URL

- Local backend URL: `http://localhost:8080`

### Authentication Summary

- `POST /register` does not require authentication.
- `POST /login?company=<company-code>` does not require authentication.
- `GET /dashboard?company=<company-code>` requires a valid JWT in the `Authorization` header.

### Tenant Resolution

- `POST /register` accepts tenant information from:
  - `tenantId` in the request body
  - `companyCode` in the request body
  - `company` query parameter as a fallback
- `POST /login` and `GET /dashboard` use the `company` query parameter.

### 1. Health Check

#### `GET /ping`

Purpose:

- confirms the backend server is running

Success response:

```text
pong
```

### 2. Register User

#### `POST /register`

Purpose:

- creates a new user under a selected tenant

Request body:

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "secret123",
  "tenantId": "acme"
}
```

Accepted fields:

- `name`: required string
- `email`: required valid email address
- `password`: required string with minimum length 6
- `tenantId`: optional tenant/company code
- `companyCode`: optional alternative tenant/company code

Success response `200 OK`:

```json
{
  "dummy": {
    "name": "Jane Doe",
    "email": "jane@example.com",
    "userId": 3,
    "tenantId": 1,
    "companyCode": "acme"
  },
  "message": "user registered"
}
```

Possible errors:

- `400 Bad Request`
  - invalid JSON
  - invalid email format
  - password too short
  - missing tenant identifier
- `404 Not Found`
  - tenant not found
- `409 Conflict`
  - email already exists for this tenant
- `500 Internal Server Error`
  - password hashing or database failure

### 3. Login User

#### `POST /login?company=<company-code>`

Purpose:

- authenticates a user for a tenant and returns a JWT token

Request body:

```json
{
  "email": "jane@example.com",
  "password": "secret123"
}
```

Success response `200 OK`:

```json
{
  "message": "user logged in",
  "token": "<jwt>"
}
```

Possible errors:

- `400 Bad Request`
  - invalid email format
  - password too short
  - missing `company` query parameter
- `401 Unauthorized`
  - invalid credentials
- `404 Not Found`
  - tenant not found
- `500 Internal Server Error`
  - token generation or database failure

### 4. Get Dashboard

#### `GET /dashboard?company=<company-code>`

Purpose:

- returns tenant-branded dashboard data for an authenticated user

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

Possible errors:

- `400 Bad Request`
  - missing `company` query parameter
- `401 Unauthorized`
  - missing or invalid authorization header
  - invalid token
- `403 Forbidden`
  - token tenant mismatch
- `404 Not Found`
  - tenant not found

## Test Execution Summary

The following automated tests were executed successfully during Sprint 2 verification:

- Backend tests:
  - `go test ./...`
  - `go test -v ./controllers ./middleware`
  - result: all named backend controller and middleware tests passed
- Frontend unit tests:
  - `npm test -- --watch=false`
  - result: `TOTAL: 16 SUCCESS`
- Cypress test:
  - `npm run cypress:run`
  - result: `3 passing`

## Conclusion

Sprint 2 successfully delivered the required integration between frontend and backend, added automated tests for both layers, and documented the backend API. The CardFlex application now supports tenant-aware registration, login, JWT-protected dashboard access, frontend validation, backend validation, and test-backed functionality across the full stack.
