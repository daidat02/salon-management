'use client';

import * as React from 'react';
import { DashDialog, DashDialogFooter } from '@/components/dashboard/DashDialog';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { DashInput } from '@/components/dashboard/DashInput';
import { DashSelect } from '@/components/dashboard/DashSelect';
import { useCreateCustomer } from '@/hooks/use-create-customer';
import { getApiErrorMessage } from '@/services/apiClient';
import { useSearchParams, useRouter } from 'next/navigation';
import { toast } from '@/components/ui/toast';

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
};

export function CustomerCreateDialog({ open, onOpenChange, onSuccess }: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [gender, setGender] = React.useState('female');
  const [error, setError] = React.useState<string | null>(null);
  const { mutate, isPending } = useCreateCustomer();

  React.useEffect(() => {
    if (!open) {
      setError(null);
      setGender('female');
    }
  }, [open]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    const form = e.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form).entries()) as Record<string, string>;
    const payload = {
      full_name: data.full_name?.trim(),
      phone: data.phone?.trim().replace(/\s+/g, ''),
      gender,
      birth_date: data.birth_date || null,
      note: data.note?.trim() || '',
    } as const;
    if (!payload.full_name || !payload.phone) {
      setError('Vui lòng nhập họ tên và số điện thoại');
      return;
    }
    mutate(payload as any, {
      onSuccess: () => {
        toast.add({
          title: 'Tạo khách hàng thành công',
          description: `Đã thêm ${payload.full_name} vào hệ thống`,
          type: 'success',
        });
        onSuccess?.();
        const params = new URLSearchParams(searchParams.toString());
        params.set('page', '1');
        router.push(`?${params.toString()}`);
        onOpenChange(false);
        form.reset();
        setGender('female');
      },
      onError: (err) => {
        const msg = getApiErrorMessage(err);
        setError(msg);
        toast.add({
          title: 'Tạo khách hàng thất bại',
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
      title="Thêm mới khách hàng"
      description="Nhập đầy đủ thông tin để lưu vào hệ thống quản lý Salon"
      footer={
        <DashDialogFooter
          onCancel={() => onOpenChange(false)}
          onConfirm={() =>
            (document.getElementById('customer-create-form') as HTMLFormElement)?.requestSubmit()
          }
          cancelText="Hủy"
          confirmText="Lưu khách hàng"
          isLoading={isPending}
          confirmDisabled={isPending}
        />
      }
    >
      <form id="customer-create-form" onSubmit={handleSubmit} className="space-y-4 text-xs">
        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-700">
            {error}
          </div>
        )}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Họ và tên <span className="text-error">*</span>
            </Label>
            <DashInput name="full_name" placeholder="Ví dụ: Nguyễn Thị Mai" required />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Số điện thoại <span className="text-error">*</span>
            </Label>
            <div className="relative">
              <span className="absolute left-3 top-1/2 -translate-y-1/2 font-medium text-slate-500 text-xs">
                +84
              </span>
              <DashInput name="phone" placeholder="0901 234 567" required className="pl-12" />
            </div>
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Email</Label>
            <DashInput name="email" type="email" placeholder="khachhang@example.com" />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Ngày sinh</Label>
            <DashInput name="birth_date" type="date" className="text-slate-600" />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-2">Giới tính</Label>
            <div className="flex h-9 items-center gap-4">
              <label className="flex items-center gap-1.5 cursor-pointer">
                <input
                  type="radio"
                  name="gender"
                  value="female"
                  checked={gender === 'female'}
                  onChange={() => setGender('female')}
                  className="text-primary focus:ring-primary"
                />
                <span>Nữ</span>
              </label>
              <label className="flex items-center gap-1.5 cursor-pointer">
                <input
                  type="radio"
                  name="gender"
                  value="male"
                  checked={gender === 'male'}
                  onChange={() => setGender('male')}
                  className="text-primary focus:ring-primary"
                />
                <span>Nam</span>
              </label>
              <label className="flex items-center gap-1.5 cursor-pointer">
                <input
                  type="radio"
                  name="gender"
                  value="other"
                  checked={gender === 'other'}
                  onChange={() => setGender('other')}
                  className="text-primary focus:ring-primary"
                />
                <span>Khác</span>
              </label>
            </div>
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Nhóm / Phân hạng</Label>
            <DashSelect
              placeholder="Chọn nhóm"
              defaultValue="new"
              options={[
                { label: 'Khách hàng mới (Mặc định)', value: 'new' },
                { label: 'Thành viên thân thiết', value: 'loyal' },
                { label: 'Khách hàng VIP Hạng Bạc', value: 'vip_silver' },
                { label: 'Khách hàng VIP Hạng Vàng', value: 'vip_gold' },
              ]}
            />
          </div>
        </div>
        <div>
          <Label className="block font-semibold text-slate-700 mb-1">Địa chỉ liên hệ</Label>
          <DashInput name="address" placeholder="Số nhà, đường, quận/huyện, tỉnh/thành phố" />
        </div>
        <div>
          <Label className="block font-semibold text-slate-700 mb-1">
            Ghi chú sở thích & Lưu ý dịch vụ
          </Label>
          <Textarea
            name="note"
            placeholder="Đặc điểm chất tóc, da đầu nhạy cảm, màu nhuộm yêu thích, thức uống ưa thích..."
            rows={3}
            className="bg-surface border-outline focus-visible:border-primary"
          />
        </div>
      </form>
    </DashDialog>
  );
}
