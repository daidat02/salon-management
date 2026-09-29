import { Customer } from '@/types/customer';
import { ApiResponse, PaginationResult } from '@/types/response';
import apiClient from './apiClient';

export type GetCustomersParams = {
  page?: number;
  pageSize?: number;
  search?: string;
};

export const getCustomersByOrgID = async ({
  page = 1,
  pageSize = 10,
  search = '',
}: GetCustomersParams = {}) => {
  const res = await apiClient.get<ApiResponse<PaginationResult<Customer>>>(`/customers`, {
    params: {
      page,
      page_size: pageSize,
      search: search || undefined,
    },
  });
  // Backend wraps pagination.Result inside SuccessResponse.Data
  // => res.data.data = { data: Customer[], page, page_size, total, total_pages }
  return res.data.data;
};

export type CreateCustomerPayload = {
  full_name: string;
  phone: string;
  gender?: Customer['gender'];
  birth_date?: string | null;
  note?: string;
};

export const createCustomer = async (payload: CreateCustomerPayload) => {
  const res = await apiClient.post<ApiResponse<Customer>>('/customers', payload);
  return res.data.data;
};
