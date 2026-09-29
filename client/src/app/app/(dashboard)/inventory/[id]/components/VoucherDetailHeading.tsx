'use client';

import Link from 'next/link';
import { ArrowLeft, Download, Printer } from 'lucide-react';
import PageHeading from '@/components/dashboard/PageHeading';
import type { InventoryDocument } from '@/types/inventory';

type VoucherDetailHeadingProps = {
  document: InventoryDocument;
};

const TYPE_LABEL: Record<string, string> = {
  import: 'Phiếu nhập',
  export: 'Phiếu xuất',
  adjust: 'Phiếu điều chỉnh',
  sale: 'Phiếu bán hàng',
  return: 'Phiếu trả hàng',
};

export default function VoucherDetailHeading({ document }: VoucherDetailHeadingProps) {
  return (
    <div className="space-y-3">
      <Link
        href="/app/inventory"
        className="inline-flex items-center gap-1.5 rounded-md p-1.5 text-xs font-medium text-on-surface-variant transition-colors hover:bg-surface-container hover:text-on-surface"
      >
        <ArrowLeft className="h-4 w-4" />
        <span>Quay lại kho hàng</span>
      </Link>
      <PageHeading
        breadcrumbs={[
          { label: 'Trang chủ', href: '#' },
          { label: 'Quản lý kho', href: '#' },
          { label: 'Nhập xuất kho', href: '/app/inventory' },
          { label: document.document_code },
        ]}
        title={`Chi tiết ${TYPE_LABEL[document.type] ?? 'phiếu kho'} ${document.document_code}`}
        badge={TYPE_LABEL[document.type] ?? document.type}
        badgeVariant="soft"
        description={
          document.created_at && (
            <p className="mt-0.5 text-xs text-on-surface-variant">
              Tạo lúc{' '}
              {new Date(document.created_at).toLocaleString('vi-VN', {
                dateStyle: 'short',
                timeStyle: 'short',
              })}
              {document.created_by_name ? ` bởi ${document.created_by_name}` : ''}
            </p>
          )
        }
        actions={[
          { label: 'In phiếu kho', icon: Printer, variant: 'outline' },
          { label: 'Xuất Excel', icon: Download, variant: 'outline' },
        ]}
      />
    </div>
  );
}
