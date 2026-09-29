import {
  Flower2,
  Package,
  Palette,
  Scissors,
  ShoppingBag,
  Sparkles,
  WandSparkles,
  type LucideIcon,
} from 'lucide-react';
import type { CategoryType } from '@/types/category';

export type CategoryTab = 'service' | 'product';

export function getCategoryIcon(name: string, type: CategoryType): LucideIcon {
  const n = name.toLowerCase();
  if (n.includes('cắt')) return Scissors;
  if (n.includes('nhuộm')) return Palette;
  if (n.includes('uốn')) return WandSparkles;
  if (n.includes('gội') || n.includes('massage') || n.includes('spa')) return Sparkles;
  if (n.includes('da ') || n.includes('chăm sóc da') || n.includes('facial')) return Flower2;
  if (n.includes('sản phẩm')) return ShoppingBag;
  return type === 'product' ? Package : Scissors;
}

export function getCategoryTypeLabel(type: CategoryType) {
  return type === 'service' ? 'Dịch vụ' : 'Sản phẩm';
}
