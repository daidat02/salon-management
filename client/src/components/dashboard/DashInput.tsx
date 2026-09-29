'use client';

import * as React from 'react';
import { Input } from '@/components/ui/input';
import { cn } from 'cn';

export interface DashInputProps extends React.ComponentProps<typeof Input> {}

export function DashInput({ className, ...props }: DashInputProps) {
  return (
    <Input
      className={cn(
        'h-9 bg-surface border-outline text-sm placeholder:text-muted-foreground',
        'focus-visible:border-primary focus-visible:ring-1 focus-visible:ring-primary/30',
        className,
      )}
      {...props}
    />
  );
}
