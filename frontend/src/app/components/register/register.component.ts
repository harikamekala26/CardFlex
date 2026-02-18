import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {
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
      name: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  onSubmit(): void {
    this.errorMessage = '';
    this.successMessage = '';

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
    this.authService.register(this.form.getRawValue() as { name: string; email: string; password: string }, company).subscribe({
      next: () => {
        this.submitting = false;
        this.successMessage = 'Registration successful. Please log in.';
        this.router.navigate(['/login'], { queryParams: { company } });
      },
      error: (err) => {
        this.submitting = false;
        this.errorMessage = err.error?.error ?? 'Registration failed';
      }
    });
  }
}
