import type {
  OrderDetail,
  OrderDetailItem,
  OrderDetailMaterial,
} from '@/types/order';
import type { Product } from '@/types/product';
import type { Service } from '@/types/service';
import type { Staff } from '@/types/staff';

export function formatVND(value: number): string {
  return `${value.toLocaleString('vi-VN')} đ`;
}

export function formatDateTime(value: string | null): string {
  if (!value) return '—';
  return new Date(value).toLocaleString('vi-VN', { dateStyle: 'short', timeStyle: 'short' });
}

export type EnrichedServiceLine = {
  id: string;
  name: string;
  desc: string;
  duration: string;
  staffName: string;
  staffRole: string;
  quantity: number;
  lineTotal: number;
};

export type EnrichedMaterialLine = {
  id: string;
  name: string;
  sku: string;
  qtyText: string;
  unitCost: number;
  lineCost: number;
};

export type EnrichedRetailLine = {
  id: string;
  name: string;
  desc: string;
  sku: string;
  qtyText: string;
  sellTotal: number;
  sellUnit: number;
  costTotal: number;
  costUnit: number;
  profit: number;
};

export type StaffCommission = {
  staffId: string;
  name: string;
  role: string;
  rateText: string;
  amount: number;
};

export function buildServiceLines(
  items: OrderDetailItem[],
  services: Map<string, Service>,
  staff: Map<string, Staff>,
): EnrichedServiceLine[] {
  return items
    .filter((it) => it.item_type === 'service')
    .map((it) => {
      const svc = it.service_id ? services.get(it.service_id) : undefined;
      const st = it.staff_id ? staff.get(it.staff_id) : undefined;
      return {
        id: it.id,
        name: svc?.name ?? 'Dịch vụ',
        desc: svc?.description ?? '',
        duration: svc ? `${svc.duration_minutes} phút` : '—',
        staffName: st ? `Stylist ${st.full_name}` : '—',
        staffRole: st?.position ?? '',
        quantity: it.quantity,
        lineTotal: it.line_total,
      };
    });
}

export function buildMaterialLines(
  materials: OrderDetailMaterial[],
  products: Map<string, Product>,
): EnrichedMaterialLine[] {
  return materials.map((m) => {
    const p = products.get(m.product_id);
    return {
      id: m.id,
      name: p?.name ?? 'Vật tư',
      sku: p?.sku ?? '—',
      qtyText: `${m.quantity} ${p?.unit ?? ''}`.trim(),
      unitCost: m.unit_cost_at_time,
      lineCost: m.quantity * m.unit_cost_at_time,
    };
  });
}

export function buildRetailLines(
  items: OrderDetailItem[],
  products: Map<string, Product>,
): EnrichedRetailLine[] {
  return items
    .filter((it) => it.item_type === 'product')
    .map((it) => {
      const p = it.product_id ? products.get(it.product_id) : undefined;
      const costUnit = p?.cost_price ?? 0;
      const costTotal = it.quantity * costUnit;
      return {
        id: it.id,
        name: p?.name ?? 'Sản phẩm',
        desc: p ? `Đơn vị: ${p.unit}` : '',
        sku: p?.sku ?? '—',
        qtyText: `x${it.quantity}${p ? ` ${p.unit}` : ''}`,
        sellTotal: it.line_total,
        sellUnit: it.quantity > 0 ? it.line_total / it.quantity : 0,
        costTotal,
        costUnit,
        profit: it.line_total - costTotal,
      };
    });
}

export type OrderProfit = {
  serviceRevenue: number;
  materialCost: number;
  serviceProfit: number;
  serviceMargin: number;
  retailRevenue: number;
  retailCost: number;
  retailProfit: number;
  retailMargin: number;
  totalProfit: number;
  totalRevenue: number;
  totalMargin: number;
};

export function calcOrderProfit(detail: OrderDetail, retailLines: EnrichedRetailLine[]): OrderProfit {
  const serviceRevenue = detail.items
    .filter((it) => it.item_type === 'service')
    .reduce((s, it) => s + it.line_total, 0);
  const materialCost = detail.materials.reduce((s, m) => s + m.quantity * m.unit_cost_at_time, 0);
  const serviceProfit = serviceRevenue - materialCost;
  const retailRevenue = retailLines.reduce((s, l) => s + l.sellTotal, 0);
  const retailCost = retailLines.reduce((s, l) => s + l.costTotal, 0);
  const retailProfit = retailRevenue - retailCost;
  const totalProfit = serviceProfit + retailProfit;
  const totalRevenue = serviceRevenue + retailRevenue;
  const pct = (part: number, whole: number) => (whole > 0 ? (part / whole) * 100 : 0);
  return {
    serviceRevenue,
    materialCost,
    serviceProfit,
    serviceMargin: pct(serviceProfit, serviceRevenue),
    retailRevenue,
    retailCost,
    retailProfit,
    retailMargin: pct(retailProfit, retailRevenue),
    totalProfit,
    totalRevenue,
    totalMargin: pct(totalProfit, totalRevenue),
  };
}

export function buildStaffCommissions(
  items: OrderDetailItem[],
  staff: Map<string, Staff>,
): StaffCommission[] {
  const revenueByStaff = new Map<string, number>();
  for (const it of items) {
    if (it.item_type !== 'service' || !it.staff_id) continue;
    revenueByStaff.set(it.staff_id, (revenueByStaff.get(it.staff_id) ?? 0) + it.line_total);
  }
  return Array.from(revenueByStaff.entries()).map(([staffId, revenue]) => {
    const st = staff.get(staffId);
    const rate = st?.commission_rate ?? 0;
    return {
      staffId,
      name: st?.full_name ?? 'Nhân viên',
      role: st?.position ?? '',
      rateText: `${rate}% doanh thu đảm nhận`,
      amount: Math.round((revenue * rate) / 100),
    };
  });
}

export function customerTier(totalSpent: number): { label: string; sub: string } {
  if (totalSpent >= 10000000) return { label: 'Hạng Vàng', sub: 'Khách hàng thân thiết' };
  if (totalSpent >= 5000000) return { label: 'Hạng Bạc', sub: 'Khách hàng thân thiết' };
  return { label: 'Hạng Thường', sub: 'Khách hàng mới' };
}
