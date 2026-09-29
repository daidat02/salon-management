'use client';

import type { InventoryDocument } from '@/types/inventory';

function Field({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div>
      <div className="mb-1.5 block text-xs font-medium text-on-surface-variant">{label}</div>
      <div className="rounded-md bg-surface-container px-3 py-2 text-xs font-medium text-on-surface">
        {value || '—'}
      </div>
    </div>
  );
}

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

export default function VoucherDetailInfoCard({ document }: { document: InventoryDocument }) {
  const partner =
    document.supplier_name ??
    (document.order_code ? `Đơn ${document.order_code}` : null);

  return (
    <section className="space-y-5 rounded-xl bg-surface p-5 shadow-sm">
      <SectionTitle>1. Thông tin phiếu kho</SectionTitle>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        <Field label="Mã chứng từ / Phiếu kho" value={document.document_code} />
        <Field
          label="Ngày giờ ghi nhận"
          value={
            document.created_at
              ? new Date(document.created_at).toLocaleString('vi-VN', {
                  dateStyle: 'short',
                  timeStyle: 'short',
                })
              : null
          }
        />
        <Field label="Nhà cung cấp / Đơn hàng" value={partner} />
      </div>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <Field label="Người tạo phiếu" value={document.created_by_name} />
        <Field label="Ghi chú & Diễn giải chứng từ" value={document.note} />
      </div>
      {document.reason && (
        <div className="grid grid-cols-1 gap-4">
          <Field label="Lý do nhập / xuất hàng" value={document.reason} />
        </div>
      )}
    </section>
  );
}
