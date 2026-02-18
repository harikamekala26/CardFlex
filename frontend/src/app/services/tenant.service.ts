import { Injectable } from '@angular/core';

import { DEFAULT_PARTNER, PARTNER_CONFIGS } from '../config/partner-config';
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
    this.tenant = PARTNER_CONFIGS[normalized] ?? {
      ...DEFAULT_PARTNER,
      companyCode: normalized,
      name: `${normalized.toUpperCase()} Cards`
    };

    console.log('Tenant company from URL:', this.tenant.companyCode);
  }

  getTenant(): Tenant | null {
    return this.tenant;
  }

  getCompanyCode(): string | null {
    return this.tenant?.companyCode ?? null;
  }

  getAllPartners(): Tenant[] {
    return Object.values(PARTNER_CONFIGS);
  }
}
