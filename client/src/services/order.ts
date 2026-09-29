import { Order, OrderDetail } from '@/types/order';
import { PosItem, PosItemsParams } from '@/types/pos';
import { ApiResponse, PaginationResult } from '@/types/response';
import apiClient from './apiClient';
import { API_ENDPOINTS } from '@/constants';
const ORDERS_ENDPOINT = API_ENDPOINTS.ORDERS;

// API đơn hàng trả entity Go không có json tag (keys PascalCase),
// map về shape chuẩn của app tại đây.
type RawOrder = {
  ID: string;
  OrganizationID: string;
  Code: string;
  CustomerID: string | null;
  AppointmentID: string | null;
  SubtotalAmount: number;
  DiscountAmount: number;
  TotalAmount: number;
  Status: Order['status'];
  CreatedAt: string | null;
  UpdatedAt: string | null;
};

function mapOrder(raw: RawOrder): Order {
  return {
    id: raw.ID,
    organization_id: raw.OrganizationID,
    code: raw.Code,
    customer_id: raw.CustomerID,
    appointment_id: raw.AppointmentID,
    subtotal_amount: raw.SubtotalAmount,
    discount_amount: raw.DiscountAmount,
    total_amount: raw.TotalAmount,
    status: raw.Status,
    created_at: raw.CreatedAt,
    updated_at: raw.UpdatedAt,
  };
}

export const getOrderDetail = async (id: string) => {
  const res = await apiClient.get<ApiResponse<OrderDetail>>(
    `${ORDERS_ENDPOINT.GET_ORDERS}/${id}`,
  );

  const detail = res.data.data;
  return {
    ...detail,
    items: detail.items || [],
    materials: detail.materials || [],
    payments: detail.payments || [],
  };
};
export const getPosItems = async ({
  search,
  category,
  cursorType,
  cursorName,
  cursorId,
  pageSize,
}: PosItemsParams = {}) => {
  const res = await apiClient.get<ApiResponse<PosItem[]>>(ORDERS_ENDPOINT.GET_POS_ITEMS, {
    params: {
      search: search || undefined,
      category: category || undefined,
      type: cursorType || undefined,
      name: cursorName || undefined,
      id: cursorId || undefined,
      page_size: pageSize,
    },
  });

  return res.data.data || [];
};

export const getOrdersByOrgID = async ({
  status,
  page,
  pageSize,
}: {
  status?: string;
  page: number;
  pageSize: number;
}) => {
  const res = await apiClient.get<ApiResponse<PaginationResult<RawOrder>>>(
    ORDERS_ENDPOINT.GET_ORDERS,
    {
      params: {
        status: status || undefined,
        page: page,
        page_size: pageSize,
      },
    },
  );

  const paged = res.data.data;
  return {
    ...paged,
    data: (paged.data || []).map(mapOrder),
  };
};
