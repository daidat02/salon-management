import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { RootState } from '../store';
import { API_BASE } from '../../constants/index';
import { loginUser } from '@/services/auth';

// Giữ re-export để code cũ dùng API_BASE từ slice không gãy.
export { API_BASE };

export type AuthUser = {
  id: string;
  email?: string;
  full_name?: string;
  phone_number: string;
  role: string;
  status: string;
};

type AuthStatus = 'idle' | 'loading' | 'succeeded' | 'failed';

type AuthState = {
  user: AuthUser | null;
  /** Chỉ giữ trong memory — KHÔNG persist (xem blacklist trong store.ts). */
  accessToken: string | null;
  status: AuthStatus;
  error: string | null;
};

const initialState: AuthState = {
  user: null,
  accessToken: null,
  status: 'idle',
  error: null,
};

// POST /login { phone, password } -> { data: { user, access_token } }
// + cookie refresh_token httpOnly do Go set (apiClient withCredentials).
// Đi qua apiClient để ăn request/response interceptor dùng chung.

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setCredentials(state, action: PayloadAction<{ user: AuthUser; accessToken: string }>) {
      state.user = action.payload.user;
      state.accessToken = action.payload.accessToken;
      state.status = 'succeeded';
      state.error = null;
    },
    setAccessToken(state, action: PayloadAction<string>) {
      state.accessToken = action.payload;
      state.status = 'succeeded';
    },
    logout(state) {
      state.user = null;
      state.accessToken = null;
      state.status = 'idle';
      state.error = null;
    },
    clearAuthError(state) {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(loginUser.pending, (state) => {
        state.status = 'loading';
        state.error = null;
      })
      .addCase(loginUser.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.user = action.payload.user;
        state.accessToken = action.payload.accessToken;
      })
      .addCase(loginUser.rejected, (state, action) => {
        state.status = 'failed';
        state.error = (action.payload as string | undefined) ?? 'Đăng nhập thất bại';
      });
  },
});

export const { setCredentials, setAccessToken, logout, clearAuthError } = authSlice.actions;
export default authSlice.reducer;

// Selectors
export const selectAuthUser = (state: RootState) => state.auth.user;
export const selectAccessToken = (state: RootState) => state.auth.accessToken;
export const selectAuthStatus = (state: RootState) => state.auth.status;
export const selectAuthError = (state: RootState) => state.auth.error;
export const selectIsAuthenticated = (state: RootState) =>
  state.auth.user !== null && state.auth.accessToken !== null;
