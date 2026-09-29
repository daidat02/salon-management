'use client';

import * as React from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import DashToolbar from '@/components/dashboard/DashToolbar';

export default function OrdersToolbar() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const status = searchParams.get('status') || '';
  const [inputValue, setInputValue] = React.useState(search);

  React.useEffect(() => {
    setInputValue(search);
  }, [search]);

  React.useEffect(() => {
    if (inputValue === search) return;
    const t = setTimeout(() => {
      const params = new URLSearchParams(searchParams.toString());
      if (inputValue.trim()) params.set('search', inputValue.trim());
      else params.delete('search');
      params.delete('page');
      router.push(`?${params.toString()}`, { scroll: false });
    }, 400);
    return () => clearTimeout(t);
  }, [inputValue, search, router, searchParams]);

  const handleStatusChange = (value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    if (!value) params.delete('status');
    else params.set('status', value);
    params.delete('page');
    router.push(`?${params.toString()}`, { scroll: false });
  };

  return (
    <DashToolbar
      searchPlaceholder="Tìm mã đơn, tên khách, SĐT..."
      searchValue={inputValue}
      onSearchChange={setInputValue}
      filters={[
        {
          id: 'status',
          placeholder: 'Tất cả trạng thái',
          value: status || undefined,
          widthClass: 'sm:w-44',
          options: [
            { value: '', label: 'Tất cả trạng thái' },
            { value: 'pending_payment', label: 'Chưa thanh toán' },
            { value: 'paid', label: 'Đã thanh toán' },
            { value: 'serving', label: 'Đang phục vụ' },
            { value: 'completed', label: 'Hoàn tất' },
            { value: 'cancelled', label: 'Đã hủy' },
            { value: 'refunded', label: 'Đã hoàn tiền' },
          ],
          onValueChange: handleStatusChange,
        },
      ]}
      onRefresh={() => router.refresh()}
      onReset={() => router.push('/app/orders', { scroll: false })}
    />
  );
}
