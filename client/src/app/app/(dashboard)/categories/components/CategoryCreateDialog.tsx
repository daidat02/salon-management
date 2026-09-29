'use client';

import * as React from 'react';
import { DashDialog, DashDialogFooter } from '@/components/dashboard/DashDialog';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { DashInput } from '@/components/dashboard/DashInput';
import { DashSelect } from '@/components/dashboard/DashSelect';
import { useCreateCategory } from '@/hooks/use-create-category';
import { useUpdateCategory } from '@/hooks/use-update-category';
import { getApiErrorMessage } from '@/services/apiClient';
import { toast } from '@/components/ui/toast';
import type { Category } from '@/types/category';

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  initial?: Category | null;
  onSuccess?: () => void;
};

export function CategoryCreateDialog({ open, onOpenChange, initial, onSuccess }: Props) {
  const isEdit = !!initial;
  const [type, setType] = React.useState<'service' | 'product'>(initial?.type ?? 'service');
  const [isActive, setIsActive] = React.useState(initial?.is_active ?? true);
  const [error, setError] = React.useState<string | null>(null);
  const { mutate: createMutate, isPending: isCreating } = useCreateCategory();
  const { mutate: updateMutate, isPending: isUpdating } = useUpdateCategory();
  const isPending = isCreating || isUpdating;

  const [lastOpen, setLastOpen] = React.useState(open);
  if (open !== lastOpen) {
    setLastOpen(open);
    if (open) {
      setError(null);
      setType(initial?.type ?? 'service');
      setIsActive(initial?.is_active ?? true);
    }
  }
  const formKey = `${open}-${initial?.id ?? 'new'}`;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    const form = e.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form).entries()) as Record<string, string>;
    const payload = {
      type,
      name: data.name?.trim(),
      description: data.description?.trim() || '',
      is_active: isActive,
    };

    if (!payload.name) {
      setError('Vui lòng nhập tên danh mục');
      return;
    }

    const done = {
      onSuccess: () => {
        toast.add({
          title: isEdit ? 'Cập nhật danh mục thành công' : 'Tạo danh mục thành công',
          description: payload.name,
          type: 'success',
        });
        onSuccess?.();
        onOpenChange(false);
      },
      onError: (err: unknown) => {
        const msg = getApiErrorMessage(err);
        setError(msg);
        toast.add({ title: 'Lưu danh mục thất bại', description: msg, type: 'error' });
      },
    };

    if (isEdit && initial) updateMutate({ id: initial.id, payload }, done);
    else createMutate(payload, done);
  };

  return (
    <DashDialog
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? 'Chỉnh sửa danh mục' : 'Thêm mới danh mục'}
      description="Nhập thông tin danh mục để phân loại dịch vụ và sản phẩm"
      footer={
        <DashDialogFooter
          onCancel={() => onOpenChange(false)}
          onConfirm={() => (document.getElementById('category-form') as HTMLFormElement)?.requestSubmit()}
          cancelText="Hủy"
          confirmText={isEdit ? 'Lưu thay đổi' : 'Lưu danh mục'}
          isLoading={isPending}
          confirmDisabled={isPending}
        />
      }
    >
      <form key={formKey} id="category-form" onSubmit={handleSubmit} className="space-y-4 text-xs">
        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-700">
            {error}
          </div>
        )}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Tên danh mục <span className="text-error">*</span>
            </Label>
            <DashInput name="name" defaultValue={initial?.name ?? ''} placeholder="VD: Cắt tóc" required />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Phân nhóm</Label>
            <DashSelect
              value={type}
              onValueChange={(v) => setType(v as 'service' | 'product')}
              placeholder="Chọn nhóm"
              options={[
                { label: 'Dịch vụ', value: 'service' },
                { label: 'Sản phẩm', value: 'product' },
              ]}
            />
          </div>
        </div>
        <div>
          <Label className="block font-semibold text-slate-700 mb-1">Mô tả</Label>
          <Textarea
            name="description"
            defaultValue={initial?.description ?? ''}
            placeholder="Mô tả ngắn về danh mục..."
            rows={3}
            className="bg-surface border-outline focus-visible:border-primary"
          />
        </div>
        <div className="flex items-center justify-between rounded-lg border border-outline bg-surface-container-low px-3 py-2">
          <div>
            <p className="text-xs font-semibold text-slate-700">Đang hoạt động</p>
            <p className="text-[11px] text-on-surface-variant">Tắt để tạm ẩn khỏi POS</p>
          </div>
          <Switch checked={isActive} onCheckedChange={setIsActive} />
        </div>
      </form>
    </DashDialog>
  );
}
