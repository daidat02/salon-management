import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { getOrderDetail } from '@/services/order';

export default function useOrderDetail(id: string | undefined) {
  return useQuery({
    queryKey: ['order-detail', id],
    queryFn: () => getOrderDetail(id as string),
    enabled: !!id,
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
