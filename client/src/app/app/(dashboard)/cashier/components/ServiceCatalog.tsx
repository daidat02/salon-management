'use client';

import CatalogToolbar from './CatalogToolbar';
import CatalogGrid from './CatalogGrid';
import useCatalogData from './useCatalogData';
import { dotFor, type CatalogItem } from './catalogTypes';

type Props = {
  onAdd: (item: CatalogItem) => void;
};

export default function ServiceCatalog({ onAdd }: Props) {
  const {
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
  } = useCatalogData();

  return (
    <section className="flex min-w-0 flex-1 flex-col gap-3 bg-transparent">
      <CatalogToolbar
        search={search}
        onSearchChange={setSearch}
        categories={categories}
        activeCat={activeCat}
        onCategoryChange={setActiveCat}
      />
      <CatalogGrid
        items={items}
        isLoading={isPosLoading}
        isError={isPosError}
        error={posError}
        sentinelRef={sentinelRef}
        isFetchingNextPage={isFetchingNextPage}
        hasNextPage={hasNextPage}
        onAdd={onAdd}
      />
    </section>
  );
}

export { dotFor };
export type { CatalogItem };
