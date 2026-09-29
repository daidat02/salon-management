import { useInfiniteQuery, keepPreviousData } from '@tanstack/react-query';
import { getPosItems } from '@/services/order';
import type { PosItem, PosItemsParams } from '@/types/pos';

type Cursor = { cursorType: string; cursorName: string; cursorId: string };

export type PosItemsBaseParams = Omit<PosItemsParams, 'cursorType' | 'cursorName' | 'cursorId'>;

export default function usePosItems(params: PosItemsBaseParams = {}) {
  const { search, category, pageSize } = params;
  return useInfiniteQuery<PosItem[], Error>({
    queryKey: ['pos-items', { search, category, pageSize }],
    queryFn: ({ pageParam }) =>
      getPosItems({ search, category, pageSize, ...((pageParam as Cursor | undefined) ?? {}) }),
    initialPageParam: undefined as Cursor | undefined,
    getNextPageParam: (lastPage) => {
      const size = pageSize ?? 20;
      if (lastPage.length < size) return undefined;
      const last = lastPage[lastPage.length - 1];
      return { cursorType: last.type, cursorName: last.name, cursorId: last.id };
    },
    placeholderData: keepPreviousData,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
