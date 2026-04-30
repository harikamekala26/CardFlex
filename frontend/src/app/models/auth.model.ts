export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  message: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token?: string;
  message?: string;
}

export interface PaymentRequest {
  amount: number;
}

export interface PaymentResponse {
  message: string;
  updatedBalance: number;
  transactionId: number;
  amount: number;
  timestamp: string;
}

export interface ProfileResponse {
  name: string;
  email: string;
}
