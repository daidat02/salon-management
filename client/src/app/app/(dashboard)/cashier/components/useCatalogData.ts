'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import usePosItems from '@/hooks/use-pos-items';
import useCategory from '@/hooks/use-category';
import type { CatalogItem } from './catalogTypes';

export default function useCatalogData() {
  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [activeCat, setActiveCat] = useState('');

  useEffect(() => {
    const t = setTimeout(() => setDebouncedSearch(search.trim()), 400);
    return () => clearTimeout(t);
  }, [search]);

  const { data: categoryData } = useCategory({ page: 1, pageSize: 100 });
  const categories = useMemo(
    () => (categoryData?.data ?? []).filter((c) => c && c.is_active),
    [categoryData],
  );
  const categoryNameById = useMemo(() => {
    const map = new Map<string, string>();
    for (const c of categories) map.set(c.id, c.name);
    return map;
  }, [categories]);

  const {
    data: posPages,
    isLoading: isPosLoading,
    isError: isPosError,
    error: posError,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = usePosItems({ search: debouncedSearch, category: activeCat || undefined, pageSize: 24 });

  const posItems = useMemo(() => (posPages?.pages ?? []).flat(), [posPages]);

  // Cuộn tới cuối grid thì tải thêm (cursor pagination).
  const sentinelRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    const el = sentinelRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
          fetchNextPage();
        }
      },
      { rootMargin: '200px' },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

  // Mô tả/thời lượng/đơn vị lấy thẳng từ API pos-items (backend JOIN sẵn),
  // không cần gọi thêm services/products.
  const items: CatalogItem[] = useMemo(
    () =>
      (posItems ?? []).map((it) => {
        const category = (it.category_id && categoryNameById.get(it.category_id)) || 'Khác';
        if (it.type === 'service') {
          return {
            id: it.id,
            name: it.name,
            desc: it.description ?? '',
            price: it.price,
            duration: it.duration_minutes ? `${it.duration_minutes}p` : undefined,
            image: '',
            type: 'service' as const,
            category,
          };
        }
        return {
          id: it.id,
          name: it.name,
          desc: '',
          price: it.price,
          stock: it.stock_quantity ?? 0,
          image: '',
          type: 'product' as const,
          category,
        };
      }),
    [posItems, categoryNameById],
  );

  return {
    search,
    setSearch,
    categories,
    activeCat,
    setActiveCat,
    items,
    isPosLoading,
    isPosError,
    posError,
    sentinelRef,
    isFetchingNextPage,
    hasNextPage,
  };
}
