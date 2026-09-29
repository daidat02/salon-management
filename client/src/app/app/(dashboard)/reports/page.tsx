import type { Metadata } from "next";
import { Download, FileDown, FileText } from "lucide-react";
import PageHeading from "@/components/dashboard/PageHeading";

export const metadata: Metadata = { title: "Báo cáo | SALON ADMIN" };

export default function ReportsPage() {
  return (
    <div className="space-y-4">
      <PageHeading
        breadcrumbs={[
          { label: "Trang chủ", href: "#" },
          { label: "Báo cáo & Thống kê", active: true },
        ]}
        title="Báo cáo & Phân tích Doanh thu"
        subtitle="Theo dõi số liệu tài chính, hiệu suất phục vụ stylist và phân tích tăng trưởng salon"
        actions={[
          { label: "Lọc", icon: Download, variant: "outline" },
          { label: "Xuất Excel", icon: FileDown, variant: "outline" },
          { label: "Tải báo cáo", icon: FileText, variant: "primary" },
        ]}
      />
      {/* Dải filter ngày — sẽ tách thành component riêng khi dựng chi tiết */}
      <div className="flex flex-wrap items-center gap-2 text-xs text-on-surface-variant">
        <span className="rounded-lg border border-outline bg-surface px-3 py-1.5">
          01/10 - 31/10/2024
        </span>
        <span className="rounded-lg bg-primary-container px-3 py-1.5 font-semibold text-primary">
          Tháng này
        </span>
      </div>
      <div className="rounded-custom border border-dashed border-border bg-surface p-8 text-center text-sm text-on-surface-variant">
        KPI + biểu đồ + bảng sẽ được dựng từ{" "}
        <code>preview/admin_preview/content/bao-cao.html</code>
      </div>
    </div>
  );
}
