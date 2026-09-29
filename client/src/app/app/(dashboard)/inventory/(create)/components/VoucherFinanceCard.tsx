'use client';

import { ReceiptText } from 'lucide-react';
import { DashInput } from '@/components/dashboard/DashInput';
import { VAT_RATE } from './voucherData';
import { formatVND, numberToVietnameseWords, type VoucherTotals } from './voucherUtils';

type VoucherFinanceCardProps = {
  totals: VoucherTotals;
  shippingFee: number;
  onShippingFeeChange: (fee: number) => void;
};

export default function VoucherFinanceCard({
  totals,
  shippingFee,
  onShippingFeeChange,
}: VoucherFinanceCardProps) {
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
            {totals.totalQty} sản phẩm ({totals.itemCount} loại)
          </span>
        </div>
        <div className="flex items-center justify-between text-on-surface-variant">
          <span>Tổng tiền hàng (gốc)</span>
          <span className="font-medium text-on-surface">{formatVND(totals.gross)} đ</span>
        </div>
        <div className="flex items-center justify-between text-on-surface-variant">
          <span className="flex items-center gap-1">
            <span>Chiết khấu NCC</span>
            <span className="text-[10px] font-medium text-success">(Khuyến mãi)</span>
          </span>
          <span className="font-medium text-error">-{formatVND(totals.discount)} đ</span>
        </div>
        <div className="flex items-center justify-between text-on-surface-variant">
          <span>Tiền thuế VAT ({Math.round(VAT_RATE * 100)}%)</span>
          <span className="font-medium text-on-surface">{formatVND(totals.vat)} đ</span>
        </div>
        <div className="flex items-center justify-between text-on-surface-variant">
          <span>Phí bốc xếp / Vận chuyển</span>
          <div className="w-24">
            <DashInput
              type="number"
              min={0}
              value={shippingFee}
              onChange={(e) => onShippingFeeChange(Math.max(0, Number(e.target.value)))}
              className="bg-surface-container px-2 py-1 text-right text-xs font-medium"
            />
          </div>
        </div>

        <div className="mt-3 space-y-1 rounded-lg bg-surface-container p-3.5">
          <div className="text-[11px] font-semibold tracking-wider text-on-surface-variant uppercase">
            Tổng thanh toán
          </div>
          <div className="font-headline text-2xl font-black tracking-tight text-primary">
            {formatVND(totals.grand)} <span className="text-base font-semibold">đ</span>
          </div>
          <div className="text-[11px] text-on-surface-variant italic">
            {numberToVietnameseWords(totals.grand)}
          </div>
        </div>
      </div>
    </section>
  );
}
