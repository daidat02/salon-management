import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { getServicesByOrgID } from '@/services/service';
type UseServiceParams = {
  page: number;
  pageSize: number;
  search?: string;
  status?: string;
  sort?: string;
};

export default function useService({ page, pageSize, search, status, sort }: UseServiceParams) {
  return useQuery({
    queryKey: ['services', { page, pageSize, search, status, sort }],
    queryFn: async () => getServicesByOrgID({ page, pageSize, search, status, sort }),
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000,
  });
}
