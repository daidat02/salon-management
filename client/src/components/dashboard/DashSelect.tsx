'use client';

import * as React from 'react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { cn } from 'cn';

export type DashSelectOption = string | { label: string; value: string };

function normalizeOption(opt: DashSelectOption): { label: string; value: string } {
  if (typeof opt === 'string') return { label: opt, value: opt };
  return opt;
}

type DashSelectProps = {
  value?: string;
  defaultValue?: string;
  onValueChange?: (value: string) => void;
  placeholder?: string;
  options: DashSelectOption[];
  disabled?: boolean;
  className?: string;
  triggerClassName?: string;
};

export function DashSelect({
  value,
  defaultValue,
  onValueChange,
  placeholder = 'Chọn...',
  options,
  disabled,
  className,
  triggerClassName,
}: DashSelectProps) {
  const normalized = options.map(normalizeOption);

  return (
    <Select
      value={value}
      items={normalized}
      defaultValue={defaultValue}
      onValueChange={(v) => onValueChange?.(v ?? '')}
      disabled={disabled}
    >
      <SelectTrigger
        className={cn(
          // Đồng bộ với DashInput: h-9, cùng border, bg, text, focus
          'h-9 w-full bg-surface border-outline text-sm',
          'data-[size=default]:h-9! data-[size=sm]:h-9!',
          'focus-visible:border-primary focus-visible:ring-1 focus-visible:ring-primary/30',
          'data-placeholder:text-muted-foreground',
          triggerClassName,
        )}
      >
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {normalized.map(({ label, value }) => (
          <SelectItem key={value} value={value} className="text-sm">
            {label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
