'use client';

import { useState } from 'react';
import ServiceCatalog, { dotFor, type CatalogItem } from './ServiceCatalog';
import InvoicePanel, { type InvoiceItem } from './InvoicePanel';

export default function PosContent() {
  const [cart, setCart] = useState<InvoiceItem[]>([]);

  const handleAdd = (item: CatalogItem) => {
    setCart((prev) => {
      const existing = prev.find((p) => p.id === item.id);
      if (existing) {
        return prev.map((p) => (p.id === item.id ? { ...p, quantity: p.quantity + 1 } : p));
      }
      return [...prev, { ...item, quantity: 1, dotColor: dotFor(item.category) } as InvoiceItem];
    });
  };

  const handleQuantity = (id: string, delta: number) => {
    setCart((prev) =>
      prev
        .map((p) => (p.id === id ? { ...p, quantity: Math.max(1, p.quantity + delta) } : p))
        .filter((p) => p.quantity > 0),
    );
  };

  const handleRemove = (id: string) => setCart((prev) => prev.filter((p) => p.id !== id));
  const handleClear = () => setCart([]);

  const handleCreateOrder = () => {
    if (cart.length === 0) return;
    console.log('Creating order with items:', cart);
    alert(
      `Tạo đơn thành công ${cart.length} mục — tạm tính ${cart.reduce((s, i) => s + i.price * i.quantity, 0).toLocaleString('vi-VN')} đ — sang Đơn hàng để thanh toán`,
    );
    setCart([]);
  };

  return (
    <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:gap-4">
      <ServiceCatalog onAdd={handleAdd} />
      <InvoicePanel
        items={cart}
        onQuantity={handleQuantity}
        onRemove={handleRemove}
        onClear={handleClear}
        onCheckout={handleCreateOrder}
      />
    </div>
  );
}
