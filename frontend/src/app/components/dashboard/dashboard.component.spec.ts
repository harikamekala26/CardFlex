import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { provideRouter, Router, RouterLink } from '@angular/router';
import { Observable, of, throwError } from 'rxjs';

import { DashboardComponent } from './dashboard.component';
import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { DashboardApiResponse } from '../../models/dashboard.model';

describe('DashboardComponent', () => {
  let fixture: ComponentFixture<DashboardComponent>;
  let component: DashboardComponent;
  let authService: {
    getDashboard: (companyCode: string) => unknown;
    logout: (companyCode?: string | null) => void;
  };
  let tenantService: {
    getCompanyCode: () => string | null;
  };
  let router: Router;

  const accountSummaryResponse: DashboardApiResponse = {
    tenant: {
      name: 'Acme Card',
      companyCode: 'acme',
      themeColor: '#0B6E4F'
    },
    accountSummary: {
      maskedCardNumber: '**** **** **** 4821',
      creditLimit: 12000,
      availableBalance: 8250,
      currency: 'USD'
    },
    transactions: [
      {
        date: '2026-02-14',
        merchant: 'Grocery Mart',
        amount: -82.41,
        status: 'Posted'
      },
      {
        date: '2026-02-20',
        merchant: 'Coffee Stand',
        amount: -17.59,
        status: 'Posted'
      },
      {
        date: '2026-03-03',
        merchant: 'Payment',
        amount: 50,
        status: 'Posted'
      }
    ]
  };

  beforeEach(async () => {
    authService = {
      getDashboard: () => of(accountSummaryResponse),
      logout: () => undefined
    };
    tenantService = {
      getCompanyCode: () => 'acme'
    };

    await TestBed.configureTestingModule({
      imports: [DashboardComponent],
      providers: [
        provideRouter([]),
        { provide: AuthService, useValue: authService },
        { provide: TenantService, useValue: tenantService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(DashboardComponent);
    component = fixture.componentInstance;
    router = TestBed.inject(Router);
  });

  it('renders the backend account summary response', () => {
    const getDashboardSpy = spyOn(authService, 'getDashboard').and.callThrough();

    fixture.detectChanges();

    expect(getDashboardSpy).toHaveBeenCalledWith('acme');
    expect(component.data?.accountSummary.maskedCardNumber).toBe('**** **** **** 4821');
    expect(component.paymentUsed).toBe(3750);
    expect(component.utilizationPercent).toBe(31);
    expect(fixture.nativeElement.textContent).toContain('Account Summary');
    expect(fixture.nativeElement.textContent).toContain('Acme Card Dashboard');
    expect(fixture.nativeElement.textContent).toContain('**** **** **** 4821');
    expect(fixture.nativeElement.textContent).toContain('Grocery Mart');
  });

  it('links to the payment page with the active tenant company', () => {
    fixture.detectChanges();

    const paymentLink = fixture.debugElement
      .queryAll(By.css('a'))
      .find((element) => element.nativeElement.textContent.includes('Make a Payment'));

    expect(paymentLink?.injector.get(RouterLink).queryParams).toEqual({ company: 'acme' });
  });

  it('renders a monthly spending chart from dashboard transactions', () => {
    fixture.detectChanges();

    expect(component.spendingSummary).toEqual([
      { month: 'Feb 2026', amount: 100, percent: 100 },
      { month: 'Mar 2026', amount: 50, percent: 50 }
    ]);
    expect(fixture.nativeElement.textContent).toContain('Monthly Transaction Volume');
    expect(fixture.nativeElement.textContent).toContain('Feb 2026');
    expect(fixture.nativeElement.textContent).toContain('Mar 2026');
    expect(fixture.debugElement.queryAll(By.css('.spending-bar')).length).toBe(2);
  });

  it('shows a loading state while the dashboard request is in flight', () => {
    spyOn(authService, 'getDashboard').and.returnValue(new Observable());

    fixture.detectChanges();

    expect(component.loading).toBeTrue();
    expect(component.data).toBeNull();
    expect(fixture.nativeElement.textContent).toContain('Preparing your tenant dashboard');
  });

  it('normalizes legacy card payloads so the current layout still works', () => {
    spyOn(authService, 'getDashboard').and.returnValue(
      of({
        tenant: accountSummaryResponse.tenant,
        card: accountSummaryResponse.accountSummary,
        transactions: accountSummaryResponse.transactions
      })
    );

    fixture.detectChanges();

    expect(component.data?.accountSummary.creditLimit).toBe(12000);
    expect(component.data?.transactions).toEqual(accountSummaryResponse.transactions);
    expect(fixture.nativeElement.textContent).toContain('Grocery Mart');
  });

  it('shows an empty state when the backend returns no transactions', () => {
    spyOn(authService, 'getDashboard').and.returnValue(
      of({
        ...accountSummaryResponse,
        transactions: []
      })
    );

    fixture.detectChanges();

    expect(component.hasTransactions).toBeFalse();
    expect(component.spendingSummary).toEqual([]);
    expect(fixture.nativeElement.textContent).toContain('No transactions to display');
    expect(fixture.nativeElement.textContent).toContain('No transactions available');
    expect(fixture.nativeElement.textContent).not.toContain('Grocery Mart');
  });

  it('shows a missing-tenant error before calling the API', () => {
    spyOn(tenantService, 'getCompanyCode').and.returnValue(null);
    const getDashboardSpy = spyOn(authService, 'getDashboard').and.callThrough();

    fixture.detectChanges();

    expect(getDashboardSpy).not.toHaveBeenCalled();
    expect(component.errorMessage).toBe('Missing tenant company in URL (?company=...)');
  });

  it('preserves the tenant company on the return-to-login link for recoverable errors', () => {
    spyOn(authService, 'getDashboard').and.returnValue(throwError(() => new Error('temporarily unavailable')));

    fixture.detectChanges();

    const loginLink = fixture.debugElement
      .queryAll(By.css('a'))
      .find((element) => element.nativeElement.textContent.includes('Return to Login'));

    expect(component.errorMessage).toBe('temporarily unavailable');
    expect(loginLink?.injector.get(RouterLink).queryParams).toEqual({ company: 'acme' });
  });

  it('logs out and routes to login when dashboard access is unauthorized', () => {
    const navigateSpy = spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));
    const getDashboardSpy = spyOn(authService, 'getDashboard').and.returnValue(
      throwError(() => ({ status: 401, message: 'Unauthorized' }))
    );
    const logoutSpy = spyOn(authService, 'logout').and.callThrough();

    fixture.detectChanges();

    expect(getDashboardSpy).toHaveBeenCalledWith('acme');
    expect(logoutSpy).toHaveBeenCalledWith('acme');
    expect(navigateSpy).toHaveBeenCalledWith(['/login'], {
      queryParams: { company: 'acme' }
    });
  });

  it('also logs out and routes to login when dashboard access is forbidden for the tenant', () => {
    const navigateSpy = spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));
    const logoutSpy = spyOn(authService, 'logout').and.callThrough();
    spyOn(authService, 'getDashboard').and.returnValue(
      throwError(() => ({ status: 403, message: 'token tenant mismatch' }))
    );

    fixture.detectChanges();

    expect(logoutSpy).toHaveBeenCalledWith('acme');
    expect(navigateSpy).toHaveBeenCalledWith(['/login'], {
      queryParams: { company: 'acme' }
    });
  });

  it('shows readable API errors from the auth service layer', () => {
    spyOn(authService, 'getDashboard').and.returnValue(throwError(() => new Error('account not found')));

    fixture.detectChanges();

    expect(component.errorMessage).toBe('account not found');
    expect(fixture.nativeElement.textContent).toContain('Dashboard Unavailable');
    expect(fixture.nativeElement.textContent).toContain('account not found');
  });
});
