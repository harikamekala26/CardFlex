import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { ActivatedRoute, convertToParamMap, provideRouter, Router, RouterLink } from '@angular/router';
import { Observable, of, throwError } from 'rxjs';

import { ProfileComponent } from './profile.component';
import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';

describe('ProfileComponent', () => {
  let fixture: ComponentFixture<ProfileComponent>;
  let component: ProfileComponent;
  let router: Router;
  let authService: {
    getProfile: (companyCode: string) => unknown;
    logout: (companyCode?: string | null) => void;
  };
  let tenantService: {
    getCompanyCode: () => string | null;
    getResolvedTenant: () => {
      name: string;
      theme: { logoUrl: string };
    };
  };
  let activatedRoute: {
    snapshot: {
      queryParamMap: ReturnType<typeof convertToParamMap>;
    };
  };

  beforeEach(async () => {
    authService = {
      getProfile: () => of({ name: 'Jane Doe', email: 'jane@example.com' }),
      logout: () => undefined
    };
    tenantService = {
      getCompanyCode: () => 'acme',
      getResolvedTenant: () => ({
        name: 'Acme Card',
        theme: {
          logoUrl: 'https://example.com/acme-logo.png'
        }
      })
    };
    activatedRoute = {
      snapshot: {
        queryParamMap: convertToParamMap({ company: 'acme' })
      }
    };

    await TestBed.configureTestingModule({
      imports: [ProfileComponent],
      providers: [
        provideRouter([]),
        { provide: AuthService, useValue: authService },
        { provide: TenantService, useValue: tenantService },
        { provide: ActivatedRoute, useValue: activatedRoute }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ProfileComponent);
    component = fixture.componentInstance;
    router = TestBed.inject(Router);
  });

  it('renders the authenticated profile returned by the backend', () => {
    const getProfileSpy = spyOn(authService, 'getProfile').and.callThrough();

    fixture.detectChanges();

    expect(getProfileSpy).toHaveBeenCalledWith('acme');
    expect(component.profile).toEqual({ name: 'Jane Doe', email: 'jane@example.com' });
    expect(fixture.nativeElement.textContent).toContain('Jane Doe');
    expect(fixture.nativeElement.textContent).toContain('jane@example.com');
    expect(fixture.nativeElement.textContent).toContain('Acme Card');
  });

  it('preserves the active tenant on the dashboard link', () => {
    fixture.detectChanges();

    const dashboardLink = fixture.debugElement
      .queryAll(By.css('a'))
      .find((element) => element.nativeElement.textContent.includes('Back to dashboard'));

    expect(dashboardLink?.injector.get(RouterLink).queryParams).toEqual({ company: 'acme' });
  });

  it('shows a loading spinner while the profile request is in flight', () => {
    spyOn(authService, 'getProfile').and.returnValue(new Observable());

    fixture.detectChanges();

    expect(component.loading).toBeTrue();
    expect(component.profile).toBeNull();
    expect(fixture.debugElement.query(By.css('.spinner'))).not.toBeNull();
    expect(fixture.nativeElement.textContent).toContain('Loading your profile');
  });

  it('shows a retryable error state when loading fails', () => {
    const getProfileSpy = spyOn(authService, 'getProfile').and.returnValues(
      throwError(() => new Error('profile service unavailable')),
      of({ name: 'Jane Doe', email: 'jane@example.com' })
    );

    fixture.detectChanges();

    expect(component.errorMessage).toBe('profile service unavailable');
    expect(fixture.nativeElement.textContent).toContain('Profile Unavailable');
    expect(fixture.nativeElement.textContent).toContain('Try Again');

    const retryButton = fixture.debugElement.query(By.css('button'));
    retryButton.triggerEventHandler('click');
    fixture.detectChanges();

    expect(getProfileSpy).toHaveBeenCalledTimes(2);
    expect(component.profile?.email).toBe('jane@example.com');
    expect(component.errorMessage).toBe('');
  });

  it('shows a missing-tenant error before calling the API', () => {
    activatedRoute.snapshot.queryParamMap = convertToParamMap({});
    spyOn(tenantService, 'getCompanyCode').and.returnValue(null);
    const getProfileSpy = spyOn(authService, 'getProfile').and.callThrough();

    fixture.detectChanges();

    expect(getProfileSpy).not.toHaveBeenCalled();
    expect(component.errorMessage).toBe('Missing tenant company in URL (?company=...)');
  });

  it('logs out and routes to login when profile access is unauthorized', () => {
    const navigateSpy = spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));
    const logoutSpy = spyOn(authService, 'logout').and.callThrough();
    spyOn(authService, 'getProfile').and.returnValue(throwError(() => ({ status: 401, message: 'Unauthorized' })));

    fixture.detectChanges();

    expect(logoutSpy).toHaveBeenCalledWith('acme');
    expect(navigateSpy).toHaveBeenCalledWith(['/login'], {
      queryParams: { company: 'acme' }
    });
  });
});
