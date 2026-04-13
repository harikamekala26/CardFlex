import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { of, throwError } from 'rxjs';

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
    expect(component.data?.transactions.length).toBe(1);
    expect(fixture.nativeElement.textContent).toContain('Grocery Mart');
  });

  it('shows a missing-tenant error before calling the API', () => {
    spyOn(tenantService, 'getCompanyCode').and.returnValue(null);
    const getDashboardSpy = spyOn(authService, 'getDashboard').and.callThrough();

    fixture.detectChanges();

    expect(getDashboardSpy).not.toHaveBeenCalled();
    expect(component.errorMessage).toBe('Missing tenant company in URL (?company=...)');
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

  it('shows readable API errors from the auth service layer', () => {
    spyOn(authService, 'getDashboard').and.returnValue(throwError(() => new Error('account not found')));

    fixture.detectChanges();

    expect(component.errorMessage).toBe('account not found');
  });
});
