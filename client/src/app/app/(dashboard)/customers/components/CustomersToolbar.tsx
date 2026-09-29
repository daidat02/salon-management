'use client';

import * as React from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import DashToolbar from '@/components/dashboard/DashToolbar';
import { useEffect } from 'react';

export default function CustomersToolbar() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const [inputValue, setInputValue] = React.useState(search);

  useEffect(() => {
    setInputValue(search);
  }, [search]);

  // debounce 400ms để không push router mỗi phím gõ
  useEffect(() => {
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

  const handleSearchChange = (v: string) => {
    setInputValue(v);
  };

  const handleFilterChange = (filterId: string, value: string, placeholder: string) => {
    const params = new URLSearchParams(searchParams.toString());

    // Nếu chọn lại placeholder "Tất cả..." thì xóa param thay vì set "Tất cả..."
    if (!value || value === placeholder || value.startsWith('Tất cả')) params.delete(filterId);
    else params.set(filterId, value);

    params.delete('page');
    router.push(`?${params.toString()}`, { scroll: false });
  };
  return (
    <DashToolbar
      searchPlaceholder="Tìm tên, SĐT, mã KH..."
      searchValue={inputValue}
      onSearchChange={handleSearchChange}
      filters={[
        {
          id: 'rank',
          placeholder: 'Tất cả hạng',
          value: searchParams.get('rank') || undefined,
          widthClass: 'sm:w-40',
          options: [
            { label: 'Tất cả hạng', value: '' },
            { label: 'Khách hàng VIP', value: 'VIP' },
            { label: 'Thành viên thân thiết', value: 'thân thiết' },
            { label: 'Khách mới', value: 'mới' },
          ],
          onValueChange: (value) => handleFilterChange('rank', value, 'Tất cả hạng'),
        },
        {
          id: 'status',
          placeholder: 'Tất cả trạng thái',
          value: searchParams.get('status') || undefined,
          widthClass: 'sm:w-40',
          options: [
            { label: 'Tất cả trạng thái', value: '' },
            { label: 'Đang hoạt động', value: 'active' },
            { label: 'Tạm ngưng', value: 'inactive' },
          ],
          onValueChange: (value) => handleFilterChange('status', value, 'Tất cả trạng thái'),
        },
        {
          id: 'time',
          placeholder: 'Tất cả thời gian',
          value: searchParams.get('time') || undefined,
          widthClass: 'sm:w-44',
          options: [
            { label: 'Tất cả thời gian', value: '' },
            { label: 'Hôm nay', value: 'today' },
            { label: 'Tháng này', value: 'this_month' },
            { label: 'Quý này', value: 'this_quarter' },
          ],
          onValueChange: (value) => handleFilterChange('time', value, 'Tất cả thời gian'),
        },
      ]}
      onRefresh={() => router.refresh()}
      onReset={() => router.push('?', { scroll: false })}
    />
  );
}
