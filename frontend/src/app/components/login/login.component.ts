import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

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
  readonly form;

  constructor(
    private readonly fb: FormBuilder,
    private readonly authService: AuthService,
    private readonly tenantService: TenantService,
    private readonly router: Router
  ) {
    this.form = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required]]
    });
  }

  onSubmit(): void {
    this.errorMessage = '';

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const company = this.tenantService.getCompanyCode();
    if (!company) {
      this.errorMessage = 'Missing tenant company in URL (?company=...)';
      return;
    }

    this.submitting = true;
    this.authService.login(this.form.getRawValue() as { email: string; password: string }, company).subscribe({
      next: ({ token }) => {
        this.submitting = false;
        this.authService.setToken(token);
        this.router.navigate(['/dashboard'], { queryParams: { company } });
      },
      error: (err) => {
        this.submitting = false;
        this.errorMessage = err.error?.error ?? 'Login failed';
      }
    });
  }
}
