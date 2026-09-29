import type { Metadata } from 'next';
import { Crown, FileDown, UserPlus, Users, Wallet } from 'lucide-react';
import PageHeading from '@/components/dashboard/PageHeading';
import StatCard from '@/components/dashboard/StatCard';
import CustomersTable from './components/CustomersTable';
import CustomersToolbar from './components/CustomersToolbar';
import CustomerHeading from './components/CustomerHeading';
export const metadata: Metadata = { title: 'Khách hàng | SALON ADMIN' };

export default function CustomersPage() {
  return (
    <div className="space-y-4">
      <CustomerHeading />

      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          label="Tổng khách hàng"
          value="1.240"
          sub={
            <span className="inline-flex items-center gap-1 font-medium text-success">
              <span className="h-1.5 w-1.5 rounded-full bg-success" /> +12% tháng này
            </span>
          }
          icon={Users}
          iconWrapClassName="bg-blue-50 text-primary"
        />
        <StatCard
          label="Khách hàng VIP"
          value="186"
          sub={
            <span className="inline-flex items-center gap-1 font-medium text-amber-700">
              <span className="h-1.5 w-1.5 rounded-full bg-amber-500" /> 15% tổng tệp khách
            </span>
          }
          icon={Crown}
          iconWrapClassName="bg-amber-50 text-amber-600"
        />
        <StatCard
          label="Khách mới tháng này"
          value="84"
          sub={<span className="font-medium text-success">↗ +8% so với tháng trước</span>}
          icon={UserPlus}
          iconWrapClassName="bg-emerald-50 text-success"
        />
        <StatCard
          label="Điểm tích lũy lưu hành"
          value="348.500"
          sub={
            <span>
              Quy đổi: <span className="font-semibold text-on-surface">~34.850.000 đ</span>
            </span>
          }
          icon={Wallet}
          iconWrapClassName="bg-indigo-50 text-indigo-600"
        />
      </section>

      <CustomersToolbar />
      <CustomersTable />
    </div>
  );
}
