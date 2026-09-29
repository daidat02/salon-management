export type CategoryType = 'service' | 'product';

export interface Category {
  id: string;
  organization_id: string;
  type: CategoryType;
  name: string;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
