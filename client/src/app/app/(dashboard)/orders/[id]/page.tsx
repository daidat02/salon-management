import type { Metadata } from 'next';
import OrderDetailContent from './components/OrderDetailContent';

export const metadata: Metadata = { title: 'Chi tiết đơn hàng | SALON ADMIN' };

export default function OrderDetailPage() {
  return <OrderDetailContent />;
}
