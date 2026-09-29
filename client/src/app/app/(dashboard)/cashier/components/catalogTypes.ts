export type CatalogItem = {
  id: string;
  name: string;
  desc: string;
  price: number;
  duration?: string;
  stock?: number;
  image: string;
  type: 'service' | 'product';
  category: string;
};

const DOT_MAP: Record<string, string> = {
  'Nhuộm & Phục hồi': 'bg-purple-500',
  'Cắt & Tạo kiểu': 'bg-blue-500',
  'Uốn & Duỗi': 'bg-amber-500',
  'Gội & Dưỡng sinh': 'bg-teal-500',
  'Sản phẩm chăm sóc tóc': 'bg-amber-500',
};

export function dotFor(categoryName: string) {
  return DOT_MAP[categoryName] ?? 'bg-slate-400';
}

export function tagVariantFor(item: { type: string; category: string }) {
  if (item.type !== 'service') return 'warning' as const;
  if (item.category.includes('Nhuộm')) return 'purple' as const;
  if (item.category.includes('Cắt')) return 'blue' as const;
  if (item.category.includes('Uốn')) return 'amber' as const;
  return 'teal' as const;
}
