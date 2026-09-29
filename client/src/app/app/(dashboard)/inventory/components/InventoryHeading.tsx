'use client';

import type { MouseEvent } from 'react';
import { ArchiveRestore, FileDown, History, SquarePlus } from 'lucide-react';
import PageHeading from '@/components/dashboard/PageHeading';
import { INVENTORY_HISTORY_ANCHOR } from './InventoryHistoryTable';
import { useRouter } from 'next/navigation';

function scrollToHistory(e: MouseEvent) {
  e.preventDefault();
  document
    .getElementById(INVENTORY_HISTORY_ANCHOR)
    ?.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

export default function InventoryHeading() {
  const router = useRouter();
  return (
    <PageHeading
      breadcrumbs={[
        { label: 'Trang chủ', href: '#' },
        { label: 'Kho Hàng', active: true },
      ]}
      title="Kho Hàng"
      badge="86 mặt hàng SKU"
      description={
        <p className="mt-0.5 text-xs text-on-surface-variant">
          Quản lý danh mục vật tư tiêu hao, mỹ phẩm bán lẻ và tồn kho kỹ thuật.
        </p>
      }
      actions={[
        {
          label: 'Xem lịch sử',
          icon: History,
          variant: 'link',
          onClick: (e) => scrollToHistory(e),
        },
        { label: 'Xuất file Excel', icon: FileDown, variant: 'outline' },
        {
          label: 'Xuất kho',
          icon: ArchiveRestore,
          variant: 'outline',
          onClick: (e) => router.push('/app/inventory/outbound'),
        },
        {
          label: 'Nhập kho',
          icon: SquarePlus,
          variant: 'primary',
          onClick: (e) => router.push('/app/inventory/inbound'),
        },
      ]}
    />
  );
}
