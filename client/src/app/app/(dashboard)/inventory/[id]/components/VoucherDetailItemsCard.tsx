'use client';

import * as React from 'react';
import DashTable, { type DashColumn } from '@/components/dashboard/DashTable';
import type { InventoryDocumentItem } from '@/types/inventory';
import { formatStockQuantity } from '../../components/inventoryUtils';

function formatVND(n: number): string {
  return n.toLocaleString('vi-VN');
}

export default function VoucherDetailItemsCard({ items }: { items: InventoryDocumentItem[] }) {
  const rows = React.useMemo(
    () => items.map((item, i) => ({ ...item, stt: i + 1, stripe: i % 2 === 1 })),
    [items],
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
      className: 'min-w-[220px]',
      cell: (row) => (
        <div className="flex items-center gap-3">
          <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-md border border-outline bg-surface-container text-sm font-bold text-primary">
            {(row.product_name || '?').charAt(0).toUpperCase()}
          </span>
          <div className="min-w-0">
            <div className="truncate font-semibold text-on-surface">{row.product_name}</div>
            <div className="flex items-center gap-2 text-[11px] text-on-surface-variant">
              <span>
                SKU: <strong className="font-mono text-on-surface">{row.sku}</strong>
              </span>
              {row.net_amount ? (
                <>
                  <span>•</span>
                  <span>
                    {row.net_amount}
                    {row.net_unit}
                  </span>
                </>
              ) : null}
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
          {row.display_unit ?? row.unit}
        </span>
      ),
    },
    {
      id: 'quantity',
      header: 'Số lượng',
      className: 'w-20 text-center font-semibold',
      cell: (row) => formatStockQuantity(row.display_quantity ?? row.quantity),
    },
    {
      id: 'unit_price',
      header: 'Đơn giá (đ)',
      className: 'text-right font-medium',
      cell: (row) => formatVND(row.display_unit_price ?? row.unit_price),
    },
    {
      id: 'total_price',
      header: 'Thành tiền (đ)',
      className: 'text-right font-semibold text-on-surface',
      cell: (row) => formatVND(row.total_price),
    },
  ];

  return (
    <section className="space-y-4 rounded-xl bg-surface p-5 shadow-sm">
      <div className="flex items-center gap-2">
        <span className="inline-block h-4 w-1 rounded-full bg-primary" />
        <h2 className="font-headline text-sm font-semibold tracking-wide text-on-surface uppercase">
          2. Danh mục sản phẩm & Vật tư
        </h2>
        <span className="rounded-full bg-surface-container px-2 py-0.5 text-xs font-medium text-on-surface-variant">
          {items.length} mục
        </span>
      </div>
      <DashTable
        columns={columns}
        data={rows}
        getRowId={(row) => row.id}
        selectable={false}
        emptyText="Phiếu không có dòng sản phẩm nào."
        rowClassName={(row) => (row.stripe ? 'bg-surface-container/20' : undefined)}
      />
    </section>
  );
}
