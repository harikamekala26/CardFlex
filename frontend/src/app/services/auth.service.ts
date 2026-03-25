import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { environment } from '../../environments/environment';
import { DashboardData } from '../models/dashboard.model';

@Injectable({ providedIn: 'root' })
export class AuthService {
  constructor(private readonly http: HttpClient) {}

  register(
    payload: { name: string; email: string; password: string },
    companyCode: string
  ): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${environment.apiBaseUrl}/register`, {
      ...payload,
      tenantId: companyCode
    });
  }

  login(payload: { email: string; password: string }, companyCode: string): Observable<{ token?: string; message?: string }> {
    return this.http.post<{ token?: string; message?: string }>(`${environment.apiBaseUrl}/login?company=${companyCode}`, payload);
  }

  getDashboard(companyCode: string): Observable<DashboardData> {
    return this.http.get<DashboardData>(`${environment.apiBaseUrl}/dashboard?company=${companyCode}`);
  }

  setToken(token: string): void {
    localStorage.setItem('cardflex_token', token);
  }

  getToken(): string | null {
    return localStorage.getItem('cardflex_token');
  }

  logout(): void {
    localStorage.removeItem('cardflex_token');
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }
}
