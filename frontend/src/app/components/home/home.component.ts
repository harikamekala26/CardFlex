import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

import { TenantService } from '../../services/tenant.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
  constructor(private readonly tenantService: TenantService) {}

  get companyCode(): string | null {
    return this.tenantService.getCompanyCode();
  }
}
