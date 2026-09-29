'use client';

import * as React from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import DashTable from '@/components/dashboard/DashTable';
import { TableSkeleton } from '@/components/dashboard/DashboardSkeletons';
import { getApiErrorMessage } from '@/services/apiClient';
import useOrders from '@/hooks/use-order';
import useCustomer from '@/hooks/use-customer';
import { buildOrderColumns, type CustomerLookup } from './columns';

const PAGE_SIZE = 10;

export default function OrdersTable() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const status = searchParams.get('status') || '';
  const pageFromUrl = parseInt(searchParams.get('page') || '1', 10) || 1;

  const page = pageFromUrl || 1;

  const { data, isLoading, isError, error } = useOrders({
    page,
    pageSize: PAGE_SIZE,
    status,
  });

  // Lookup tên + SĐT khách (API đơn chỉ trả customer_id).
  const { data: customerData } = useCustomer({ page: 1, pageSize: 100 });
  const customerMap = React.useMemo(() => {
    const map = new Map<string, { name: string; phone: string }>();
    for (const c of customerData?.data ?? []) {
      if (c) map.set(c.id, { name: c.full_name, phone: c.phone });
    }
    return map;
  }, [customerData]);
  const lookup: CustomerLookup = React.useCallback(
    (id) => (id ? customerMap.get(id) : undefined),
    [customerMap],
  );
  const columns = React.useMemo(() => buildOrderColumns(lookup), [lookup]);

  const serverRows = (data?.data || []).filter(
    (r): r is NonNullable<typeof r> => r !== null && r !== undefined,
  );
  const serverTotal = data?.total || 0;

  // Backend chưa hỗ trợ search đơn → lọc theo mã đơn trong trang hiện tại.
  const query = search.trim().toLowerCase();
  const displayed = query
    ? serverRows.filter((row) => row.code.toLowerCase().includes(query))
    : serverRows;

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('page', newPage.toString());
    router.push(`?${params.toString()}`, { scroll: false });
  };

  if (isError) {
    return (
      <div className="rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
        Không tải được danh sách đơn hàng: {getApiErrorMessage(error)} — thử đăng nhập lại hoặc
        kiểm tra kết nối API.
      </div>
    );
  }

  if (isLoading && serverRows.length === 0) {
    return <TableSkeleton />;
  }

  return (
    <DashTable
      columns={columns}
      data={displayed}
      getRowId={(r) => r.id}
      pagination={{
        page: query ? 1 : page,
        pageSize: PAGE_SIZE,
        total: query ? displayed.length : serverTotal,
        itemLabel: 'đơn hàng',
        onPageChange: handlePageChange,
      }}
    />
  );
}
