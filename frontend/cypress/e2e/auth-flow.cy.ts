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

  it('shows an API error when tenant registration fails', () => {
    cy.intercept('POST', '**/register', {
      statusCode: 409,
      body: { error: 'email already exists for this tenant' }
    }).as('registerRequest');

    cy.visit('/register?company=wells-fargo');

    cy.get('input[formcontrolname="name"]').type('Jane Doe');
    cy.get('input[formcontrolname="email"]').type('jane@example.com');
    cy.get('input[formcontrolname="password"]').type('secret123');
    cy.contains('button', 'Create Account').click();

    cy.wait('@registerRequest');
    cy.contains('email already exists for this tenant').should('be.visible');
    cy.location('pathname').should('eq', '/register');
    cy.location('search').should('include', 'company=wells-fargo');
  });

  it('redirects protected dashboard routes to tenant-aware login when no session exists', () => {
    cy.visit('/dashboard?company=chase-bank');

    cy.location('pathname').should('eq', '/login');
    cy.location('search').should('include', 'company=chase-bank');
    cy.contains('Signing in for tenant').should('contain.text', 'chase-bank');
  });

  it('clears the tenant session and keeps tenant context on logout', () => {
    cy.intercept('GET', 'http://localhost:8080/dashboard?company=capital-one', {
      statusCode: 200,
      body: {
        tenant: {
          name: 'Capital One',
          companyCode: 'capital-one',
          themeColor: '#003B95'
        },
        accountSummary: {
          maskedCardNumber: '**** **** **** 4821',
          creditLimit: 12000,
          availableBalance: 8250,
          currency: 'USD'
        },
        transactions: []
      }
    }).as('dashboardRequest');

    cy.visit('/login?company=capital-one', {
      onBeforeLoad(window) {
        window.localStorage.setItem('cardflex_sessions', JSON.stringify({ 'capital-one': 'mock-jwt' }));
      }
    });

    cy.location('pathname').should('eq', '/dashboard');
    cy.wait('@dashboardRequest');
    cy.contains('button', 'Logout').click();

    cy.location('pathname').should('eq', '/login');
    cy.location('search').should('include', 'company=capital-one');
    cy.window()
      .its('localStorage.cardflex_sessions')
      .then((storedSessions) => {
        expect(JSON.parse(storedSessions as string)).to.deep.equal({});
    });
    cy.contains('Signing in for tenant').should('contain.text', 'capital-one');
  });

  it('completes a full payment flow from login through confirmation', () => {
    let dashboardRequests = 0;

    cy.intercept('POST', 'http://localhost:8080/login?company=chase-bank', {
      statusCode: 200,
      body: { token: 'payment-flow-jwt', message: 'user logged in' }
    }).as('loginRequest');

    cy.intercept('GET', 'http://localhost:8080/dashboard?company=chase-bank', (req) => {
      dashboardRequests += 1;
      req.reply({
        statusCode: 200,
        body: {
          tenant: {
            name: 'Chase Bank',
            companyCode: 'chase-bank',
            themeColor: '#0A2A66'
          },
          accountSummary: {
            maskedCardNumber: '**** **** **** 4821',
            creditLimit: 12000,
            availableBalance: dashboardRequests === 1 ? 8250 : 8500,
            currency: 'USD'
          },
          features: {
            paymentsEnabled: true,
            profileEnabled: true
          },
          transactions: [
            { date: '2026-02-14', merchant: 'Grocery Mart', amount: -82.41, status: 'Posted' },
            { date: '2026-04-29', merchant: 'Payment', amount: 250, status: 'Posted' }
          ]
        }
      });
    }).as('dashboardRequest');

    cy.intercept('POST', 'http://localhost:8080/payment?company=chase-bank', {
      statusCode: 200,
      body: {
        message: 'payment recorded successfully',
        updatedBalance: 8500,
        transactionId: 42,
        amount: 250,
        timestamp: '2026-04-29T10:00:00Z'
      }
    }).as('paymentRequest');

    cy.visit('/login?company=chase-bank');

    cy.get('input[formcontrolname="email"]').type('demo+chase-bank@cardflex.local');
    cy.get('input[formcontrolname="password"]').type('secret123');
    cy.contains('button', 'Sign In').click();

    cy.wait('@loginRequest');
    cy.wait('@dashboardRequest')
      .its('request.headers.authorization')
      .should('eq', 'Bearer payment-flow-jwt');

    cy.contains('Chase Bank Dashboard').should('be.visible');
    cy.contains('a', 'Make a Payment').click();

    cy.location('pathname').should('eq', '/payment');
    cy.location('search').should('include', 'company=chase-bank');
    cy.get('input[formcontrolname="amount"]').type('250');
    cy.contains('button', 'Submit Payment').click();

    cy.wait('@paymentRequest').then(({ request }) => {
      expect(request.headers.authorization).to.eq('Bearer payment-flow-jwt');
      expect(request.body).to.deep.equal({ amount: 250 });
    });
    cy.contains('payment recorded successfully').should('be.visible');
    cy.location('pathname', { timeout: 3000 }).should('eq', '/dashboard');
    cy.wait('@dashboardRequest');
    cy.contains('Payment').should('be.visible');
  });

  it('navigates to the profile page and displays authenticated user data', () => {
    cy.intercept('GET', 'http://localhost:8080/dashboard?company=chase-bank', {
      statusCode: 200,
      body: {
        tenant: {
          name: 'Chase Bank',
          companyCode: 'chase-bank',
          themeColor: '#0A2A66'
        },
        accountSummary: {
          maskedCardNumber: '**** **** **** 4821',
          creditLimit: 12000,
          availableBalance: 8250,
          currency: 'USD'
        },
        features: {
          paymentsEnabled: true,
          profileEnabled: true
        },
        transactions: []
      }
    }).as('dashboardRequest');

    cy.intercept('GET', 'http://localhost:8080/profile?company=chase-bank', {
      statusCode: 200,
      body: {
        name: 'Jane Profile',
        email: 'jane.profile@example.com'
      }
    }).as('profileRequest');

    cy.visit('/?company=chase-bank', {
      onBeforeLoad(window) {
        window.localStorage.setItem('cardflex_sessions', JSON.stringify({ 'chase-bank': 'profile-flow-jwt' }));
      }
    });

    cy.contains('a', 'Dashboard').click();
    cy.wait('@dashboardRequest');
    cy.contains('a', 'Profile').should('have.attr', 'href').and('include', 'company=chase-bank');
    cy.contains('a', 'Profile').click();

    cy.location('pathname').should('eq', '/profile');
    cy.location('search').should('include', 'company=chase-bank');
    cy.wait('@profileRequest')
      .its('request.headers.authorization')
      .should('eq', 'Bearer profile-flow-jwt');
    cy.contains('Jane Profile').should('be.visible');
    cy.contains('jane.profile@example.com').should('be.visible');
    cy.contains('Chase Bank').should('be.visible');
  });

  it('hides and blocks profile access when the tenant profile feature is disabled', () => {
    cy.intercept('GET', 'http://localhost:8080/profile?company=capital-one', {
      statusCode: 403,
      body: { error: 'profiles are disabled for this tenant' }
    }).as('profileRequest');

    cy.visit('/?company=capital-one', {
      onBeforeLoad(window) {
        window.localStorage.setItem('cardflex_sessions', JSON.stringify({ 'capital-one': 'profile-disabled-jwt' }));
      }
    });

    cy.contains('a', 'Dashboard').should('be.visible');
    cy.contains('a', 'Profile').should('not.exist');

    cy.visit('/profile?company=capital-one', {
      onBeforeLoad(window) {
        window.localStorage.setItem('cardflex_sessions', JSON.stringify({ 'capital-one': 'profile-disabled-jwt' }));
      }
    });

    cy.wait('@profileRequest')
      .its('request.headers.authorization')
      .should('eq', 'Bearer profile-disabled-jwt');
    cy.contains('profiles are disabled for this tenant').should('be.visible');
    cy.location('pathname').should('eq', '/profile');
  });
});
