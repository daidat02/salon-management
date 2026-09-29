import type { Metadata } from "next";
import PageHeading from "@/components/dashboard/PageHeading";

export const metadata: Metadata = { title: "Cài đặt salon | SALON ADMIN" };

export default function SettingsPage() {
  return (
    <div className="space-y-4">
      <PageHeading
        title="Cài đặt salon"
        subtitle="Thông tin chi nhánh, giờ làm việc và phân quyền — điều hướng chi tiết bên trái."
      />
      <div className="rounded-custom border border-dashed border-border bg-surface p-8 text-center text-sm text-on-surface-variant">
        Layout 2 cột + form sẽ được dựng từ{" "}
        <code>preview/admin_preview/content/cai-dat.html</code>
      </div>
    </div>
  );
}
