import { authClient } from './api';
import { LoginRequest, RegisterRequest, LoginResponse, User } from '@/types';

export const authService = {
  async register(data: RegisterRequest): Promise<LoginResponse> {
    const response = await authClient.post<LoginResponse>('/auth/register', data);
    if (response.data.access_token) {
      localStorage.setItem('access_token', response.data.access_token);
      localStorage.setItem('refresh_token', response.data.refresh_token);
    }
    return response.data;
  },

  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await authClient.post<LoginResponse>('/auth/login', data);
    if (response.data.access_token) {
      localStorage.setItem('access_token', response.data.access_token);
      localStorage.setItem('refresh_token', response.data.refresh_token);
    }
    return response.data;
  },

  async logout(): Promise<void> {
    const refreshToken = localStorage.getItem('refresh_token');
    try {
      await authClient.post('/auth/logout', { refresh_token: refreshToken });
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  },

  async getProfile(): Promise<User> {
    const response = await authClient.get<{ user: User }>('/auth/me');
    return response.data.user;
  },

  async changePassword(oldPassword: string, newPassword: string): Promise<void> {
    await authClient.post('/auth/change-password', {
      old_password: oldPassword,
      new_password: newPassword,
    });
  },

  getAccessToken(): string | null {
    return localStorage.getItem('access_token');
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('access_token');
  },
};
