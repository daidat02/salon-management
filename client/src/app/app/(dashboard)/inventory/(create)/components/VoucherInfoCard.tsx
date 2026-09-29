'use client';

import * as React from 'react';
import { format, isValid, parse } from 'date-fns';
import { vi } from 'date-fns/locale';
import { Plus, RefreshCw } from 'lucide-react';
import { Label } from '@/components/ui/label';
import { DashInput } from '@/components/dashboard/DashInput';
import { DashSelect } from '@/components/dashboard/DashSelect';
import { DashDatePicker } from '@/components/dashboard/DashDatePicker';
import {
  VOUCHER_META,
  inboundReasons,
  keepers,
  outboundReasons,
  randomVoucherCode,
  suppliers,
  warehouses,
  type VoucherFormState,
} from './voucherData';

type VoucherInfoCardProps = {
  form: VoucherFormState;
  onChange: (patch: Partial<VoucherFormState>) => void;
};

function SectionTitle({ children }: { children: string }) {
  return (
    <div className="flex items-center gap-2">
      <span className="inline-block h-4 w-1 rounded-full bg-primary" />
      <h2 className="font-headline text-sm font-semibold tracking-wide text-on-surface uppercase">
        {children}
      </h2>
    </div>
  );
}

/** Đồng hồ hệ thống chạy live (ngày/giờ/phút), tránh lệch hydration bằng cách mount ở client. */
function SystemClock() {
  const [now, setNow] = React.useState<Date | null>(null);

  React.useEffect(() => {
    const t = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(t);
  }, []);

  return (
    <span className="text-xs text-on-surface-variant">
      Thời gian hệ thống:{' '}
      <strong className="text-on-surface tabular-nums">
        {now ? format(now, 'dd/MM/yyyy HH:mm', { locale: vi }) : '…'}
      </strong>
    </span>
  );
}

export default function VoucherInfoCard({ form, onChange }: VoucherInfoCardProps) {
  const isInbound = form.type === 'inbound';
  const reasons = isInbound ? inboundReasons : outboundReasons;
  const meta = VOUCHER_META[form.type];

  const parsedDate = (() => {
    if (!form.date) return undefined;
    const d = parse(form.date, 'yyyy-MM-dd', new Date());
    return isValid(d) ? d : undefined;
  })();

  return (
    <section className="space-y-5 rounded-xl bg-surface p-5 shadow-sm">
      <div className="flex items-center justify-between gap-3">
        <SectionTitle>1. Thông tin phiếu kho</SectionTitle>
        <SystemClock />
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        <div>
          <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
            Mã chứng từ / Phiếu kho
          </Label>
          <div className="relative">
            <DashInput
              value={form.code}
              placeholder="Nhập mã phiếu kho..."
              onChange={(e) => onChange({ code: e.target.value })}
              className="bg-surface-container text-primary shadow-inner"
            />
            <button
              type="button"
              title="Tạo mã ngẫu nhiên mới"
              onClick={() => onChange({ code: randomVoucherCode(form.type) })}
              className="absolute top-1/2 right-2 -translate-y-1/2 text-on-surface-variant transition-colors hover:text-primary"
            >
              <RefreshCw className="h-4 w-4" />
            </button>
          </div>
        </div>
        <div className="flex-1">
          <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
            Người tạo phiếu
          </Label>
          <div className="flex h-9 items-center gap-2 truncate rounded-md bg-surface-container px-3 py-2 text-xs font-medium text-on-surface">
            <span className="h-2 w-2 shrink-0 rounded-full bg-success" />
            <span className="truncate">Minh Anh (Quản lý)</span>
          </div>
        </div>
        <div>
          <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
            Ngày giờ ghi nhận
          </Label>
          <DashDatePicker
            value={parsedDate}
            onValueChange={(d) => onChange({ date: d ? format(d, 'yyyy-MM-dd') : '' })}
            placeholder="Chọn ngày ghi nhận..."
            triggerClassName="bg-surface-container"
          />
        </div>

        {/* <div>
          <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
            Kho lưu trữ / Chi nhánh
          </Label>
          <DashInput
            value={form.warehouse}
            onChange={(e) => onChange({ warehouse: e.target.value })}
            placeholder="Nhập kho lưu trữ"
            className="bg-surface-container"
          />
        </div> */}
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div>
          <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
            {meta.reasonLabel}
          </Label>
          <DashSelect
            value={form.reason}
            onValueChange={(v) => onChange({ reason: v })}
            options={reasons}
            placeholder="Chọn lý do nhập / xuất hàng"
            triggerClassName="bg-surface-container"
          />
        </div>
        {isInbound ? (
          <div>
            <div className="mb-1.5 flex items-center justify-between">
              <Label className="text-xs font-medium text-on-surface-variant">
                Nhà Cung Cấp / Nhãn hàng
              </Label>
              <button
                type="button"
                className="flex items-center gap-0.5 text-[11px] font-semibold text-primary hover:underline"
              >
                <Plus className="h-[13px] w-[13px]" />
                Thêm NCC mới
              </button>
            </div>
            <DashSelect
              value={form.supplier}
              onValueChange={(v) => onChange({ supplier: v })}
              options={suppliers}
              placeholder="Chọn nhà cung cấp / nhãn hàng"
              triggerClassName="bg-surface-container"
            />
          </div>
        ) : (
          <div>
            <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
              Người nhận / Khách hàng
            </Label>
            <DashInput
              value={form.recipient}
              onChange={(e) => onChange({ recipient: e.target.value })}
              placeholder="Nhập tên người nhận, SĐT hoặc mã khách hàng..."
              className="bg-surface-container"
            />
          </div>
        )}
      </div>

      <div className="flex flex-col gap-4 md:flex-row md:items-center">
        <div className="flex-3">
          <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
            Ghi chú & Diễn giải chứng từ
          </Label>
          <DashInput
            value={form.note}
            onChange={(e) => onChange({ note: e.target.value })}
            placeholder="Nhập ghi chú cho phiếu kho..."
            className="bg-surface-container"
          />
        </div>

        {/* <div className="grid grid-cols-2 gap-3">
          <div>
            <Label className="mb-1.5 block text-xs font-medium text-on-surface-variant">
              Thủ kho giao / nhận
            </Label>
            <DashSelect
              value={form.keeper}
              onValueChange={(v) => onChange({ keeper: v })}
              options={keepers}
              triggerClassName="bg-surface-container"
            />
          </div>
        </div> */}
      </div>
    </section>
  );
}
