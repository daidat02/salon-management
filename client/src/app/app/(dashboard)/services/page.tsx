import type { Metadata } from 'next';
import { BadgeCheck, PauseCircle, Scissors, Wallet } from 'lucide-react';
import StatCard from '@/components/dashboard/StatCard';
import ServicesTable from './components/ServicesTable';
import ServiceHeading from './components/ServiceHeading';
import ServiceToolbar from './components/ServiceToolbar';

export const metadata: Metadata = { title: 'Dịch vụ | SALON ADMIN' };

export default function ServicesPage() {
  return (
    <div className="space-y-4">
      <ServiceHeading />

      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Tổng dịch vụ"
          value={
            <span className="flex items-baseline justify-between">
              <span>24</span>
              <span className="rounded-full bg-primary-container px-2 py-0.5 text-[11px] font-semibold text-primary">
                +2 tháng này
              </span>
            </span>
          }
          sub="Phục vụ trên 5 danh mục chính"
          icon={Scissors}
          iconWrapClassName="bg-surface-container text-primary"
        />
        <StatCard
          label="Đang hoạt động"
          value={
            <span className="flex items-baseline justify-between">
              <span>21</span>
              <span className="inline-flex items-center gap-1 rounded-full border border-emerald-100 bg-emerald-50 px-2 py-0.5 text-[11px] font-semibold text-success">
                <span className="h-1.5 w-1.5 rounded-full bg-success" /> 87.5% tỷ lệ
              </span>
            </span>
          }
          sub="Sẵn sàng nhận khách đặt lịch"
          icon={BadgeCheck}
          iconWrapClassName="bg-emerald-50 text-success"
        />
        <StatCard
          label="Tạm ngưng / Ẩn"
          value={
            <span className="flex items-baseline justify-between">
              <span>3</span>
              <span className="inline-flex items-center gap-1 rounded-full border border-amber-100 bg-amber-50 px-2 py-0.5 text-[11px] font-semibold text-warning">
                <span className="h-1.5 w-1.5 rounded-full bg-warning" /> Cần nhập liệu
              </span>
            </span>
          }
          sub="Hết thuốc đặc trị hoặc tái cơ cấu"
          icon={PauseCircle}
          iconWrapClassName="bg-amber-50 text-warning"
        />
        <StatCard
          label="Doanh thu dịch vụ tháng"
          value="186.400.000 đ"
          sub={<span className="font-semibold text-success">↗ +15.8% so với tháng trước</span>}
          icon={Wallet}
          iconWrapClassName="bg-primary-container text-primary"
        />
      </section>

      <ServiceToolbar />
      <ServicesTable />
    </div>
  );
}
