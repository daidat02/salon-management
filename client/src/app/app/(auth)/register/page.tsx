"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  ChevronDown,
  CircleCheck,
  Eye,
  EyeOff,
  Scissors,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

const API_BASE =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

const ROLES = [
  { value: "manager", label: "Quản lý salon" },
  { value: "receptionist", label: "Lễ tân" },
  { value: "owner", label: "Chủ salon" },
];

// POST /register { full_name, phone_number, email, password, role }
// Backend tự set status=active.
export default function RegisterPage() {
  const router = useRouter();
  const [fullName, setFullName] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [role, setRole] = useState("manager");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [done, setDone] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    if (password !== confirm) {
      setError("Mật khẩu nhập lại chưa khớp.");
      return;
    }
    if (password.length < 6) {
      setError("Mật khẩu phải từ 6 ký tự trở lên.");
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          full_name: fullName.trim(),
          phone_number: phone.trim(),
          email: email.trim(),
          password,
          role,
        }),
      });
      const json = await res.json().catch(() => null);
      if (!res.ok)
        throw new Error(json?.message ?? json?.error ?? "Đăng ký thất bại");
      setDone(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Đăng ký thất bại");
    } finally {
      setLoading(false);
    }
  }

  if (done) {
    return (
      <div className="rounded-custom border border-border bg-surface p-6 text-center shadow-sm sm:p-8">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-emerald-50 text-emerald-600">
          <CircleCheck className="h-6 w-6" />
        </div>
        <h1 className="mt-4 text-xl font-bold tracking-tight text-body">
          Đăng ký thành công
        </h1>
        <p className="mt-1 text-sm text-on-surface-variant">
          Tài khoản <span className="font-semibold text-body">{phone}</span> đã
          được tạo. Hãy đăng nhập để bắt đầu.
        </p>
        <Button
          onClick={() => router.push("/app/login")}
          className="mt-6 h-10 w-full text-sm"
        >
          Đăng nhập ngay
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2.5 lg:hidden">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-brand text-white shadow-sm">
          <Scissors className="h-5 w-5" />
        </div>
        <span className="text-base font-bold tracking-tight text-slate-900">
          SALON ADMIN
        </span>
      </div>

      <div className="rounded-custom border border-border bg-surface p-6 shadow-sm sm:p-8">
        <h1 className="text-xl font-bold tracking-tight text-body">
          Tạo tài khoản
        </h1>
        <p className="mt-1 text-sm text-on-surface-variant">
          Đăng ký tài khoản quản lý salon của bạn.
        </p>

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="fullName">Họ và tên</Label>
            <Input
              id="fullName"
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              placeholder="vd: Nguyễn Minh Anh"
              required
            />
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="phone">Số điện thoại</Label>
              <Input
                id="phone"
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="vd: 0901234567"
                required
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="role">Vai trò</Label>
              <div className="relative">
                <select
                  id="role"
                  value={role}
                  onChange={(e) => setRole(e.target.value)}
                  className="h-8 w-full appearance-none rounded-lg border border-input bg-transparent pr-8 pl-2.5 text-sm text-body transition-colors outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50"
                >
                  {ROLES.map((r) => (
                    <option key={r.value} value={r.value}>
                      {r.label}
                    </option>
                  ))}
                </select>
                <ChevronDown className="pointer-events-none absolute top-1/2 right-2.5 h-4 w-4 -translate-y-1/2 text-slate-400" />
              </div>
            </div>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="vd: minhanh@salon.vn"
              required
            />
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="password">Mật khẩu</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="new-password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Tối thiểu 6 ký tự"
                  className="pr-10"
                  required
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => setShowPassword((v) => !v)}
                  aria-label={showPassword ? "Ẩn mật khẩu" : "Hiện mật khẩu"}
                  className="absolute top-1/2 right-1.5 -translate-y-1/2 text-slate-400 hover:text-slate-700"
                >
                  {showPassword ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="confirm">Nhập lại mật khẩu</Label>
              <Input
                id="confirm"
                type={showPassword ? "text" : "password"}
                autoComplete="new-password"
                value={confirm}
                onChange={(e) => setConfirm(e.target.value)}
                placeholder="Nhập lại mật khẩu"
                required
              />
            </div>
          </div>

          {error && (
            <p className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-600">
              {error}
            </p>
          )}

          <Button type="submit" disabled={loading} className="h-10 w-full text-sm">
            {loading ? "Đang tạo tài khoản..." : "Đăng ký"}
          </Button>
        </form>

        <p className="mt-5 text-center text-xs text-on-surface-variant">
          Đã có tài khoản?{" "}
          <Link
            href="/app/login"
            className="font-semibold text-brand hover:text-brand-hover"
          >
            Đăng nhập
          </Link>
        </p>
      </div>
    </div>
  );
}
