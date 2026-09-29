'use client';

import * as React from 'react';
import { DashDialog, DashDialogFooter } from '@/components/dashboard/DashDialog';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { DashInput } from '@/components/dashboard/DashInput';
import { DashSelect } from '@/components/dashboard/DashSelect';
import { useCreateProduct } from '@/hooks/use-create-product';
import useCategory from '@/hooks/use-category';
import { getApiErrorMessage } from '@/services/apiClient';
import { toast } from '@/components/ui/toast';

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
};

export function ProductCreateDialog({ open, onOpenChange, onSuccess }: Props) {
  const [productType, setProductType] = React.useState<'retail' | 'material' | 'both'>('retail');
  const [categoryId, setCategoryId] = React.useState('');
  const [netUnit, setNetUnit] = React.useState('');
  const [minStock, setMinStock] = React.useState('');
  const [isActive, setIsActive] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const { mutate, isPending } = useCreateProduct();
  const { data: categoryData, isLoading: isCategoriesLoading } = useCategory({
    page: 1,
    pageSize: 100,
  });
  const productCategories = (categoryData?.data ?? []).filter(
    (c) => c.type === 'product' && c.is_active,
  );

  React.useEffect(() => {
    if (!open) {
      setError(null);
      setProductType('retail');
      setCategoryId('');
      setNetUnit('');
      setMinStock('');
      setIsActive(true);
    }
  }, [open]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    const form = e.target as HTMLFormElement;
    const data = Object.fromEntries(new FormData(form).entries()) as Record<string, string>;

    const payload = {
      sku: data.sku?.trim(),
      name: data.name?.trim(),
      unit: data.unit?.trim(),
      net_unit: netUnit || '',
      net_amount: data.net_amount ? parseFloat(data.net_amount) : 0,
      product_type: productType,
      category_id: categoryId || undefined,
      cost_price: data.cost_price_first ? parseFloat(data.cost_price_first) : 0,
      sell_price: data.sell_price_first ? parseFloat(data.sell_price_first) : 0,
      stock_quantity: data.stock_first ? parseFloat(data.stock_first) : 0,
      min_stock: data.min_stock ? parseFloat(data.min_stock) : 0,
      is_active: isActive,
    };
    console.log('Submitting product create payload:', payload);
    if (!payload.sku || !payload.name) {
      setError('Vui lòng nhập SKU và tên sản phẩm');
      return;
    }

    mutate(payload as any, {
      onSuccess: () => {
        toast.add({
          title: 'Tạo sản phẩm thành công',
          description: `Đã thêm ${payload.name} (${payload.sku})`,
          type: 'success',
        });
        onSuccess?.();
        onOpenChange(false);
        form.reset();
        setProductType('retail');
        setCategoryId('');
        setNetUnit('');
        setMinStock('');
        setIsActive(true);
      },
      onError: (err) => {
        const msg = getApiErrorMessage(err);
        setError(msg);
        toast.add({ title: 'Tạo sản phẩm thất bại', description: msg, type: 'error' });
      },
    });
  };

  return (
    <DashDialog
      open={open}
      onOpenChange={onOpenChange}
      title="Thêm mới sản phẩm"
      description="Nhập thông tin sản phẩm để quản lý tồn kho và bán lẻ"
      footer={
        <DashDialogFooter
          onCancel={() => onOpenChange(false)}
          onConfirm={() =>
            (document.getElementById('product-create-form') as HTMLFormElement)?.requestSubmit()
          }
          cancelText="Hủy"
          confirmText="Lưu sản phẩm"
          isLoading={isPending}
          confirmDisabled={isPending}
        />
      }
    >
      <form id="product-create-form" onSubmit={handleSubmit} className="space-y-4 text-xs">
        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-700">
            {error}
          </div>
        )}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Mã SKU</Label>
            <DashInput name="sku" placeholder="VD: SKU-TN-1001" />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">
              Tên sản phẩm <span className="text-error">*</span>
            </Label>
            <DashInput name="name" placeholder="VD: Dầu gội Keratin 500ml" required />
          </div>

          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Danh mục</Label>
            <DashSelect
              value={categoryId}
              onValueChange={setCategoryId}
              placeholder={isCategoriesLoading ? 'Đang tải danh mục...' : 'Chọn danh mục'}
              disabled={isCategoriesLoading || productCategories.length === 0}
              options={productCategories.map((c) => ({ label: c.name, value: c.id }))}
            />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Loại sản phẩm</Label>
            <DashSelect
              value={productType}
              onValueChange={(v) => setProductType(v as any)}
              placeholder="Chọn loại"
              options={[
                { label: 'Bán lẻ', value: 'retail' },
                { label: 'Vật tư tiêu hao', value: 'material' },
                { label: 'Cả hai', value: 'both' },
              ]}
            />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Đơn vị</Label>
            <DashInput name="unit" placeholder="VD: chai" />
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <Label className="block font-semibold text-slate-700 mb-1">Dung tích (net)</Label>
              <DashInput name="net_amount" placeholder="VD: 500" />
            </div>
            <div>
              <Label className="block font-semibold text-slate-700 mb-1">Đơn vị tịnh</Label>
              <DashSelect
                value={netUnit}
                onValueChange={setNetUnit}
                placeholder="Chọn đơn vị"
                options={[
                  { label: 'ml', value: 'ml' },
                  { label: 'g (gam)', value: 'g' },
                ]}
              />
            </div>
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Giá Nhập Ban Đầu</Label>
            <DashInput type="number" name="cost_price_first" placeholder="VD: 500000" />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Giá Bán Ban Đầu</Label>
            <DashInput type="number" name="sell_price_first" placeholder="VD: 500000" />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Tồn Kho Ban đầu</Label>
            <DashInput type="number" name="stock_first" placeholder="VD: 50" />
          </div>
          <div>
            <Label className="block font-semibold text-slate-700 mb-1">Ngưỡng Tối Thiểu</Label>
            <DashInput type="number" name="min_stock" placeholder="VD: 50" />
          </div>
        </div>
        <div className="flex items-center justify-between rounded-lg border border-outline bg-surface-container-low px-3 py-2">
          <div>
            <p className="text-xs font-semibold text-slate-700">Đang kinh doanh</p>
            <p className="text-[11px] text-on-surface-variant">Tắt nếu ngừng bán / hết hàng</p>
          </div>
          <Switch checked={isActive} onCheckedChange={setIsActive} />
        </div>
        <span className="text-xs text-on-surface-variant px-2">
          Ghi chú: Hệ thống sẽ tự động sinh mã SKU nếu bạn bỏ trống.
        </span>
      </form>
    </DashDialog>
  );
}
