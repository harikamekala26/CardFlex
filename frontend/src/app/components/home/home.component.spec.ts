import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { provideRouter, RouterLink } from '@angular/router';

import { TenantService } from '../../services/tenant.service';
import { HomeComponent } from './home.component';

describe('HomeComponent', () => {
  let fixture: ComponentFixture<HomeComponent>;
  let component: HomeComponent;
  let tenantService: {
    getCompanyCode: () => string | null;
    getResolvedTenant: () => { name: string };
  };

  beforeEach(async () => {
    tenantService = {
      getCompanyCode: () => 'capital-one',
      getResolvedTenant: () => ({ name: 'Capital One' })
    };

    await TestBed.configureTestingModule({
      imports: [HomeComponent],
      providers: [provideRouter([]), { provide: TenantService, useValue: tenantService }]
    }).compileComponents();

    fixture = TestBed.createComponent(HomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('renders the resolved tenant name and company code', () => {
    expect(component.tenantName).toBe('Capital One');
    expect(component.companyCode).toBe('capital-one');
    expect(fixture.nativeElement.textContent).toContain('Welcome to Capital One');
    expect(fixture.nativeElement.textContent).toContain('capital-one');
  });

  it('preserves the tenant company in register and login navigation links', () => {
    const links = fixture.debugElement.queryAll(By.css('a'));

    expect(links[0].injector.get(RouterLink).queryParams).toEqual({ company: 'capital-one' });
    expect(links[1].injector.get(RouterLink).queryParams).toEqual({ company: 'capital-one' });
  });
});
