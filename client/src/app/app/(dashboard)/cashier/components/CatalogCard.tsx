'use client';

import { Plus, Timer } from 'lucide-react';
import { Tag } from '@/components/ui/tag';
import { tagVariantFor, type CatalogItem } from './catalogTypes';

type CatalogCardProps = {
  item: CatalogItem;
  onAdd: (item: CatalogItem) => void;
};

export default function CatalogCard({ item, onAdd }: CatalogCardProps) {
  return (
    <div
      onClick={() => onAdd(item)}
      className="bg-surface border-outline hover:border-primary/60 group flex cursor-pointer flex-col justify-between rounded-md border p-3 transition-all duration-200 hover:shadow-md active:scale-[0.98]"
    >
      <div>
        <div className="mb-2 flex items-start justify-between">
          <Tag
            variant={tagVariantFor(item)}
            size="sm"
            shape="rounded"
            className="text-[10px] tracking-wide uppercase"
          >
            {item.type === 'service' ? 'Dịch vụ' : 'Sản phẩm'}
          </Tag>
          <span className="text-on-surface-variant flex items-center gap-0.5 text-[11px]">
            {item.type === 'service' ? (
              item.duration ? (
                <>
                  <Timer className="h-3 w-3" /> {item.duration}
                </>
              ) : null
            ) : (
              <span className="bg-emerald-50 text-emerald-600 rounded px-1.5 py-0.2 text-[10px] font-medium">
                Kho: {item.stock ?? 0}
              </span>
            )}
          </span>
        </div>
        <h3 className="group-hover:text-primary mb-1 line-clamp-2 text-xs font-semibold text-on-surface transition-colors">
          {item.name}
        </h3>
        <p className="text-on-surface-variant mb-2 line-clamp-1 text-[11px]">
          {item.desc || item.category}
        </p>
      </div>
      <div className="mt-auto flex items-center justify-between border-t border-slate-100 pt-2">
        <span className="text-primary text-xs font-bold">
          {item.price.toLocaleString('vi-VN')} đ
        </span>
        <span className="bg-primary-container text-primary group-hover:bg-primary flex h-6 w-6 items-center justify-center rounded group-hover:text-white transition-colors">
          <Plus className="h-4 w-4" />
        </span>
      </div>
    </div>
  );
}
