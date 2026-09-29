import {
  VAT_RATE,
  buildUnitOptions,
  priceForUnit,
  toSuggestion,
  type UnitOption,
  type VoucherItem,
} from './voucherData';
import type { Product } from '@/types/product';

export type VoucherLine = VoucherItem & {
  name: string;
  sku: string;
  spec: string;
  unit: string;
  netUnit?: string;
  netAmount?: string;
  costPrice: number;
  /** Tồn kho hiện tại (đơn vị nhỏ nhất) — dùng để chặn xuất vượt tồn ở UI. */
  stock: number;
  /** Đơn vị xuất đang chọn + hệ số quy đổi + đơn vị chuẩn. */
  selectedUnit: string;
  unitFactor: number;
  baseUnit: string;
  unitOptions: UnitOption[];
  /** Số lượng đã quy đổi ra đơn vị nhỏ nhất (= qty × factor), dùng để gửi API. */
  convertedQty: number;
  /** Đơn giá theo đơn vị xuất đang chọn (tự tính theo hệ số quy đổi). */
  price: number;
  /**
   * Đơn giá quy về 1 đơn vị nhỏ nhất (= price / factor, làm tròn 2 thập phân).
   * Dùng để gửi API vì quantity gửi đi cũng ở đơn vị nhỏ nhất.
   */
  pricePerNet: number;
  gross: number;
  discountAmount: number;
  total: number;
};

export function buildLines(items: VoucherItem[], products: Product[]): VoucherLine[] {
  const byId = new Map(products.map((p) => [p.id, p]));
  return items.map((item) => {
    const p = byId.get(item.productId);
    const s = p ? toSuggestion(p) : null;
    const { baseUnit, options } = p
      ? buildUnitOptions(p)
      : { baseUnit: '', options: [] as UnitOption[] };
    const selectedUnit = item.selectedUnit ?? p?.unit ?? '';
    const unitFactor = options.find((o) => o.value === selectedUnit)?.factor ?? 1;
    const price = p ? priceForUnit(p, selectedUnit) : 0;
    const gross = item.qty * price;
    const discountAmount = Math.round((gross * item.discount) / 100);
    const total = Math.round((gross - discountAmount) * (1 + VAT_RATE));
    return {
      ...item,
      name: s?.name ?? 'Sản phẩm không xác định',
      sku: s?.sku ?? '—',
      spec: s?.spec ?? '',
      unit: s?.unit ?? '',
      costPrice: s?.costPrice ?? 0,
      stock: p?.stock_quantity ?? 0,
      selectedUnit,
      unitFactor,
      baseUnit,
      unitOptions: options,
      convertedQty: item.qty * unitFactor,
      price,
      pricePerNet: unitFactor > 0 ? Math.round((price / unitFactor) * 100) / 100 : price,
      gross,
      discountAmount,
      total,
    };
  });
}

/** Dòng text quy đổi chuẩn: "1,500 ml". */
export function formatConvertedQty(line: Pick<VoucherLine, 'convertedQty' | 'baseUnit'>): string {
  return `${line.convertedQty.toLocaleString('vi-VN')} ${line.baseUnit}`.trim();
}

/** Dòng chi tiết quy đổi: "(1.5 x 1000 ml)". */
export function formatConversionDetail(
  line: Pick<VoucherLine, 'qty' | 'unitFactor' | 'baseUnit'>,
): string {
  return `(${line.qty} x ${line.unitFactor.toLocaleString('vi-VN')} ${line.baseUnit})`.trim();
}

/** Các dòng phiếu xuất đang vượt tồn kho (so theo đơn vị nhỏ nhất). */
export function findOverStockLines(lines: VoucherLine[]): VoucherLine[] {
  return lines.filter((l) => l.convertedQty > l.stock);
}

export type VoucherTotals = {
  totalQty: number;
  itemCount: number;
  gross: number;
  discount: number;
  vat: number;
  shippingFee: number;
  grand: number;
};

export function calcTotals(lines: VoucherLine[], shippingFee: number): VoucherTotals {
  const gross = lines.reduce((s, l) => s + l.gross, 0);
  const discount = lines.reduce((s, l) => s + l.discountAmount, 0);
  const vat = Math.round((gross - discount) * VAT_RATE);
  return {
    totalQty: lines.reduce((s, l) => s + l.qty, 0),
    itemCount: lines.length,
    gross,
    discount,
    vat,
    shippingFee,
    grand: gross - discount + vat + shippingFee,
  };
}

export function formatVND(n: number): string {
  return n.toLocaleString('vi-VN');
}

const DIGITS = ['không', 'một', 'hai', 'ba', 'bốn', 'năm', 'sáu', 'bảy', 'tám', 'chín'];

function readThreeDigits(n: number): string {
  const tram = Math.floor(n / 100);
  const chuc = Math.floor((n % 100) / 10);
  const donvi = n % 10;
  const parts: string[] = [];
  if (tram > 0) {
    parts.push(`${DIGITS[tram]} trăm`);
  }
  if (chuc === 0) {
    if (donvi > 0) parts.push(tram > 0 ? `linh ${DIGITS[donvi]}` : DIGITS[donvi]);
  } else if (chuc === 1) {
    parts.push('mười');
    if (donvi === 5) parts.push('lăm');
    else if (donvi > 0) parts.push(DIGITS[donvi]);
  } else {
    parts.push(`${DIGITS[chuc]} mươi`);
    if (donvi === 1) parts.push('mốt');
    else if (donvi === 5) parts.push('lăm');
    else if (donvi > 0) parts.push(DIGITS[donvi]);
  }
  return parts.join(' ');
}

/** Đọc số tiền VNĐ thành chữ (mức triệu/tỷ, đủ cho tổng phiếu kho). */
export function numberToVietnameseWords(n: number): string {
  if (!Number.isFinite(n) || n < 0) return '';
  if (n === 0) return 'Không đồng chẵn.';
  const ty = Math.floor(n / 1_000_000_000);
  const trieu = Math.floor((n % 1_000_000_000) / 1_000_000);
  const nghin = Math.floor((n % 1_000_000) / 1000);
  const tram = n % 1000;
  const parts: string[] = [];
  if (ty > 0) parts.push(`${readThreeDigits(ty)} tỷ`);
  if (trieu > 0) parts.push(`${readThreeDigits(trieu)} triệu`);
  if (nghin > 0) parts.push(`${readThreeDigits(nghin)} nghìn`);
  if (tram > 0) parts.push(readThreeDigits(tram));
  const text = parts.join(' ');
  return `${text.charAt(0).toUpperCase()}${text.slice(1)} đồng chẵn.`;
}
