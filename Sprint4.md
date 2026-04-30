# Sprint 4 Report
Project: CardFlex - Multi-Tenant Credit Card Account Management Portal  
Team: Resolvers

## Sprint 4 Overview

Sprint 4 focused on completing the remaining account-management workflows, strengthening tenant-specific feature control, expanding backend unit test coverage, and updating documentation for the final project submission.

The main Sprint 4 outcomes were:

- frontend support for tenant-controlled feature availability
- payment UI and dashboard feature gating
- protected backend payment and profile endpoints
- account initialization for registered users
- expanded frontend, Cypress, and backend test coverage
- updated backend API documentation in `README.md`

## Work Completed In Sprint 4

### Frontend Work Completed

- Added dashboard support for the backend `features` response.
- Hid the dashboard payment action when `paymentsEnabled` is false for the tenant.
- Maintained payment component support for submitting payments through the backend API.
- Added and maintained frontend unit tests for dashboard, payment, layout, auth, tenant services, auth interceptor, login, register, and home components.
- Maintained Cypress coverage for tenant-aware authentication and protected-route flows.

### Backend Work Completed

- Added tenant feature-flag support using a JSON-backed `features` field on tenant records.
- Seeded example feature configurations for `chase-bank`, `capital-one`, and `wells-fargo`.
- Updated `GET /dashboard` to return a `features` map so the frontend can show or hide feature-specific UI.
- Added payment feature gating so `POST /payment` returns `403 Forbidden` when `payments_enabled` is false for the tenant.
- Implemented account initialization for newly registered users so the dashboard can load immediately after registration.
- Added login-time account backfill for existing users who were created before account initialization existed.
- Added and expanded backend unit tests for payment, profile, dashboard, auth, middleware, and migrations.
- Updated the root `README.md` Backend API Reference with Sprint 4 endpoints and error documentation.

## Sprint 4 Feature Summary

### Frontend Features

#### Tenant-Aware Experience

The frontend continues to preserve tenant context with the `company` query parameter. Dashboard and payment navigation use the active tenant company so users stay inside the correct issuer workspace.

The dashboard consumes the backend feature map and uses it to conditionally show payment access. For example, `wells-fargo` hides the payment action because payments are disabled for that tenant.

#### Payment UI Gating

The dashboard payment link is displayed only when the active tenant has `paymentsEnabled` set to true. This prevents users from entering the payment workflow when the backend has disabled payments for that issuer.

#### Payment Flow UI

The payment page supports submitting a card payment for the authenticated tenant account. It handles amount input, validation messaging, submission state, success messaging, and backend error display.

### Backend Features

#### Tenant Feature Flags

Tenant records now support boolean feature flags. The backend stores flags in snake_case, for example:

```json
{
  "payments_enabled": true,
  "profile_enabled": true
}
```

The dashboard API returns a frontend-friendly camelCase map:

```json
{
  "features": {
    "paymentsEnabled": true,
    "profileEnabled": true
  }
}
```

Seeded Sprint 4 examples:

- `chase-bank`: payments enabled, profile enabled
- `capital-one`: payments enabled, profile disabled
- `wells-fargo`: payments disabled, profile enabled

#### Payment Endpoint

The protected `POST /payment` endpoint records a card payment for the authenticated user and tenant. It validates:

- tenant context
- JWT authentication
- token tenant match
- `payments_enabled` feature flag
- valid positive amount
- available account balance
- tenant-scoped account lookup

Successful payments update the account balance and create a transaction record.

#### Profile Endpoint

The protected `GET /profile` endpoint returns the authenticated user's profile for the active tenant:

```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

The endpoint validates tenant context, JWT authentication, token tenant match, and user existence.

#### Account Initialization

Registration now creates a default account for each new user:

- masked card number: `**** **** **** 0000`
- credit limit: `5000`
- available balance: `5000`
- currency: `USD`

Login also backfills this default account for existing users who do not already have an account.

## Frontend Unit Tests

Frontend unit test files included in the Sprint 4 codebase:

- `frontend/src/app/services/auth.service.spec.ts`
- `frontend/src/app/services/tenant.service.spec.ts`
- `frontend/src/app/interceptors/auth.interceptor.spec.ts`
- `frontend/src/app/components/login/login.component.spec.ts`
- `frontend/src/app/components/register/register.component.spec.ts`
- `frontend/src/app/components/home/home.component.spec.ts`
- `frontend/src/app/components/layout/layout.component.spec.ts`
- `frontend/src/app/components/dashboard/dashboard.component.spec.ts`
- `frontend/src/app/components/payment/payment.component.spec.ts`

Important Sprint 4 frontend coverage includes:

- dashboard rendering with backend account data
- dashboard loading, empty, and error states
- tenant-aware dashboard and payment links
- hidden payment action when tenant payments are disabled
- payment form rendering
- payment validation and submission
- payment error display
- tenant-aware auth service request behavior
- auth interceptor tenant session behavior
- layout navigation and logout behavior

Run frontend unit tests:

```bash
cd frontend
npm install
npm test
```

Note: the local checkout used during Sprint 4 had incomplete `node_modules`, so frontend test execution requires `npm install` before running `npm test`.

## Cypress Tests

Cypress test file:

- `frontend/cypress/e2e/auth-flow.cy.ts`

Cypress coverage includes:

- tenant-aware home page navigation
- successful registration flow
- successful login and dashboard navigation
- failed registration error flow
- protected dashboard redirect to login when no session exists
- logout flow with preserved tenant context and cleared tenant session

Run Cypress tests:

```bash
cd frontend
npm run cypress:run
```

For Cypress execution, run the backend at `http://localhost:8080` and the frontend at `http://localhost:4200`.

## Backend Unit Tests

Backend unit tests were expanded for Sprint 4. Run all backend tests with:

```bash
cd backend
GOCACHE=/private/tmp/cardflex-go-build-cache go test -v ./...
```

Expected package results:

```text
ok      cardflex-backend/controllers
ok      cardflex-backend/middleware
ok      cardflex-backend/migrations
```

### Auth Controller Tests

File:

- `backend/controllers/auth_controller_test.go`

Tests:

- `TestRegisterReturnsDummyResponseForTenant`
- `TestRegisterSupportsCompanyQueryFallback`
- `TestRegisterRejectsDuplicateEmailForTenant`
- `TestRegisterRequiresTenantIdentifier`
- `TestLoginReturnsDummyResponseForTenant`
- `TestLoginValidationRejectsInvalidInputs`

Sprint 4 additions include verification that registration creates a default account and login backfills a missing account.

### Dashboard Controller Tests

File:

- `backend/controllers/dashboard_controller_test.go`

Tests:

- `TestGetDashboardReturnsTenantScopedData`
- `TestGetDashboardReturnsEmptyTransactionsWhenNoTransactionsExist`
- `TestGetDashboardReturnsNotFoundWhenAccountMissing`
- `TestGetDashboardReturnsUnauthorizedWhenClaimsMissing`
- `TestGetDashboardReturnsUnauthorizedWhenTenantMissing`

Sprint 4 additions include verification that dashboard responses include feature flags such as `paymentsEnabled` and `profileEnabled`.

### Payment Controller Tests

File:

- `backend/controllers/payment_controller_test.go`

Tests:

- `TestRecordPayment_Success`
- `TestRecordPayment_InvalidAmount`
- `TestRecordPayment_InvalidJSONBody`
- `TestRecordPayment_ExceedsBalance`
- `TestRecordPayment_AccountNotFound`
- `TestRecordPayment_MissingTenant`
- `TestRecordPayment_MissingClaims`
- `TestRecordPayment_MissingJWT`
- `TestRecordPayment_RejectsTenantMismatch`
- `TestRecordPayment_PaymentsDisabled`

Coverage:

- successful payment response
- account balance update
- transaction record creation
- invalid amount validation
- malformed JSON handling
- amount greater than available balance
- missing account handling
- missing tenant context
- missing authentication claims
- missing JWT through middleware
- tenant mismatch through middleware
- disabled payment feature flag
- no balance mutation or transaction creation on failed payments

### Profile Controller Tests

File:

- `backend/controllers/profile_controller_test.go`

Tests:

- `TestGetProfileReturnsAuthenticatedUser`
- `TestGetProfileRequiresCompanyQuery`
- `TestGetProfileRequiresValidJWT`
- `TestGetProfileRejectsInvalidJWT`
- `TestGetProfileRejectsTenantMismatch`
- `TestGetProfileReturnsNotFoundWhenTenantMissing`
- `TestGetProfileReturnsNotFoundWhenUserMissing`
- `TestGetProfileRejectsInvalidUserIDClaim`

Coverage:

- successful authenticated profile response
- response shape limited to `name` and `email`
- missing company query parameter
- missing JWT
- invalid JWT
- token tenant mismatch
- missing tenant
- missing user
- malformed user ID claim

### Middleware Tests

File:

- `backend/middleware/auth_middleware_test.go`

Tests:

- `TestJWTAuthAllowsProtectedRouteWithValidToken`
- `TestJWTAuthRejectsExpiredToken`
- `TestJWTAuthRejectsTokenForDifferentTenant`
- `TestJWTAuthRejectsMissingAuthorizationHeader`

### Migration Tests

Files:

- `backend/migrations/account_migration_test.go`
- `backend/migrations/tenant_migration_test.go`
- `backend/migrations/transaction_migration_test.go`

Tests:

- `TestMigrateAccountsCreatesTenantScopedAccountTable`
- `TestMigrateTenantsSeedsFeatureFlags`
- `TestMigrateTransactionsSeedsSampleDashboardData`

## Frontend Demo Notes

Seeded tenant URLs used for frontend walkthroughs:

```text
http://localhost:4200/login?company=chase-bank
http://localhost:4200/login?company=wells-fargo
http://localhost:4200/login?company=capital-one
```

Frontend demo behavior:

- `chase-bank` shows payment access.
- `wells-fargo` hides payment access because payments are disabled.
- `capital-one` demonstrates a different seeded feature configuration.

## Backend API Documentation Updates

The root `README.md` Backend API Reference was updated for Sprint 4.

Documented endpoints now include:

- `POST /register`
- `POST /login?company=<company-code>`
- `GET /dashboard?company=<company-code>`
- `GET /profile?company=<company-code>`
- `POST /payment?company=<company-code>`

Sprint 4 documentation updates include:

- required `Authorization: Bearer <jwt>` headers for protected endpoints
- dashboard `features` response shape
- profile request and response documentation
- payment request and response documentation
- payment feature-flag gating behavior
- common error codes for profile and payment
- registration default account creation
- login account backfill behavior

## Backend Demo Notes

Seeded demo users are created for tenant walkthroughs with this password:

```text
secret123
```

Examples:

- `demo+chase-bank@cardflex.local`
- `demo+wells-fargo@cardflex.local`
- `demo+capital-one@cardflex.local`

## Sprint 4 Summary

Sprint 4 completed the remaining backend and frontend work needed for a stronger final CardFlex workflow. The project now supports tenant-controlled feature availability, protected payment and profile endpoints, account initialization for registered users, detailed backend unit tests, Cypress and frontend unit test coverage, and updated API documentation for final project presentation.
