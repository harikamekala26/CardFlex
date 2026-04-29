import { Component, OnInit } from '@angular/core';
import { CommonModule, CurrencyPipe, DatePipe } from '@angular/common';
import { Router, RouterLink } from '@angular/router';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { DashboardData, DashboardTransaction, normalizeDashboardData } from '../../models/dashboard.model';

interface SpendingSummaryItem {
  month: string;
  amount: number;
  percent: number;
}

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

  get spendingSummary(): SpendingSummaryItem[] {
    const transactions = this.data?.transactions ?? [];
    const monthlySpending = new Map<string, { date: Date; amount: number }>();

    for (const transaction of transactions) {
      const date = new Date(transaction.date);

      if (Number.isNaN(date.getTime())) {
        continue;
      }

      const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
      const existing = monthlySpending.get(monthKey);
      const spendAmount = Math.abs(transaction.amount);

      monthlySpending.set(monthKey, {
        date: existing?.date ?? new Date(date.getFullYear(), date.getMonth(), 1),
        amount: (existing?.amount ?? 0) + spendAmount
      });
    }

    const items = Array.from(monthlySpending.values()).sort((first, second) => first.date.getTime() - second.date.getTime());
    const maxAmount = Math.max(...items.map((item) => item.amount), 0);

    return items.map((item) => ({
      month: item.date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' }),
      amount: item.amount,
      percent: maxAmount > 0 ? Math.round((item.amount / maxAmount) * 100) : 0
    }));
  }

  trackTransaction(_: number, transaction: DashboardTransaction): string {
    return `${transaction.date}-${transaction.merchant}-${transaction.amount}`;
  }

  trackSpendingMonth(_: number, item: SpendingSummaryItem): string {
    return item.month;
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
