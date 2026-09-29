import type { Metadata } from 'next';
import CashierHeading from './components/CashierHeading';
import OrdersTable from '../orders/components/OrdersTable';

export const metadata: Metadata = { title: 'Thu ngân | SALON ADMIN' };

export default function CashierPage() {
  return (
    <div className="space-y-4">
      <CashierHeading />
      <OrdersTable />
    </div>
  );
}
