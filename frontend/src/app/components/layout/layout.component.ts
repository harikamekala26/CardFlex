import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterLink, RouterOutlet } from '@angular/router';

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
    private readonly route: ActivatedRoute,
    private readonly router: Router,
    private readonly tenantService: TenantService,
    protected readonly authService: AuthService
  ) {}

  ngOnInit(): void {
    this.route.queryParamMap.subscribe((params) => {
      const company = params.get('company');
      this.tenantService.setFromCompanyCode(company);
      this.applyTenantBranding();
    });
  }

  get companyCode(): string | null {
    return this.tenantService.getCompanyCode();
  }

  get tenantName(): string {
    return this.tenantService.getTenant()?.name ?? 'CardFlex';
  }

  get activeTenant() {
    return this.tenantService.getTenant();
  }

  get footerPartners() {
    return this.activeTenant ? [this.activeTenant] : this.tenantService.getAllPartners();
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/login'], {
      queryParams: { company: this.companyCode ?? undefined }
    });
  }

  private applyTenantBranding(): void {
    const tenant = this.activeTenant;
    if (!tenant) {
      document.title = 'CardFlex';
      document.documentElement.style.setProperty('--tenant-color', '#00539C');
      document.documentElement.style.setProperty('--tenant-secondary-color', '#8FB4D8');
      return;
    }

    document.title = `${tenant.name} | CardFlex`;
    document.documentElement.style.setProperty('--tenant-color', tenant.primaryColor);
    document.documentElement.style.setProperty('--tenant-secondary-color', tenant.secondaryColor);
  }
}
