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
http://localhost:4200/?company=citibank
```

### Team: Resolvers

### Members

- Pooja Aslesha Kunchepu (Back-end)
- Ashrita Yanala (Back-end)
- Harika Mekala (Front-end)
- Venkata Ratna Chandu Gembali (Front-end)
