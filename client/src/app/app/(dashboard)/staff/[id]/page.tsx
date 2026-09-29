import PageHeading from '@/components/dashboard/PageHeading';
import { Edit, Lock } from 'lucide-react';
import React from 'react';

const StaffDetailPage = () => {
  return (
    <div className="space-y-4">
      <PageHeading
        breadcrumbs={[
          { label: 'Trang chủ', href: '#' },
          { label: 'Nhân viên', active: false, href: '/app/staff' },
          { label: 'Nguyễn Văn A', active: true },
        ]}
        title="Nguyễn Văn A"
        subtitle="Chi tiết thông tin nhân viên, lịch sử làm việc và hiệu suất"
        actions={[
          { label: 'Khóa Tài khoản', icon: Lock, variant: 'outline' },
          { label: 'Chỉnh sửa hồ sơ', icon: Edit, variant: 'primary' },
        ]}
      />
    </div>
  );
};

export default StaffDetailPage;
