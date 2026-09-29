'use client';

import { Suspense, useState } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { Eye, EyeOff, LockKeyhole, Phone, Scissors } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAppDispatch, useAppSelector } from '@/store/hook';
import { clearAuthError, selectAuthError, selectAuthStatus } from '@/store/slices/authSlice';
import { loginUser } from '@/services/auth';

// POST /login { phone, password } -> { data: { user, access_token } }
// + cookie refresh_token httpOnly do Go set.
export default function LoginPage() {
  return (
    <Suspense
      fallback={<div className="text-center text-sm text-on-surface-variant">Đang tải...</div>}
    >
      <LoginForm />
    </Suspense>
  );
}

function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const dispatch = useAppDispatch();
  const status = useAppSelector(selectAuthStatus);
  const error = useAppSelector(selectAuthError);
  const loading = status === 'loading';
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [remember, setRemember] = useState(true);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    dispatch(clearAuthError());
    const result = await dispatch(loginUser({ phone, password }));
    if (loginUser.fulfilled.match(result)) {
      // Mirror token vào sessionStorage khi "ghi nhớ" (tab-scoped).
      // Token chính nằm trong Redux memory (không persist — chống XSS).
      if (remember) sessionStorage.setItem('access_token', result.payload.accessToken);
      router.push(searchParams.get('next') ?? '/app/overview');
    }
  }

  return (
    <div className="space-y-6">
      {/* Brand gọn cho mobile */}
      <div className="flex items-center gap-2.5 lg:hidden">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-brand text-white shadow-sm">
          <Scissors className="h-5 w-5" />
        </div>
        <span className="text-base font-bold tracking-tight text-slate-900">SALON ADMIN</span>
      </div>

      <div className="rounded-custom border border-border bg-surface p-6 shadow-sm sm:p-8">
        <h1 className="text-xl font-bold tracking-tight text-body">Đăng nhập</h1>
        <p className="mt-1 text-sm text-on-surface-variant">
          Chào Mừng Trở Lại Với Hệ Thống Quản Lý Salon.
        </p>

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="phone">Số điện thoại</Label>
            <div className="relative">
              <Phone className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-slate-400" />
              <Input
                id="phone"
                type="tel"
                autoComplete="username"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="vd: 0901234567"
                className="pl-9"
                required
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="password">Mật khẩu</Label>
            <div className="relative">
              <LockKeyhole className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-slate-400" />
              <Input
                id="password"
                type={showPassword ? 'text' : 'password'}
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Nhập mật khẩu"
                className="pr-10 pl-9"
                required
              />
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                onClick={() => setShowPassword((v) => !v)}
                aria-label={showPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'}
                className="absolute top-1/2 right-1.5 -translate-y-1/2 text-slate-400 hover:text-slate-700"
              >
                {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </Button>
            </div>
          </div>

          <div className="flex items-center justify-between text-xs">
            <label className="flex cursor-pointer items-center gap-2 text-on-surface-variant">
              <Checkbox checked={remember} onCheckedChange={(v) => setRemember(v === true)} />
              Ghi nhớ đăng nhập
            </label>
            <Link
              href="/app/forgot-password"
              className="font-semibold text-brand hover:text-brand-hover"
            >
              Quên mật khẩu?
            </Link>
          </div>

          {error && (
            <p className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-600">
              {error}
            </p>
          )}

          <Button type="submit" disabled={loading} className="h-10 w-full text-sm">
            {loading ? 'Đang đăng nhập...' : 'Đăng nhập'}
          </Button>
        </form>
      </div>
    </div>
  );
}
