import { Component, OnInit } from '@angular/core';
import { CommonModule, CurrencyPipe, DatePipe } from '@angular/common';
import { Router, RouterLink } from '@angular/router';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { DashboardData, DashboardTransaction, normalizeDashboardData } from '../../models/dashboard.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, CurrencyPipe, DatePipe, RouterLink],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent implements OnInit {
  data: DashboardData | null = null;
  loading = true;
  errorMessage = '';

  constructor(
    private readonly authService: AuthService,
    private readonly tenantService: TenantService,
    private readonly router: Router
  ) {}

  ngOnInit(): void {
    this.loadDashboard();
  }

  get companyCode(): string | null {
    return this.tenantService.getCompanyCode();
  }

  get paymentUsed(): number {
    if (!this.data) {
      return 0;
    }

    return this.data.accountSummary.creditLimit - this.data.accountSummary.availableBalance;
  }

  get utilizationPercent(): number {
    if (!this.data?.accountSummary.creditLimit) {
      return 0;
    }

    return Math.round((this.paymentUsed / this.data.accountSummary.creditLimit) * 100);
  }

  get hasTransactions(): boolean {
    return (this.data?.transactions.length ?? 0) > 0;
  }

  trackTransaction(_: number, transaction: DashboardTransaction): string {
    return `${transaction.date}-${transaction.merchant}-${transaction.amount}`;
  }

  reload(): void {
    this.loadDashboard();
  }

  private loadDashboard(): void {
    const company = this.tenantService.getCompanyCode();
    this.loading = true;
    this.errorMessage = '';

    if (!company) {
      this.loading = false;
      this.errorMessage = 'Missing tenant company in URL (?company=...)';
      return;
    }

    this.authService.getDashboard(company).subscribe({
      next: (response) => {
        this.loading = false;
        this.data = normalizeDashboardData(response);
      },
      error: (err) => {
        this.loading = false;
        this.data = null;
        const status = typeof err === 'object' && err !== null && 'status' in err ? Number(err.status) : undefined;
        const message =
          err instanceof Error
            ? err.message
            : typeof err === 'object' && err !== null && 'message' in err && typeof err.message === 'string'
              ? err.message
              : 'Failed to load dashboard';
        this.errorMessage = message;

        if (status === 401 || status === 403) {
          this.authService.logout(company);
          void this.router.navigate(['/login'], {
            queryParams: { company }
          });
        }
      }
    });
  }
}
