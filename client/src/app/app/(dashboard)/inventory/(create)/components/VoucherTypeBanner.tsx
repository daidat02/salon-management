'use client';

import { ArrowDownToLine, ArrowUpToLine, Lock } from 'lucide-react';
import { cn } from 'cn';
import { VOUCHER_META, type VoucherType } from './voucherData';

const DESCRIPTIONS: Record<VoucherType, string> = {
  inbound: 'Nhập hàng mới từ nhà cung cấp — tồn kho và giá vốn cập nhật ngay khi lưu phiếu.',
  outbound: 'Xuất hàng khỏi kho — số lượng xuất không được vượt tồn kho hiện tại.',
};

/** Banner khóa loại phiếu (mỗi page nhập/xuất riêng, không cho đổi loại giữa chừng). */
export default function VoucherTypeBanner({ type }: { type: VoucherType }) {
  const meta = VOUCHER_META[type];
  const Icon = type === 'inbound' ? ArrowDownToLine : ArrowUpToLine;

  return (
    <div
      className={cn(
        'flex items-center gap-3 rounded-xl border p-3.5 shadow-sm',
        meta.bannerClassName,
      )}
    >
      <span
        className={cn(
          'flex h-10 w-10 shrink-0 items-center justify-center rounded-lg shadow-sm',
          meta.bannerIconClassName,
        )}
      >
        <Icon className="h-[22px] w-[22px]" />
      </span>
      <span className="min-w-0 flex-1">
        <span className="flex items-center gap-2">
          <span className="text-sm font-semibold text-on-surface">
            {type === 'inbound' ? 'Phiếu Nhập Kho (Inbound)' : 'Phiếu Xuất Kho (Outbound)'}
          </span>
          <span className="inline-flex items-center gap-1 rounded-full bg-surface px-2 py-0.5 text-[10px] font-semibold text-on-surface-variant">
            <Lock className="h-3 w-3" />
            Đã khóa
          </span>
        </span>
        <span className="block truncate text-xs text-on-surface-variant">
          {DESCRIPTIONS[type]}
        </span>
      </span>
    </div>
  );
}
