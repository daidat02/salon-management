'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  Archive,
  Calendar,
  ChartColumn,
  CirclePlus,
  LayoutDashboard,
  LogOut,
  Package,
  Palette,
  Receipt,
  Settings,
  ShoppingCart,
  Sparkles,
  Tags,
  UserCheck,
  Users,
  Wallet,
  type LucideIcon,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar';
import { cn } from 'cn';

type NavItem = {
  key: string;
  label: string;
  href: string;
  icon: LucideIcon;
  /** Pill nổi bật, vd. "12" ở Lịch hẹn */
  badge?: string;
  /** Chữ mờ bên phải, vd. "1.2k" ở Khách hàng */
  meta?: string;
};

// Nav chính — bám preview/admin_preview/index.html (dòng 164–223)
const MAIN_NAV: NavItem[] = [
  { key: 'overview', label: 'Tổng quan', href: '/app/overview', icon: LayoutDashboard },
  { key: 'cashier', label: 'Thu ngân', href: '/app/cashier', icon: ShoppingCart },
  {
    key: 'appointments',
    label: 'Lịch hẹn',
    href: '/app/appointments',
    icon: Calendar,
    badge: '12',
  },
  { key: 'orders', label: 'Đơn hàng', href: '/app/orders', icon: Receipt },
  { key: 'customers', label: 'Khách hàng', href: '/app/customers', icon: Users, meta: '1.2k' },
  { key: 'categories', label: 'Danh mục', href: '/app/categories', icon: Tags },
  { key: 'products', label: 'Sản phẩm', href: '/app/products', icon: Package },
  { key: 'services', label: 'Dịch vụ', href: '/app/services', icon: Sparkles },
  { key: 'inventory', label: 'Kho hàng', href: '/app/inventory', icon: Archive },
  { key: 'staff', label: 'Nhân viên', href: '/app/staff', icon: UserCheck },
  { key: 'cashflow', label: 'Thu & Chi', href: '/app/cashflow', icon: Wallet },
  { key: 'reports', label: 'Báo cáo', href: '/app/reports', icon: ChartColumn },
  { key: 'palette', label: 'Bảng màu', href: '/app/palette', icon: Palette },
];

type DashSidebarProps = {
  /** Active mặc định khi pathname chưa khớp item nào (preview đang highlight Lịch hẹn) */
  defaultActiveKey?: string;
  onCreateOrder?: () => void;
  onLogout?: () => void;
};

/**
 * Sidebar dashboard — dựng lại từ preview/admin_preview/index.html (dòng 156–236)
 * bằng shadcn Sidebar (SidebarMenuButton isActive + render thành Next Link).
 * Màu active/border/badge ăn theo --sidebar-accent / --brand (token app).
 */
export default function DashSidebar({
  defaultActiveKey = 'appointments',
  onCreateOrder,
  onLogout,
}: DashSidebarProps) {
  const pathname = usePathname();

  // Sửa lại logic nhận diện active theo prefix (startsWith) để support các trang con (create, [id],...)
  const activeKey =
    [...MAIN_NAV, { key: 'settings', href: '/app/settings' }]
      .sort((a, b) => b.href.length - a.href.length) // Ưu tiên match đường dẫn dài hơn trước
      .find((item) => pathname === item.href || pathname.startsWith(item.href + '/'))?.key ??
    defaultActiveKey;

  return (
    <Sidebar
      collapsible="none"
      className="h-full min-h-0 w-64 shrink-0 border-r border-border bg-white select-none max-md:hidden"
    >
      {/* Vùng scroll: CTA + nav */}
      <SidebarContent className="custom-scrollbar gap-6 p-4 [scrollbar:thin]">
        {/* Quick CTA */}
        <Button
          onClick={onCreateOrder}
          className="h-10 w-full rounded-xl text-xs font-medium shadow-sm active:scale-[0.99]"
        >
          <CirclePlus />
          <span>Đặt lịch mới</span>
        </Button>

        {/* Navigation Links */}
        <SidebarMenu className="gap-1">
          {MAIN_NAV.map((item) => {
            const Icon = item.icon;
            const isActive = item.key === activeKey;
            return (
              <SidebarMenuItem key={item.key}>
                <SidebarMenuButton
                  isActive={isActive}
                  render={<Link href={item.href} />}
                  className={cn(
                    'h-auto rounded-xl px-3.5 py-2.5 text-xs font-medium text-slate-600 hover:text-slate-900 [&_svg]:size-4 [&_svg]:text-slate-400',
                    isActive &&
                      'border-l-[3px] border-l-brand bg-brand-light font-semibold text-brand hover:text-brand [&_svg]:text-brand',
                  )}
                >
                  <Icon />
                  <span className="flex-1">{item.label}</span>
                  {item.badge && (
                    <span className="rounded-full bg-brand px-2 py-0.5 text-[11px] font-bold text-white">
                      {item.badge}
                    </span>
                  )}
                  {item.meta && (
                    <span className="text-[11px] font-medium text-slate-400">{item.meta}</span>
                  )}
                </SidebarMenuButton>
              </SidebarMenuItem>
            );
          })}
        </SidebarMenu>
      </SidebarContent>

      {/* Footer Menu */}
      <SidebarFooter className="shrink-0 gap-1 border-t border-border p-4">
        <SidebarMenu className="gap-1">
          <SidebarMenuItem>
            <SidebarMenuButton
              isActive={activeKey === 'settings'}
              render={<Link href="/app/settings" />}
              className="h-auto rounded-xl px-3.5 py-2 text-xs font-medium text-slate-600 hover:text-slate-900 [&_svg]:size-4 [&_svg]:text-slate-400"
            >
              <Settings />
              <span className="flex-1">Cài đặt salon</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              onClick={onLogout}
              className="h-auto rounded-xl px-3.5 py-2 text-xs font-medium text-red-600 hover:bg-red-50 hover:text-red-600 [&_svg]:size-4 [&_svg]:text-red-500"
            >
              <LogOut />
              <span className="flex-1">Đăng xuất</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
