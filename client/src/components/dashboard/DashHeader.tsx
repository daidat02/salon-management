'use client';

import { useEffect, useRef, useState } from 'react';
import { Bell, CircleQuestionMark, Plus, Scissors, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

type DashHeaderProps = {
  branchName?: string;
  userName?: string;
  userRole?: string;
  avatarUrl?: string;
  hasNotification?: boolean;
  searchPlaceholder?: string;
  onCreateOrder?: () => void;
  onSearchChange?: (value: string) => void;
  onNotificationClick?: () => void;
  onHelpClick?: () => void;
  onProfileClick?: () => void;
};

const DEFAULT_AVATAR =
  'https://lh3.googleusercontent.com/aida-public/AB6AXuDWatsTUy3FswNW9YhcwNpZeuPjy_8D1tsLZnG-u6hb-Ymb0748_O91x04uvaxXTjSgHMSoTKmL1g-SzIHEzEEsrD1QjzwB2C-dxUSUxowi5-gxlKK5uOcslAQHTDha0tGU7ieQsW3AFR4s91R4TUoYXq-NBfcaen4ItsDytNT0TKBPA1KdNfja8EOcv2KNI2_ufCzxWk8zYOgAlr_BNGvLTv4t-fr-PYnUsrRa5Iq9ptA_awm94XR-FA';

/**
 * Header dashboard — dựng lại từ preview/admin_preview/index.html (dòng 102–153).
 * Nút chính dùng shadcn Button (variant default → --primary = brand #4F7CAC).
 */
export default function DashHeader({
  branchName = 'Hệ Thống Điều Hành Salon',
  userName = 'Minh Anh',
  userRole = 'Quản lý salon',
  avatarUrl = DEFAULT_AVATAR,
  hasNotification = true,
  searchPlaceholder = 'Tìm theo tên, SĐT, mã lịch hẹn... (⌘K)',
  onCreateOrder,
  onSearchChange,
  onNotificationClick,
  onHelpClick,
  onProfileClick,
}: DashHeaderProps) {
  const [query, setQuery] = useState('');
  const searchRef = useRef<HTMLInputElement>(null);

  // Bấm ⌘K / Ctrl+K để focus vào ô search (giống preview)
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault();
        searchRef.current?.focus();
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);

  return (
    <header className="z-30 flex h-16 shrink-0 items-center justify-between border-b border-border bg-surface px-6 select-none">
      {/* Brand Title & Search */}
      <div className="flex items-center gap-8">
        <div className="flex w-56 items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-brand text-lg font-bold text-white shadow-sm">
            <Scissors className="h-5 w-5" />
          </div>
          <div>
            <span className="block text-base leading-tight font-bold tracking-tight text-slate-900">
              SALON OS
            </span>
            <span className="text-[11px] leading-none font-medium text-slate-400">
              {branchName}
            </span>
          </div>
        </div>

        {/* Quick Global Search */}
        <div className="relative hidden w-80 md:block">
          <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400">
            <Search className="h-4 w-4" />
          </span>
          <Input
            ref={searchRef}
            value={query}
            onChange={(e) => {
              setQuery(e.target.value);
              onSearchChange?.(e.target.value);
            }}
            type="text"
            placeholder={searchPlaceholder}
            className="bg-slate-50 py-1.5 pr-12 pl-9 text-xs hover:bg-slate-100/80 focus-visible:border-brand focus-visible:bg-white focus-visible:ring-2 focus-visible:ring-brand/20"
          />
          <span className="absolute inset-y-0 right-0 flex items-center pr-2.5">
            <kbd className="rounded border border-slate-200 bg-white px-1.5 py-0.5 font-mono text-[10px] text-slate-400 shadow-2xs">
              ⌘K
            </kbd>
          </span>
        </div>
      </div>

      {/* Actions & User Profile */}
      <div className="flex items-center gap-4">
        {/* Nút chính — shadcn Button, màu brand qua --primary */}
        <Button
          size="sm"
          onClick={onCreateOrder}
          className="hidden text-xs font-semibold shadow-sm sm:inline-flex"
        >
          <Plus />
          <span>Tạo đơn</span>
        </Button>

        {/* Notification Bell */}
        <div className="relative">
          <Button
            variant="ghost"
            size="icon"
            onClick={onNotificationClick}
            aria-label="Thông báo"
            className="relative size-9 rounded-lg text-slate-500 hover:text-slate-800"
          >
            <Bell className="h-4 w-4" />
            {hasNotification && (
              <span className="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-red-500 ring-2 ring-white" />
            )}
          </Button>
        </div>

        {/* Help */}
        <Button
          variant="ghost"
          size="icon"
          onClick={onHelpClick}
          aria-label="Trợ giúp"
          className="size-9 rounded-lg text-slate-500 hover:text-slate-800"
        >
          <CircleQuestionMark className="h-4 w-4" />
        </Button>

        {/* User Avatar */}
        <div
          onClick={onProfileClick}
          className="group flex cursor-pointer items-center gap-3 border-l border-border pl-2"
        >
          <div className="h-9 w-9 overflow-hidden rounded-full border border-border ring-2 ring-transparent transition-all group-hover:ring-brand/30">
            <img alt={userName} src={avatarUrl} className="h-full w-full object-cover" />
          </div>
          <div className="hidden text-left lg:block">
            <div className="text-xs leading-tight font-semibold text-slate-800">{userName}</div>
            <div className="text-[11px] leading-tight text-slate-400">{userRole}</div>
          </div>
        </div>
      </div>
    </header>
  );
}
