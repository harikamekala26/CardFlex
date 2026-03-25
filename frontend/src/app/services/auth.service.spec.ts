import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';

import { environment } from '../../environments/environment';
import { AuthService } from './auth.service';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    localStorage.clear();

    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()]
    });

    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
    localStorage.clear();
  });

  it('calls register with tenantId in the request body', () => {
    service.register({ name: 'Jane', email: 'jane@example.com', password: 'secret123' }, 'chase-bank').subscribe();

    const request = httpMock.expectOne(`${environment.apiBaseUrl}/register`);
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual({
      name: 'Jane',
      email: 'jane@example.com',
      password: 'secret123',
      tenantId: 'chase-bank'
    });

    request.flush({ message: 'user registered' });
  });

  it('calls login with the company query parameter', () => {
    service.login({ email: 'jane@example.com', password: 'secret123' }, 'capital-one').subscribe();

    const request = httpMock.expectOne(
      (req) => req.url === `${environment.apiBaseUrl}/login` && req.params.get('company') === 'capital-one'
    );

    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual({
      email: 'jane@example.com',
      password: 'secret123'
    });

    request.flush({ token: 'jwt-token', message: 'user logged in' });
  });

  it('stores and reads tenant-scoped sessions', () => {
    service.setSession('wells-fargo', 'tenant-token');

    expect(service.getToken('wells-fargo')).toBe('tenant-token');
    expect(service.isAuthenticated('wells-fargo')).toBeTrue();
    expect(service.getToken('capital-one')).toBeNull();
  });

  it('maps backend validation errors into readable errors', () => {
    let errorMessage = '';

    service.login({ email: 'jane@example.com', password: 'secret123' }, 'wells-fargo').subscribe({
      next: () => fail('expected login to fail'),
      error: (error: Error) => {
        errorMessage = error.message;
      }
    });

    const request = httpMock.expectOne(
      (req) => req.url === `${environment.apiBaseUrl}/login` && req.params.get('company') === 'wells-fargo'
    );
    request.flush({ error: 'invalid credentials' }, { status: 401, statusText: 'Unauthorized' });

    expect(errorMessage).toBe('invalid credentials');
  });

  it('maps network failures into a backend unavailable message', () => {
    let errorMessage = '';

    service.register({ name: 'Jane', email: 'jane@example.com', password: 'secret123' }, 'chase-bank').subscribe({
      next: () => fail('expected register to fail'),
      error: (error: Error) => {
        errorMessage = error.message;
      }
    });

    const request = httpMock.expectOne(`${environment.apiBaseUrl}/register`);
    request.error(new ProgressEvent('error'), { status: 0, statusText: 'Unknown Error' });

    expect(errorMessage).toBe('Backend is unreachable. Verify the API server and base URL.');
  });
});
