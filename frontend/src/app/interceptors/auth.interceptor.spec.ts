import { TestBed } from '@angular/core/testing';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';

import { environment } from '../../environments/environment';
import { AuthService } from '../services/auth.service';
import { TenantService } from '../services/tenant.service';
import { authInterceptor } from './auth.interceptor';

describe('authInterceptor', () => {
  let authService: AuthService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    localStorage.clear();

    TestBed.configureTestingModule({
      providers: [
        provideHttpClient(withInterceptors([authInterceptor])),
        provideHttpClientTesting(),
        {
          provide: TenantService,
          useValue: {
            getCompanyCode: () => 'acme'
          }
        }
      ]
    });

    authService = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
    localStorage.clear();
  });

  it('adds the tenant session authorization header to payment requests', () => {
    authService.setSession('acme', 'tenant-token');

    authService.makePayment(75, 'acme').subscribe();

    const request = httpMock.expectOne(
      (req) => req.url === `${environment.apiBaseUrl}/payment` && req.params.get('company') === 'acme'
    );

    expect(request.request.headers.get('Authorization')).toBe('Bearer tenant-token');
    request.flush({
      message: 'payment recorded successfully',
      updatedBalance: 7925,
      transactionId: 7,
      amount: 75,
      timestamp: '2026-04-29T10:00:00Z'
    });
  });
});
