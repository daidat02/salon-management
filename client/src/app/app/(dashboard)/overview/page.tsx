import type { Metadata } from "next";
import { Calendar, Clock, TrendingUp, Users } from "lucide-react";
import PageHeading from "@/components/dashboard/PageHeading";
import StatCard from "@/components/dashboard/StatCard";

export const metadata: Metadata = { title: "Tổng quan | SALON ADMIN" };

export default function OverviewPage() {
  return (
    <div className="space-y-4">
      <PageHeading
        title="Chào buổi sáng, Minh Anh 👋"
        subtitle={
          <p className="mt-1 text-sm text-on-surface-variant">
            Hôm nay bạn có <span className="font-semibold text-primary">12 lịch hẹn</span> và{" "}
            <span className="font-semibold text-warning">8 đơn hàng</span> cần xử lý
          </p>
        }
      />
      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard label="Doanh thu hôm nay" value="28.500.000₫" sub={<span className="font-medium text-success">↗ +12% so với hôm qua</span>} icon={TrendingUp} iconWrapClassName="bg-[#EFF6FF] text-primary" />
        <StatCard label="Lợi nhuận ròng" value="12.800.000₫" sub={<span className="font-medium text-success">↗ +8% so với hôm qua</span>} icon={TrendingUp} iconWrapClassName="bg-emerald-50 text-success" />
        <StatCard label="Lịch hẹn hôm nay" value={<>12 <span className="text-xs font-normal text-on-surface-variant">lịch</span></>} sub="4 đang chờ tiếp nhận" icon={Calendar} />
        <StatCard label="Khách hàng mới" value="8" sub="tháng này +12%" icon={Users} iconWrapClassName="bg-blue-50 text-primary" />
      </section>
      <section className="grid grid-cols-1 gap-4 lg:grid-cols-3">
        <div className="rounded-xl border border-outline bg-surface p-5 shadow-sm lg:col-span-2">
          <h3 className="text-sm font-semibold text-on-surface">Doanh thu 7 ngày qua</h3>
          <div className="mt-4 flex h-40 items-center justify-center rounded-lg bg-surface-container-low text-xs text-on-surface-variant">Biểu đồ sẽ nối API — preview tong-quan.html</div>
        </div>
        <div className="rounded-xl border border-outline bg-surface p-5 shadow-sm">
          <h3 className="text-sm font-semibold text-on-surface">Lịch hẹn sắp tới</h3>
          <div className="mt-4 space-y-2 text-xs text-on-surface-variant">
            <div className="flex items-center gap-2"><Clock className="h-4 w-4 text-primary" />09:30 — Chị Thu Thảo (Balayage)</div>
            <div className="flex items-center gap-2"><Clock className="h-4 w-4 text-primary" />10:00 — Chị Kim Oanh (Gội dưỡng sinh)</div>
            <div className="flex items-center gap-2"><Clock className="h-4 w-4 text-primary" />13:00 — Chị Lan Hương (Nối mi)</div>
          </div>
        </div>
      </section>
    </div>
  );
}
