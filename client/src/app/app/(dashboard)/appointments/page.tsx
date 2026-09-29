import type { Metadata } from "next";
import { Clock, FileDown, Plus, Users, Wallet } from "lucide-react";
import PageHeading from "@/components/dashboard/PageHeading";
import StatCard from "@/components/dashboard/StatCard";
import DashToolbar from "@/components/dashboard/DashToolbar";
import AppointmentsView from "./components/AppointmentsView";

export const metadata: Metadata = { title: "Lịch hẹn | SALON ADMIN" };

export default function AppointmentsPage() {
  return (
    <div className="space-y-4">
      <PageHeading
        breadcrumbs={[
          { label: "Trang chủ", href: "#" },
          { label: "Lịch hẹn", active: true },
        ]}
        title="Lịch hẹn"
        subtitle="Quản lý lịch hẹn, phân công kỹ thuật viên và theo dõi trạng thái phục vụ"
        actions={[
          { label: "Xuất Excel / CSV", icon: FileDown, variant: "outline" },
          { label: "Đặt lịch mới", icon: Plus, variant: "primary" },
        ]}
      />
      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard label="Tổng hôm nay" value={<>12 <span className="text-xs font-normal text-on-surface-variant">lịch</span></>} sub={<span className="font-medium text-warning">● 4 đang chờ tiếp nhận</span>} icon={Clock} />
        <StatCard label="Đang phục vụ" value={<>3 <span className="text-xs font-normal text-on-surface-variant">khách</span></>} sub={<span className="font-medium text-primary">● 3 thợ đang bận</span>} icon={Users} iconWrapClassName="bg-blue-50 text-primary" />
        <StatCard label="Đã hoàn thành" value={<>6 <span className="text-xs font-normal text-on-surface-variant">lịch</span></>} sub={<span className="font-medium text-success">● 50% chỉ tiêu ngày</span>} icon={Clock} iconWrapClassName="bg-emerald-50 text-success" />
        <StatCard label="Đã hủy / Vắng" value={<>1 <span className="text-xs font-normal text-on-surface-variant">lịch</span></>} sub={<span className="font-medium text-error">● Tỷ lệ hủy 8.3%</span>} icon={Wallet} iconWrapClassName="bg-red-50 text-error" />
      </section>
      <DashToolbar
        searchPlaceholder="Tìm theo tên, SĐT, mã lịch..."
        filters={[
          {
            id: "status",
            placeholder: "Tất cả trạng thái",
            widthClass: "sm:w-44",
            options: ["Tất cả trạng thái", "Chờ xác nhận", "Đã xác nhận", "Đang phục vụ", "Đã hoàn thành", "Đã hủy"],
          },
          {
            id: "staff",
            placeholder: "Tất cả nhân viên",
            widthClass: "sm:w-40",
            options: ["Tất cả nhân viên", "Hương Stylist", "Tuấn Stylist", "Linh Barber", "Ngọc Spa", "Trang Eyelash"],
          },
          {
            id: "service",
            placeholder: "Tất cả dịch vụ",
            widthClass: "sm:w-44",
            options: ["Tất cả dịch vụ", "Cắt & Tạo kiểu", "Nhuộm & Phục hồi", "Uốn & Duỗi", "Gội & Dưỡng sinh"],
          },
        ]}
      />
      <AppointmentsView />
    </div>
  );
}
