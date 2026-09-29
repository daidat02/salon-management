'use client';

import * as React from 'react';
import { DashDialog, DashDialogFooter } from '@/components/dashboard/DashDialog';
import { Label } from '@/components/ui/label';
import { DashInput } from '@/components/dashboard/DashInput';
import { DashSelect } from '@/components/dashboard/DashSelect';
import { useCreateStaff } from '@/hooks/use-create-staff';
import { getApiErrorMessage } from '@/services/apiClient';
import { useSearchParams, useRouter } from 'next/navigation';
import { toast } from '@/components/ui/toast';

const POSITION_OPTIONS = [
  { label: 'Master Stylist', value: 'Master Stylist' },
  { label: 'Stylist', value: 'Stylist' },
  { label: 'Kỹ thuật viên', value: 'Kỹ thuật viên' },
  { label: 'Lễ tân', value: 'Lễ tân' },
  { label: 'Quản lý salon', value: 'Quản lý salon' },
];

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
};

export function StaffCreateDialog({ open, onOpenChange, onSuccess }: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [position, setPosition] = React.useState('');
  const [status, setStatus] = React.useState<'active' | 'inactive'>('active');
  const [error, setError] = React.useState<string | null>(null);
  const { mutate, isPending } = useCreateStaff();

  React.useEffect(() => {
    if (!open) {
      setError(null);
      setPosition('');
      setStatus('active');
    }
  }, [open ]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    const form = e.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form).entries()) as Record<string, string>;
    const commission = Number(data.commission_rate);
    const payload = {
      code: data.code?.trim(),
      full_name: data.full_name?.trim(),
      phone: data.phone?.trim().replace(/\s+/g, ''),
      position,
      commission_rate: Number.isNaN(commission) ? 0 : commission,
      status,
      hire_date: data.hire_date || new Date().toISOString().slice(0, 10),
    } as const;
    if (!payload.code || !payload.full_name || !payload.phone || !payload.position) {
      setError('Vui lòng nhập mã NV, họ tên, số điện thoại và vị trí');
      return;
    }
    if (payload.commission_rate < 0 || payload.commission_rate > 100) {
      setError('Tỷ lệ hoa hồng phải nằm trong khoảng từ 0 đến 100');
      return;
    }
    mutate(payload, {
      onSuccess: () => {
        toast.add({
          title: 'Tạo nhân viên thành công',
          description: `Đã thêm ${payload.full_name} vào hệ thống`,
          type: 'success',
        });
        onSuccess?.();
        const params = new URLSearchParams(searchParams.toString());
        params.set('page', '1');
        router.push(`?${params.toString()}`);
        onOpenChange(false);
        form.reset();
        setPosition('');
        setStatus('active');
      },
      onError: (err) => {
        const msg = getApiErrorMessage(err);
        setError(msg);
        toast.add({
          title: 'Tạo nhân viên thất bại',
          description: msg,
          type: 'error',
        });
      },
    });
  };

  return (
    <DashDialog
      open={open}
      onOpenChange={onOpenChange}
      title="Thêm mới nhân viên"
      description="Nhập đầy đủ thông tin để lưu vào hệ thống quản lý Salon"
      footer={
        <DashDialogFooter
          onCancel={() => onOpenChange(false)}
          onConfirm={() =>
            (document.getElementById('staff-create-form') as HTMLFormElement)?.requestSubmit()
          }
          cancelText="Hủy"
          confirmText="Lưu nhân viên"
          isLoading={isPending}
          confirmDisabled={isPending}
        />
      }
    >
      <form id="staff-create-form" onSubmit={handleSubmit} className="space-y-4 text-xs">
        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-700">
            {error}
          </div>
        )}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Mã nhân viên <span className="text-error">*</span>
            </Label>
            <DashInput name="code" placeholder="Ví dụ: NV-01" required />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Họ và tên <span className="text-error">*</span>
            </Label>
            <DashInput name="full_name" placeholder="Ví dụ: Nguyễn Anh Tuấn" required />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Số điện thoại <span className="text-error">*</span>
            </Label>
            <DashInput name="phone" placeholder="0901 234 567" required />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Vị trí <span className="text-error">*</span>
            </Label>
            <DashSelect
              placeholder="Chọn vị trí"
              value={position}
              onValueChange={setPosition}
              options={POSITION_OPTIONS}
            />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Hoa hồng (%)
            </Label>
            <DashInput
              name="commission_rate"
              type="number"
              min={0}
              max={100}
              step={0.5}
              placeholder="Ví dụ: 10"
            />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Ngày vào làm</Label>
            <DashInput name="hire_date" type="date" className="text-slate-600" />
          </div>
          <div className="sm:col-span-2">
            <Label className="block font-semibold text-slate-700 mb-2">Trạng thái</Label>
            <div className="flex h-9 items-center gap-4">
              <label className="flex items-center gap-1.5 cursor-pointer">
                <input
                  type="radio"
                  name="status"
                  value="active"
                  checked={status === 'active'}
                  onChange={() => setStatus('active')}
                  className="text-primary focus:ring-primary"
                />
                <span>Đang hoạt động</span>
              </label>
              <label className="flex items-center gap-1.5 cursor-pointer">
                <input
                  type="radio"
                  name="status"
                  value="inactive"
                  checked={status === 'inactive'}
                  onChange={() => setStatus('inactive')}
                  className="text-primary focus:ring-primary"
                />
                <span>Tạm ngưng</span>
              </label>
            </div>
          </div>
        </div>
      </form>
    </DashDialog>
  );
}
