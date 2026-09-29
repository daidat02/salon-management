'use client';

import * as React from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { cn } from 'cn';

type DashDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: string;
  children: React.ReactNode;
  footer?: React.ReactNode;
  maxWidth?: string;
};

export function DashDialog({
  open,
  onOpenChange,
  title,
  description,
  children,
  footer,
  maxWidth = 'max-w-[640px]',
}: DashDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        showCloseButton={false}
        className={cn(
          'bg-surface border-outline p-0 rounded-2xl shadow-xl border overflow-hidden',
          'w-full max-w-[calc(100%-2rem)] sm:max-w-160',
          maxWidth,
        )}
      >
        {/* Header - code.html:919-929 */}
        <div className="flex items-center justify-between px-6 pt-6 pb-4 border-b border-outline">
          <div>
            <DialogTitle className="text-lg font-bold font-headline text-slate-900 leading-none">
              {title}
            </DialogTitle>
            {description && (
              <DialogDescription className="text-xs text-on-surface-variant mt-1">
                {description}
              </DialogDescription>
            )}
          </div>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => onOpenChange(false)}
            className="text-on-surface-variant hover:text-slate-800 hover:bg-surface-container rounded-lg"
            aria-label="Đóng"
          >
            <span className="sr-only">Đóng</span>
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
            >
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </Button>
        </div>

        {/* Body - code.html:931-994 */}
        <div className="px-6 py-5 space-y-4 text-xs max-h-[70vh] overflow-y-auto custom-scroll">
          {children}
        </div>

        {/* Footer - code.html:996-1003 */}
        {footer && (
          <div className="flex items-center justify-end gap-3 px-6 py-4 border-t border-outline bg-surface">
            {footer}
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}

// Convenience footer buttons matching code.html
export function DashDialogFooter({
  onCancel,
  onConfirm,
  cancelText = 'Hủy',
  confirmText = 'Lưu khách hàng',
  confirmDisabled,
  isLoading,
}: {
  onCancel: () => void;
  onConfirm: () => void;
  cancelText?: string;
  confirmText?: string;
  confirmDisabled?: boolean;
  isLoading?: boolean;
}) {
  return (
    <>
      <Button
        variant="outline"
        onClick={onCancel}
        className="border-outline text-slate-600 hover:bg-surface-container px-4 py-2 text-xs font-semibold"
      >
        {cancelText}
      </Button>
      <Button
        onClick={onConfirm}
        disabled={confirmDisabled || isLoading}
        className="bg-primary hover:bg-[#436b95] text-on-primary px-5 py-2 text-xs font-semibold shadow-sm"
      >
        {isLoading ? 'Đang lưu...' : confirmText}
      </Button>
    </>
  );
}
