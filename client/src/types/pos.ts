export type PosItemType = 'product' | 'service';

export interface PosItem {
  id: string;
  organization_id: string;
  category_id: string | null;
  type: PosItemType;
  name: string;
  description: string | null;
  duration_minutes: number | null;
  unit: string | null;
  price: number;
  stock_quantity: number | null;
  is_active: boolean;
}

export interface PosItemsParams {
  search?: string;
  category?: string;
  /** Con trỏ phân trang cursor: type/name/id của dòng cuối trang trước */
  cursorType?: string;
  cursorName?: string;
  cursorId?: string;
  pageSize?: number;
}
