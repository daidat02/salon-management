import { InventoryDocument, InventoryDocumentDetail, InventoryTransaction } from '@/types/inventory';
import { ApiResponse, PaginationResult } from '@/types/response';
import apiClient from './apiClient';
import { API_ENDPOINTS } from '@/constants';
const INVENTORY_ENDPOINT = API_ENDPOINTS.INVENTORY;

export const getInventoryDocuments = async ({
  search,
  sort,
  filter,
  page,
  pageSize,
}: {
  search?: string;
  sort?: string;
  filter?: string;
  page: number;
  pageSize: number;
}) => {
  const res = await apiClient.get<ApiResponse<PaginationResult<InventoryDocument>>>(
    INVENTORY_ENDPOINT.GET_DOCUMENTS,
    {
      params: {
        search: search || undefined,
        sort: sort || undefined,
        filter: filter || undefined,
        page: page,
        page_size: pageSize,
      },
    },
  );

  return res.data.data;
};

export type CreateDocumentItemPayload = {
  product_id: string;
  quantity: number;
  unit_price: number;
};

export type CreateDocumentPayload = {
  supplier_id?: string | null;
  order_id?: string | null;
  document_code: string;
  type: 'import' | 'export' | 'adjust' | 'sale' | 'return';
  note?: string | null;
  reason?: string | null;
  items: CreateDocumentItemPayload[];
};

export const createInventoryDocument = async (payload: CreateDocumentPayload) => {
  const res = await apiClient.post<ApiResponse<InventoryDocument>>(
    INVENTORY_ENDPOINT.GET_DOCUMENTS,
    payload,
  );

  return res.data.data;
};

export const getInventoryDocument = async (id: string) => {
  const res = await apiClient.get<ApiResponse<InventoryDocumentDetail>>(
    `${INVENTORY_ENDPOINT.GET_DOCUMENTS}${id}`,
  );

  return res.data.data;
};

export const getInventoryTransactions = async ({
  productId,
  page,
  pageSize,
}: {
  productId?: string;
  page: number;
  pageSize: number;
}) => {
  const res = await apiClient.get<ApiResponse<PaginationResult<InventoryTransaction>>>(
    INVENTORY_ENDPOINT.GET_TRANSACTIONS,
    {
      params: {
        product_id: productId || undefined,
        page: page,
        page_size: pageSize,
      },
    },
  );

  return res.data.data;
};
