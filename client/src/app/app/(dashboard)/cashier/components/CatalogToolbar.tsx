'use client';

import { Filter, ScanBarcode, Search, Store } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { cn } from 'cn';
import type { Category } from '@/types/category';

type CatalogToolbarProps = {
  search: string;
  onSearchChange: (value: string) => void;
  categories: Category[];
  activeCat: string;
  onCategoryChange: (id: string) => void;
};

export default function CatalogToolbar({
  search,
  onSearchChange,
  categories,
  activeCat,
  onCategoryChange,
}: CatalogToolbarProps) {
  return (
    <div className="bg-surface flex flex-col gap-3 rounded-md border border-outline p-4 shadow-sm">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Store className="text-primary h-6 w-6" />
          <h2 className="text-lg font-bold text-on-surface">Thu ngân & Bán hàng</h2>
        </div>
      </div>

      <div className="flex gap-2">
        <div className="relative flex-1">
          <Search className="text-on-surface-variant pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2" />
          <Input
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder="Tìm tên dịch vụ, liệu trình, sản phẩm... (F2 hoặc Cmd+K)"
            className="bg-surface-container-low border-outline h-8 rounded-md border py-2 pr-14 pl-9 text-sm text-on-surface placeholder:text-xs focus-visible:border-primary focus-visible:ring-1 focus-visible:ring-primary"
          />
        </div>
        <button className="bg-surface-container hover:bg-slate-200 border-outline flex items-center gap-1.5 rounded-md border px-3 py-2 text-xs font-medium text-on-surface transition-colors">
          <Filter className="h-4 w-4" /> Bộ lọc
        </button>
      </div>

      <div className="flex items-center gap-2 overflow-x-auto pt-1 pb-0.5 no-scrollbar">
        <button
          onClick={() => onCategoryChange('')}
          className={cn(
            'whitespace-nowrap rounded-full px-3.5 py-1.5 text-xs font-medium transition-all',
            activeCat === ''
              ? 'bg-primary text-on-primary font-semibold shadow-sm'
              : 'bg-surface-container-low hover:bg-slate-200 text-on-surface border border-outline',
          )}
        >
          Tất cả
        </button>
        {categories.map((cat) => {
          const isActive = cat.id === activeCat;
          return (
            <button
              key={cat.id}
              onClick={() => onCategoryChange(isActive ? '' : cat.id)}
              className={cn(
                'whitespace-nowrap rounded-full px-3.5 py-1.5 text-xs font-medium transition-all',
                isActive
                  ? 'bg-primary text-on-primary font-semibold shadow-sm'
                  : 'bg-surface-container-low hover:bg-slate-200 text-on-surface border border-outline',
              )}
            >
              {cat.name}
            </button>
          );
        })}
      </div>
    </div>
  );
}
