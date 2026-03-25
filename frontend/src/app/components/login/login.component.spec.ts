import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { of, throwError } from 'rxjs';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { LoginComponent } from './login.component';

describe('LoginComponent', () => {
  let fixture: ComponentFixture<LoginComponent>;
  let component: LoginComponent;
  let authService: Partial<AuthService>;
  let tenantService: Partial<TenantService>;
  let router: Router;

  beforeEach(async () => {
    authService = {
      login: () => of({ token: 'mock-token', message: 'user logged in' }),
      setSession: () => undefined,
      isAuthenticated: () => false
    };

    tenantService = {
      getCompanyCode: () => 'capital-one'
    };

    await TestBed.configureTestingModule({
      imports: [LoginComponent],
      providers: [
        provideRouter([]),
        { provide: AuthService, useValue: authService },
        { provide: TenantService, useValue: tenantService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(LoginComponent);
    component = fixture.componentInstance;
    router = TestBed.inject(Router);
    fixture.detectChanges();
  });

  it('shows a validation error when the form is invalid', () => {
    component.onSubmit();

    expect(component.errorMessage).toBe('Enter a valid email and a password with at least 6 characters.');
  });

  it('shows an error when the tenant company is missing', () => {
    spyOn(tenantService, 'getCompanyCode').and.returnValue(null);

    component.form.setValue({
      email: 'jane@example.com',
      password: 'secret123'
    });
    component.onSubmit();

    expect(component.errorMessage).toBe('Missing tenant company in URL (?company=...)');
  });

  it('stores the tenant session and navigates after successful login', () => {
    const setSessionSpy = spyOn(authService, 'setSession');
    const navigateSpy = spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));
    spyOn(authService, 'login').and.returnValue(of({ token: 'jwt-token', message: 'user logged in' }));

    component.form.setValue({
      email: 'jane@example.com',
      password: 'secret123'
    });
    component.onSubmit();

    expect(setSessionSpy).toHaveBeenCalledWith('capital-one', 'jwt-token');
    expect(navigateSpy).toHaveBeenCalledWith(['/dashboard'], {
      queryParams: { company: 'capital-one' }
    });
  });

  it('shows backend login errors returned by the API layer', () => {
    spyOn(authService, 'login').and.returnValue(throwError(() => new Error('invalid credentials')));

    component.form.setValue({
      email: 'jane@example.com',
      password: 'secret123'
    });
    component.onSubmit();

    expect(component.errorMessage).toBe('invalid credentials');
  });
});
