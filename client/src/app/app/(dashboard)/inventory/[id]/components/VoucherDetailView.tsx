'use client';

import { useParams } from 'next/navigation';
import useInventoryDocument from '@/hooks/use-inventory-document';
import { getApiErrorMessage } from '@/services/apiClient';
import { DashboardSkeleton } from '@/components/dashboard/DashboardSkeletons';
import VoucherTypeBanner from '../../(create)/components/VoucherTypeBanner';
import type { VoucherType } from '../../(create)/components/voucherData';
import VoucherDetailHeading from './VoucherDetailHeading';
import VoucherDetailInfoCard from './VoucherDetailInfoCard';
import VoucherDetailItemsCard from './VoucherDetailItemsCard';
import VoucherDetailFinanceCard from './VoucherDetailFinanceCard';

function toVoucherType(type: string): VoucherType {
  return type === 'export' || type === 'sale' ? 'outbound' : 'inbound';
}

/** Trang chi tiết phiếu kho — layout giống form tạo nhưng chỉ xem, không nhập. */
export default function VoucherDetailView() {
  const params = useParams<{ id: string }>();
  const id = params?.id;
  const { data, isLoading, isError, error } = useInventoryDocument(id);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <DashboardSkeleton />
      </div>
    );
  }

  if (isError || !data) {
    return (
      <div className="space-y-4">
        <div className="rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
          Không tải được chi tiết phiếu kho: {getApiErrorMessage(error)} — thử đăng nhập lại hoặc
          kiểm tra kết nối API.
        </div>
      </div>
    );
  }

  const { document, items } = data;

  return (
    <div className="space-y-4">
      <VoucherDetailHeading document={document} />
      <VoucherTypeBanner type={toVoucherType(document.type)} />

      <div className="grid grid-cols-1 items-start gap-6 lg:grid-cols-12">
        <div className="space-y-6 lg:col-span-8">
          <VoucherDetailInfoCard document={document} />
          <VoucherDetailItemsCard items={items} />
        </div>

        <div className="space-y-6 lg:col-span-4">
          <VoucherDetailFinanceCard document={document} items={items} />
        </div>
      </div>
    </div>
  );
}
