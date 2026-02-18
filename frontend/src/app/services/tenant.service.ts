import { Injectable } from '@angular/core';

import { Tenant } from '../models/tenant.model';

@Injectable({ providedIn: 'root' })
export class TenantService {
  private tenant: Tenant | null = null;

  private readonly tenantMap: Record<string, Omit<Tenant, 'companyCode'>> = {
    acme: {
      name: 'Acme Card',
      themeColor: '#0B6E4F',
      contactNumber: '+1 (800) 555-0111',
      address: '200 Market Street, San Francisco, CA 94105',
      billingDetails: 'Billing Cycle: 1st monthly | Due Date: 25th monthly'
    },
    nova: {
      name: 'Nova Finance',
      themeColor: '#C84B31',
      contactNumber: '+1 (800) 555-0182',
      address: '500 W Madison St, Chicago, IL 60661',
      billingDetails: 'Billing Cycle: 5th monthly | Due Date: 28th monthly'
    },
    prime: {
      name: 'Prime Credit',
      themeColor: '#00539C',
      contactNumber: '+1 (800) 555-0199',
      address: '1201 3rd Ave, Seattle, WA 98101',
      billingDetails: 'Billing Cycle: 10th monthly | Due Date: 2nd monthly'
    }
  };

  setFromCompanyCode(companyCode: string | null): void {
    if (!companyCode) {
      this.tenant = null;
      return;
    }

    const normalized = companyCode.toLowerCase();
    const fromMap = this.tenantMap[normalized] ?? {
      name: `${normalized.toUpperCase()} Cards`,
      themeColor: '#00539C',
      contactNumber: '+1 (800) 555-0100',
      address: '100 Main St, New York, NY 10001',
      billingDetails: 'Billing Cycle: 1st monthly | Due Date: 25th monthly'
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

  getAllPartners(): Tenant[] {
    return Object.entries(this.tenantMap).map(([companyCode, tenant]) => ({
      companyCode,
      ...tenant
    }));
  }
}