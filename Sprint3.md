# Sprint 3 Report
Project: CardFlex - Multi-Tenant Credit Card Account Management Portal  
Team: Resolvers

## Sprint 3 Overview

Sprint 3 focused on turning the dashboard experience into a real tenant-aware account workspace backed by database data, then strengthening submission quality with broader frontend and backend test coverage plus sprint documentation.

All 8 Sprint 3 issues were completed:

- 4 backend Sprint 3 issues
- 4 frontend Sprint 3 issues

## Sprint 3 Issues Completed

### Backend Issues

- `#39` `[Backend][Sprint-3] Add Account model and migration`
- `#40` `[Backend][Sprint-3] Add Transaction model and seed sample data`
- `#41` `[Backend][Sprint-3] Replace hardcoded dashboard response with database-backed data`
- `#42` `[Backend][Sprint-3] Add dashboard backend unit tests, validation, and API documentation`

### Frontend Issues

- `#43` `[Frontend][Sprint-3] Refactor dashboard to display real backend account data`
- `#44` `[Frontend][Sprint-3] Add transaction history, loading, empty, and error states`
- `#45` `[Frontend][Sprint-3] Add dashboard and tenant-flow unit tests`
- `#46` `[Frontend][Sprint-3] Expand Cypress coverage, verify test setup, and update Sprint3.md`

## Work Completed In Sprint 3

### Backend Work Completed

- Added the `Account` model and database migration for tenant-scoped card account information.
- Added the `Transaction` model and migration for transaction history tied to tenant accounts.
- Seeded sample account and transaction data for tenant demo flows.
- Replaced the previous hardcoded dashboard response with real database-backed dashboard data.
- Updated the backend dashboard flow to return tenant information, account summary data, and transaction history from persisted records.
- Added backend dashboard unit tests and validation coverage.
- Completed backend API documentation for the Sprint 3 dashboard flow.

### Frontend Work Completed

- Refactored the dashboard to consume real backend account data instead of static placeholder content.
- Added transaction history rendering in the dashboard UI.
- Implemented dashboard loading, empty, and error states for a production-style tenant flow.
- Preserved tenant-aware behavior using the `company` query parameter across dashboard navigation and auth-related routes.
- Added frontend unit tests for dashboard rendering and tenant-aware flow behavior.
- Expanded Cypress coverage beyond the happy path with failed auth, protected-route redirect, and logout coverage.
- Verified frontend unit test and Cypress execution from the project scripts used by the team.
- Fixed Cypress setup reliability by clearing `ELECTRON_RUN_AS_NODE` in the frontend scripts.

## Sprint 3 Dashboard Outcome

Sprint 3 delivered the complete dashboard data flow across both layers:

- backend models and migrations for account and transaction data
- seeded tenant sample data for demo and testing
- database-backed dashboard API responses
- frontend dashboard rendering of account summary and transaction history
- loading, empty, unauthorized, and recoverable error states
- tenant-aware redirects and navigation preservation
- test-backed frontend and backend submission quality

## Frontend Unit Tests

Files added or expanded during Sprint 3:

- `frontend/src/app/components/dashboard/dashboard.component.spec.ts`
  - successful dashboard rendering
  - loading state while dashboard data is in flight
  - legacy dashboard payload normalization
  - empty transaction state
  - missing-tenant error handling
  - tenant-aware return-to-login link
  - unauthorized redirect to login with preserved tenant context
  - forbidden tenant redirect to login with preserved tenant context
  - readable backend error messaging
- `frontend/src/app/components/home/home.component.spec.ts`
  - tenant name and company code rendering
  - tenant-aware register and login links
- `frontend/src/app/components/layout/layout.component.spec.ts`
  - tenant sync from query parameters
  - branding application from resolved tenant data
  - tenant-aware navigation links for unauthenticated users
  - dashboard link visibility for authenticated users
  - tenant-aware logout navigation

Supporting frontend specs also remain in place for:

- `frontend/src/app/services/auth.service.spec.ts`
- `frontend/src/app/services/tenant.service.spec.ts`
- `frontend/src/app/components/login/login.component.spec.ts`
- `frontend/src/app/components/register/register.component.spec.ts`

Frontend unit execution result:

- `31 SUCCESS`

## Cypress Tests

File:

- `frontend/cypress/e2e/auth-flow.cy.ts`

Coverage:

- tenant-aware home page navigation
- successful registration flow
- successful login and dashboard navigation
- failed registration error flow
- protected dashboard redirect to login when no session exists
- logout flow with preserved tenant context and cleared tenant session

Cypress execution result:

- `6 passing`

## Backend Tests

Sprint 3 backend verification included:

- dashboard controller test coverage
- account and transaction migration test coverage
- tenant-aware dashboard validation behavior

Relevant backend test files:

- `backend/controllers/dashboard_controller_test.go`
- `backend/migrations/account_migration_test.go`
- `backend/migrations/transaction_migration_test.go`
- existing authentication and middleware tests remained part of backend verification

## Test Execution

Run backend tests:

```bash
cd backend
go test ./...
```

Run frontend unit tests:

```bash
cd frontend
npm test
```

Run Cypress end-to-end tests:

```bash
cd frontend
npm run cypress:run
```

For Cypress execution, keep the local backend on `http://localhost:8080` and frontend on `http://localhost:4200` running.

## Sprint 3 Summary

Sprint 3 completed all planned backend and frontend dashboard work for this phase. The team delivered persistent account and transaction data, a database-backed dashboard API, a frontend dashboard that renders real tenant data with robust states, stronger unit and end-to-end test coverage, and a submission-ready demo flow for presentation.
