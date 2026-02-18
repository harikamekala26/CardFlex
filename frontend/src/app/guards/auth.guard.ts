import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';

import { AuthService } from '../services/auth.service';
import { TenantService } from '../services/tenant.service';

export const authGuard: CanActivateFn = (route) => {
  const authService = inject(AuthService);
  const router = inject(Router);
  const tenantService = inject(TenantService);

  if (authService.isAuthenticated()) {
    return true;
  }

  return router.createUrlTree(['/login'], {
    queryParams: { company: route.queryParamMap.get('company') ?? tenantService.getCompanyCode() ?? undefined }
  });
};
