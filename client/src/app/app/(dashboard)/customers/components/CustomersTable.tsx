'use client';

import * as React from 'react';
import DashTable from '@/components/dashboard/DashTable';
import { customerColumns } from './columns';
import type { Customer } from '@/types/customer';
import { useRouter, useSearchParams } from 'next/navigation';
import useCustomer from '@/hooks/use-customer';
import { getApiErrorMessage } from '@/services/apiClient';
import { CustomerDrawer } from './CustomerDrawer';
import { TableSkeleton } from '@/components/dashboard/DashboardSkeletons';
import { useState } from 'react';

type CustomersTableProps = {
  customers?: Customer[];
  total?: number;
  currentPage?: number;
  pageSize?: number;
};

export default function CustomersTable({
  customers: initialCustomers,
  total: initialTotal,
  currentPage: initialPage,
  pageSize: initialPageSize = 20,
}: CustomersTableProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const pageFromUrl = parseInt(searchParams.get('page') || '1', 10) || 1;

  const page = initialPage ?? pageFromUrl;
  const effectivePage = Number.isNaN(pageFromUrl) ? page : pageFromUrl;

  const { data, isLoading, isError, error } = useCustomer({
    page: effectivePage,
    pageSize: initialPageSize,
    search,
  });

  // Dùng data từ hook, fallback về initial props khi chưa có (SSR hydration)
  const customers = data?.data ?? initialCustomers ?? [];
  const total = data?.total ?? initialTotal ?? 0;
  const pageSize = data?.page_size ?? initialPageSize;

  const [drawerOpen, setDrawerOpen] = useState(false);
  const handleRowClick = () => {
    setDrawerOpen(true);
  };

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('page', newPage.toString());
    router.push(`?${params.toString()}`, { scroll: false });
  };

  if (isError) {
    return (
      <div className="rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
        Không tải được danh sách khách hàng: {getApiErrorMessage(error)} — thử đăng nhập lại hoặc
        kiểm tra kết nối API.
      </div>
    );
  }

  // Hiển thị skeleton khi lần đầu chưa có data (delay 1s demo)
  if (isLoading && customers.length === 0) {
    return <TableSkeleton rows={5} cols={7} />;
  }

  return (
    <>
      <DashTable
        columns={customerColumns}
        data={customers}
        getRowId={(r) => r.id}
        onRowClick={handleRowClick}
        pagination={{
          page: effectivePage,
          pageSize,
          total,
          itemLabel: 'khách hàng',
          onPageChange: handlePageChange,
        }}
      />
      <CustomerDrawer open={drawerOpen} onOpenChange={setDrawerOpen} />
    </>
  );
}
