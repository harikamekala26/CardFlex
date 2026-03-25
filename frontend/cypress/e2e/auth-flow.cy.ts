describe('CardFlex auth flow', () => {
  it('shows tenant-aware navigation on the home page', () => {
    cy.visit('/?company=chase-bank');

    cy.contains('Welcome to Chase Bank').should('be.visible');
    cy.contains('Tenant company detected from URL:').should('contain.text', 'chase-bank');
    cy.contains('a', 'Register').should('have.attr', 'href').and('include', 'company=chase-bank');
    cy.contains('a', 'Login').should('have.attr', 'href').and('include', 'company=chase-bank');
  });

  it('submits registration and routes the user to login', () => {
    cy.intercept('POST', '**/register', {
      statusCode: 200,
      body: { message: 'user registered' }
    }).as('registerRequest');

    cy.visit('/register?company=chase-bank');

    cy.get('input[formcontrolname="name"]').type('Jane Doe');
    cy.get('input[formcontrolname="email"]').type('jane@example.com');
    cy.get('input[formcontrolname="password"]').type('secret123');
    cy.contains('button', 'Create Account').click();

    cy.wait('@registerRequest')
      .its('request.body')
      .should('deep.equal', {
        name: 'Jane Doe',
        email: 'jane@example.com',
        password: 'secret123',
        tenantId: 'chase-bank'
      });

    cy.contains('user registered').should('be.visible');
    cy.location('pathname').should('eq', '/login');
    cy.location('search').should('include', 'company=chase-bank');
  });

  it('submits login and navigates to the dashboard', () => {
    cy.intercept('POST', '**/login?company=capital-one', {
      statusCode: 200,
      body: { token: 'mock-jwt', message: 'user logged in' }
    }).as('loginRequest');

    cy.intercept('GET', '**/dashboard?company=capital-one', {
      statusCode: 200,
      body: {
        tenant: {
          name: 'Capital One',
          companyCode: 'capital-one',
          themeColor: '#003B95'
        },
        card: {
          maskedCardNumber: '**** **** **** 4821',
          creditLimit: 12000,
          availableBalance: 8250,
          currency: 'USD'
        },
        transactions: [
          { date: '2026-02-14', merchant: 'Grocery Mart', amount: -82.41, status: 'Posted' }
        ]
      }
    }).as('dashboardRequest');

    cy.visit('/login?company=capital-one');

    cy.get('input[formcontrolname="email"]').type('jane@example.com');
    cy.get('input[formcontrolname="password"]').type('secret123');
    cy.contains('button', 'Sign In').click();

    cy.wait('@loginRequest')
      .its('request.body')
      .should('deep.equal', {
        email: 'jane@example.com',
        password: 'secret123'
      });

    cy.wait('@dashboardRequest');
    cy.location('pathname').should('eq', '/dashboard');
    cy.contains('Capital One Dashboard').should('be.visible');
  });
});
