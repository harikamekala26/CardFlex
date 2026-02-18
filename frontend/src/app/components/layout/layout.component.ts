import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NavigationEnd, Router, RouterLink, RouterOutlet } from '@angular/router';
import { filter } from 'rxjs';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

@Component({
  selector: 'app-layout',
  standalone: true,
  imports: [CommonModule, RouterOutlet, RouterLink],
  templateUrl: './layout.component.html',
  styleUrl: './layout.component.css'
})
export class LayoutComponent implements OnInit {
  constructor(
    private readonly router: Router,
    private readonly tenantService: TenantService,
    protected readonly authService: AuthService
  ) {}

  ngOnInit(): void {
    this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe(() => this.syncTenantFromUrl());

    this.syncTenantFromUrl();
  }

  get companyCode(): string | null {
    return this.tenantService.getCompanyCode();
  }

  get tenantName(): string {
    return this.tenantService.getTenant()?.name ?? 'CardFlex';
  }

  get footerPartners() {
    const activeTenant = this.tenantService.getTenant();
    return activeTenant ? [activeTenant] : this.tenantService.getAllPartners();
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/login'], {
      queryParams: { company: this.companyCode ?? undefined }
    });
  }

  private syncTenantFromUrl(): void {
    const tree = this.router.parseUrl(this.router.url);
    const company = tree.queryParams['company'] ?? null;
    this.tenantService.setFromCompanyCode(company);

    const tenant = this.tenantService.getTenant();
    if (!tenant) {
      document.title = 'CardFlex';
      document.documentElement.style.setProperty('--tenant-color', '#00539C');
      return;
    }

    document.title = `${tenant.name} | CardFlex`;
    document.documentElement.style.setProperty('--tenant-color', tenant.themeColor);
  }
}