import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { getCustomersByOrgID } from '@/services/customer';

type UseCustomersParams = {
  page: number;
  pageSize: number;
  search?: string;
  status?: string;
};

export default function useCustomer({ page, pageSize, search, status }: UseCustomersParams) {
  return useQuery({
    queryKey: ['customers', { page, pageSize, search, status }],
    queryFn: () => getCustomersByOrgID({ page, pageSize, search }),
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
