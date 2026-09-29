'use client';

import * as React from 'react';
import { ArrowUpToLine, Pencil, Plus, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Tag } from '@/components/ui/tag';
import { useUpdateCategory } from '@/hooks/use-update-category';
import useService from '@/hooks/use-service';
import { getApiErrorMessage } from '@/services/apiClient';
import { toast } from '@/components/ui/toast';
import type { Category } from '@/types/category';
import { getCategoryIcon, getCategoryTypeLabel } from './data';

type Props = {
  category: Category | null;
  position: number;
  onClose: () => void;
  onEdit: (category: Category) => void;
};

export default function CategoryDetail({ category, position, onClose, onEdit }: Props) {
  const { mutate: updateMutate } = useUpdateCategory();
  const { data: servicesData } = useService({ page: 1, pageSize: 100 });
  const services = (servicesData?.data ?? []).filter(
    (s) => s && category && s.category_id === category.id,
  );

  if (!category) {
    return (
      <div className="bg-surface border border-dashed border-outline rounded-xl p-8 text-center text-sm text-on-surface-variant shadow-sm">
        Chọn một danh mục bên trái để xem chi tiết
      </div>
    );
  }

  const iconElement = React.createElement(getCategoryIcon(category.name, category.type), {
    className: 'h-4 w-4',
  });

  const handleToggle = (checked: boolean) => {
    updateMutate(
      { id: category.id, payload: { is_active: checked } },
      {
        onSuccess: () => {
          toast.add({
            title: checked ? 'Đã kích hoạt danh mục' : 'Đã tạm ẩn danh mục',
            description: category.name,
            type: 'success',
          });
        },
        onError: (err) => {
          toast.add({ title: 'Cập nhật thất bại', description: getApiErrorMessage(err), type: 'error' });
        },
      },
    );
  };

  return (
    <div className="bg-surface border border-outline rounded-xl p-5 shadow-sm space-y-5 sticky top-20">
      <div className="flex items-center justify-between pb-3 border-b border-outline">
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-primary-container text-primary flex items-center justify-center">
            {iconElement}
          </div>
          <div>
            <h2 className="font-headline font-semibold text-base text-on-surface">
              Chi tiết danh mục: {category.name}
            </h2>
            <p className="text-[11px] text-on-surface-variant">
              Cập nhật {new Date(category.updated_at).toLocaleDateString('vi-VN')}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Tag variant="primary" size="sm" shape="rounded">
            {services.length} {category.type === 'service' ? 'dịch vụ' : 'sản phẩm'}
          </Tag>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={onClose}
            title="Đóng panel"
            className="text-on-surface-variant hover:text-on-surface"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div className="space-y-3 bg-surface-container-low p-3.5 rounded-lg border border-outline/60 text-xs">
        <div className="grid grid-cols-2 gap-3">
          <div>
            <span className="text-on-surface-variant block mb-1">Tên danh mục</span>
            <span className="font-semibold text-on-surface">{category.name}</span>
          </div>
          <div>
            <span className="text-on-surface-variant block mb-1">Mã danh mục</span>
            <span className="font-mono font-semibold bg-white px-2 py-0.5 rounded border border-outline text-primary inline-block">
              {category.id.slice(0, 8).toUpperCase()}
            </span>
          </div>
        </div>
        <div className="grid grid-cols-2 gap-3 pt-2 border-t border-outline/50">
          <div>
            <span className="text-on-surface-variant block mb-1">Vị trí hiển thị</span>
            <span className="font-semibold text-on-surface flex items-center gap-1">
              <ArrowUpToLine className="h-3.5 w-3.5 text-primary" />
              Thứ tự #{position} (Ưu tiên trên POS)
            </span>
          </div>
          <div>
            <span className="text-on-surface-variant block mb-1">Phân nhóm</span>
            <span className="font-semibold text-on-surface">
              {getCategoryTypeLabel(category.type)} trực tiếp
            </span>
          </div>
        </div>
        <div className="pt-2 border-t border-outline/50 flex items-center justify-between">
          <div>
            <span className="text-on-surface-variant block mb-0.5">Trạng thái hoạt động</span>
            <span className="text-[11px] text-success font-medium">
              {category.is_active ? 'Đang kích hoạt và hiển thị cho nhân viên' : 'Đang tạm ẩn'}
            </span>
          </div>
          <Switch checked={category.is_active} onCheckedChange={handleToggle} />
        </div>
        <div className="pt-2 border-t border-outline/50">
          <span className="text-on-surface-variant block mb-1">Mô tả hiển thị</span>
          <p className="text-on-surface font-normal italic">
            {category.description ? `"${category.description}"` : 'Chưa có mô tả'}
          </p>
        </div>
      </div>

      <div className="space-y-2.5">
        <div className="flex items-center justify-between">
          <h3 className="text-xs font-bold uppercase tracking-wider text-on-surface">
            {category.type === 'service' ? 'Dịch vụ' : 'Sản phẩm'} tiêu biểu trong nhóm (
            {Math.min(services.length, 4)}/{services.length})
          </h3>
        </div>
        {services.length === 0 ? (
          <p className="text-xs text-on-surface-variant border border-dashed border-outline rounded-lg p-4 text-center">
            Chưa có {category.type === 'service' ? 'dịch vụ' : 'sản phẩm'} nào trong nhóm này
          </p>
        ) : (
          <div className="divide-y divide-outline border border-outline rounded-lg bg-surface overflow-hidden text-xs">
            {services.slice(0, 4).map((s, i) => (
              <div
                key={s.id}
                className="p-3 flex items-center justify-between hover:bg-surface-container-low transition-colors"
              >
                <div className="flex items-center gap-2.5">
                  <div className="w-7 h-7 rounded-md bg-slate-100 flex items-center justify-center text-on-surface-variant font-bold text-xs">
                    {String(i + 1).padStart(2, '0')}
                  </div>
                  <div>
                    <h4 className="font-semibold text-on-surface">{s.name}</h4>
                    <p className="text-[11px] text-on-surface-variant">
                      {s.duration_minutes} phút • Giá:{' '}
                      {Number(s.price).toLocaleString('vi-VN')} đ
                    </p>
                  </div>
                </div>
                <Tag variant="success" size="sm" shape="rounded">
                  Đang mở
                </Tag>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="pt-3 border-t border-outline flex items-center gap-2.5">
        <Button
          variant="outline"
          className="flex-1 py-2 px-3 text-xs font-semibold shadow-sm"
          onClick={() => onEdit(category)}
        >
          <Plus className="h-4 w-4 text-primary" />
          Thêm {category.type === 'service' ? 'dịch vụ' : 'sản phẩm'} vào nhóm
        </Button>
        <Button
          onClick={() => onEdit(category)}
          className="py-2 px-4 text-xs font-semibold bg-primary text-white shadow-sm"
        >
          <Pencil className="h-4 w-4" />
          Chỉnh sửa danh mục
        </Button>
      </div>
    </div>
  );
}
