import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
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
    private readonly tenantService: TenantService
  ) {
    this.form = this.fb.group({
      email: ['', [Validators.required, Validators.email, Validators.pattern(EMAIL_PATTERN)]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  onSubmit(): void {
    this.errorMessage = '';
    this.successMessage = '';

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      this.errorMessage = 'Enter a valid email and a password with at least 6 characters.';
      return;
    }

    const company = this.tenantService.getCompanyCode();
    if (!company) {
      this.errorMessage = 'Missing tenant company in URL (?company=...)';
      return;
    }

    this.submitting = true;
    this.authService.login(this.form.getRawValue() as { email: string; password: string }, company).subscribe({
      next: ({ token, message }) => {
        this.submitting = false;
        if (token) {
          this.authService.setToken(token);
        }
        this.successMessage = message || 'user logged in';
        this.form.reset();
      },
      error: (err) => {
        this.submitting = false;
        this.errorMessage = err.error?.error ?? 'Login failed';
      }
    });
  }
}
