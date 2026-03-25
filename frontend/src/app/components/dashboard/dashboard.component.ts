import { Component, OnInit } from '@angular/core';
import { CommonModule, CurrencyPipe, DatePipe } from '@angular/common';
import { Router, RouterLink } from '@angular/router';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { DashboardData } from '../../models/dashboard.model';

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

    return this.data.card.creditLimit - this.data.card.availableBalance;
  }

  get utilizationPercent(): number {
    if (!this.data?.card.creditLimit) {
      return 0;
    }

    return Math.round((this.paymentUsed / this.data.card.creditLimit) * 100);
  }

  trackTransaction(_: number, transaction: DashboardData['transactions'][number]): string {
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
        this.data = response;
      },
      error: (err) => {
        this.loading = false;
        this.data = null;
        this.errorMessage = err.error?.error ?? 'Failed to load dashboard';

        if (err.status === 401 || err.status === 403) {
          this.authService.logout(company);
          void this.router.navigate(['/login'], {
            queryParams: { company }
          });
        }
      }
    });
  }
}
