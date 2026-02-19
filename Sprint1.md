# Sprint 1 Report
Project: CardFlex - Multi-Tenant Credit Card Account Management Portal  
Team: Resolvers  

---

## User Stories

1. As a credit card customer, I want the system to detect my issuing company from the URL parameter so that I automatically access the correct branded portal.

2. As a credit card customer, I want the portal interface (logo, theme, branding) to change based on my issuer so that it feels like my bank’s official platform.

3. As the CardFlex platform, we want a backend server running so that we can provide a shared multi-tenant service architecture.

4. As the platform, we want tenant and user data models implemented so that the system can support multiple credit card issuers and their customers.

5. As the platform, we want a database connection configured so that tenant and user information can be stored persistently.

6. As the development team, we want a structured Angular frontend with a reusable base layout so that tenant-specific pages can be built consistently.

---

## Sprint 1 Planned Issues

### Backend
- [Sprint-1] Setup Go server
- [Sprint-1] Setup database connection
- [Sprint-1] Implement Tenant model
- [Sprint-1] Implement User model

### Frontend
- [Sprint-1] Initialize Angular project structure
- [Sprint-1] Create base layout template
- [Sprint-1] Implement tenant detection via URL parameter
- [Sprint-1] Implement tenant theme configuration

---

## Completed Issues

All Sprint-1 issues planned for both frontend and backend were successfully completed.

The backend server is operational, connected to the database, and includes Tenant and User models.  
The frontend correctly detects tenants via URL parameters and dynamically applies tenant-specific themes.

---

## Not Completed & Reasons

All Sprint-1 planned issues were completed. No Sprint-1 tasks were left unfinished.

**Note:**
Since this was our first sprint and we were working with new technologies (Go and Angular), we intentionally kept Sprint-1 focused on foundational setup and core architecture. More complex features such as authentication and API integration were planned for Sprint-2 to ensure realistic sprint planning and stable implementation.
