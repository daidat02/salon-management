'use client';

import PageHeading from '@/components/dashboard/PageHeading';
import { FileDown, Plus } from 'lucide-react';
import { useState } from 'react';
import { ServiceCreateDialog } from './ServiceCreateDialog';

const ServiceHeading = () => {
  const [open, setOpen] = useState(false);
  return (
    <>
      <PageHeading
        breadcrumbs={[
          { label: 'Trang chủ', href: '#' },
          { label: 'Dịch vụ', active: true },
        ]}
        title="Dịch vụ"
        badge="24 dịch vụ"
        subtitle="Quản lý danh sách dịch vụ làm đẹp, bảng giá và thời lượng phục vụ"
        actions={[
          { label: 'Xuất Excel', icon: FileDown, variant: 'outline' },
          { label: 'Thêm dịch vụ', icon: Plus, variant: 'primary', onClick: () => setOpen(true) },
        ]}
      />
      <ServiceCreateDialog open={open} onOpenChange={setOpen} />
    </>
  );
};

export default ServiceHeading;
