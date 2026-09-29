/** Types + dữ liệu cho page Tạo phiếu nhập / xuất kho. */

import type { Product } from '@/types/product';

export type VoucherType = 'inbound' | 'outbound';

/** Map loại phiếu UI -> type của API POST /inventory/documents. */
export const VOUCHER_TYPE_MAP: Record<VoucherType, 'import' | 'export'> = {
  inbound: 'import',
  outbound: 'export',
};

export type VoucherItem = {
  productId: string;
  qty: number;
  /** Chiết khấu % của dòng */
  discount: number;
  /** Đơn giá ghi đè (mặc định lấy theo loại phiếu: nhập -> giá vốn, xuất -> giá bán) */
  unitPrice?: number;
  /** Đơn vị xuất đã chọn (mặc định = unit của sản phẩm) */
  selectedUnit?: string;
};

/** Một nấc đơn vị quy đổi của sản phẩm (factor = số đơn vị nhỏ nhất trong 1 đơn vị này). */
export type UnitOption = {
  value: string;
  label: string;
  factor: number;
};

/**
 * Dựng danh sách đơn vị xuất từ unit + net_unit của sản phẩm, theo mẫu quy đổi chuẩn.
 * VD: unit='chai', net_unit='ml', net_amount=1000 -> [ml ×1, chai ×1000].
 * Không có dữ liệu net -> chỉ 1 option duy nhất với factor = 1.
 */
export function buildUnitOptions(p: Product): { baseUnit: string; options: UnitOption[] } {
  const baseUnit = p.net_unit || p.unit;
  const hasNet = !!p.net_unit && p.net_unit !== p.unit && !!p.net_amount && p.net_amount > 0;

  if (!hasNet) {
    return { baseUnit, options: [{ value: p.unit, label: p.unit, factor: 1 }] };
  }

  const factor = p.net_amount as number;
  return {
    baseUnit,
    options: [
      { value: p.net_unit as string, label: `${p.net_unit} (Đơn vị chuẩn)`, factor: 1 },
      { value: p.unit, label: `${p.unit} (${factor.toLocaleString('vi-VN')} ${baseUnit})`, factor },
    ],
  };
}

export type VoucherSuggestion = {
  id: string;
  name: string;
  sku: string;
  spec: string;
  unit: string;
  costPrice: number;
  sellPrice: number;
};

/** Chuẩn hoá Product từ API thành gợi ý hiển thị trong ô tìm kiếm. */
export function toSuggestion(p: Product): VoucherSuggestion {
  const specParts = [p.net_amount ? `${p.net_amount}${p.net_unit ?? ''}` : null].filter(Boolean);
  return {
    id: p.id,
    name: p.name,
    sku: p.sku,
    spec: specParts.length > 0 ? specParts.join(' • ') : `Đơn vị: ${p.unit}`,
    unit: p.unit,
    costPrice: p.cost_price,
    sellPrice: p.sell_price,
  };
}

/** Hệ số của đơn vị lớn (unit) quy ra đơn vị nhỏ nhất; =1 nếu không có dữ liệu net. */
export function largeFactorOf(p: Product): number {
  return p.net_unit && p.net_unit !== p.unit && p.net_amount && p.net_amount > 0
    ? p.net_amount
    : 1;
}

/** Giá gốc (chưa quy đổi) của sản phẩm: luôn lấy giá vốn (cost_price). */
export function baseLargePrice(p: Product): number {
  return p.cost_price;
}

/**
 * Đơn giá theo đơn vị đang chọn: giá vốn (theo unit lớn) chia theo hệ số.
 * VD: cost 120.000đ/chai, net 1000ml -> chọn ml = 120đ, chọn chai = 120.000đ.
 */
export function priceForUnit(p: Product, unitValue: string): number {
  const { options } = buildUnitOptions(p);
  const factor = options.find((o) => o.value === unitValue)?.factor ?? 1;
  const large = largeFactorOf(p);
  if (large <= 0) return 0;
  return Math.round((baseLargePrice(p) * factor) / large);
}

/**
 * Tính lại đơn giá khi đổi đơn vị, chỉ dùng dữ liệu của dòng
 * (costPrice + unitOptions + unit lớn) nên không phụ thuộc list tìm kiếm.
 */
export function recalcLinePrice(
  costPrice: number,
  unitOptions: UnitOption[],
  largeUnit: string,
  nextUnit: string,
): number {
  const next = unitOptions.find((o) => o.value === nextUnit);
  const large = unitOptions.find((o) => o.value === largeUnit);
  const largeFactor = large?.factor ?? 1;
  if (!next || largeFactor <= 0) return Math.round(costPrice);
  return Math.round((costPrice * next.factor) / largeFactor);
}

/** Đặc trưng giao diện từng loại phiếu (dùng chung form, khóa theo page). */
export type VoucherMeta = {
  title: string;
  breadcrumb: string;
  badge: string;
  submitLabel: string;
  reasonLabel: string;
  bannerClassName: string;
  bannerIconClassName: string;
};

export const VOUCHER_META: Record<VoucherType, VoucherMeta> = {
  inbound: {
    title: 'Tạo Phiếu Nhập Kho',
    breadcrumb: 'Nhập kho',
    badge: 'Phiếu nhập',
    submitLabel: 'Hoàn tất & Nhập kho',
    reasonLabel: 'Lý do nhập hàng',
    bannerClassName: 'border-emerald-200 bg-emerald-50/60',
    bannerIconClassName: 'bg-emerald-100 text-emerald-700',
  },
  outbound: {
    title: 'Tạo Phiếu Xuất Kho',
    breadcrumb: 'Xuất kho',
    badge: 'Phiếu xuất',
    submitLabel: 'Hoàn tất & Xuất kho',
    reasonLabel: 'Lý do xuất hàng',
    bannerClassName: 'border-amber-200 bg-amber-50/60',
    bannerIconClassName: 'bg-amber-100 text-amber-700',
  },
};

export type PaymentStatus = 'paid' | 'partial' | 'unpaid';

export type VoucherFormState = {
  type: VoucherType;
  code: string;
  date: string;
  warehouse: string;
  reason: string;
  supplier: string;
  /** Người nhận hàng — chỉ dùng cho phiếu xuất. */
  recipient: string;
  keeper: string;
  note: string;
  items: VoucherItem[];
  shippingFee: number;
  paymentMethod: string;
  paymentStatus: PaymentStatus;
  vatInvoiceCode: string;
  autoAddStock: boolean;
  autoCashbook: boolean;
  trackExpiry: boolean;
};

export const VAT_RATE = 0.08;

export const warehouses = [
  'Kho Chi nhánh Trung tâm (Quận 1)',
  'Kho Salon Thảo Điền (Quận 2)',
  'Kho Tổng Vận Hành (Bình Thạnh)',
];

export const inboundReasons = [
  'Nhập hàng định kỳ từ Nhà Cung Cấp',
  'Nhập bổ sung đơn hàng khẩn cấp',
  'Nhập hoàn trả dư liệu trình salon',
  'Nhập kiểm kê chênh lệch tăng',
];

export const outboundReasons = [
  'Xuất cho Stylist làm dịch vụ',
  'Bán lẻ tại quầy POS',
  'Điều chuyển sang chi nhánh khác',
  'Xuất tiêu hủy / hết hạn',
];

export const suppliers = [
  "L'Oréal Professionnel VN (MST: 0309988771)",
  'Kérastase Paris Flagship VN',
  'Olaplex Vietnam Official',
  'Công ty TNHH Mỹ Phẩm Tóc Châu Âu',
];

export const keepers = ['Lê Tuấn (Thủ kho chính)', 'Trần Hoàng (Kho phụ)', 'Thu Hà (Thu ngân)'];

export const paymentMethods = [
  'Chuyển khoản ngân hàng (NCC)',
  'Tiền mặt tại quầy salon',
  'Ghi nợ công nợ (Hạn thanh toán 30 ngày)',
  'Ví điện tử công ty',
];

export function makeInitialFormState(type: VoucherType): VoucherFormState {
  return {
    type,
    code: '',
    date: new Date().toISOString().split('T')[0],
    warehouse: '',
    reason: type === 'inbound' ? inboundReasons[0] : outboundReasons[0],
    supplier: '',
    recipient: '',
    keeper: keepers[0],
    note: '',
    items: [],
    shippingFee: 0,
    paymentMethod: paymentMethods[0],
    paymentStatus: 'paid',
    vatInvoiceCode: '',
    autoAddStock: true,
    autoCashbook: true,
    trackExpiry: true,
  };
}

/** Giữ tương thích ngược cho form chưa truyền type (mặc định phiếu nhập). */
export const initialFormState: VoucherFormState = makeInitialFormState('inbound');

export function randomVoucherCode(type: VoucherType): string {
  const prefix = type === 'inbound' ? 'PNK' : 'PXK';
  const num = Math.floor(1000 + Math.random() * 9000);
  const year = new Date().getFullYear();
  return `${prefix}-${year}-${num}`;
}
