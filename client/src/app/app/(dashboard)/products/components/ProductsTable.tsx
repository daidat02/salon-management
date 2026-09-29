'use client';

import * as React from 'react';
import DashTable from '@/components/dashboard/DashTable';
import { productColumns } from './columns';
import { useRouter, useSearchParams } from 'next/navigation';
import useProduct from '@/hooks/use-product';
import { TableSkeleton } from '@/components/dashboard/DashboardSkeletons';
import { getApiErrorMessage } from '@/services/apiClient';

export default function ProductsTable() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const status = searchParams.get('status') || '';
  const sort = searchParams.get('sort') || '';
  const brand = searchParams.get('brand') || '';
  const pageFromUrl = parseInt(searchParams.get('page') || '1', 10) || 1;

  const page = pageFromUrl || 1;
  const initialPageSize = 20;

  // brand hiện chưa có API riêng, map tạm vào search để lọc theo tên thương hiệu
  const effectiveSearch = brand ? (search ? `${search} ${brand}` : brand) : search;

  const { data, isLoading, isError, error } = useProduct({
    page,
    pageSize: initialPageSize,
    search: effectiveSearch,
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
        Không tải được danh sách sản phẩm: {getApiErrorMessage(error)} — thử đăng nhập lại hoặc kiểm
        tra kết nối API.
      </div>
    );
  }

  if (isLoading && products.length === 0) {
    return <TableSkeleton />;
  }

  return (
    <DashTable
      columns={productColumns}
      data={products}
      getRowId={(r) => r.id}
      pagination={{
        page,
        pageSize: initialPageSize,
        total: total,
        itemLabel: 'sản phẩm',
        onPageChange: handlePageChange,
      }}
    />
  );
}
