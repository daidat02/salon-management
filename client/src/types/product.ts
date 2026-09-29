export interface Product {
  id: string;
  organization_id: string;
  category_id: string;
  sku: string;
  name: string;
  unit: string;
  net_unit: string | null;
  net_amount: number | null;
  product_type: 'retail' | 'material' | 'both';
  show_on_web: boolean;
  cost_price: number;
  sell_price: number;
  stock_quantity: number;
  min_stock: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
