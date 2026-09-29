import Link from 'next/link';
import { Tag } from '@/components/ui/tag';
import type { DashColumn } from '@/components/dashboard/DashTable';
import type { Order, OrderStatus } from '@/types/order';

export type CustomerLookup = (
  customerId: string | null,
) => { name: string; phone: string } | undefined;

const statusMeta: Record<
  OrderStatus,
  { label: string; variant: 'success' | 'warning' | 'purple' | 'teal' | 'blue' | 'gray' | 'slate' }
> = {
  paid: { label: 'Đã thanh toán', variant: 'success' },
  pending_payment: { label: 'Chưa thanh toán', variant: 'warning' },
  draft: { label: 'Nháp', variant: 'slate' },
  serving: { label: 'Đang phục vụ', variant: 'blue' },
  completed: { label: 'Hoàn tất', variant: 'teal' },
  cancelled: { label: 'Đã hủy', variant: 'gray' },
  refunded: { label: 'Đã hoàn tiền', variant: 'purple' },
};

function formatVND(value: number): string {
  return `${value.toLocaleString('vi-VN')} đ`;
}

export function buildOrderColumns(lookup: CustomerLookup): DashColumn<Order>[] {
  return [
    {
      id: 'code',
      header: 'Mã đơn',
      className: 'font-mono font-medium text-primary min-w-[110px]',
      cell: (row) => (
        <Link href={`/app/orders/${row.id}`} className="hover:underline">
          {row.code}
        </Link>
      ),
    },
    {
      id: 'customer',
      header: 'Khách hàng',
      className: 'min-w-[180px]',
      cell: (row) => {
        const customer = lookup(row.customer_id);
        return (
          <div>
            <div className="font-semibold text-on-surface">{customer?.name ?? 'Khách lẻ'}</div>
            <div className="text-[11px] text-on-surface-variant">{customer?.phone ?? '—'}</div>
          </div>
        );
      },
    },
    {
      id: 'notes',
      header: 'Dịch Vụ / Sản phẩm',
      className: 'min-w-[180px]',
      cell: (row) => {
        const customer = lookup(row.customer_id);
        return (
          <div>
            <div className="text-[11px] text-on-surface-variant">{customer?.phone ?? '—'}</div>
          </div>
        );
      },
    },
    {
      id: 'total',
      header: 'Tổng tiền',
      className: 'text-right font-semibold text-on-surface whitespace-nowrap',
      cell: (row) => formatVND(row.total_amount),
    },
    {
      id: 'status',
      header: 'Trạng thái',
      className: 'text-center',
      cell: (row) => {
        const meta = statusMeta[row.status] ?? { label: row.status, variant: 'slate' as const };
        return (
          <Tag variant={meta.variant} shape="pill" size="md" dot>
            {meta.label}
          </Tag>
        );
      },
    },
    {
      id: 'created_at',
      header: 'Ngày tạo',
      className: 'whitespace-nowrap text-on-surface-variant text-right',
      cell: (row) =>
        row.created_at
          ? new Date(row.created_at).toLocaleString('vi-VN', {
              dateStyle: 'short',
              timeStyle: 'short',
            })
          : '—',
    },
  ];
}
