import type { Metadata } from 'next';
import VoucherCreateForm from '../components/VoucherCreateForm';

export const metadata: Metadata = { title: 'Tạo phiếu nhập kho | SALON ADMIN' };

export default function CreateInboundPage() {
  return <VoucherCreateForm initialType="inbound" />;
}
