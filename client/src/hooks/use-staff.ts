import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { getStaffsByOrgID } from '@/services/staff';

type UseStaffsParams = {
  page: number;
  pageSize: number;
  search?: string;
  status?: string;
  sort?: string;
};

export default function useStaff({ page, pageSize, search, status, sort }: UseStaffsParams) {
  return useQuery({
    queryKey: ['staff', { page, pageSize, search, status, sort }],
    queryFn: () => getStaffsByOrgID({ page, pageSize, search, status, sort }),
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
