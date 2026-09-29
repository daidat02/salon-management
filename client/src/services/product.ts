import { Product } from '@/types/product';
import { ApiResponse, PaginationResult } from '@/types/response';
import apiClient from './apiClient';
import { API_ENDPOINTS } from '@/constants';
const PRODUCTS_ENDPOINT = API_ENDPOINTS.PRODUCTS;

export const getProductsByOrgID = async ({
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
  const res = await apiClient.get<ApiResponse<PaginationResult<Product>>>(
    PRODUCTS_ENDPOINT.GET_PRODUCTS,
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

export type CreateProductPayload = {
  sku: string;
  name: string;
  unit?: string;
  net_unit?: string;
  product_type?: Product['product_type'];
  category_id?: string | null;
  is_active?: boolean;
};

export const createProduct = async (payload: CreateProductPayload) => {
  const res = await apiClient.post<ApiResponse<Product>>(PRODUCTS_ENDPOINT.GET_PRODUCTS, payload);
  return res.data.data;
};
