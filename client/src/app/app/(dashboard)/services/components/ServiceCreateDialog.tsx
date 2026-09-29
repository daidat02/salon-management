'use client';

import * as React from 'react';
import { DashDialog, DashDialogFooter } from '@/components/dashboard/DashDialog';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { DashInput } from '@/components/dashboard/DashInput';
import { DashSelect } from '@/components/dashboard/DashSelect';
import { useCreateService } from '@/hooks/use-create-service';
import { getApiErrorMessage } from '@/services/apiClient';
import { toast } from '@/components/ui/toast';

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
};

export function ServiceCreateDialog({ open, onOpenChange, onSuccess }: Props) {
  const [isActive, setIsActive] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const { mutate, isPending } = useCreateService();

  React.useEffect(() => {
    if (!open) {
      setError(null);
      setIsActive(true);
    }
  }, [open]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    const form = e.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form).entries()) as Record<string, string>;

    const payload = {
      name: data.name?.trim(),
      description: data.description?.trim() || '',
      duration_minutes: Number(data.duration_minutes || 0),
      buffer_minutes: Number(data.buffer_minutes || 0),
      price: Number(data.price || 0),
      is_active: isActive,
    };

    if (!payload.name || !payload.duration_minutes || !payload.price) {
      setError('Vui lòng nhập tên, thời lượng và giá');
      return;
    }

    mutate(payload as any, {
      onSuccess: () => {
        toast.add({
          title: 'Tạo dịch vụ thành công',
          description: `Đã thêm ${payload.name}`,
          type: 'success',
        });
        onSuccess?.();
        onOpenChange(false);
        form.reset();
        setIsActive(true);
      },
      onError: (err) => {
        const msg = getApiErrorMessage(err);
        setError(msg);
        toast.add({ title: 'Tạo dịch vụ thất bại', description: msg, type: 'error' });
      },
    });
  };

  return (
    <DashDialog
      open={open}
      onOpenChange={onOpenChange}
      title="Thêm mới dịch vụ"
      description="Nhập thông tin dịch vụ để quản lý bảng giá và thời lượng"
      footer={
        <DashDialogFooter
          onCancel={() => onOpenChange(false)}
          onConfirm={() => (document.getElementById('service-create-form') as HTMLFormElement)?.requestSubmit()}
          cancelText="Hủy"
          confirmText="Lưu dịch vụ"
          isLoading={isPending}
          confirmDisabled={isPending}
        />
      }
    >
      <form id="service-create-form" onSubmit={handleSubmit} className="space-y-4 text-xs">
        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-700">
            {error}
          </div>
        )}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div className="sm:col-span-2">
            <Label className="block font-semibold text-slate-700 mb-1">
              Tên dịch vụ <span className="text-error">*</span>
            </Label>
            <DashInput name="name" placeholder="VD: Cắt tóc nam cao cấp" required />
          </div>
          <div className="sm:col-span-2">
            <Label className="block font-semibold text-slate-700 mb-1">Mô tả</Label>
            <Textarea name="description" placeholder="Mô tả chi tiết dịch vụ..." rows={2} className="bg-surface border-outline focus-visible:border-primary" />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Thời lượng (phút) <span className="text-error">*</span>
            </Label>
            <DashInput name="duration_minutes" type="number" placeholder="45" required />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Thời gian đệm (phút)</Label>
            <DashInput name="buffer_minutes" type="number" placeholder="10" />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Giá (đ) <span className="text-error">*</span>
            </Label>
            <DashInput name="price" type="number" placeholder="150000" required />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Danh mục</Label>
            <DashSelect
              placeholder="Chọn danh mục"
              options={[
                { label: 'Cắt tóc', value: 'cat_cut' },
                { label: 'Nhuộm', value: 'cat_dye' },
                { label: 'Uốn & Duỗi', value: 'cat_perm' },
                { label: 'Gội đầu & Chăm sóc', value: 'cat_wash' },
              ]}
            />
          </div>
        </div>
        <div className="flex items-center justify-between rounded-lg border border-outline bg-surface-container-low px-3 py-2">
          <div>
            <p className="text-xs font-semibold text-slate-700">Đang hoạt động</p>
            <p className="text-[11px] text-on-surface-variant">Tắt nếu tạm ngưng dịch vụ</p>
          </div>
          <Switch checked={isActive} onCheckedChange={setIsActive} />
        </div>
      </form>
    </DashDialog>
  );
}
