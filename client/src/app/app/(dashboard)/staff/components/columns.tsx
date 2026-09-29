import { Tag } from "@/components/ui/tag"
import type { DashColumn } from "@/components/dashboard/DashTable"
import type { Staff } from "@/types/staff"

function formatDate(value: string | null) {
  if (!value) return '—';
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return value;
  return d.toLocaleDateString('vi-VN');
}

function getInitials(name: string) {
  return name
    .trim()
    .split(/\s+/)
    .slice(-2)
    .map((p) => p[0]?.toUpperCase() ?? '')
    .join('');
}

// Columns bám trực tiếp Staff API interface — không qua StaffItem trung gian.
// Mọi accessor/cell đều đọc field gốc: full_name, code, position, phone, commission_rate, hire_date, status.
export const staffColumns: DashColumn<Staff>[] = [
  {
    id: "staff",
    header: "Nhân viên",
    className: "min-w-[200px]",
    cell: (row) => (
      <div className="flex items-center gap-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary-container text-xs font-bold text-primary">
          {getInitials(row.full_name)}
        </div>
        <div className="min-w-0">
          <div className="truncate font-semibold text-on-surface">{row.full_name}</div>
          <div className="text-[11px] text-on-surface-variant">{row.code}</div>
        </div>
      </div>
    ),
  },
  {
    id: "position",
    header: "Vị trí",
    cell: (row) => <Tag variant="slate" size="sm" shape="rounded">{row.position}</Tag>,
  },
  { id: "phone", header: "SĐT", className: "font-mono text-xs whitespace-nowrap", accessorKey: "phone" },
  {
    id: "commission_rate",
    header: "Hoa hồng",
    className: "text-center font-medium whitespace-nowrap",
    cell: (row) => <span className="font-semibold text-on-surface">{row.commission_rate}%</span>,
  },
  {
    id: "hire_date",
    header: "Ngày vào làm",
    className: "whitespace-nowrap",
    cell: (row) => (
      <span className="text-xs text-on-surface-variant">{formatDate(row.hire_date)}</span>
    ),
  },
  {
    id: "status",
    header: "Trạng thái",
    className: "text-center",
    cell: (row) => (
      <Tag variant={row.status === "active" ? "success" : "gray"} shape="pill" size="md" dot>
        {row.status === "active" ? "Hoạt động" : "Tạm ngưng"}
      </Tag>
    ),
  },
]
