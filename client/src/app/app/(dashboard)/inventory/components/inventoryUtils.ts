import type { Product } from '@/types/product';

export type InventoryStatus = 'low' | 'in_stock' | 'out_of_stock' | 'warning';

/**
 * Hệ số quy đổi từ đơn vị tồn kho (nhỏ) sang đơn vị hiển thị (lớn).
 * VD: stock tính bằng g, net_amount = 500 (g/chai) => 1000g hiển thị thành 2 chai.
 */
export function getConversionFactor(product: Product): number {
  return product.net_amount && product.net_amount > 0 ? product.net_amount : 1;
}

/** Tồn kho đã quy đổi về đơn vị hiển thị. */
export function getDisplayStock(product: Product): number {
  return product.stock_quantity / getConversionFactor(product);
}

/** Ngưỡng tối thiểu đã quy đổi về cùng đơn vị hiển thị. */
export function getDisplayMinStock(product: Product): number {
  return product.min_stock / getConversionFactor(product);
}

/** Format số lượng tồn: tối đa 2 thập phân, bỏ số 0 thừa (5 thay vì 5.00). */
export function formatStockQuantity(value: number): string {
  if (!Number.isFinite(value)) return '0';
  return value.toLocaleString('vi-VN', { maximumFractionDigits: 2 });
}

/**
 * Logic trạng thái tồn kho — giữ nguyên process như mock data cũ:
 * - 0 = hết hàng
 * - <= min = sắp hết
 * - còn lại: progress < 40% = cảnh báo, ngược lại = còn hàng an toàn
 */
export function getInventoryStatus(product: Product): InventoryStatus {
  const stockQuantity = getDisplayStock(product);
  if (stockQuantity <= 0) return 'out_of_stock';
  if (stockQuantity <= getDisplayMinStock(product)) return 'low';
  if (getStockProgress(product) < 40) return 'warning';
  return 'in_stock';
}

/**
 * Process tính progress bar giống mock data:
 * progress = stock / min * 100 (clamp 0-100).
 * VD mock: stock 3/min 15 -> 20%, 24/30 -> 80%, 2/20 -> 10%.
 * Nếu min = 0 (chưa set ngưỡng) thì coi như đầy 100% khi còn hàng.
 */
export function getStockProgress(product: Product): number {
  const stockQuantity = getDisplayStock(product);
  const minStock = getDisplayMinStock(product);

  if (minStock <= 0) return stockQuantity > 0 ? 100 : 0;

  const pct = (stockQuantity / minStock) * 100;
  return Math.max(0, Math.min(100, Math.round(pct)));
}
export function formatVND(value: number): string {
  return `${value.toLocaleString('vi-VN')} đ`;
}

/** Đơn vị hiển thị tồn kho: ưu tiên `unit`, fallback `net_unit`. */
export function getStockUnit(product: Product): string {
  return product.unit || product.net_unit || '';
}
