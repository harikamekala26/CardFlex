import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AbstractControl, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';

import { AuthService } from '../../services/auth.service';
import { RegisterRequest } from '../../models/auth.model';
import { TenantService } from '../../services/tenant.service';

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
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
      email: ['', [Validators.required, Validators.email, Validators.pattern(EMAIL_PATTERN)]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  get nameControl(): AbstractControl | null {
    return this.form.get('name');
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

  getFieldError(controlName: 'name' | 'email' | 'password'): string {
    const control = this.form.get(controlName);

    if (!control?.errors || !this.showFieldError(control)) {
      return '';
    }

    if (control.errors['required']) {
      return `${this.getFieldLabel(controlName)} is required.`;
    }

    if (control.errors['minlength']) {
      if (controlName === 'password') {
        return 'Password must be at least 6 characters.';
      }

      return 'Name must be at least 2 characters.';
    }

    if (control.errors['email'] || control.errors['pattern']) {
      return 'Enter a valid email address.';
    }

    return 'Please review this field.';
  }

  onSubmit(): void {
    this.errorMessage = '';
    this.successMessage = '';

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      this.errorMessage = 'Enter a valid name, email, and password (minimum 6 characters).';
      return;
    }

    const company = this.companyCode;
    if (!company) {
      this.errorMessage = 'Missing tenant company in URL (?company=...)';
      return;
    }

    this.submitting = true;
    this.authService.register(this.form.getRawValue() as RegisterRequest, company).subscribe({
      next: (response) => {
        this.submitting = false;
        this.successMessage = response.message || 'User registered successfully. You can sign in now.';
        this.form.reset();
        this.form.markAsPristine();
        this.form.markAsUntouched();
        window.setTimeout(() => {
          void this.router.navigate(['/login'], {
            queryParams: { company: company ?? undefined }
          });
        }, 900);
      },
      error: (err: Error) => {
        this.submitting = false;
        this.errorMessage = err.message || 'Registration failed';
      }
    });
  }

  private getFieldLabel(controlName: 'name' | 'email' | 'password'): string {
    switch (controlName) {
      case 'name':
        return 'Name';
      case 'email':
        return 'Email';
      case 'password':
        return 'Password';
    }
  }
}
