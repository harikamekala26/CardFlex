import { Component, OnInit } from '@angular/core';
import { CommonModule, CurrencyPipe, DatePipe } from '@angular/common';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { DashboardData } from '../../models/dashboard.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, CurrencyPipe, DatePipe],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent implements OnInit {
  data: DashboardData | null = null;
  errorMessage = '';

  constructor(
    private readonly authService: AuthService,
    private readonly tenantService: TenantService
  ) {}

  ngOnInit(): void {
    const company = this.tenantService.getCompanyCode();
    if (!company) {
      this.errorMessage = 'Missing tenant company in URL (?company=...)';
      return;
    }

    this.authService.getDashboard(company).subscribe({
      next: (response) => {
        this.data = response;
      },
      error: (err) => {
        this.errorMessage = err.error?.error ?? 'Failed to load dashboard';
      }
    });
  }
}
