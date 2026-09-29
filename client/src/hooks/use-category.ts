import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { getCategoriesByOrgID } from '@/services/category';
type UseCategoryParams = {
  page: number;
  pageSize: number;
};

export default function useCategory({ page, pageSize }: UseCategoryParams) {
  return useQuery({
    queryKey: ['categories', { page, pageSize }],
    queryFn: async () => getCategoriesByOrgID({ page, pageSize }),
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000,
  });
}
