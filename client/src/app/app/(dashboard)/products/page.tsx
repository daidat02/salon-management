import type { Metadata } from 'next';
import { FileDown, Package, Plus, TriangleAlert, Wallet } from 'lucide-react';
import PageHeading from '@/components/dashboard/PageHeading';
import StatCard from '@/components/dashboard/StatCard';
import DashToolbar from '@/components/dashboard/DashToolbar';
import ProductsTable from './components/ProductsTable';
import ProductToolbar from './components/ProductToolbar';
import ProductHeading from './components/ProductHeading';

export const metadata: Metadata = { title: 'Sản phẩm | SALON ADMIN' };

export default function ProductsPage() {
  return (
    <div className="space-y-4">
      <ProductHeading />

      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Tổng sản phẩm"
          value={
            <span className="flex items-baseline justify-between">
              <span>18</span>
              <span className="rounded-full bg-primary-container px-2 py-0.5 text-[11px] font-semibold text-primary">
                +3 tháng này
              </span>
            </span>
          }
          sub="Thuộc 4 thương hiệu chính"
          icon={Package}
          iconWrapClassName="bg-surface-container text-primary"
        />
        <StatCard
          label="Đang bán"
          value={
            <span className="flex items-baseline justify-between">
              <span>15</span>
              <span className="inline-flex items-center gap-1 rounded-full border border-emerald-100 bg-emerald-50 px-2 py-0.5 text-[11px] font-semibold text-success">
                <span className="h-1.5 w-1.5 rounded-full bg-success" /> 83.3% tỷ lệ
              </span>
            </span>
          }
          sub="Sẵn sàng bán tại quầy"
          icon={Package}
          iconWrapClassName="bg-emerald-50 text-success"
        />
        <StatCard
          label="Sắp hết hàng"
          value={
            <span className="flex items-baseline justify-between">
              <span>3</span>
              <span className="inline-flex items-center gap-1 rounded-full border border-amber-100 bg-amber-50 px-2 py-0.5 text-[11px] font-semibold text-warning">
                <span className="h-1.5 w-1.5 rounded-full bg-warning" /> Cần nhập hàng
              </span>
            </span>
          }
          sub="Dưới mức tồn an toàn"
          icon={TriangleAlert}
          iconWrapClassName="bg-amber-50 text-warning"
        />
        <StatCard
          label="Doanh thu bán lẻ tháng"
          value="48.600.000 đ"
          sub={<span className="font-semibold text-success">↗ +9.2% so với tháng trước</span>}
          icon={Wallet}
          iconWrapClassName="bg-primary-container text-primary"
        />
      </section>

      <ProductToolbar />
      <ProductsTable />
    </div>
  );
}
