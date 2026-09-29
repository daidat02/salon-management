'use client';

import * as React from 'react';
import { cn } from 'cn';
import DashTable from '@/components/dashboard/DashTable';
import { inventoryColumns } from './columns';
import { getInventoryStatus } from './inventoryUtils';
import { useRouter, useSearchParams } from 'next/navigation';
import useProduct from '@/hooks/use-product';
import { TableSkeleton } from '@/components/dashboard/DashboardSkeletons';
import { getApiErrorMessage } from '@/services/apiClient';

export default function InventoryTable() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const status = searchParams.get('status') || '';
  const sort = searchParams.get('sort') || '';
  const pageFromUrl = parseInt(searchParams.get('page') || '1', 10) || 1;

  const page = pageFromUrl || 1;
  const pageSize = 20;

  const { data, isLoading, isError, error } = useProduct({
    page,
    pageSize,
    search,
    status,
    sort,
  });

  const products = (data?.data || []).filter(
    (r): r is NonNullable<typeof r> => r !== null && r !== undefined,
  );
  const total = data?.total || 0;

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('page', newPage.toString());
    router.push(`?${params.toString()}`, { scroll: false });
  };

  if (isError) {
    return (
      <div className="rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
        Không tải được tồn kho: {getApiErrorMessage(error)} — thử đăng nhập lại hoặc kiểm tra kết
        nối API.
      </div>
    );
  }

  if (isLoading && products.length === 0) {
    return <TableSkeleton />;
  }

  return (
    <DashTable
      columns={inventoryColumns}
      data={products}
      getRowId={(r) => r.id}
      rowClassName={(row) => {
        const status = getInventoryStatus(row);
        return cn(
          status === 'low' && 'bg-red-50/10',
          status === 'out_of_stock' && 'bg-slate-50/50',
        );
      }}
      pagination={{
        page,
        pageSize,
        total,
        itemLabel: 'mặt hàng SKU',
        onPageChange: handlePageChange,
      }}
    />
  );
}
