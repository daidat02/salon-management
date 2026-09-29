'use client';

import { CloudUpload, FileText, Tag, Wallet, X } from 'lucide-react';
import { cn } from 'cn';
import { Label } from '@/components/ui/label';
import { DashInput } from '@/components/dashboard/DashInput';
import { DashSelect } from '@/components/dashboard/DashSelect';
import { paymentMethods, type PaymentStatus, type VoucherFormState } from './voucherData';

type VoucherPaymentCardProps = {
  form: VoucherFormState;
  onChange: (patch: Partial<VoucherFormState>) => void;
};

const statusOptions: { value: PaymentStatus; label: string }[] = [
  { value: 'paid', label: 'Đã trả hết' },
  { value: 'partial', label: 'Một phần' },
  { value: 'unpaid', label: 'Chưa trả' },
];

export default function VoucherPaymentCard({ form, onChange }: VoucherPaymentCardProps) {
  return (
    <section className="space-y-4 rounded-xl bg-surface p-5 shadow-sm">
      <h3 className="flex items-center gap-2 font-headline text-sm font-semibold tracking-wide text-on-surface uppercase">
        <Wallet className="h-[18px] w-[18px] text-primary" />
        <span>Thanh toán & Chứng từ</span>
      </h3>

      <div className="space-y-3.5">
        <div>
          <Label className="mb-1 block text-xs font-medium text-on-surface-variant">
            Phương thức thanh toán
          </Label>
          <DashSelect
            value={form.paymentMethod}
            onValueChange={(v) => onChange({ paymentMethod: v })}
            options={paymentMethods}
            triggerClassName="bg-surface-container font-medium"
          />
        </div>

        <div>
          <Label className="mb-1 block text-xs font-medium text-on-surface-variant">
            Trạng thái thanh toán
          </Label>
          <div className="grid grid-cols-3 gap-1.5 rounded-md bg-surface-container p-1 text-xs">
            {statusOptions.map((opt) => (
              <button
                key={opt.value}
                type="button"
                onClick={() => onChange({ paymentStatus: opt.value })}
                className={cn(
                  'rounded py-1 transition-colors',
                  form.paymentStatus === opt.value
                    ? 'bg-surface font-semibold text-primary shadow-xs'
                    : 'text-on-surface-variant hover:text-on-surface',
                )}
              >
                {opt.label}
              </button>
            ))}
          </div>
        </div>

        <div>
          <Label className="mb-1 block text-xs font-medium text-on-surface-variant">
            Mã hóa đơn VAT / Phiếu gốc
          </Label>
          <div className="relative">
            <DashInput
              value={form.vatInvoiceCode}
              onChange={(e) => onChange({ vatInvoiceCode: e.target.value })}
              className="bg-surface-container pl-8 text-xs"
              placeholder="Nhập mã hóa đơn VAT / phiếu gốc..."
            />
            <Tag className="pointer-events-none absolute top-1/2 left-2.5 h-[15px] w-[15px] -translate-y-1/2 text-on-surface-variant" />
          </div>
        </div>

        <div>
          <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
            Ảnh phiếu giao hàng / Hóa đơn đỏ
          </Label>
          <div className="cursor-pointer space-y-2 rounded-lg bg-surface-container p-3 text-center transition-colors hover:bg-surface-container/70">
            <div className="mx-auto flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-primary">
              <CloudUpload className="h-[18px] w-[18px]" />
            </div>
            <div className="text-xs font-medium text-on-surface">
              Kéo thả ảnh hoặc <span className="font-semibold text-primary">Tải lên tệp</span>
            </div>
            <div className="text-[10px] text-on-surface-variant">
              Hỗ trợ JPG, PNG, PDF (Tối đa 15MB)
            </div>
          </div>
          <div className="mt-2 flex items-center justify-between rounded-md bg-surface-container p-2 text-xs">
            <div className="flex min-w-0 items-center gap-2">
              <FileText className="h-[18px] w-[18px] shrink-0 text-primary" />
              <span className="truncate text-[11px] font-medium text-on-surface">
                Hoa_don_Loreal_24102024.pdf
              </span>
            </div>
            <button
              type="button"
              className="shrink-0 text-on-surface-variant transition-colors hover:text-error"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}
