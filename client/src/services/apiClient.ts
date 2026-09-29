import axios, { type AxiosError, type AxiosInstance, type InternalAxiosRequestConfig } from 'axios';
import { API_BASE } from '@/constants';

// Holder token in-memory — đồng bộ từ authSlice (login thunk).
// Không đọc Redux store trực tiếp vì store tạo per-Provider,
// module này đứng ngoài React tree.
let inMemoryToken: string | null = null;

export function setApiAccessToken(token: string | null) {
  inMemoryToken = token;
  // Đồng bộ sessionStorage nếu đang ở browser và trước đó đã lưu
  if (typeof window !== 'undefined') {
    if (token) {
      // Chỉ ghi nếu trước đó có lưu (tôn trọng "ghi nhớ" ở login)
      // Nhưng để refresh luôn khả dụng sau reload, ta lưu luôn khi có token mới
      // Nếu muốn chỉ lưu khi đã từng lưu, kiểm tra hasSession trước.
      // Ở đây lưu luôn để lần sau reload không cần refresh ngay.
      try {
        const prev = sessionStorage.getItem('access_token');
        // Nếu trước đó đã có hoặc token mới khác null, cập nhật
        if (prev || token) sessionStorage.setItem('access_token', token);
      } catch {}
    } else {
      try {
        sessionStorage.removeItem('access_token');
      } catch {}
    }
  }
  // Thử đồng bộ Redux nếu store đã được gắn (tránh circular import)
  try {
    // @ts-ignore
    if (typeof window !== 'undefined' && (window as any).__REDUX_STORE__?.dispatch) {
      // @ts-ignore
      const { setAccessToken } = require('@/store/slices/authSlice');
      if (token) (window as any).__REDUX_STORE__.dispatch(setAccessToken(token));
    }
  } catch {}
}

function resolveToken(): string | null {
  if (inMemoryToken) return inMemoryToken;
  // Fallback: mirror tab-scoped do trang login ghi khi tick "ghi nhớ".
  if (typeof window !== 'undefined') {
    try {
      return sessionStorage.getItem('access_token');
    } catch {
      return null;
    }
  }
  return null;
}

export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE,
  timeout: 15000,
  // Nhận/gửi cookie refresh_token httpOnly do Go set.
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

// Request interceptor: gắn Bearer token khi có.
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = resolveToken();
    if (token) {
      config.headers.set('Authorization', `Bearer ${token}`);
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// --- Refresh token queue ---
let isRefreshing = false;
type FailedQueueItem = {
  resolve: (token: string) => void;
  reject: (err: unknown) => void;
};
let failedQueue: FailedQueueItem[] = [];

function processQueue(error: unknown, token: string | null = null) {
  failedQueue.forEach((prom) => {
    if (error) prom.reject(error);
    else prom.resolve(token!);
  });
  failedQueue = [];
}

function clearAuthAndRedirect() {
  inMemoryToken = null;
  if (typeof window !== 'undefined') {
    try {
      sessionStorage.removeItem('access_token');
      localStorage.removeItem('access_token');
    } catch {}
    const next = window.location.pathname + window.location.search;
    if (!window.location.pathname.startsWith('/app/login')) {
      // Cố ý reload cả trang (không dùng router): interceptor đứng ngoài React tree
      // eslint-disable-next-line @next/next/no-location-assign-relative-destination
      window.location.href = `/app/login?next=${encodeURIComponent(next)}`;
    }
  }
}

// Response interceptor: 401 -> thử refresh token rồi retry (chỉ trên browser, server không có cookie refresh_token)
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // Chỉ xử lý 401 và có config
    if (error.response?.status !== 401 || !originalRequest) {
      return Promise.reject(error);
    }

    // Không retry nếu chính refresh endpoint bị 401 -> refresh_token hết hạn
    if (originalRequest.url?.includes('/refresh-token')) {
      clearAuthAndRedirect();
      return Promise.reject(error);
    }

    // Tránh loop vô hạn
    if (originalRequest._retry) {
      clearAuthAndRedirect();
      return Promise.reject(error);
    }

    // Server component không có window/sessionStorage/cookie refresh_token khả dụng -> không thể refresh, trả về luôn
    if (typeof window === 'undefined') {
      return Promise.reject(error);
    }

    // Nếu đang refresh, queue lại
    if (isRefreshing) {
      return new Promise<string>((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      })
        .then((token) => {
          originalRequest.headers.set('Authorization', `Bearer ${token}`);
          return apiClient(originalRequest);
        })
        .catch((err) => Promise.reject(err));
    }

    originalRequest._retry = true;
    isRefreshing = true;

    try {
      // Dùng axios thuần (không qua apiClient) để tránh interceptor loop
      const res = await axios.post(`${API_BASE}/refresh-token`, {}, { withCredentials: true });
      // Handler Go trả: SuccessResponse{ Data: { access_token: "..." } }
      const newToken =
        (res.data?.data?.access_token as string | undefined) ??
        (res.data?.access_token as string | undefined) ??
        null;

      if (!newToken) throw new Error('No access_token from refresh');

      setApiAccessToken(newToken);
      processQueue(null, newToken);

      originalRequest.headers.set('Authorization', `Bearer ${newToken}`);
      return apiClient(originalRequest);
    } catch (refreshError) {
      processQueue(refreshError, null);
      clearAuthAndRedirect();
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  },
);

type ApiErrorBody = {
  message?: string;
  error?: string;
};

// Helper bóc message lỗi theo envelope chuẩn của backend Go.
export function getApiErrorMessage(
  error: unknown,
  fallback = 'Có lỗi xảy ra, thử lại sau',
): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as ApiErrorBody | undefined;
    return data?.message ?? data?.error ?? fallback;
  }
  return error instanceof Error ? error.message : fallback;
}

export default apiClient;
