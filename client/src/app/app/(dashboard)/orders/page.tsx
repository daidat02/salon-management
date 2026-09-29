import type { Metadata } from 'next';
import OrdersHeading from './components/OrdersHeading';
import OrdersToolbar from './components/OrdersToolbar';
import OrdersTable from './components/OrdersTable';

export const metadata: Metadata = { title: 'Đơn hàng | SALON ADMIN' };

export default function DonHangPage() {
  return (
    <div className="space-y-4">
      <OrdersHeading />
      <OrdersToolbar />
      <OrdersTable />
    </div>
  );
}
