import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { getInventoryDocument } from '@/services/inventory';

export default function useInventoryDocument(id: string | undefined) {
  return useQuery({
    queryKey: ['inventory-document', id],
    queryFn: async () => getInventoryDocument(id as string),
    enabled: !!id,
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
