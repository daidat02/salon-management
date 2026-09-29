import type { Metadata } from 'next';
import InventoryHeading from './components/InventoryHeading';
import InventoryStats from './components/InventoryStats';
import InventoryAlert from './components/InventoryAlert';
import InventoryToolbar from './components/InventoryToolbar';
import InventoryTable from './components/InventoryTable';
import InventoryHistoryTable from './components/InventoryHistoryTable';

export const metadata: Metadata = { title: 'Kho hàng | SALON ADMIN' };

export default function KhoHangPage() {
  return (
    <div className="space-y-4">
      <InventoryHeading />
      <InventoryStats />
      <InventoryAlert />
      <InventoryToolbar />
      <InventoryTable />
      <InventoryHistoryTable />
    </div>
  );
}
