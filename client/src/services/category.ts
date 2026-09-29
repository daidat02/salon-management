import { Category } from '@/types/category';
import { ApiResponse, PaginationResult } from '@/types/response';
import apiClient from './apiClient';
import { API_ENDPOINTS } from '@/constants';
const CATEGORIES_ENDPOINT = API_ENDPOINTS.CATEGORIES;

export const getCategoriesByOrgID = async ({
  page,
  pageSize,
}: {
  page: number;
  pageSize: number;
}) => {
  const res = await apiClient.get<ApiResponse<PaginationResult<Category>>>(
    CATEGORIES_ENDPOINT.GET_CATEGORIES,
    {
      params: {
        page: page,
        page_size: pageSize,
      },
    },
  );

  return res.data.data;
};

export type CreateCategoryPayload = {
  type: Category['type'];
  name: string;
  description?: string;
  is_active?: boolean;
};

export const createCategory = async (payload: CreateCategoryPayload) => {
  const res = await apiClient.post<ApiResponse<Category>>(
    CATEGORIES_ENDPOINT.GET_CATEGORIES,
    payload,
  );
  return res.data.data;
};

export type UpdateCategoryPayload = Partial<CreateCategoryPayload>;

export const updateCategory = async (id: string, payload: UpdateCategoryPayload) => {
  const res = await apiClient.put<ApiResponse<Category>>(
    `${CATEGORIES_ENDPOINT.GET_CATEGORIES}/${id}`,
    payload,
  );
  return res.data.data;
};

export const deleteCategory = async (id: string) => {
  const res = await apiClient.delete<ApiResponse<{ id: string }>>(
    `${CATEGORIES_ENDPOINT.GET_CATEGORIES}/${id}`,
  );
  return res.data.data;
};
