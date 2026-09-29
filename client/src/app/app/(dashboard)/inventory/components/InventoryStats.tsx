import { Package, TriangleAlert, Truck, Wallet } from 'lucide-react';
import StatCard from '@/components/dashboard/StatCard';

export default function InventoryStats() {
  return (
    <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <StatCard
        label="Tổng SKU"
        value={
          <>
            86 <span className="text-xs font-normal text-on-surface-variant">mặt hàng</span>
          </>
        }
        sub="58 sản phẩm bán lẻ • 28 vật tư"
        icon={Package}
      />
      <StatCard
        label="Sắp hết hàng"
        value={
          <>
            07{' '}
            <span className="ml-1 rounded bg-red-100 px-2 py-0.5 text-xs font-semibold text-error">
              Báo động
            </span>
          </>
        }
        sub={<span className="text-red-600/80">Cần tạo đơn bổ sung trong 48h</span>}
        icon={TriangleAlert}
        iconWrapClassName="bg-red-100 text-error"
        className="border-red-200 bg-red-50/20"
      />
      <StatCard
        label="Tổng giá trị tồn kho"
        value={
          <>
            124.500.000 <span className="text-xs font-normal text-on-surface-variant">đ</span>
          </>
        }
        sub={
          <span className="inline-flex items-center gap-1 font-medium text-success">
            ↗ +4.2% so với tháng trước
          </span>
        }
        icon={Wallet}
        iconWrapClassName="bg-emerald-50 text-success"
      />
      <StatCard
        label="Chi phí nhập tháng này"
        value={
          <>
            32.800.000 <span className="text-xs font-normal text-on-surface-variant">đ</span>
          </>
        }
        sub="6 đợt nhập hàng đã đối soát"
        icon={Truck}
        iconWrapClassName="bg-blue-50 text-primary"
      />
    </section>
  );
}
