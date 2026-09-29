'use client';

import type { RefObject } from 'react';
import CatalogCard from './CatalogCard';
import { getApiErrorMessage } from '@/services/apiClient';
import type { CatalogItem } from './catalogTypes';

type CatalogGridProps = {
  items: CatalogItem[];
  isLoading: boolean;
  isError: boolean;
  error: unknown;
  sentinelRef: RefObject<HTMLDivElement | null>;
  isFetchingNextPage: boolean;
  hasNextPage: boolean | undefined;
  onAdd: (item: CatalogItem) => void;
};

function SkeletonCards({ count }: { count: number }) {
  return (
    <div className="grid grid-cols-3 gap-3">
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="bg-surface animate-pulse rounded-md border border-outline p-3">
          <div className="bg-surface-container mb-2 h-4 w-1/3 rounded" />
          <div className="bg-surface-container mb-2 h-4 w-full rounded" />
          <div className="bg-surface-container h-4 w-2/3 rounded" />
        </div>
      ))}
    </div>
  );
}

export default function CatalogGrid({
  items,
  isLoading,
  isError,
  error,
  sentinelRef,
  isFetchingNextPage,
  hasNextPage,
  onAdd,
}: CatalogGridProps) {
  return (
    <div className="pr-1">
      {isLoading ? (
        <SkeletonCards count={6} />
      ) : isError ? (
        <div className="rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
          Không tải được danh sách bán hàng: {getApiErrorMessage(error)} — thử đăng nhập lại hoặc
          kiểm tra kết nối API.
        </div>
      ) : (
        <div className="grid grid-cols-3 gap-3">
          {items.map((item) => (
            <CatalogCard key={item.id} item={item} onAdd={onAdd} />
          ))}
        </div>
      )}
      {!isLoading && !isError && items.length === 0 && (
        <div className="py-10 text-center text-sm text-on-surface-variant">
          Không tìm thấy sản phẩm nào
        </div>
      )}
      {/* Sentinel cuộn vô hạn + trạng thái */}
      <div ref={sentinelRef} aria-hidden className="h-1" />
      {isFetchingNextPage && <SkeletonCards count={3} />}
      {!hasNextPage && items.length > 0 && (
        <div className="py-4 text-center text-xs text-on-surface-variant">
          Đã hiển thị tất cả {items.length} món
        </div>
      )}
    </div>
  );
}
