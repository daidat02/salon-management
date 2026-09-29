'use client';
import PageHeading from '@/components/dashboard/PageHeading';
import { FileDown, UserPlus } from 'lucide-react';
import { useState } from 'react';
import { CustomerCreateDialog } from './CustomerCreateDialog';

const CustomerHeading = () => {
  const [openCreateDialog, setOpenCreateDialog] = useState(false);

  return (
    <>
      <PageHeading
        breadcrumbs={[
          { label: 'Trang chủ', href: '#' },
          { label: 'Khách hàng', active: true },
        ]}
        title="Khách hàng"
        subtitle="Quản lý 1.240 khách hàng thân thiết và khách hàng mới"
        actions={[
          {
            label: 'Xuất Excel',
            icon: FileDown,
            variant: 'outline',
            onClick: () => console.log('Export Excel'),
          },
          {
            label: 'Thêm khách hàng',
            icon: UserPlus,
            variant: 'primary',
            onClick: () => setOpenCreateDialog(true),
          },
        ]}
      />

      <CustomerCreateDialog open={openCreateDialog} onOpenChange={setOpenCreateDialog} />
    </>
  );
};

export default CustomerHeading;
