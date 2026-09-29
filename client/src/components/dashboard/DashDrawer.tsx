'use client';

import * as React from 'react';
import { X, MoreVertical } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Sheet, SheetContent } from '@/components/ui/sheet';
import { cn } from 'cn';

type DashDrawerProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  subtitle?: string;
  badge?: React.ReactNode;
  children: React.ReactNode;
  footer?: React.ReactNode;
  widthClass?: string;
};

export function DashDrawer({
  open,
  onOpenChange,
  title,
  subtitle,
  badge,
  children,
  footer,
  widthClass = 'max-w-[680px] w-[680px]',
}: DashDrawerProps) {
  // Map width sang class có !important để thắng Sheet mặc định w-3/4 + sm:max-w-sm
  // Thêm comment safelist để Tailwind JIT sinh ra các class này
  // w-[480px] max-w-[480px] !w-[480px] !max-w-[480px] sm:!max-w-[480px] data-[side=right]:!w-[480px] w-[680px] max-w-[680px] !w-[680px] !max-w-[680px] sm:!max-w-[680px] data-[side=right]:!w-[680px] w-[720px] max-w-[720px] !w-[720px] !max-w-[720px] sm:!max-w-[720px] data-[side=right]:!w-[720px] w-[800px] max-w-[800px] !w-[800px] !max-w-[800px] sm:!max-w-[800px] data-[side=right]:!w-[800px]
  const widthMap: Record<string, string> = {
    '480px':
      '!w-[480px] !max-w-[480px] sm:!max-w-[480px] data-[side=right]:!w-[480px] data-[side=right]:!max-w-[480px]',
    '680px':
      '!w-[680px] !max-w-[680px] sm:!max-w-[680px] data-[side=right]:!w-[680px] data-[side=right]:!max-w-[680px]',
    '720px':
      '!w-[720px] !max-w-[720px] sm:!max-w-[720px] data-[side=right]:!w-[720px] data-[side=right]:!max-w-[720px]',
    '800px':
      '!w-[800px] !max-w-[800px] sm:!max-w-[800px] data-[side=right]:!w-[800px] data-[side=right]:!max-w-[800px]',
    '900px':
      '!w-[900px] !max-w-[900px] sm:!max-w-[900px] data-[side=right]:!w-[900px] data-[side=right]:!max-w-[900px]',
  };
  const match = widthClass.match(/(\d+)px/);
  const widthPx = match ? `${match[1]}px` : '680px';
  const widthClasses = widthMap[widthPx] ?? `!w-[${widthPx}] !max-w-[${widthPx}]`;

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side="right"
        showCloseButton={false}
        className={cn(
          'p-0 gap-0 bg-surface border-l border-outline flex flex-col h-dvh max-h-dvh',
          widthClasses,
        )}
      >
        {/* Header - code.html:732-748 */}
        <div className="flex items-center justify-between border-b border-outline bg-surface-container-low p-5 shrink-0">
          <div>
            <div className="flex items-center gap-2">
              {badge}
              <h3 className="text-base font-bold text-slate-900">{title}</h3>
            </div>
            {subtitle && <p className="text-[11px] text-on-surface-variant mt-0.5">{subtitle}</p>}
          </div>
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon-sm"
              className="border border-outline bg-white text-on-surface-variant hover:text-slate-800 rounded-lg"
            >
              <MoreVertical className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => onOpenChange(false)}
              className="border border-outline bg-white text-on-surface-variant hover:text-slate-800 rounded-lg"
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        </div>

        {/* Body */}
        <div className="flex-1 overflow-y-auto p-5 space-y-6 custom-scroll min-h-0">{children}</div>

        {/* Footer */}
        {footer && (
          <div className="border-t border-outline bg-surface p-4 flex items-center justify-end gap-2.5 shrink-0">
            {footer}
          </div>
        )}
      </SheetContent>
    </Sheet>
  );
}
