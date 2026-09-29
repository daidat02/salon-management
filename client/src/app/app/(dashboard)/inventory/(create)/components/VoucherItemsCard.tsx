'use client';

import * as React from 'react';
import { ListFilter, PlusCircle, QrCode, Search, Trash2 } from 'lucide-react';
import { cn } from 'cn';
import DashTable, { type DashColumn } from '@/components/dashboard/DashTable';
import ProductPickerDialog from './ProductPickerDialog';
import type { Category } from '@/types/category';
import type { Product } from '@/types/product';
import {
  recalcLinePrice,
  toSuggestion,
  type VoucherItem,
  type VoucherSuggestion,
  type VoucherType,
} from './voucherData';
import {
  formatConversionDetail,
  formatConvertedQty,
  formatVND,
  type VoucherLine,
} from './voucherUtils';

type VoucherItemsCardProps = {
  lines: VoucherLine[];
  voucherType: VoucherType;
  /** Kết quả tìm kiếm sản phẩm từ API (đã lọc theo query ở form cha). */
  products: Product[];
  isProductsLoading?: boolean;
  query: string;
  onQueryChange: (q: string) => void;
  /** Danh mục loại product để lọc gợi ý. */
  categories: Category[];
  categoryFilter: string;
  onCategoryFilterChange: (id: string) => void;
  onUpdateItem: (productId: string, patch: Partial<VoucherItem>) => void;
  onRemoveItem: (productId: string) => void;
  onAddItem: (productId: string) => void;
  onAddItems: (products: Product[]) => void;
  onPriceChange?: (productId: string, newPrice: number) => void;
};

function toFilteredSuggestions(
  products: Product[],
  addedIds: Set<string>,
  categoryFilter: string,
): VoucherSuggestion[] {
  return products
    .filter((p) => !addedIds.has(p.id))
    .filter((p) => !categoryFilter || p.category_id === categoryFilter)
    .slice(0, 6)
    .map(toSuggestion);
}

export default function VoucherItemsCard({
  lines,
  voucherType,
  products,
  isProductsLoading,
  query,
  onQueryChange,
  categories,
  categoryFilter,
  onUpdateItem,
  onRemoveItem,
  onAddItem,
  onAddItems,
  onPriceChange,
}: VoucherItemsCardProps) {
  const [focused, setFocused] = React.useState(false);
  const [scanMode, setScanMode] = React.useState(false);
  const [pickerOpen, setPickerOpen] = React.useState(false);

  const addedIds = React.useMemo(() => new Set(lines.map((l) => l.productId)), [lines]);
  const suggestions = React.useMemo(
    () => toFilteredSuggestions(products, addedIds, categoryFilter),
    [products, addedIds, categoryFilter],
  );

  const totalQty = lines.reduce((s, l) => s + l.qty, 0);

  const rows = React.useMemo(
    () => lines.map((line, i) => ({ ...line, stt: i + 1, stripe: i % 2 === 1 })),
    [lines],
  );

  const columns: DashColumn<(typeof rows)[number]>[] = [
    {
      id: 'stt',
      header: 'STT',
      className: 'w-10 text-center font-medium text-on-surface-variant',
      accessorKey: 'stt',
    },
    {
      id: 'product',
      header: 'Sản phẩm / Quy cách',
      className: 'min-w-[200px]',
      cell: (row) => (
        <div className="flex items-center gap-3">
          <div className="min-w-0">
            <div className="truncate font-semibold text-on-surface">{row.name}</div>
            <div className="flex items-center gap-2 text-[11px] text-on-surface-variant">
              <span>
                SKU: <strong className="font-mono text-on-surface">{row.sku}</strong>
              </span>
              <span>•</span>
              <span>{row.spec}</span>
            </div>
          </div>
        </div>
      ),
    },
    {
      id: 'unit',
      header: 'Đơn vị xuất',
      className: 'min-w-[140px]',
      cell: (row) => (
        <div onClick={(e) => e.stopPropagation()}>
          <select
            value={row.selectedUnit}
            onChange={(e) => {
              const nextUnit = e.target.value;
              // Đổi đơn vị -> tính lại đơn giá theo hệ số quy đổi.
              // Dùng dữ liệu của dòng (row) thay vì tra list tìm kiếm,
              // vì list tìm kiếm có thể không còn chứa sản phẩm này.
              const nextPrice = recalcLinePrice(row.costPrice, row.unitOptions, row.unit, nextUnit);
              onUpdateItem(row.productId, { selectedUnit: nextUnit, unitPrice: nextPrice });
            }}
            className="h-8 w-full rounded bg-surface-container px-1.5 py-1 text-xs font-semibold text-on-surface focus:outline-none focus:ring-1 focus:ring-primary"
          >
            {row.unitOptions.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>
      ),
    },

    {
      id: 'qty',
      header: 'Số lượng',
      className: 'w-20 text-center',
      cell: (row) => {
        const overStock = voucherType === 'outbound' && row.convertedQty > row.stock;
        return (
          <div className="flex flex-col items-center gap-0.5">
            <input
              type="number"
              min={0}
              step="0.01"
              value={row.qty}
              onClick={(e) => e.stopPropagation()}
              onChange={(e) =>
                onUpdateItem(row.productId, { qty: Math.max(0, Number(e.target.value)) })
              }
              className={cn(
                'w-16 rounded bg-surface-container px-1.5 py-1 text-center text-xs font-semibold focus:ring-1 focus:outline-none',
                overStock ? 'text-error ring-1 ring-error focus:ring-error' : 'focus:ring-primary',
              )}
            />
            {voucherType === 'outbound' && (
              <span
                className={cn(
                  'text-[10px]',
                  overStock ? 'font-semibold text-error' : 'text-on-surface-variant',
                )}
              >
                Tồn: {row.stock.toLocaleString('vi-VN')} {row.baseUnit}
              </span>
            )}
          </div>
        );
      },
    },
    {
      id: 'converted',
      header: 'Quy đổi chuẩn',
      className: 'min-w-[130px] whitespace-nowrap',
      cell: (row) => (
        <div>
          <span className="font-semibold text-primary">{formatConvertedQty(row)}</span>
          <span className="block text-xs text-on-surface-variant">
            {formatConversionDetail(row)}
          </span>
        </div>
      ),
    },
    {
      id: 'price',
      header: 'Đơn giá (đ)',
      className: 'text-right font-medium',
      cell: (row) =>
        voucherType === 'outbound' ? (
          <span className="font-semibold text-on-surface">{formatVND(row.price)}</span>
        ) : (
          <input
            type="number"
            min={0}
            value={row.price}
            onClick={(e) => e.stopPropagation()}
            onChange={(e) => onPriceChange?.(row.productId, Math.max(0, Number(e.target.value)))}
            className="w-20 rounded bg-surface-container px-1.5 py-1 text-center text-xs font-semibold focus:ring-1 focus:ring-primary focus:outline-none"
          />
        ),
    },
    {
      id: 'discount',
      header: 'CK %',
      className: 'w-16 text-center',
      cell: (row) => (
        <input
          type="number"
          min={0}
          max={100}
          value={row.discount}
          onClick={(e) => e.stopPropagation()}
          onChange={(e) =>
            onUpdateItem(row.productId, {
              discount: Math.min(100, Math.max(0, Number(e.target.value))),
            })
          }
          className="w-12 rounded bg-surface-container px-1 py-1 text-center text-xs text-on-surface-variant focus:ring-1 focus:ring-primary focus:outline-none"
        />
      ),
    },
    // {
    //   id: 'vat',
    //   header: 'VAT %',
    //   className: 'w-16 text-center font-medium text-on-surface-variant',
    //   cell: () => '8%',
    // },
    {
      id: 'total',
      header: 'Thành tiền (đ)',
      className: 'text-right font-semibold text-on-surface',
      cell: (row) => formatVND(row.total),
    },
    {
      id: 'actions',
      header: '',
      className: 'w-10 text-center',
      cell: (row) => (
        <button
          type="button"
          title="Xóa dòng"
          onClick={() => onRemoveItem(row.productId)}
          className="rounded p-1 text-on-surface-variant transition-colors hover:bg-error/10 hover:text-error"
        >
          <Trash2 className="h-4 w-4" />
        </button>
      ),
    },
  ];

  return (
    <section className="space-y-4 rounded-xl ">
      <div className="flex flex-col gap-3 bg-surface p-5 shadow-sm  rounded-xl">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between ">
          <div className="flex items-center gap-2">
            <span className="inline-block h-4 w-1 rounded-full bg-primary" />
            <h2 className="font-headline text-sm font-semibold tracking-wide text-on-surface uppercase">
              2. Danh mục sản phẩm & Vật tư
            </h2>
            <span className="rounded-full bg-surface-container px-2 py-0.5 text-xs font-medium text-on-surface-variant">
              {lines.length} mục
            </span>
          </div>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => setPickerOpen(true)}
              className="flex items-center gap-1.5 rounded-md bg-surface-container px-3 py-1.5 text-xs font-medium text-on-surface transition-all hover:bg-outline-variant"
            >
              <ListFilter className="h-4 w-4 text-primary" />
              <span>Chọn nhiều từ danh mục</span>
            </button>
            <button
              type="button"
              onClick={() => setScanMode((v) => !v)}
              className={cn(
                'flex items-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-semibold transition-all',
                scanMode
                  ? 'bg-primary text-on-primary'
                  : 'bg-primary-container/80 text-primary hover:bg-primary-container',
              )}
            >
              <QrCode className="h-4 w-4" />
              <span>{scanMode ? 'Đang quét mã...' : 'Bật máy quét mã'}</span>
            </button>
          </div>
        </div>
        <div className="flex flex-col gap-2 sm:flex-row">
          <div className="relative w-full">
            <Search className="absolute top-1/2 left-3.5 h-[18px] w-[18px] -translate-y-1/2 text-on-surface-variant" />
            <input
              value={query}
              onChange={(e) => onQueryChange(e.target.value)}
              onFocus={() => setFocused(true)}
              onBlur={() => setTimeout(() => setFocused(false), 150)}
              placeholder="Quét mã vạch hoặc nhập tìm thuốc nhuộm, dầu gội, oxy, tinh dầu dưỡng..."
              className="w-full rounded-lg bg-surface-container py-2.5 pr-24 pl-10 text-xs text-on-surface shadow-inner placeholder:text-on-surface-variant focus:ring-1 focus:ring-primary focus:outline-none"
            />
            <div className="absolute top-1/2 right-2 flex -translate-y-1/2 items-center gap-1">
              <span className="rounded bg-surface px-1.5 py-0.5 text-[10px] font-semibold text-on-surface-variant uppercase shadow-xs">
                F2: Tìm nhanh
              </span>
            </div>
            {focused && (
              <div className="absolute inset-x-0 top-full z-10 mt-1 overflow-hidden rounded-lg border border-outline bg-surface shadow-md">
                {isProductsLoading ? (
                  <div className="px-3 py-2 text-xs text-on-surface-variant">
                    Đang tìm sản phẩm…
                  </div>
                ) : suggestions.length === 0 ? (
                  <div className="px-3 py-2 text-xs text-on-surface-variant">
                    Không tìm thấy sản phẩm phù hợp.
                  </div>
                ) : (
                  suggestions.map((p) => (
                    <button
                      key={p.id}
                      type="button"
                      onMouseDown={(e) => e.preventDefault()}
                      onClick={() => {
                        onAddItem(p.id);
                        onQueryChange('');
                      }}
                      className="flex w-full items-center justify-between gap-3 px-3 py-2 text-left text-xs transition-colors hover:bg-surface-container"
                    >
                      <span className="min-w-0">
                        <span className="block truncate font-semibold text-on-surface">
                          {p.name}
                        </span>
                        <span className="font-mono text-[11px] text-on-surface-variant">
                          SKU: {p.sku} • {formatVND(p.costPrice)} đ
                        </span>
                      </span>
                      <PlusCircle className="h-4 w-4 shrink-0 text-primary" />
                    </button>
                  ))
                )}
              </div>
            )}
          </div>
        </div>
      </div>

      <DashTable
        columns={columns}
        data={rows}
        getRowId={(row) => row.productId}
        selectable={false}
        emptyText="Chưa có sản phẩm nào — tìm kiếm ở trên để thêm vào phiếu."
        rowClassName={(row) => (row.stripe ? 'bg-surface-container/20' : undefined)}
        footer={
          <div className="text-xs font-medium text-on-surface-variant">
            Tổng cộng:{' '}
            <span className="font-bold text-on-surface">
              {totalQty} đơn vị ({lines.length} mặt hàng)
            </span>
          </div>
        }
      />
      <ProductPickerDialog
        key={pickerOpen ? 'open' : 'closed'}
        open={pickerOpen}
        onOpenChange={setPickerOpen}
        categories={categories}
        addedIds={addedIds}
        onConfirm={onAddItems}
      />
    </section>
  );
}
