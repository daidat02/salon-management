export type OrderStatus =
  | 'draft'
  | 'pending_payment'
  | 'paid'
  | 'serving'
  | 'completed'
  | 'cancelled'
  | 'refunded';

export interface Order {
  id: string;
  organization_id: string;
  code: string;
  customer_id: string | null;
  appointment_id: string | null;
  subtotal_amount: number;
  discount_amount: number;
  total_amount: number;
  status: OrderStatus;
  created_at: string | null;
  updated_at: string | null;
}

export interface OrderDetailItem {
  id: string;
  item_type: 'service' | 'product';
  service_id: string | null;
  product_id: string | null;
  staff_id: string | null;
  quantity: number;
  unit_price: number;
  line_total: number;
}

export interface OrderDetailMaterial {
  id: string;
  order_item_id: string;
  product_id: string;
  quantity: number;
  unit_cost_at_time: number;
}

export interface OrderDetailPayment {
  id: string;
  order_id: string;
  method: string;
  provider: string | null;
  provider_txn_id: string | null;
  amount: number;
  status: string;
  paid_at: string | null;
  created_at: string | null;
}

export interface OrderDetail {
  order: Order;
  items: OrderDetailItem[];
  materials: OrderDetailMaterial[];
  payments: OrderDetailPayment[];
}
