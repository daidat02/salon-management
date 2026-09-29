'use client';

import { Skeleton } from '@/components/ui/skeleton';

export function StatSkeleton() {
  return (
    <div className="flex flex-col justify-between rounded-xl border border-outline bg-surface px-4 py-3.5 shadow-sm">
      <div className="flex items-center justify-between gap-3">
        <Skeleton className="h-3.5 w-24" />
        <Skeleton className="h-7 w-7 rounded-lg" />
      </div>
      <div className="mt-2 space-y-1.5">
        <Skeleton className="h-6 w-16" />
        <Skeleton className="h-3 w-32" />
      </div>
    </div>
  );
}

export function StatsSkeleton({ count = 4 }: { count?: number }) {
  return (
    <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {Array.from({ length: count }).map((_, i) => (
        <StatSkeleton key={i} />
      ))}
    </section>
  );
}

export function ToolbarSkeleton() {
  return (
    <section className="bg-surface flex flex-col gap-3 rounded-md border border-outline p-3.5 shadow-sm lg:flex-row lg:items-center lg:justify-between">
      <div className="flex w-full flex-1 flex-col items-center gap-2.5 sm:flex-row lg:w-auto">
        <Skeleton className="h-8 w-full sm:w-72 rounded-md" />
        <Skeleton className="h-8 w-full sm:w-44 rounded-md" />
        <Skeleton className="h-8 w-full sm:w-40 rounded-md" />
        <Skeleton className="h-8 w-full sm:w-44 rounded-md" />
      </div>
      <div className="flex shrink-0 items-center gap-2 self-end lg:self-auto">
        <Skeleton className="h-7 w-7 rounded-md" />
        <Skeleton className="h-7 w-24 rounded-md" />
      </div>
    </section>
  );
}

export function TableSkeleton({ rows = 5, cols = 7 }: { rows?: number; cols?: number }) {
  return (
    <div className="bg-surface overflow-hidden rounded-md border border-outline shadow-sm">
      <div className="overflow-x-auto">
        <table className="w-full text-left">
          <thead>
            <tr className="bg-surface-container border-b border-outline">
              {Array.from({ length: cols }).map((_, i) => (
                <th key={i} className="px-3 py-3">
                  <Skeleton className="h-3 w-full max-w-20" />
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-outline">
            {Array.from({ length: rows }).map((_, r) => (
              <tr key={r}>
                {Array.from({ length: cols }).map((_, c) => (
                  <td key={c} className="px-3 py-3">
                    {c === 0 ? (
                      <div className="flex items-center gap-3">
                        <Skeleton className="h-8 w-8 rounded-full" />
                        <div className="space-y-1.5">
                          <Skeleton className="h-3 w-24" />
                          <Skeleton className="h-2.5 w-20" />
                        </div>
                      </div>
                    ) : c === 1 ? (
                      <Skeleton className="h-3 w-32" />
                    ) : (
                      <Skeleton className="h-3 w-full max-w-22.5" />
                    )}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="flex flex-col items-center justify-between gap-3 border-t border-outline bg-surface px-6 py-3.5 sm:flex-row">
        <Skeleton className="h-3 w-48" />
        <div className="flex items-center gap-1">
          <Skeleton className="h-7 w-16 rounded" />
          <Skeleton className="h-7 w-7 rounded" />
          <Skeleton className="h-7 w-7 rounded" />
          <Skeleton className="h-7 w-7 rounded" />
          <Skeleton className="h-7 w-16 rounded" />
        </div>
      </div>
    </div>
  );
}

export function PageHeadingSkeleton() {
  return (
    <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div className="space-y-2">
        <Skeleton className="h-3 w-40" />
        <Skeleton className="h-6 w-48" />
        <Skeleton className="h-3 w-64" />
      </div>
      <div className="flex items-center gap-2.5">
        <Skeleton className="h-8 w-28 rounded-md" />
        <Skeleton className="h-8 w-36 rounded-md" />
      </div>
    </div>
  );
}

export function DashboardSkeleton() {
  return (
    <div className="space-y-4">
      <PageHeadingSkeleton />
      <StatsSkeleton count={4} />
      <ToolbarSkeleton />
      <TableSkeleton rows={6} cols={7} />
    </div>
  );
}
