'use client';

import type { LucideIcon } from 'lucide-react';
import { CheckCircle2, FlaskConical, Scissors, ShoppingBag, Timer, TrendingUp } from 'lucide-react';
import DashTable, { type DashColumn } from '@/components/dashboard/DashTable';
import type {
  EnrichedMaterialLine,
  EnrichedRetailLine,
  EnrichedServiceLine,
} from './orderDetailUtils';
import { formatVND } from './orderDetailUtils';

function CardHeader({
  icon: Icon,
  iconClass,
  title,
  sub,
  badge,
}: {
  icon: LucideIcon;
  iconClass: string;
  title: string;
  sub: string;
  badge?: React.ReactNode;
}) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-3">
        <div className={`flex h-9 w-9 items-center justify-center rounded-lg ${iconClass}`}>
          <Icon className="h-5 w-5" />
        </div>
        <div>
          <h3 className="text-base font-semibold text-on-surface">{title}</h3>
          <p className="text-xs text-on-surface-variant">{sub}</p>
        </div>
      </div>
      {badge}
    </div>
  );
}

const serviceColumns: DashColumn<EnrichedServiceLine>[] = [
  {
    id: 'name',
    header: 'Tên dịch vụ',
    className: 'min-w-[200px]',
    cell: (row) => (
      <div>
        <div className="font-medium text-on-surface">{row.name}</div>
        {row.desc && <div className="mt-0.5 text-xs text-on-surface-variant">{row.desc}</div>}
      </div>
    ),
  },
  {
    id: 'duration',
    header: 'Thời lượng',
    cell: (row) => (
      <span className="inline-flex items-center gap-1.5 text-xs font-medium text-on-surface-variant">
        <Timer className="h-3.5 w-3.5" />
        {row.duration}
      </span>
    ),
  },
  {
    id: 'staff',
    header: 'Stylist thực hiện',
    cell: (row) => (
      <div>
        <div className="text-sm font-medium text-on-surface">{row.staffName}</div>
        {row.staffRole && <div className="text-xs font-medium text-primary">{row.staffRole}</div>}
      </div>
    ),
  },
  {
    id: 'quantity',
    header: 'Số lượng',
    className: 'text-center',
    cell: (row) => <span className="text-sm font-medium text-on-surface">x{row.quantity}</span>,
  },
  {
    id: 'lineTotal',
    header: 'Đơn giá & Thành tiền',
    className: 'text-right',
    cell: (row) => (
      <span className="text-sm font-semibold text-on-surface">{formatVND(row.lineTotal)}</span>
    ),
  },
];

const materialColumns: DashColumn<EnrichedMaterialLine>[] = [
  {
    id: 'name',
    header: 'Tên vật tư kỹ thuật',
    className: 'min-w-[200px]',
    cell: (row) => <span className="font-medium text-on-surface">{row.name}</span>,
  },
  {
    id: 'sku',
    header: 'Mã SKU',
    cell: (row) => <span className="font-mono text-xs text-on-surface-variant">{row.sku}</span>,
  },
  {
    id: 'qtyText',
    header: 'Số lượng tiêu hao',
    className: 'text-center',
    cell: (row) => <span className="text-xs font-medium text-on-surface">{row.qtyText}</span>,
  },
  {
    id: 'unitCost',
    header: 'Giá vốn định mức',
    className: 'text-right',
    cell: (row) => <span className="text-xs text-on-surface">{formatVND(row.unitCost)}</span>,
  },
  {
    id: 'lineCost',
    header: 'Thành tiền vốn',
    className: 'text-right',
    cell: (row) => (
      <span className="text-sm font-semibold text-on-surface">{formatVND(row.lineCost)}</span>
    ),
  },
];

const retailColumns: DashColumn<EnrichedRetailLine>[] = [
  {
    id: 'name',
    header: 'Tên sản phẩm',
    className: 'min-w-[200px]',
    cell: (row) => (
      <div>
        <div className="font-medium text-on-surface">{row.name}</div>
        {row.desc && <div className="text-xs text-on-surface-variant">{row.desc}</div>}
      </div>
    ),
  },
  {
    id: 'sku',
    header: 'Mã SKU',
    cell: (row) => <span className="font-mono text-xs text-on-surface-variant">{row.sku}</span>,
  },
  {
    id: 'qtyText',
    header: 'Số lượng',
    className: 'text-center',
    cell: (row) => <span className="font-medium text-on-surface">{row.qtyText}</span>,
  },
  {
    id: 'sellTotal',
    header: 'Giá bán lẻ',
    className: 'text-right',
    cell: (row) => (
      <div>
        <div className="font-medium text-on-surface">{formatVND(row.sellTotal)}</div>
        <div className="text-[11px] text-on-surface-variant">({formatVND(row.sellUnit)}/đv)</div>
      </div>
    ),
  },
  {
    id: 'costTotal',
    header: 'Giá vốn nhập',
    className: 'text-right',
    cell: (row) => (
      <div className="text-xs text-on-surface-variant">
        <div>{formatVND(row.costTotal)}</div>
        <div className="text-[11px]">({formatVND(row.costUnit)}/đv)</div>
      </div>
    ),
  },
  {
    id: 'profit',
    header: 'Lợi nhuận gộp',
    className: 'text-right',
    cell: (row) => <span className="font-semibold text-success">+{formatVND(row.profit)}</span>,
  },
];

export default function OrderDetailMain({
  services,
  materials,
  retail,
}: {
  services: EnrichedServiceLine[];
  materials: EnrichedMaterialLine[];
  retail: EnrichedRetailLine[];
}) {
  const serviceTotal = services.reduce((s, l) => s + l.lineTotal, 0);
  const materialTotal = materials.reduce((s, l) => s + l.lineCost, 0);
  const retailRevenue = retail.reduce((s, l) => s + l.sellTotal, 0);
  const retailCost = retail.reduce((s, l) => s + l.costTotal, 0);
  const retailProfit = retailRevenue - retailCost;

  return (
    <>
      <section className="space-y-3 ">
        <div className="bg-surface rounded-lg p-4 shadow-sm">
          <CardHeader
            icon={Scissors}
            iconClass="bg-primary-container text-primary"
            title="Dịch vụ trong đơn"
            sub="Quy trình thực hiện bởi đội ngũ Stylist chính"
            badge={
              <span className="rounded-md bg-primary-container px-2.5 py-1 text-xs font-semibold text-primary">
                {services.length} dịch vụ
              </span>
            }
          />
        </div>
        <DashTable
          columns={serviceColumns}
          data={services}
          getRowId={(row) => row.id}
          selectable={false}
          emptyText="Đơn không có dịch vụ."
          footer={
            <div className="flex w-full items-center justify-between">
              <span className="text-xs font-medium text-on-surface-variant">
                Tổng giá trị dịch vụ thực hiện
              </span>
              <div className="text-base font-bold tracking-tight text-primary">
                {formatVND(serviceTotal)}
              </div>
            </div>
          }
        />
      </section>

      <section className="space-y-3 ">
        <div className="bg-surface rounded-lg p-4 shadow-sm">
          <CardHeader
            icon={FlaskConical}
            iconClass="bg-orange-50 text-warning"
            title="Vật tư đã sử dụng (Tiêu hao nội bộ)"
            sub="Khấu trừ trực tiếp từ kho hóa chất salon"
            badge={
              <span className="inline-flex items-center gap-1 rounded-md border border-emerald-200 bg-emerald-50 px-2.5 py-1 text-xs font-semibold text-success">
                <CheckCircle2 className="h-3.5 w-3.5" />
                Tự động xuất kho
              </span>
            }
          />
        </div>
        <DashTable
          columns={materialColumns}
          data={materials}
          getRowId={(row) => row.id}
          selectable={false}
          emptyText="Đơn không tiêu hao vật tư."
          footer={
            <div className="flex w-full items-center justify-between">
              <span className="text-xs font-medium text-on-surface-variant">
                Tổng chi phí vật tư tiêu hao:
              </span>
              <div className="text-sm font-bold text-on-surface">{formatVND(materialTotal)}</div>
            </div>
          }
        />
      </section>

      <section className="space-y-3 ">
        <div className="bg-surface rounded-lg p-4 shadow-sm">
          <CardHeader
            icon={ShoppingBag}
            iconClass="bg-blue-50 text-primary"
            title="Sản phẩm bán lẻ mang về"
            sub="Chăm sóc tóc tại nhà chính hãng"
            badge={
              <span className="rounded-md bg-primary-container px-2.5 py-1 text-xs font-semibold text-primary">
                {retail.length} sản phẩm
              </span>
            }
          />
        </div>

        <DashTable
          columns={retailColumns}
          data={retail}
          getRowId={(row) => row.id}
          selectable={false}
          emptyText="Đơn không có sản phẩm bán lẻ."
          footer={
            <div className="flex w-full flex-wrap items-center justify-between gap-3 text-xs">
              <div className="text-on-surface-variant">
                Doanh thu bán lẻ:{' '}
                <span className="font-semibold text-on-surface">{formatVND(retailRevenue)}</span>
              </div>
              <div className="text-on-surface-variant">
                Tổng giá vốn:{' '}
                <span className="font-semibold text-on-surface">{formatVND(retailCost)}</span>
              </div>
              <div className="flex items-center gap-1 font-semibold text-success">
                <TrendingUp className="h-3.5 w-3.5" />
                Lãi gộp sản phẩm: {formatVND(retailProfit)}
              </div>
            </div>
          }
        />
      </section>
    </>
  );
}
