import { CircleCheck, Scissors } from "lucide-react";

// Layout dùng chung cho các trang auth (login / register / forgot-password).
// Route group (auth) không ảnh hưởng URL: vẫn là /app/login, /app/register...
// Không dùng DashHeader/DashSidebar — full màn hình split-screen theo
// design system SALON ADMIN (primary #4F7CAC, Inter).
const HIGHLIGHTS = [
  "Lịch hẹn & nhắc lịch tự động",
  "Thu ngân, đơn hàng, tồn kho",
  "Báo cáo doanh thu theo thời gian thực",
];

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="theme-app grid min-h-screen bg-background font-sans text-body lg:grid-cols-[1.1fr_1fr]">
      {/* Panel thương hiệu */}
      <div className="relative hidden flex-col justify-between overflow-hidden bg-brand p-10 text-white lg:flex">
        <div
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              "radial-gradient(600px 300px at 20% 10%, rgba(255,255,255,0.18), transparent 60%), radial-gradient(500px 400px at 90% 90%, rgba(0,0,0,0.25), transparent 60%)",
          }}
        />
        <div className="relative flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white/15 ring-1 ring-white/30">
            <Scissors className="h-5 w-5" />
          </div>
          <div>
            <p className="text-base leading-tight font-bold tracking-tight">
              SALON ADMIN
            </p>
            <p className="text-xs text-white/70">Chi nhánh Trung tâm</p>
          </div>
        </div>

        <div className="relative space-y-6">
          <h1 className="text-3xl leading-tight font-bold tracking-tight">
            Quản lý salon
            <br />
            gọn nhẹ, chuẩn salon.
          </h1>
          <ul className="space-y-3">
            {HIGHLIGHTS.map((item) => (
              <li key={item} className="flex items-center gap-2.5 text-sm text-white/90">
                <CircleCheck className="h-4 w-4 shrink-0 text-white" />
                {item}
              </li>
            ))}
          </ul>
        </div>

        <p className="relative text-xs text-white/60">
          © 2026 Salon Project — Dành cho quản lý & nhân viên salon
        </p>
      </div>

      {/* Vùng form */}
      <div className="custom-scrollbar flex min-h-0 items-center justify-center overflow-y-auto p-6 sm:p-10">
        <div className="w-full max-w-sm">{children}</div>
      </div>
    </div>
  );
}
