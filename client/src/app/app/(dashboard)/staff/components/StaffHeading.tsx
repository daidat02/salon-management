'use client';
import PageHeading from '@/components/dashboard/PageHeading';
import { FileDown, UserPlus } from 'lucide-react';
import { useState } from 'react';
import { StaffCreateDialog } from './StaffCreateDialog';

const StaffHeading = () => {
  const [openCreateDialog, setOpenCreateDialog] = useState(false);

  return (
    <>
      <PageHeading
        breadcrumbs={[
          { label: 'Trang chủ', href: '#' },
          { label: 'Nhân viên', active: true },
        ]}
        title="Nhân viên"
        badge="18 nhân viên"
        subtitle="Quản lý hồ sơ, lịch làm việc, chuyên môn kỹ thuật và hiệu suất phục vụ"
        actions={[
          {
            label: 'Xuất Excel',
            icon: FileDown,
            variant: 'outline',
            onClick: () => console.log('Export Excel'),
          },
          {
            label: 'Thêm nhân viên',
            icon: UserPlus,
            variant: 'primary',
            onClick: () => setOpenCreateDialog(true),
          },
        ]}
      />

      <StaffCreateDialog open={openCreateDialog} onOpenChange={setOpenCreateDialog} />
    </>
  );
};

export default StaffHeading;
