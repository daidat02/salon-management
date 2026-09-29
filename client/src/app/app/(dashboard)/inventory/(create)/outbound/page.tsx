import type { Metadata } from 'next';
import VoucherCreateForm from '../components/VoucherCreateForm';

export const metadata: Metadata = { title: 'Tạo phiếu xuất kho | SALON ADMIN' };

export default function CreateOutboundPage() {
  return <VoucherCreateForm initialType="outbound" />;
}
