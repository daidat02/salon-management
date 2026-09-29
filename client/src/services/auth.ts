import { AuthUser } from '@/store/slices/authSlice';
import apiClient, { getApiErrorMessage, setApiAccessToken } from './apiClient';
import { createAsyncThunk } from '@reduxjs/toolkit';
import { API_ENDPOINTS } from '@/constants';

const { AUTH } = API_ENDPOINTS;
export const loginUser = createAsyncThunk(
  'auth/login',
  async (payload: { phone: string; password: string }, { rejectWithValue }) => {
    try {
      const res = await apiClient.post(AUTH.LOGIN, payload);
      const user = res.data?.data?.user as AuthUser;
      const accessToken = (res.data?.data?.access_token ?? '') as string;
      setApiAccessToken(accessToken || null);
      console.log('loginUser: user', user);
      return { user, accessToken };
    } catch (err) {
      return rejectWithValue(getApiErrorMessage(err, 'Đăng nhập thất bại'));
    }
  },
);
