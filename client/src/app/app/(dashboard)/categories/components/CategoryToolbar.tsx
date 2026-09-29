'use client';

import { DashInput } from '@/components/dashboard/DashInput';
import { DashSelect } from '@/components/dashboard/DashSelect';
import { Button } from '@/components/ui/button';
import useCategory from '@/hooks/use-category';
import { cn } from 'cn';
import { RefreshCw, Search, ShoppingBag, Sparkles } from 'lucide-react';
import { useRouter, useSearchParams } from 'next/navigation';
import React, { useEffect } from 'react';

const CategoryToolbar = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const search = searchParams.get('search') || '';
  const status = searchParams.get('status') || '';
  const tab = searchParams.get('tab') || 'service';
  const [inputValue, setInputValue] = React.useState(search);
  const [lastSearch, setLastSearch] = React.useState(search);
  if (search !== lastSearch) {
    setLastSearch(search);
    setInputValue(search);
  }

  const { data } = useCategory({ page: 1, pageSize: 100 });
  const categories = data?.data ?? [];
  const serviceCount = categories.filter((c) => c.type === 'service').length;
  const productCount = categories.filter((c) => c.type === 'product').length;

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

  const handleStatusChange = (value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    if (!value) params.delete('status');
    else params.set('status', value);
    params.delete('page');
    router.push(`?${params.toString()}`, { scroll: false });
  };

  const handleTabChange = (value: 'service' | 'product') => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('tab', value);
    params.delete('page');
    router.push(`?${params.toString()}`, { scroll: false });
  };

  const handleReset = () => {
    router.push('/app/categories?tab=service', { scroll: false });
  };

  return (
    <section className="bg-surface rounded-xl border border-outline p-4 shadow-sm flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
      <div className="flex items-center p-1 bg-surface-container rounded-lg border border-outline/70 w-fit">
        <button
          onClick={() => handleTabChange('service')}
          className={cn(
            'flex items-center gap-2 px-3.5 py-1.5 rounded-md text-xs transition-all',
            tab === 'service'
              ? 'bg-primary text-white font-semibold shadow-sm'
              : 'text-on-surface-variant hover:text-on-surface font-medium',
          )}
        >
          <Sparkles className="h-3.5 w-3.5" />
          <span>Danh mục Dịch vụ ({serviceCount})</span>
        </button>
        <button
          onClick={() => handleTabChange('product')}
          className={cn(
            'flex items-center gap-2 px-3.5 py-1.5 rounded-md text-xs transition-all',
            tab === 'product'
              ? 'bg-primary text-white font-semibold shadow-sm'
              : 'text-on-surface-variant hover:text-on-surface font-medium',
          )}
        >
          <ShoppingBag className="h-3.5 w-3.5" />
          <span>Danh mục Sản phẩm ({productCount})</span>
        </button>
      </div>

      <div className="flex items-center gap-2.5 flex-1 md:justify-end">
        <div className="relative flex-1 md:w-64 md:flex-none md:max-w-xs">
          <span className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-on-surface-variant">
            <Search className="h-4 w-4" />
          </span>
          <DashInput
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder="Tìm theo tên danh mục..."
            className="pl-9 text-xs"
          />
        </div>
        <DashSelect
          value={status}
          onValueChange={handleStatusChange}
          placeholder="Tất cả trạng thái"
          triggerClassName="w-40"
          options={[
            { value: '', label: 'Tất cả trạng thái' },
            { value: 'active', label: 'Hoạt động' },
            { value: 'inactive', label: 'Tạm ẩn' },
          ]}
        />
        <Button
          variant="outline"
          size="icon"
          onClick={handleReset}
          title="Đặt lại bộ lọc"
          className="border-outline text-on-surface-variant hover:text-on-surface hover:bg-surface-container rounded-lg shrink-0"
        >
          <RefreshCw className="h-4 w-4" />
        </Button>
      </div>
    </section>
  );
};

export default CategoryToolbar;
