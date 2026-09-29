export interface Customer {
  id: string;
  organization_id: string;
  full_name: string;
  phone: string;
  gender: 'male' | 'female' | 'other';
  birth_date: string | null;
  note: string;
  total_spent: number;
  total_visits: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
}
