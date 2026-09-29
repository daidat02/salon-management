'use client';

import * as React from 'react';
import { Search } from 'lucide-react';
import { cn } from 'cn';
import { DashDialog, DashDialogFooter } from '@/components/dashboard/DashDialog';
import { DashSelect } from '@/components/dashboard/DashSelect';
import DashTable, { type DashColumn } from '@/components/dashboard/DashTable';
import useProduct from '@/hooks/use-product';
import type { Category } from '@/types/category';
import type { Product } from '@/types/product';
import { formatVND } from './voucherUtils';

type ProductPickerDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  categories: Category[];
  /** ID sản phẩm đã có trong phiếu — ẩn khỏi danh sách chọn. */
  addedIds: Set<string>;
  onConfirm: (products: Product[]) => void;
};

const pickerColumns: DashColumn<Product>[] = [
  {
    id: 'product',
    header: 'Sản phẩm / Quy cách',
    className: 'min-w-[220px]',
    cell: (row) => (
      <div className="flex items-center gap-3">
        <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-outline bg-surface-container text-sm font-bold text-primary">
          {row.name.charAt(0).toUpperCase()}
        </span>
        <div className="min-w-0">
          <div className="truncate font-semibold text-on-surface">{row.name}</div>
          <div className="text-[11px] text-on-surface-variant">
            SKU: <strong className="font-mono text-on-surface">{row.sku}</strong>
          </div>
        </div>
      </div>
    ),
  },
  {
    id: 'unit',
    header: 'ĐVT',
    className: 'text-center',
    cell: (row) => (
      <span className="rounded bg-surface-container px-2 py-0.5 text-[11px] font-medium text-on-surface">
        {row.unit}
      </span>
    ),
  },
  {
    id: 'cost_price',
    header: 'Giá vốn (đ)',
    className: 'text-right font-medium whitespace-nowrap',
    cell: (row) => formatVND(row.cost_price),
  },
  {
    id: 'stock_quantity',
    header: 'Tồn kho',
    className: 'text-center',
    cell: (row) => (
      <span
        className={cn(
          'font-semibold',
          row.stock_quantity <= row.min_stock ? 'text-error' : 'text-success',
        )}
      >
        {row.stock_quantity} {row.unit}
      </span>
    ),
  },
];

export default function ProductPickerDialog({
  open,
  onOpenChange,
  categories,
  addedIds,
  onConfirm,
}: ProductPickerDialogProps) {
  const [query, setQuery] = React.useState('');
  const [categoryId, setCategoryId] = React.useState('');
  const [selectedIds, setSelectedIds] = React.useState<Set<string>>(new Set());

  const { data, isLoading } = useProduct({ page: 1, pageSize: 100 });
  const all = React.useMemo(
    () => (data?.data ?? []).filter((p): p is NonNullable<typeof p> => p != null),
    [data],
  );

  const rows = React.useMemo(() => {
    const q = query.trim().toLowerCase();
    return all.filter(
      (p) =>
        !addedIds.has(p.id) &&
        (!categoryId || p.category_id === categoryId) &&
        (q === '' ||
          p.name.toLowerCase().includes(q) ||
          p.sku.toLowerCase().includes(q)),
    );
  }, [all, addedIds, categoryId, query]);

  const handleConfirm = () => {
    if (selectedIds.size === 0) return;
    // Trả về object đầy đủ để form cache lại (list tìm kiếm của form có thể không chứa SP này).
    onConfirm(all.filter((p) => selectedIds.has(p.id)));
    onOpenChange(false);
  };

  return (
    <DashDialog
      open={open}
      onOpenChange={onOpenChange}
      title="Chọn sản phẩm từ danh mục"
      description="Tích chọn các sản phẩm cần thêm vào phiếu kho."
      maxWidth="max-w-[720px]"
      footer={
        <DashDialogFooter
          onCancel={() => onOpenChange(false)}
          onConfirm={handleConfirm}
          cancelText="Đóng"
          confirmText={
            selectedIds.size > 0 ? `Thêm ${selectedIds.size} sản phẩm` : 'Thêm sản phẩm'
          }
          confirmDisabled={selectedIds.size === 0}
        />
      }
    >
      <div className="flex flex-col gap-2 sm:flex-row">
        <div className="relative w-full">
          <Search className="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-on-surface-variant" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Tìm theo tên, SKU..."
            className="w-full rounded-lg bg-surface-container py-2 pr-3 pl-9 text-xs text-on-surface shadow-inner placeholder:text-on-surface-variant focus:ring-1 focus:ring-primary focus:outline-none"
          />
        </div>
        <DashSelect
          value={categoryId}
          onValueChange={setCategoryId}
          placeholder="Tất cả danh mục"
          options={[
            { label: 'Tất cả danh mục', value: '' },
            ...categories.map((c) => ({ label: c.name, value: c.id })),
          ]}
          triggerClassName="bg-surface-container sm:w-52"
        />
      </div>

      {isLoading ? (
        <div className="py-10 text-center text-xs text-on-surface-variant">
          Đang tải danh sách sản phẩm…
        </div>
      ) : (
        <DashTable
          columns={pickerColumns}
          data={rows}
          getRowId={(row) => row.id}
          selectedIds={selectedIds}
          onSelectionChange={setSelectedIds}
          emptyText="Không tìm thấy sản phẩm phù hợp."
        />
      )}
    </DashDialog>
  );
}
