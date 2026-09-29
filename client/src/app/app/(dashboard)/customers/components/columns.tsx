import { Tag } from '@/components/ui/tag';
import type { DashColumn } from '@/components/dashboard/DashTable';
import type { Customer } from '@/types/customer';

function formatCurrency(value: number) {
  return value.toLocaleString('vi-VN') + ' đ';
}

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

function getGenderLabel(gender: Customer['gender']) {
  switch (gender) {
    case 'male':
      return 'Nam';
    case 'female':
      return 'Nữ';
    default:
      return 'Khác';
  }
}

// Columns bám trực tiếp Customer API interface — không qua CustomerItem trung gian.
// Mọi accessor/cell đều đọc field gốc: full_name, phone, gender, birth_date, note, total_spent, total_visits, deleted_at.
export const customerColumns: DashColumn<Customer>[] = [
  {
    id: 'customer',
    header: 'Khách hàng',
    className: 'min-w-[220px]',
    cell: (row) => (
      <div className="flex items-center gap-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary-container text-xs font-bold text-primary">
          {getInitials(row.full_name)}
        </div>
        <div className="min-w-0">
          <div className="truncate font-semibold text-on-surface">{row.full_name}</div>
          <div className="text-[11px] text-on-surface-variant">{row.phone}</div>
        </div>
      </div>
    ),
  },
  {
    id: 'gender',
    header: 'Giới tính',
    className: 'text-center',
    cell: (row) => (
      <Tag
        variant={row.gender === 'female' ? 'purple' : row.gender === 'male' ? 'blue' : 'slate'}
        shape="pill"
        size="sm"
      >
        {getGenderLabel(row.gender)}
      </Tag>
    ),
  },
  {
    id: 'birth_date',
    header: 'Ngày sinh',
    className: 'whitespace-nowrap',
    cell: (row) => (
      <span className="text-xs text-on-surface-variant">{formatDate(row.birth_date)}</span>
    ),
  },
  {
    id: 'note',
    header: 'Ghi chú',
    className: 'min-w-[180px] max-w-[260px]',
    cell: (row) => (
      <span className="line-clamp-1 text-xs text-on-surface-variant" title={row.note}>
        {row.note?.trim() ? row.note : '—'}
      </span>
    ),
  },
  {
    id: 'total_spent',
    header: 'Tổng chi tiêu',
    className: 'text-right font-semibold text-on-surface whitespace-nowrap',
    cell: (row) => formatCurrency(row.total_spent),
  },
  {
    id: 'total_visits',
    header: 'Lượt ghé',
    className: 'text-center font-medium',
    cell: (row) => <span className="font-semibold text-on-surface">{row.total_visits}</span>,
  },
  {
    id: 'status',
    header: 'Trạng thái',
    className: 'text-center',
    cell: (row) => {
      const isActive = !row.deleted_at;
      return (
        <Tag variant={isActive ? 'success' : 'gray'} shape="pill" size="md" dot>
          {isActive ? 'Hoạt động' : 'Đã xóa'}
        </Tag>
      );
    },
  },
  {
    id: 'date',
    header: 'Ngày tạo',
    className: 'whitespace-nowrap text-on-surface-variant',
    accessorKey: 'created_at',
    cell: (row) => (
      <span className="text-xs text-on-surface-variant">{formatDate(row.created_at)}</span>
    ),
  },
];
