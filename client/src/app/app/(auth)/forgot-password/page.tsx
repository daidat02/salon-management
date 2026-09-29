"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ArrowLeft, CircleCheck, Phone, Scissors } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

// Backend chưa có endpoint quên mật khẩu — flow hiện tại: nhập SĐT,
// hiển thị trạng thái đã gửi (sẵn sàng nối API /forgot-password khi có).
export default function ForgotPasswordPage() {
  const router = useRouter();
  const [phone, setPhone] = useState("");
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    // TODO: gọi POST /forgot-password khi backend có endpoint.
    await new Promise((r) => setTimeout(r, 800));
    setLoading(false);
    setSent(true);
  }

  if (sent) {
    return (
      <div className="rounded-custom border border-border bg-surface p-6 text-center shadow-sm sm:p-8">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-brand-light text-brand">
          <CircleCheck className="h-6 w-6" />
        </div>
        <h1 className="mt-4 text-xl font-bold tracking-tight text-body">
          Đã gửi yêu cầu
        </h1>
        <p className="mt-1 text-sm text-on-surface-variant">
          Nếu số <span className="font-semibold text-body">{phone}</span> tồn
          tại trong hệ thống, mã đặt lại mật khẩu sẽ được gửi qua SMS trong vài
          phút.
        </p>
        <div className="mt-6 space-y-2.5">
          <Button
            variant="outline"
            onClick={() => setSent(false)}
            className="h-10 w-full text-sm"
          >
            Gửi lại mã
          </Button>
          <Button
            onClick={() => router.push("/app/login")}
            className="h-10 w-full text-sm"
          >
            Quay lại đăng nhập
          </Button>
        </div>
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
          Quên mật khẩu
        </h1>
        <p className="mt-1 text-sm text-on-surface-variant">
          Nhập số điện thoại đã đăng ký — chúng tôi sẽ gửi mã đặt lại cho bạn.
        </p>

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="phone">Số điện thoại</Label>
            <div className="relative">
              <Phone className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-slate-400" />
              <Input
                id="phone"
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="vd: 0901234567"
                className="pl-9"
                required
              />
            </div>
          </div>

          <Button type="submit" disabled={loading} className="h-10 w-full text-sm">
            {loading ? "Đang gửi..." : "Gửi mã đặt lại"}
          </Button>
        </form>

        <Link
          href="/app/login"
          className="mt-5 inline-flex items-center gap-1.5 text-xs font-semibold text-brand hover:text-brand-hover"
        >
          <ArrowLeft className="h-3.5 w-3.5" />
          Quay lại đăng nhập
        </Link>
      </div>
    </div>
  );
}
