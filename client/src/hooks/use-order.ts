import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { getOrdersByOrgID } from '@/services/order';

type UseOrdersParams = {
  page: number;
  pageSize: number;
  status?: string;
};

export default function useOrders({ page, pageSize, status }: UseOrdersParams) {
  return useQuery({
    queryKey: ['orders', { page, pageSize, status }],
    queryFn: async () => getOrdersByOrgID({ status, page, pageSize }),
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
