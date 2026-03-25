import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AbstractControl, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  submitting = false;
  errorMessage = '';
  successMessage = '';
  readonly form;

  constructor(
    private readonly fb: FormBuilder,
    private readonly authService: AuthService,
    private readonly tenantService: TenantService,
    private readonly router: Router
  ) {
    this.form = this.fb.group({
      email: ['', [Validators.required, Validators.email, Validators.pattern(EMAIL_PATTERN)]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  get emailControl(): AbstractControl | null {
    return this.form.get('email');
  }

  get passwordControl(): AbstractControl | null {
    return this.form.get('password');
  }

  get companyCode(): string | null {
    return this.tenantService.getCompanyCode();
  }

  showFieldError(control: AbstractControl | null): boolean {
    return !!control && control.invalid && (control.touched || control.dirty);
  }

  getFieldError(controlName: 'email' | 'password'): string {
    const control = this.form.get(controlName);

    if (!control?.errors || !this.showFieldError(control)) {
      return '';
    }

    if (control.errors['required']) {
      return controlName === 'email' ? 'Email is required.' : 'Password is required.';
    }

    if (control.errors['email'] || control.errors['pattern']) {
      return 'Enter a valid email address.';
    }

    if (control.errors['minlength']) {
      return 'Password must be at least 6 characters.';
    }

    return 'Please review this field.';
  }

  onSubmit(): void {
    this.errorMessage = '';
    this.successMessage = '';

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      this.errorMessage = 'Enter a valid email and a password with at least 6 characters.';
      return;
    }

    const company = this.companyCode;
    if (!company) {
      this.errorMessage = 'Missing tenant company in URL (?company=...)';
      return;
    }

    this.submitting = true;
    this.authService.login(this.form.getRawValue() as { email: string; password: string }, company).subscribe({
      next: ({ token, message }) => {
        this.submitting = false;
        if (!token) {
          this.errorMessage = 'Login succeeded but no token was returned.';
          return;
        }

        this.authService.setToken(token);
        this.successMessage = message || 'User logged in successfully.';
        this.form.reset();
        this.form.markAsPristine();
        this.form.markAsUntouched();
        void this.router.navigate(['/dashboard'], {
          queryParams: { company }
        });
      },
      error: (err) => {
        this.submitting = false;
        this.errorMessage = err.error?.error ?? 'Login failed';
      }
    });
  }
}
