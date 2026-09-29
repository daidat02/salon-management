'use client';

import * as React from 'react';
import { GripVertical, Pencil, Trash2 } from 'lucide-react';
import { DashDialog, DashDialogFooter } from '@/components/dashboard/DashDialog';
import { Tag } from '@/components/ui/tag';
import { TableSkeleton } from '@/components/dashboard/DashboardSkeletons';
import { getApiErrorMessage } from '@/services/apiClient';
import { toast } from '@/components/ui/toast';
import { useRouter, useSearchParams } from 'next/navigation';
import useCategory from '@/hooks/use-category';
import useService from '@/hooks/use-service';
import { useDeleteCategory } from '@/hooks/use-delete-category';
import type { Category } from '@/types/category';
import { cn } from 'cn';
import CategoryDetail from './CategoryDetail';
import { CategoryCreateDialog } from './CategoryCreateDialog';
import { getCategoryIcon } from './data';

export default function CategoryGrid() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const tab = searchParams.get('tab') || 'service';
  const search = searchParams.get('search') || '';
  const status = searchParams.get('status') || '';

  const { data, isLoading, isError, error } = useCategory({ page: 1, pageSize: 100 });
  const { data: servicesData } = useService({ page: 1, pageSize: 100 });

  const [selectedId, setSelectedId] = React.useState<string | null>(null);
  const [editing, setEditing] = React.useState<Category | null>(null);
  const [editOpen, setEditOpen] = React.useState(false);
  const [deleting, setDeleting] = React.useState<Category | null>(null);
  const { mutate: deleteMutate, isPending: isDeleting } = useDeleteCategory();

  const categories = data?.data ?? [];
  const serviceCountByCategory = React.useMemo(() => {
    const services = servicesData?.data ?? [];
    const map = new Map<string, number>();
    for (const s of services) {
      if (!s) continue;
      map.set(s.category_id, (map.get(s.category_id) ?? 0) + 1);
    }
    return map;
  }, [servicesData]);

  const filtered = categories
    .filter((c) => c.type === tab)
    .filter((c) =>
      search
        ? c.name.toLowerCase().includes(search.toLowerCase()) ||
          c.description?.toLowerCase().includes(search.toLowerCase())
        : true,
    )
    .filter((c) => {
      if (status === 'active') return c.is_active;
      if (status === 'inactive') return !c.is_active;
      return true;
    });

  const selected =
    filtered.find((c) => c.id === selectedId) ?? filtered[0] ?? null;
  const selectedIndex = selected ? filtered.findIndex((c) => c.id === selected.id) + 1 : 0;

  const handleEdit = (category: Category) => {
    setEditing(category);
    setEditOpen(true);
  };

  const handleDelete = () => {
    if (!deleting) return;
    deleteMutate(deleting.id, {
      onSuccess: () => {
        toast.add({ title: 'Đã xóa danh mục', description: deleting.name, type: 'success' });
        if (selectedId === deleting.id) setSelectedId(null);
        setDeleting(null);
      },
      onError: (err) => {
        toast.add({ title: 'Xóa thất bại', description: getApiErrorMessage(err), type: 'error' });
      },
    });
  };

  if (isError) {
    return (
      <div className="rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
        Không tải được danh sách danh mục: {getApiErrorMessage(error)}
      </div>
    );
  }

  if (isLoading && categories.length === 0) {
    return <TableSkeleton rows={4} cols={4} />;
  }

  return (
    <>
      <div className="grid grid-cols-1 xl:grid-cols-12 gap-6 items-start">
        <div className="xl:col-span-7 space-y-4">
          <div className="flex items-center justify-between px-1">
            <div className="flex items-center gap-2">
              <span className="text-xs font-bold uppercase tracking-wider text-on-surface-variant">
                Danh sách hiển thị trên POS
              </span>
              <span className="text-[11px] bg-primary-container text-primary font-semibold px-2 py-0.5 rounded-full">
                {filtered.length} thẻ
              </span>
            </div>
            <span className="text-xs text-on-surface-variant hidden md:flex items-center gap-1">
              <GripVertical className="h-3.5 w-3.5" />
              Kéo thả thẻ để đổi thứ tự
            </span>
          </div>

          {filtered.length === 0 ? (
            <div className="rounded-xl border border-dashed border-outline bg-surface p-8 text-center text-sm text-on-surface-variant">
              Không tìm thấy danh mục nào. Thử đổi từ khóa hoặc{' '}
              <button
                className="font-semibold text-primary hover:underline"
                onClick={() => router.push('/app/categories?tab=service', { scroll: false })}
              >
                đặt lại bộ lọc
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {filtered.map((category) => {
                const isSelected = selected?.id === category.id;
                const iconElement = React.createElement(
                  getCategoryIcon(category.name, category.type),
                  { className: 'h-5 w-5' },
                );
                const count = serviceCountByCategory.get(category.id) ?? 0;
                return (
                  <div
                    key={category.id}
                    onClick={() => setSelectedId(category.id)}
                    className={cn(
                      'group relative bg-surface rounded-xl p-4 shadow-sm transition-all cursor-pointer',
                      isSelected
                        ? 'border-2 border-primary shadow-md ring-2 ring-primary/10'
                        : 'border border-outline hover:border-primary/50 hover:shadow-md',
                    )}
                  >
                    {isSelected && (
                      <div className="absolute -top-2.5 right-4 bg-primary text-white text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider shadow-sm">
                        Đang xem
                      </div>
                    )}
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className="w-11 h-11 rounded-xl bg-primary-container text-primary flex items-center justify-center shadow-inner">
                          {iconElement}
                        </div>
                        <div>
                          <h3 className="font-semibold text-base text-on-surface leading-snug">
                            {category.name}
                          </h3>
                          <span className="text-[11px] font-medium text-on-surface-variant">
                            Mã: {category.id.slice(0, 8).toUpperCase()}
                          </span>
                        </div>
                      </div>
                      <button
                        className="text-on-surface-variant hover:text-on-surface p-1 rounded hover:bg-surface-container"
                        title="Kéo để đổi vị trí"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <GripVertical className="h-4 w-4" />
                      </button>
                    </div>
                    <p className="text-xs text-on-surface-variant line-clamp-2 min-h-[32px] mb-3">
                      {category.description || 'Chưa có mô tả cho danh mục này.'}
                    </p>
                    <div className="pt-3 border-t border-outline flex items-center justify-between text-xs">
                      <div className="flex items-center gap-2">
                        <span className="px-2 py-0.5 rounded-md bg-blue-50 text-primary font-semibold border border-primary/20">
                          {count} {category.type === 'service' ? 'dịch vụ' : 'sản phẩm'}
                        </span>
                        <Tag variant={category.is_active ? 'success' : 'gray'} size="sm" shape="rounded" dot>
                          {category.is_active ? 'Hoạt động' : 'Tạm ẩn'}
                        </Tag>
                      </div>
                      <div className="flex items-center gap-1">
                        <button
                          className="p-1.5 text-on-surface-variant hover:text-primary hover:bg-surface-container rounded-lg transition-colors"
                          title="Chỉnh sửa"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleEdit(category);
                          }}
                        >
                          <Pencil className="h-4 w-4" />
                        </button>
                        <button
                          className="p-1.5 text-on-surface-variant hover:text-error hover:bg-surface-container rounded-lg transition-colors"
                          title="Xóa"
                          onClick={(e) => {
                            e.stopPropagation();
                            setDeleting(category);
                          }}
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>

        <div className="xl:col-span-5">
          <CategoryDetail
            category={selected}
            position={selectedIndex}
            onClose={() => setSelectedId(null)}
            onEdit={handleEdit}
          />
        </div>
      </div>

      <CategoryCreateDialog
        open={editOpen}
        onOpenChange={setEditOpen}
        initial={editing}
        onSuccess={() => setEditing(null)}
      />

      <DashDialog
        open={!!deleting}
        onOpenChange={(open) => !open && setDeleting(null)}
        title="Xóa danh mục"
        description={`Bạn có chắc muốn xóa "${deleting?.name}"? Các dịch vụ/sản phẩm trong nhóm sẽ mất liên kết danh mục.`}
        footer={
          <DashDialogFooter
            onCancel={() => setDeleting(null)}
            onConfirm={handleDelete}
            cancelText="Hủy"
            confirmText="Xóa danh mục"
            isLoading={isDeleting}
            confirmDisabled={isDeleting}
          />
        }
      >
        <p className="text-xs text-on-surface-variant">
          Hành động này không thể hoàn tác. Hãy chắc chắn nhóm không còn mục nào quan trọng.
        </p>
      </DashDialog>
    </>
  );
}
