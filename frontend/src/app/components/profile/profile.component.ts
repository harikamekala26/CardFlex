import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';

import { ProfileResponse } from '../../models/auth.model';
import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './profile.component.html',
  styleUrl: './profile.component.css'
})
export class ProfileComponent implements OnInit {
  profile: ProfileResponse | null = null;
  loading = true;
  errorMessage = '';

  constructor(
    private readonly authService: AuthService,
    private readonly tenantService: TenantService,
    private readonly route: ActivatedRoute,
    private readonly router: Router
  ) {}

  ngOnInit(): void {
    this.loadProfile();
  }

  get companyCode(): string | null {
    return this.route.snapshot.queryParamMap.get('company') ?? this.tenantService.getCompanyCode();
  }

  get tenantName(): string {
    return this.tenantService.getResolvedTenant().name;
  }

  get tenantLogo(): string {
    return this.tenantService.getResolvedTenant().theme.logoUrl;
  }

  reload(): void {
    this.loadProfile();
  }

  private loadProfile(): void {
    const company = this.companyCode;
    this.loading = true;
    this.errorMessage = '';

    if (!company) {
      this.loading = false;
      this.profile = null;
      this.errorMessage = 'Missing tenant company in URL (?company=...)';
      return;
    }

    this.authService.getProfile(company).subscribe({
      next: (profile) => {
        this.loading = false;
        this.profile = profile;
      },
      error: (err) => {
        this.loading = false;
        this.profile = null;
        const status = typeof err === 'object' && err !== null && 'status' in err ? Number(err.status) : undefined;
        const message =
          err instanceof Error
            ? err.message
            : typeof err === 'object' && err !== null && 'message' in err && typeof err.message === 'string'
              ? err.message
              : 'Failed to load profile';
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
