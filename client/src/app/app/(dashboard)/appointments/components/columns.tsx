import { Tag } from "@/components/ui/tag"
import type { DashColumn } from "@/components/dashboard/DashTable"
import type { AppointmentItem } from "./data"

export const appointmentColumns: DashColumn<AppointmentItem>[] = [
  { id: "time", header: "Giờ hẹn", className: "whitespace-nowrap font-medium" },
  {
    id: "customer",
    header: "Khách hàng",
    className: "min-w-[180px]",
    cell: (row) => (
      <div>
        <div className="font-semibold text-on-surface">{row.customer}</div>
        <div className="text-[11px] text-on-surface-variant">{row.phone}</div>
      </div>
    ),
  },
  {
    id: "service",
    header: "Dịch vụ",
    className: "min-w-[200px]",
    cell: (row) => (
      <div>
        <div className="font-medium text-on-surface">{row.service}</div>
        <div className="text-[11px] text-on-surface-variant">{row.duration}</div>
      </div>
    ),
  },
  { id: "staff", header: "Nhân viên", accessorKey: "staff" },
  {
    id: "status",
    header: "Trạng thái",
    className: "text-center",
    cell: (row) => {
      const variantMap = {
        pending: "warning",
        confirmed: "blue",
        serving: "primary",
        done: "success",
        cancelled: "gray",
      } as const
      const label: Record<string, string> = {
        pending: "Chờ xác nhận",
        confirmed: "Đã xác nhận",
        serving: "Đang phục vụ",
        done: "Hoàn thành",
        cancelled: "Đã hủy",
      }
      return (
        <Tag variant={variantMap[row.status]} shape="pill" size="md">
          {label[row.status]}
        </Tag>
      )
    },
  },
]
