export interface Service {
  id: string;
  organization_id: string;
  category_id: string;
  name: string;
  description: string;
  duration_minutes: number;
  buffer_minutes: number;
  price: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
