'use client';

import PageHeading from '@/components/dashboard/PageHeading';
import { ArrowUpDown, Download, Plus } from 'lucide-react';
import { useState } from 'react';
import { CategoryCreateDialog } from './CategoryCreateDialog';

const CategoryHeading = () => {
  const [open, setOpen] = useState(false);
  return (
    <>
      <PageHeading
        breadcrumbs={[
          { label: 'Trang chủ', href: '#' },
          { label: 'Quản lý dịch vụ', href: '#' },
          { label: 'Danh mục', active: true },
        ]}
        title="Quản lý Danh mục"
        subtitle="Tổ chức và phân loại hệ thống dịch vụ làm đẹp, liệu trình và sản phẩm bán lẻ tại Salon"
        actions={[
          { label: 'Sắp xếp thứ tự', icon: ArrowUpDown, variant: 'outline' },
          { label: 'Xuất Excel', icon: Download, variant: 'outline' },
          { label: 'Thêm danh mục', icon: Plus, variant: 'primary', onClick: () => setOpen(true) },
        ]}
      />
      <CategoryCreateDialog open={open} onOpenChange={setOpen} />
    </>
  );
};

export default CategoryHeading;
