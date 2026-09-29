import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { getInventoryTransactions } from '@/services/inventory';

type UseInventoryTransactionsParams = {
  productId?: string;
  page: number;
  pageSize: number;
};

export default function useInventoryTransactions({
  productId,
  page,
  pageSize,
}: UseInventoryTransactionsParams) {
  return useQuery({
    queryKey: ['inventory-transactions', { productId, page, pageSize }],
    queryFn: async () => getInventoryTransactions({ productId, page, pageSize }),
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
