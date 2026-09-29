// Layout riêng của dashboard — sidebar/navbar chỉ bọc các route trong
// group (dashboard), KHÔNG bọc /app/login.
import DashHeader from '@/components/dashboard/DashHeader';
import DashSidebar from '@/components/dashboard/DashSidebar';
import { SidebarProvider } from '@/components/ui/sidebar';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="theme-app flex h-screen flex-col overflow-hidden bg-background font-sans text-body">
      {/* Top header cố định — không scroll (preview/admin_preview/index.html:102) */}
      <DashHeader />

      {/* Vùng dưới header */}
      <div className="flex min-h-0 flex-1">
        {/* Sidebar shadcn — active theo pathname (preview/admin_preview/index.html:156) */}
        <SidebarProvider className="contents">
          <DashSidebar />
        </SidebarProvider>

        {/* Vùng scroll duy nhất */}
        <main className="custom-scrollbar min-h-0 min-w-0 flex-1 overflow-x-hidden overflow-y-auto bg-background">
          <div className="mx-auto w-full space-y-6 p-6">{children}</div>
        </main>
      </div>
    </div>
  );
}
