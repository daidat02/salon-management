export interface Staff {
  id: string;
  organization_id: string;
  code: string;
  full_name: string;
  phone: string;
  position: string;
  commission_rate: number;
  status: 'active' | 'inactive';
  hire_date: string;
  created_at: string | null;
  updated_at: string | null;
  deleted_at: string | null;
}
