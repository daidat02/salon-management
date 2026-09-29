import { Service } from '@/types/service';
import { ApiResponse, PaginationResult } from '@/types/response';
import apiClient from './apiClient';
import { API_ENDPOINTS } from '@/constants';
const SERVICES_ENDPOINT = API_ENDPOINTS.SERVICES;

export const getServicesByOrgID = async ({
  search,
  page,
  pageSize,
  status,
  sort,
}: {
  search?: string;
  page: number;
  pageSize: number;
  status?: string;
  sort?: string;
}) => {
  const res = await apiClient.get<ApiResponse<PaginationResult<Service>>>(
    SERVICES_ENDPOINT.GET_SERVICES,
    {
      params: {
        search: search,
        page: page,
        page_size: pageSize,
        status: status,
        sort: sort,
      },
    },
  );

  return res.data.data;
};

export type CreateServicePayload = {
  category_id?: string | null;
  name: string;
  description?: string;
  duration_minutes: number;
  buffer_minutes?: number;
  price: number;
  is_active?: boolean;
};

export const createService = async (payload: CreateServicePayload) => {
  const res = await apiClient.post<ApiResponse<Service>>(SERVICES_ENDPOINT.GET_SERVICES, payload);
  return res.data.data;
};
