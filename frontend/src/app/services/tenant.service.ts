import { Injectable } from '@angular/core';

import { DEFAULT_TENANT, TENANT_CONFIGS } from '../config/tenant-config';
import { Tenant } from '../models/tenant.model';

@Injectable({ providedIn: 'root' })
export class TenantService {
  private tenant: Tenant | null = null;

  setFromCompanyCode(companyCode: string | null): void {
    if (!companyCode) {
      this.tenant = null;
      console.log('Tenant company from URL:', null);
      return;
    }

    const normalized = companyCode.toLowerCase();
    this.tenant = TENANT_CONFIGS[normalized] ?? {
      ...DEFAULT_TENANT,
      companyCode: normalized,
      name: `${normalized.toUpperCase()} Cards`
    };

    console.log('Tenant company from URL:', this.tenant.companyCode);
  }

  isFeatureEnabled(feature: string, defaultValue = true): boolean {
    const tenant = this.getResolvedTenant();
    return tenant.features?.[feature] ?? defaultValue;
  }

  getTenant(): Tenant | null {
    return this.tenant;
  }

  getResolvedTenant(): Tenant {
    return this.tenant ?? DEFAULT_TENANT;
  }

  getCompanyCode(): string | null {
    return this.tenant?.companyCode ?? null;
  }

  getAllPartners(): Tenant[] {
    return Object.values(TENANT_CONFIGS);
  }
}
