import type { Metadata } from 'next';
import VoucherDetailView from './components/VoucherDetailView';

export const metadata: Metadata = { title: 'Chi tiết phiếu kho | SALON ADMIN' };

export default function InventoryDetailPage() {
  return <VoucherDetailView />;
}
