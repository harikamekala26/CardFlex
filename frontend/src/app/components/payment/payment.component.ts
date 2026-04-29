import { CommonModule, CurrencyPipe } from '@angular/common';
import { Component } from '@angular/core';
import { AbstractControl, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

@Component({
  selector: 'app-payment',
  standalone: true,
  imports: [CommonModule, CurrencyPipe, ReactiveFormsModule, RouterLink],
  templateUrl: './payment.component.html',
  styleUrl: './payment.component.css'
})
export class PaymentComponent {
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
      amount: [null as number | null, [Validators.required, Validators.min(0.01)]]
    });
  }

  get amountControl(): AbstractControl | null {
    return this.form.get('amount');
  }

  get companyCode(): string | null {
    return this.tenantService.getCompanyCode();
  }

  showAmountError(): boolean {
    return !!this.amountControl && this.amountControl.invalid && (this.amountControl.touched || this.amountControl.dirty);
  }

  getAmountError(): string {
    const errors = this.amountControl?.errors;

    if (!errors || !this.showAmountError()) {
      return '';
    }

    if (errors['required']) {
      return 'Payment amount is required.';
    }

    if (errors['min']) {
      return 'Payment amount must be greater than $0.00.';
    }

    return 'Enter a valid payment amount.';
  }

  onSubmit(): void {
    this.errorMessage = '';
    this.successMessage = '';

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const company = this.companyCode;
    if (!company) {
      this.errorMessage = 'Missing tenant company in URL (?company=...)';
      return;
    }

    const amount = Number(this.form.getRawValue().amount);
    this.submitting = true;

    this.authService.makePayment(amount, company).subscribe({
      next: (response) => {
        this.submitting = false;
        this.successMessage = response.message || 'Payment recorded successfully.';
        window.setTimeout(() => {
          void this.router.navigate(['/dashboard'], {
            queryParams: { company }
          });
        }, 900);
      },
      error: (err: Error) => {
        this.submitting = false;
        this.errorMessage = err.message || `Unable to submit payment for ${company}.`;
      }
    });
  }
}
