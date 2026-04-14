import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { ActivatedRoute, convertToParamMap, provideRouter, Router, RouterLink } from '@angular/router';
import { BehaviorSubject } from 'rxjs';

import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { LayoutComponent } from './layout.component';

describe('LayoutComponent', () => {
  let fixture: ComponentFixture<LayoutComponent>;
  let component: LayoutComponent;
  let router: Router;
  let authService: {
    isAuthenticated: (companyCode: string | null) => boolean;
    logout: (companyCode?: string | null) => void;
  };
  let tenantService: {
    setFromCompanyCode: (companyCode: string | null) => void;
    getCompanyCode: () => string | null;
    getResolvedTenant: () => {
      name: string;
      theme: { logoUrl: string; primaryColor: string; secondaryColor: string };
    };
    getTenant: () => {
      name: string;
      cardArt?: { frontGradient: string; backGradient: string };
    } | null;
    getAllPartners: () => Array<{ name: string }>;
  };
  let queryParamMap$: BehaviorSubject<ReturnType<typeof convertToParamMap>>;

  beforeEach(async () => {
    queryParamMap$ = new BehaviorSubject(convertToParamMap({ company: 'chase-bank' }));

    authService = {
      isAuthenticated: () => false,
      logout: () => undefined
    };

    tenantService = {
      setFromCompanyCode: () => undefined,
      getCompanyCode: () => 'chase-bank',
      getResolvedTenant: () => ({
        name: 'Chase Bank',
        theme: {
          logoUrl: 'https://example.com/chase-logo.png',
          primaryColor: '#0A2A66',
          secondaryColor: '#2E8BC0'
        }
      }),
      getTenant: () => ({
        name: 'Chase Bank',
        cardArt: {
          frontGradient: 'linear-gradient(#fff, #000)',
          backGradient: 'linear-gradient(#000, #fff)'
        }
      }),
      getAllPartners: () => [{ name: 'Chase Bank' }]
    };

    await TestBed.configureTestingModule({
      imports: [LayoutComponent],
      providers: [
        provideRouter([]),
        { provide: ActivatedRoute, useValue: { queryParamMap: queryParamMap$.asObservable() } },
        { provide: AuthService, useValue: authService },
        { provide: TenantService, useValue: tenantService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(LayoutComponent);
    component = fixture.componentInstance;
    router = TestBed.inject(Router);
    fixture.detectChanges();
  });

  it('syncs the tenant from the company query parameter and applies branding', () => {
    const setTenantSpy = spyOn(tenantService, 'setFromCompanyCode').and.callThrough();

    queryParamMap$.next(convertToParamMap({ company: 'capital-one' }));
    fixture.detectChanges();

    expect(setTenantSpy).toHaveBeenCalledWith('capital-one');
    expect(document.title).toBe('Chase Bank | CardFlex');
    expect(document.documentElement.style.getPropertyValue('--tenant-color')).toBe('#0A2A66');
    expect(document.documentElement.style.getPropertyValue('--tenant-secondary-color')).toBe('#2E8BC0');
  });

  it('preserves the tenant company in unauthenticated navigation links', () => {
    const linkTexts = fixture.debugElement.queryAll(By.css('nav a')).map((element) => ({
      text: element.nativeElement.textContent.trim(),
      params: element.injector.get(RouterLink).queryParams
    }));

    expect(linkTexts.find((link) => link.text === 'Home')?.params).toEqual({ company: 'chase-bank' });
    expect(linkTexts.find((link) => link.text === 'Register')?.params).toEqual({ company: 'chase-bank' });
    expect(linkTexts.find((link) => link.text === 'Login')?.params).toEqual({ company: 'chase-bank' });
    expect(linkTexts.find((link) => link.text === 'Dashboard')).toBeUndefined();
  });

  it('shows the dashboard link for authenticated users and preserves tenant context on logout', () => {
    spyOn(authService, 'isAuthenticated').and.returnValue(true);
    const logoutSpy = spyOn(authService, 'logout').and.callThrough();
    const navigateSpy = spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));

    fixture.detectChanges();

    const linkTexts = fixture.debugElement.queryAll(By.css('nav a')).map((element) => ({
      text: element.nativeElement.textContent.trim(),
      params: element.injector.get(RouterLink).queryParams
    }));

    expect(linkTexts.find((link) => link.text === 'Dashboard')?.params).toEqual({ company: 'chase-bank' });

    component.logout();

    expect(logoutSpy).toHaveBeenCalledWith('chase-bank');
    expect(navigateSpy).toHaveBeenCalledWith(['/login'], {
      queryParams: { company: 'chase-bank' }
    });
  });
});
