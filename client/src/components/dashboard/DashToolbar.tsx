'use client';

import * as React from 'react';
import { RefreshCw, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { cn } from 'cn';

export type ToolbarFilterOption = string | { label: string; value: string };

export type ToolbarFilter = {
  id: string;
  placeholder: string;
  options: ToolbarFilterOption[];
  value?: string;
  onValueChange?: (value: string) => void;
  widthClass?: string;
};

type DashToolbarProps = {
  searchPlaceholder?: string;
  searchValue?: string;
  onSearchChange?: (value: string) => void;
  filters?: ToolbarFilter[];
  onRefresh?: () => void;
  onReset?: () => void;
  className?: string;
};

function normalizeOption(opt: ToolbarFilterOption): { label: string; value: string } {
  if (typeof opt === 'string') return { label: opt, value: opt };
  return opt;
}
export default function DashToolbar({
  searchPlaceholder = 'Tìm kiếm...',
  searchValue,
  onSearchChange,
  filters = [],
  onRefresh,
  onReset,
  className,
}: DashToolbarProps) {
  const [internalSearch, setInternalSearch] = React.useState(searchValue ?? '');

  React.useEffect(() => {
    if (searchValue !== undefined) setInternalSearch(searchValue);
  }, [searchValue]);

  const handleSearch = (v: string) => {
    if (searchValue === undefined) setInternalSearch(v);
    onSearchChange?.(v);
  };

  return (
    <section
      className={cn(
        'bg-surface flex flex-col gap-3 rounded-md border border-outline p-3.5 shadow-sm lg:flex-row lg:items-center lg:justify-between',
        className,
      )}
    >
      <div className="flex w-full flex-1 flex-col items-center gap-2.5 sm:flex-row lg:w-auto">
        {/* Search */}
        <div className="relative w-full sm:w-72">
          <span className="text-on-surface-variant pointer-events-none absolute inset-y-0 left-0 flex items-center pl-2.5">
            <Search className="h-3.5 w-3.5" />
          </span>
          <Input
            value={internalSearch}
            onChange={(e) => handleSearch(e.target.value)}
            placeholder={searchPlaceholder}
            className="bg-surface-container border-outline placeholder:text-on-surface-variant focus-visible:border-primary focus-visible:ring-primary/30 h-8 w-full rounded-md border py-2 pr-3 pl-8 text-xs placeholder:text-xs focus-visible:ring-1"
          />
        </div>

        {/* Filters - hiển thị label, value dùng để lưu */}
        {filters.map((filter) => {
          const normalized = filter.options.map(normalizeOption);
          const valueToLabel = new Map(normalized.map((o) => [o.value, o.label] as const));
          const labelToValue = new Map(normalized.map((o) => [o.label, o.value] as const));
          const displayValue = filter.value
            ? (valueToLabel.get(filter.value) ?? filter.value)
            : null;
          return (
            <div key={filter.id} className={cn('relative w-full', filter.widthClass ?? 'sm:w-44')}>
              <Select
                value={displayValue}
                onValueChange={(label) => {
                  const actualValue = label ? (labelToValue.get(label) ?? label) : '';
                  filter.onValueChange?.(actualValue);
                }}
              >
                <SelectTrigger
                  size="default"
                  className="bg-surface-container border-outline text-on-surface focus-visible:border-primary focus-visible:ring-primary/30 h-8 w-full rounded-md border px-3 py-2 text-xs focus-visible:ring-1 data-placeholder:text-on-surface-variant"
                >
                  <SelectValue placeholder={filter.placeholder} />
                </SelectTrigger>
                <SelectContent>
                  {normalized.map(({ label }) => (
                    <SelectItem key={label} value={label} className="text-xs">
                      {label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          );
        })}
      </div>

      {/* Actions */}
      <div className="flex shrink-0 items-center gap-2 self-end lg:self-auto">
        <Button
          variant="outline"
          size="icon-sm"
          onClick={onRefresh}
          title="Làm mới dữ liệu"
          className="border-outline text-on-surface-variant hover:text-primary hover:bg-surface-container rounded-md border"
        >
          <RefreshCw className="h-4.5 w-4.5" />
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={onReset}
          className="border-outline text-on-surface-variant hover:text-primary hover:bg-surface-container rounded-md border px-3 py-1.5 text-xs font-medium"
        >
          Đặt lại bộ lọc
        </Button>
      </div>
    </section>
  );
}
