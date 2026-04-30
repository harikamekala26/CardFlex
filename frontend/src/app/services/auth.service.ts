import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpParams } from '@angular/common/http';
import { Observable, catchError, throwError } from 'rxjs';

import { environment } from '../../environments/environment';
import {
  LoginRequest,
  LoginResponse,
  PaymentResponse,
  ProfileResponse,
  RegisterRequest,
  RegisterResponse
} from '../models/auth.model';
import { DashboardApiResponse } from '../models/dashboard.model';

interface TenantSession {
  companyCode: string;
  token: string;
}

export class ApiError extends Error {
  constructor(message: string, readonly status?: number) {
    super(message);
    this.name = 'ApiError';
  }
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly storageKey = 'cardflex_sessions';
  private readonly apiBaseUrl = environment.apiBaseUrl.replace(/\/$/, '');

  constructor(private readonly http: HttpClient) {}

  register(payload: RegisterRequest, companyCode: string): Observable<RegisterResponse> {
    return this.http
      .post<RegisterResponse>(`${this.apiBaseUrl}/register`, {
        ...payload,
        tenantId: companyCode
      })
      .pipe(catchError((error) => this.handleApiError(error, 'Registration failed. Please try again.')));
  }

  login(payload: LoginRequest, companyCode: string): Observable<LoginResponse> {
    const params = new HttpParams().set('company', companyCode);

    return this.http
      .post<LoginResponse>(`${this.apiBaseUrl}/login`, payload, { params })
      .pipe(catchError((error) => this.handleApiError(error, 'Login failed. Please try again.')));
  }

  getDashboard(companyCode: string): Observable<DashboardApiResponse> {
    const params = new HttpParams().set('company', companyCode);
    return this.http
      .get<DashboardApiResponse>(`${this.apiBaseUrl}/dashboard`, { params })
      .pipe(catchError((error) => this.handleApiError(error, 'Unable to load dashboard data.')));
  }

  getProfile(companyCode: string): Observable<ProfileResponse> {
    const params = new HttpParams().set('company', companyCode);
    return this.http
      .get<ProfileResponse>(`${this.apiBaseUrl}/profile`, { params })
      .pipe(catchError((error) => this.handleApiError(error, 'Unable to load profile data.')));
  }

  makePayment(amount: number, companyCode: string): Observable<PaymentResponse> {
    const params = new HttpParams().set('company', companyCode);
    return this.http
      .post<PaymentResponse>(`${this.apiBaseUrl}/payment`, { amount }, { params })
      .pipe(catchError((error) => this.handleApiError(error, 'Unable to submit payment.')));
  }

  setSession(companyCode: string, token: string): void {
    const sessions = this.getStoredSessions();
    sessions[companyCode] = token;
    localStorage.setItem(this.storageKey, JSON.stringify(sessions));
  }

  getToken(companyCode: string | null): string | null {
    if (!companyCode) {
      return null;
    }

    return this.getStoredSessions()[companyCode] ?? null;
  }

  getSession(companyCode: string | null): TenantSession | null {
    const token = this.getToken(companyCode);
    if (!companyCode || !token) {
      return null;
    }

    return { companyCode, token };
  }

  logout(companyCode?: string | null): void {
    if (!companyCode) {
      localStorage.removeItem(this.storageKey);
      return;
    }

    const sessions = this.getStoredSessions();
    delete sessions[companyCode];
    localStorage.setItem(this.storageKey, JSON.stringify(sessions));
  }

  isAuthenticated(companyCode: string | null): boolean {
    return !!this.getToken(companyCode);
  }

  private getStoredSessions(): Record<string, string> {
    const raw = localStorage.getItem(this.storageKey);
    if (!raw) {
      return {};
    }

    try {
      const parsed = JSON.parse(raw) as Record<string, string>;
      return typeof parsed === 'object' && parsed !== null ? parsed : {};
    } catch {
      return {};
    }
  }

  private handleApiError(error: unknown, fallbackMessage: string): Observable<never> {
    if (error instanceof HttpErrorResponse) {
      if (typeof error.error?.error === 'string' && error.error.error.trim()) {
        return throwError(() => new ApiError(error.error.error, error.status));
      }

      if (typeof error.error?.message === 'string' && error.error.message.trim()) {
        return throwError(() => new ApiError(error.error.message, error.status));
      }

      if (error.status === 0) {
        return throwError(() => new ApiError('Backend is unreachable. Verify the API server and base URL.', error.status));
      }

      return throwError(() => new ApiError(fallbackMessage, error.status));
    }

    return throwError(() => new Error(fallbackMessage));
  }
}
