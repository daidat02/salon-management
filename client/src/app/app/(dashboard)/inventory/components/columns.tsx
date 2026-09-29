import { History, Pencil, SquarePlus } from 'lucide-react';
import { cn } from 'cn';
import { Tag } from '@/components/ui/tag';
import type { DashColumn } from '@/components/dashboard/DashTable';
import type { Product } from '@/types/product';
import {
  formatStockQuantity,
  formatVND,
  getDisplayMinStock,
  getDisplayStock,
  getInventoryStatus,
  getStockProgress,
  getStockUnit,
  type InventoryStatus,
} from './inventoryUtils';

function Shelves(props: React.ComponentProps<'svg'>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" {...props}>
      <path d="M3 5h18M3 12h18M3 19h18M5 5v14M19 5v14" />
    </svg>
  );
}

function StatusBadge({ status }: { status: InventoryStatus }) {
  if (status === 'low')
    return (
      <Tag variant="error" shape="pill" size="md" dot>
        Sắp hết
      </Tag>
    );
  if (status === 'in_stock')
    return (
      <Tag variant="success" shape="pill" size="md" dot>
        Còn hàng
      </Tag>
    );
  if (status === 'out_of_stock')
    return (
      <Tag variant="gray" shape="pill" size="md" dot>
        Hết hàng
      </Tag>
    );
  return (
    <Tag variant="warning" shape="pill" size="md" dot>
      Sắp hết
    </Tag>
  );
}

export const inventoryColumns: DashColumn<Product>[] = [
  {
    id: 'sku',
    header: 'Mã SKU',
    className:
      'font-mono text-primary font-semibold text-[11px] min-w-[120px] hover:underline hover:text-primary/80 transition-colors cursor-pointer',

    accessorKey: 'sku',
  },
  {
    id: 'product',
    header: 'Sản phẩm & Quy cách',
    className: 'min-w-[220px]',
    cell: (row) => (
      <div className="flex items-center gap-3">
        <div>
          <span className="hover:text-primary cursor-pointer font-semibold text-on-surface">
            {row.name}
          </span>
          <p className="text-[11px] text-on-surface-variant">
            Đơn vị: {row.unit || '—'}
            {row.net_amount ? ` • Dung tích: ${row.net_amount}${row.net_unit}` : ''}
          </p>
        </div>
      </div>
    ),
  },
  {
    id: 'category',
    header: 'Phân loại',
    cell: () => (
      <Tag variant="slate" size="sm" shape="rounded">
        Chưa phân loại
      </Tag>
    ),
  },
  {
    id: 'stock',
    header: 'Tồn kho & Mức an toàn',
    className: 'min-w-[160px]',
    cell: (row) => {
      const status = getInventoryStatus(row);
      const progress = getStockProgress(row);
      const unit = getStockUnit(row);
      return (
        <div className="flex flex-col gap-1">
          <div className="flex items-center justify-between text-[11px]">
            <span
              className={cn(
                'font-bold',
                status === 'low' || status === 'out_of_stock' ? 'text-error' : 'text-success',
              )}
            >
              {formatStockQuantity(getDisplayStock(row))} {unit}
            </span>
            <span className="text-on-surface-variant">
              Min: {formatStockQuantity(getDisplayMinStock(row))} {unit}
            </span>
          </div>
          <div className="h-1.5 w-full overflow-hidden rounded-full bg-slate-100">
            <div
              className={cn(
                'h-1.5 rounded-full',
                status === 'low'
                  ? 'bg-error'
                  : status === 'out_of_stock'
                    ? 'bg-gray-400'
                    : progress < 40
                      ? 'bg-amber-500'
                      : 'bg-success',
              )}
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>
      );
    },
  },
  {
    id: 'cost',
    header: 'Giá nhập',
    className: 'text-right font-medium',
    cell: (row) => (row.cost_price === 0 ? 'Chưa nhập' : formatVND(row.cost_price)),
  },
  {
    id: 'price',
    header: 'Giá niêm yết',
    className: 'text-right font-semibold text-primary',
    cell: (row) => (row.sell_price === 0 ? 'Dùng nội bộ' : formatVND(row.sell_price)),
  },
  {
    id: 'status',
    header: 'Trạng thái',
    className: 'text-center',
    cell: (row) => <StatusBadge status={getInventoryStatus(row)} />,
  },
  {
    id: 'actions',
    header: 'Thao tác',
    className: 'text-center',
    cell: () => (
      <div className="flex items-center justify-center gap-1 text-on-surface-variant">
        <button
          className="hover:text-primary hover:bg-primary-container rounded p-1 transition-colors"
          title="Nhập thêm"
        >
          <SquarePlus className="h-[17px] w-[17px]" />
        </button>
        <button
          className="hover:text-primary hover:bg-surface-container rounded p-1 transition-colors"
          title="Điều chỉnh kho"
        >
          <Pencil className="h-[17px] w-[17px]" />
        </button>
        <button
          className="hover:text-primary hover:bg-surface-container rounded p-1 transition-colors"
          title="Lịch sử thẻ kho"
        >
          <History className="h-[17px] w-[17px]" />
        </button>
      </div>
    ),
  },
];
