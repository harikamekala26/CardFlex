import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { provideRouter, Router, RouterLink } from '@angular/router';
import { By } from '@angular/platform-browser';
import { of, throwError } from 'rxjs';

import { PaymentComponent } from './payment.component';
import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

describe('PaymentComponent', () => {
  let fixture: ComponentFixture<PaymentComponent>;
  let component: PaymentComponent;
  let authService: {
    makePayment: (amount: number, companyCode: string) => unknown;
  };
  let tenantService: {
    getCompanyCode: () => string | null;
  };
  let router: Router;

  beforeEach(async () => {
    authService = {
      makePayment: () =>
        of({
          message: 'payment recorded successfully',
          updatedBalance: 8000,
          transactionId: 12,
          amount: 250,
          timestamp: '2026-04-29T10:00:00Z'
        })
    };
    tenantService = {
      getCompanyCode: () => 'acme'
    };

    await TestBed.configureTestingModule({
      imports: [PaymentComponent],
      providers: [
        provideRouter([]),
        { provide: AuthService, useValue: authService },
        { provide: TenantService, useValue: tenantService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(PaymentComponent);
    component = fixture.componentInstance;
    router = TestBed.inject(Router);
  });

  it('renders the payment form for the active tenant', () => {
    fixture.detectChanges();

    expect(fixture.nativeElement.textContent).toContain('Submit a card payment');
    expect(fixture.nativeElement.textContent).toContain('Paying for tenant');
    expect(fixture.nativeElement.textContent).toContain('acme');
    expect(fixture.debugElement.query(By.css('input[formControlName="amount"]'))).toBeTruthy();
    expect(fixture.debugElement.query(By.css('button[type="submit"]')).nativeElement.textContent).toContain('Submit Payment');
  });

  it('shows an inline validation error before calling the API for a missing amount', () => {
    const makePaymentSpy = spyOn(authService, 'makePayment').and.callThrough();

    fixture.detectChanges();
    component.onSubmit();
    fixture.detectChanges();

    expect(makePaymentSpy).not.toHaveBeenCalled();
    expect(fixture.nativeElement.textContent).toContain('Payment amount is required.');
  });

  it('shows an inline validation error before calling the API for an invalid amount', () => {
    const makePaymentSpy = spyOn(authService, 'makePayment').and.callThrough();

    fixture.detectChanges();
    component.form.patchValue({ amount: 0 });
    component.onSubmit();
    fixture.detectChanges();

    expect(makePaymentSpy).not.toHaveBeenCalled();
    expect(fixture.nativeElement.textContent).toContain('Payment amount must be greater than $0.00.');
  });

  it('submits a payment and redirects back to the tenant dashboard', fakeAsync(() => {
    const makePaymentSpy = spyOn(authService, 'makePayment').and.callThrough();
    const navigateSpy = spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));

    fixture.detectChanges();
    component.form.patchValue({ amount: 250 });
    component.onSubmit();
    tick(900);

    expect(makePaymentSpy).toHaveBeenCalledWith(250, 'acme');
    expect(component.successMessage).toBe('payment recorded successfully');
    expect(navigateSpy).toHaveBeenCalledWith(['/dashboard'], {
      queryParams: { company: 'acme' }
    });
  }));

  it('shows server errors without crashing', () => {
    spyOn(authService, 'makePayment').and.returnValue(throwError(() => new Error('payment amount exceeds available balance')));

    fixture.detectChanges();
    component.form.patchValue({ amount: 3000 });
    component.onSubmit();
    fixture.detectChanges();

    expect(component.submitting).toBeFalse();
    expect(fixture.nativeElement.textContent).toContain('payment amount exceeds available balance');
  });

  it('preserves the company query parameter on the dashboard link', () => {
    fixture.detectChanges();

    const dashboardLink = fixture.debugElement.query(By.css('.dashboard-link'));

    expect(dashboardLink.injector.get(RouterLink).queryParams).toEqual({ company: 'acme' });
  });
});
