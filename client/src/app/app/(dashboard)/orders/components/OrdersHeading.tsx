'use client';

import { useRouter } from 'next/navigation';
import PageHeading from '@/components/dashboard/PageHeading';
import { FileDown, Plus } from 'lucide-react';

export default function OrdersHeading() {
  const router = useRouter();

  return (
    <PageHeading
      breadcrumbs={[{ label: 'Trang chủ', href: '#' }, { label: 'Đơn hàng', active: true }]}
      title="Đơn hàng"
      subtitle="Theo dõi đơn hàng, thanh toán và lịch sử mua hàng tại salon"
      actions={[
        { label: 'Xuất file Excel', icon: FileDown, variant: 'outline' },
        {
          label: 'Tạo đơn',
          icon: Plus,
          variant: 'primary',
          onClick: () => router.push('/app/cashier/pos'),
        },
      ]}
    />
  );
}
