import { TestBed } from '@angular/core/testing';

import { DEFAULT_TENANT } from '../config/tenant-config';
import { TenantService } from './tenant.service';

describe('TenantService', () => {
  let service: TenantService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TenantService);
  });

  it('resolves a known tenant from the company code', () => {
    service.setFromCompanyCode('chase-bank');

    expect(service.getCompanyCode()).toBe('chase-bank');
    expect(service.getResolvedTenant().name).toBe('Chase Bank');
  });

  it('creates a fallback tenant for an unknown company code', () => {
    service.setFromCompanyCode('acme');

    expect(service.getCompanyCode()).toBe('acme');
    expect(service.getResolvedTenant().name).toBe('ACME Cards');
  });

  it('returns the default tenant when no company code is present', () => {
    service.setFromCompanyCode(null);

    expect(service.getTenant()).toBeNull();
    expect(service.getResolvedTenant()).toEqual(DEFAULT_TENANT);
  });

  it('reports tenant feature availability from tenant configuration', () => {
    service.setFromCompanyCode('capital-one');

    expect(service.isFeatureEnabled('paymentsEnabled')).toBeTrue();
    expect(service.isFeatureEnabled('profileEnabled')).toBeFalse();
  });
});
