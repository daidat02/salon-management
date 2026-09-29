'use client';

import * as React from 'react';
import DashTable from '@/components/dashboard/DashTable';
import { staffColumns } from './columns';
import type { Staff } from '@/types/staff';
import { useRouter, useSearchParams } from 'next/navigation';
import useStaff from '@/hooks/use-staff';
import { getApiErrorMessage } from '@/services/apiClient';
import { TableSkeleton } from '@/components/dashboard/DashboardSkeletons';

type StaffTableProps = {
  staffs?: Staff[];
  total?: number;
  currentPage?: number;
  pageSize?: number;
};

export default function StaffTable({
  staffs: initialStaffs,
  total: initialTotal,
  currentPage: initialPage,
  pageSize: initialPageSize = 20,
}: StaffTableProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const status = searchParams.get('status') || '';
  const sort = searchParams.get('sort') || 'created_at_desc';
  const pageFromUrl = parseInt(searchParams.get('page') || '1', 10) || 1;

  const page = initialPage ?? pageFromUrl;
  const effectivePage = Number.isNaN(pageFromUrl) ? page : pageFromUrl;

  const { data, isLoading, isError, error } = useStaff({
    page: effectivePage,
    pageSize: initialPageSize,
    search,
    status,
    sort,
  });

  // Dùng data từ hook, fallback về initial props khi chưa có (SSR hydration)
  const staffs = data?.data ?? initialStaffs ?? [];
  const total = data?.total ?? initialTotal ?? 0;
  const pageSize = data?.page_size ?? initialPageSize;

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('page', newPage.toString());
    router.push(`?${params.toString()}`, { scroll: false });
  };

  const handleClickRow = (id: string) => {
    router.push(`/app/staff/${id}`);
  };
  if (isError) {
    return (
      <div className="rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
        Không tải được danh sách nhân viên: {getApiErrorMessage(error)} — thử đăng nhập lại hoặc
        kiểm tra kết nối API.
      </div>
    );
  }

  // Hiển thị skeleton khi lần đầu chưa có data
  if (isLoading && staffs.length === 0) {
    return <TableSkeleton rows={5} cols={6} />;
  }

  return (
    <DashTable
      columns={staffColumns}
      data={staffs}
      getRowId={(r) => r.id}
      onRowClick={(row) => handleClickRow(row.id)}
      pagination={{
        page: effectivePage,
        pageSize,
        total,
        itemLabel: 'nhân viên',
        onPageChange: handlePageChange,
      }}
    />
  );
}
