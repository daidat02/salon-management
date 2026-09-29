'use client';
import DashToolbar from '@/components/dashboard/DashToolbar';
import React, { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

const ServiceToolbar = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const status = searchParams.get('status') || '';
  const sort = searchParams.get('sort') || '';
  const category = searchParams.get('category') || '';
  const [inputValue, setInputValue] = React.useState(search);

  useEffect(() => {
    setInputValue(search);
  }, [search]);

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

  const handleFilterChange = (filterId: string, value: string, placeholder: string) => {
    const params = new URLSearchParams(searchParams.toString());
    if (!value || value === placeholder || value.startsWith('Tất cả') || value === '') params.delete(filterId);
    else params.set(filterId, value);
    params.delete('page');
    router.push(`?${params.toString()}`, { scroll: false });
  };

  return (
    <DashToolbar
      searchPlaceholder="Tìm theo tên dịch vụ, mã dịch vụ..."
      searchValue={inputValue}
      onSearchChange={(v) => setInputValue(v)}
      filters={[
        {
          id: 'category',
          placeholder: 'Tất cả danh mục',
          value: category || undefined,
          widthClass: 'sm:w-44',
          options: ['Tất cả danh mục', 'Cắt tóc', 'Nhuộm', 'Uốn & Duỗi', 'Gội đầu & Chăm sóc', 'Combo'],
          onValueChange: (value) => handleFilterChange('category', value, 'Tất cả danh mục'),
        },
        {
          id: 'status',
          placeholder: 'Tất cả trạng thái',
          value: status || undefined,
          widthClass: 'sm:w-40',
          options: [
            { value: '', label: 'Tất cả trạng thái' },
            { value: 'active', label: 'Đang hoạt động' },
            { value: 'inactive', label: 'Tạm ngưng / Ẩn' },
          ],
          onValueChange: (value) => handleFilterChange('status', value, 'Tất cả trạng thái'),
        },
        {
          id: 'sort',
          placeholder: 'Sắp xếp: Mới nhất',
          value: sort || undefined,
          widthClass: 'sm:w-44',
          options: [
            { value: 'created_at_desc', label: 'Sắp xếp: Mới nhất' },
            { value: 'price_asc', label: 'Giá tăng dần' },
            { value: 'price_desc', label: 'Giá giảm dần' },
            { value: 'duration_asc', label: 'Thời lượng tăng dần' },
          ],
          onValueChange: (value) => handleFilterChange('sort', value, 'Sắp xếp: Mới nhất'),
        },
      ]}
      onRefresh={() => router.refresh()}
      onReset={() => router.push('/app/services', { scroll: false })}
    />
  );
};

export default ServiceToolbar;
