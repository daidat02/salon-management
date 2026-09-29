import type { Metadata } from "next";
import { CalendarOff, Star, Users, Wrench } from "lucide-react";
import StatCard from "@/components/dashboard/StatCard";
import StaffHeading from "./components/StaffHeading";
import StaffToolbar from "./components/StaffToolbar";
import StaffTable from "./components/StaffTable";

export const metadata: Metadata = { title: "Nhân viên | SALON ADMIN" };

export default function StaffPage() {
  return (
    <div className="space-y-4">
      <StaffHeading />

      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          label="Tổng nhân viên"
          value={
            <span className="flex items-baseline gap-2">
              <span>18</span>
              <span className="inline-flex items-center gap-1 rounded-full border border-emerald-100 bg-emerald-50 px-2 py-0.5 text-xs font-semibold text-success">↗ +2 tháng này</span>
            </span>
          }
          sub="12 Stylist, 4 Kỹ thuật viên, 2 Lễ tân"
          icon={Users}
          iconWrapClassName="bg-primary-container/40 text-primary"
        />
        <StatCard
          label="Stylist đang làm"
          value={
            <span className="flex items-baseline gap-2">
              <span>14</span>
              <span className="rounded-full bg-primary-container px-2 py-0.5 text-xs font-semibold text-primary">77.8% ca trực</span>
            </span>
          }
          sub={
            <span className="inline-flex items-center gap-1.5">
              <span className="h-2 w-2 animate-pulse rounded-full bg-success" />9 thợ đang bận làm tóc cho khách
            </span>
          }
          icon={Wrench}
          iconWrapClassName="bg-emerald-50 text-success"
        />
        <StatCard
          label="Nghỉ hôm nay"
          value={
            <span className="flex items-baseline gap-2">
              <span>3</span>
              <span className="rounded-full border border-amber-200 bg-amber-50 px-2 py-0.5 text-xs font-semibold text-warning">Có phép</span>
            </span>
          }
          sub="1 stylist nghỉ ốm, 2 off định kỳ"
          icon={CalendarOff}
          iconWrapClassName="bg-amber-50 text-warning"
        />
        <StatCard
          label="Rating trung bình"
          value={
            <>
              4.8 <span className="text-lg font-normal text-on-surface-variant">/ 5.0</span>
            </>
          }
          sub="Dựa trên 1.420 đánh giá của khách — ↗ +0.2 so với tháng trước"
          icon={Star}
          iconWrapClassName="bg-yellow-50 text-yellow-600"
        />
      </section>

      <StaffToolbar />
      <StaffTable />
    </div>
  );
}
