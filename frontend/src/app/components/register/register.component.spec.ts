import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of, throwError } from 'rxjs';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { RegisterComponent } from './register.component';

describe('RegisterComponent', () => {
  let fixture: ComponentFixture<RegisterComponent>;
  let component: RegisterComponent;
  let authService: Partial<AuthService>;
  let tenantService: Partial<TenantService>;

  beforeEach(async () => {
    authService = {
      register: () => of({ message: 'user registered' })
    };

    tenantService = {
      getCompanyCode: () => 'chase-bank'
    };

    await TestBed.configureTestingModule({
      imports: [RegisterComponent],
      providers: [
        provideRouter([]),
        { provide: AuthService, useValue: authService },
        { provide: TenantService, useValue: tenantService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(RegisterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('shows a validation error when the form is invalid', () => {
    component.onSubmit();

    expect(component.errorMessage).toBe('Enter a valid name, email, and password (minimum 6 characters).');
  });

  it('shows an error when the tenant company is missing', () => {
    spyOn(tenantService, 'getCompanyCode').and.returnValue(null);

    component.form.setValue({
      name: 'Jane Doe',
      email: 'jane@example.com',
      password: 'secret123'
    });
    component.onSubmit();

    expect(component.errorMessage).toBe('Missing tenant company in URL (?company=...)');
  });

  it('passes form values to the register API', () => {
    const registerSpy = spyOn(authService, 'register').and.returnValue(of({ message: 'user registered' }));

    component.form.setValue({
      name: 'Jane Doe',
      email: 'jane@example.com',
      password: 'secret123'
    });
    component.onSubmit();

    expect(registerSpy).toHaveBeenCalledWith(
      {
        name: 'Jane Doe',
        email: 'jane@example.com',
        password: 'secret123'
      },
      'chase-bank'
    );
  });

  it('shows backend register errors returned by the API layer', () => {
    spyOn(authService, 'register').and.returnValue(throwError(() => new Error('email already exists for this tenant')));

    component.form.setValue({
      name: 'Jane Doe',
      email: 'jane@example.com',
      password: 'secret123'
    });
    component.onSubmit();

    expect(component.errorMessage).toBe('email already exists for this tenant');
  });
});
