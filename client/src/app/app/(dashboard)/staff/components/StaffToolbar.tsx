'use client';

import * as React from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import DashToolbar from '@/components/dashboard/DashToolbar';
import { useEffect } from 'react';

export default function StaffToolbar() {
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
      searchPlaceholder="Tìm theo tên nhân viên, SĐT, mã NV..."
      searchValue={inputValue}
      onSearchChange={handleSearchChange}
      filters={[
        {
          id: 'status',
          placeholder: 'Tất cả trạng thái',
          value: searchParams.get('status') || undefined,
          widthClass: 'sm:w-44',
          options: [
            { label: 'Tất cả trạng thái', value: '' },
            { label: 'Đang hoạt động', value: 'active' },
            { label: 'Tạm ngưng', value: 'inactive' },
          ],
          onValueChange: (value) => handleFilterChange('status', value, 'Tất cả trạng thái'),
        },
        {
          id: 'sort',
          placeholder: 'Sắp xếp',
          value: searchParams.get('sort') || undefined,
          widthClass: 'sm:w-44',
          options: [
            { label: 'Mới nhất', value: 'created_at_desc' },
            { label: 'Tên A – Z', value: 'name_asc' },
            { label: 'Tên Z – A', value: 'name_desc' },
            { label: 'Vào làm sớm nhất', value: 'hire_date_asc' },
            { label: 'Vào làm mới nhất', value: 'hire_date_desc' },
          ],
          onValueChange: (value) => handleFilterChange('sort', value, 'Sắp xếp'),
        },
      ]}
      onRefresh={() => router.refresh()}
      onReset={() => router.push('?', { scroll: false })}
    />
  );
}
