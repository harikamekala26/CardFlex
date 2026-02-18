import { Injectable } from '@angular/core';

import { Tenant } from '../models/tenant.model';

@Injectable({ providedIn: 'root' })
export class TenantService {
  private tenant: Tenant | null = null;

  private readonly tenantMap: Record<string, Omit<Tenant, 'companyCode'>> = {
    acme: { name: 'Acme Card', themeColor: '#0B6E4F' },
    nova: { name: 'Nova Finance', themeColor: '#C84B31' },
    prime: { name: 'Prime Credit', themeColor: '#00539C' }
  };

  setFromCompanyCode(companyCode: string | null): void {
    if (!companyCode) {
      this.tenant = null;
      return;
    }

    const normalized = companyCode.toLowerCase();
    const fromMap = this.tenantMap[normalized] ?? {
      name: `${normalized.toUpperCase()} Cards`,
      themeColor: '#00539C'
    };

    this.tenant = {
      companyCode: normalized,
      ...fromMap
    };
  }

  getTenant(): Tenant | null {
    return this.tenant;
  }

  getCompanyCode(): string | null {
    return this.tenant?.companyCode ?? null;
  }
}
