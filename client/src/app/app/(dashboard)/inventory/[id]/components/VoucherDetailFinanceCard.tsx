'use client';

import { ReceiptText } from 'lucide-react';
import type { InventoryDocument, InventoryDocumentItem } from '@/types/inventory';
import {
  formatVND,
  numberToVietnameseWords,
} from '../../(create)/components/voucherUtils';
import { VAT_RATE } from '../../(create)/components/voucherData';

type VoucherDetailFinanceCardProps = {
  document: InventoryDocument;
  items: InventoryDocumentItem[];
};

export default function VoucherDetailFinanceCard({ document, items }: VoucherDetailFinanceCardProps) {
  const totalQty = items.reduce((s, it) => s + it.quantity, 0);
  // Backend chỉ lưu total_amount = Σ(qty × unit_price), không tách CK/VAT/phí ship
  // nên các dòng đó hiển thị 0đ cho khớp layout form tạo.
  const gross = items.reduce((s, it) => s + it.total_price, 0);

  return (
    <section className="space-y-4 rounded-xl bg-surface p-5 shadow-sm">
      <div className="flex items-center justify-between">
        <h3 className="flex items-center gap-2 font-headline text-sm font-semibold tracking-wide text-on-surface uppercase">
          <ReceiptText className="h-[18px] w-[18px] text-primary" />
          <span>Hạch toán tài chính</span>
        </h3>
        <span className="rounded bg-success/10 px-2 py-0.5 text-[11px] font-semibold text-success">
          Tự động tính
        </span>
      </div>

      <div className="space-y-3 pt-1 text-xs">
        <div className="flex items-center justify-between text-on-surface-variant">
          <span>Tổng lượng vật tư</span>
          <span className="font-semibold text-on-surface">
            {totalQty.toLocaleString('vi-VN')} sản phẩm ({items.length} loại)
          </span>
        </div>
        <div className="flex items-center justify-between text-on-surface-variant">
          <span>Tổng tiền hàng (gốc)</span>
          <span className="font-medium text-on-surface">{formatVND(gross)} đ</span>
        </div>
        <div className="flex items-center justify-between text-on-surface-variant">
          <span className="flex items-center gap-1">
            <span>Chiết khấu NCC</span>
            <span className="text-[10px] font-medium text-success">(Khuyến mãi)</span>
          </span>
          <span className="font-medium text-error">-{formatVND(0)} đ</span>
        </div>
        <div className="flex items-center justify-between text-on-surface-variant">
          <span>Tiền thuế VAT ({Math.round(VAT_RATE * 100)}%)</span>
          <span className="font-medium text-on-surface">{formatVND(0)} đ</span>
        </div>
        <div className="flex items-center justify-between text-on-surface-variant">
          <span>Phí bốc xếp / Vận chuyển</span>
          <span className="font-medium text-on-surface">{formatVND(0)} đ</span>
        </div>

        <div className="mt-3 space-y-1 rounded-lg bg-surface-container p-3.5">
          <div className="text-[11px] font-semibold tracking-wider text-on-surface-variant uppercase">
            Tổng thanh toán
          </div>
          <div className="font-headline text-2xl font-black tracking-tight text-primary">
            {formatVND(document.total_amount)} <span className="text-base font-semibold">đ</span>
          </div>
          <div className="text-[11px] text-on-surface-variant italic">
            {numberToVietnameseWords(document.total_amount)}
          </div>
        </div>
      </div>
    </section>
  );
}
