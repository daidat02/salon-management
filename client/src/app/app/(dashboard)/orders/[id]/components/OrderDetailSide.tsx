'use client';

import { BadgeDollarSign, ChartColumn, Phone, Star, User, Wallet } from 'lucide-react';
import { Tag } from '@/components/ui/tag';
import type { Customer } from '@/types/customer';
import type { Order, OrderDetailPayment } from '@/types/order';
import type { OrderProfit, StaffCommission } from './orderDetailUtils';
import { customerTier, formatDateTime, formatVND } from './orderDetailUtils';

function SideCard({
  icon: Icon,
  title,
  extra,
  children,
}: {
  icon: typeof User;
  title: string;
  extra?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-4 rounded-xl border border-outline bg-surface p-5 shadow-sm">
      <div className="flex items-center justify-between border-b border-outline pb-3">
        <h3 className="flex items-center gap-2 text-sm font-semibold text-on-surface">
          <Icon className="h-[18px] w-[18px] text-primary" />
          {title}
        </h3>
        {extra}
      </div>
      {children}
    </div>
  );
}

const METHOD_LABEL: Record<string, string> = {
  cash: 'Tiền mặt',
  transfer: 'Chuyển khoản',
  card: 'Thẻ',
  ewallet: 'Ví điện tử',
};

export default function OrderDetailSide({
  customer,
  commissions,
  profit,
  order,
  payments,
  cashierName,
}: {
  customer: Customer | undefined;
  commissions: StaffCommission[];
  profit: OrderProfit;
  order: Order;
  payments: OrderDetailPayment[];
  cashierName: string;
}) {
  const totalCommission = commissions.reduce((s, c) => s + c.amount, 0);
  const paidTotal = payments.reduce((s, p) => s + p.amount, 0);
  const tier = customerTier(customer?.total_spent ?? 0);
  const mainPayment = payments[0];

  return (
    <>
      <SideCard
        icon={User}
        title="Thông tin khách hàng"
        extra={
          <span className="text-xs font-medium text-primary hover:underline">Hồ sơ chi tiết</span>
        }
      >
        <div className="flex items-center gap-3.5">
          <span className="flex h-14 w-14 shrink-0 items-center justify-center rounded-full border-2 border-primary/20 bg-primary-container text-lg font-bold text-primary shadow-xs">
            {(customer?.full_name || 'K').charAt(0).toUpperCase()}
          </span>
          <div className="space-y-1">
            <h4 className="text-base font-semibold text-on-surface">
              {customer?.full_name ?? 'Khách lẻ'}
            </h4>
            {customer && (
              <p className="flex items-center gap-1 font-mono text-xs text-on-surface-variant">
                <Phone className="h-3.5 w-3.5" />
                {customer.phone}
              </p>
            )}
            <div className="flex items-center gap-1.5 pt-0.5">
              <Tag variant="warning" shape="rounded" size="sm" dot={false}>
                <Star className="h-3 w-3" />
                {tier.label}
              </Tag>
              <span className="rounded-full border border-outline bg-surface-container px-2 py-0.5 text-[11px] font-medium text-on-surface-variant">
                {tier.sub}
              </span>
            </div>
          </div>
        </div>
        {customer && (
          <div className="grid grid-cols-3 gap-2 border-t border-outline pt-2 text-center">
            <div className="rounded-lg bg-surface-container-low p-2">
              <div className="text-[11px] text-on-surface-variant">Tích lũy</div>
              <div className="mt-0.5 text-xs font-bold text-on-surface">
                {(customer.total_spent / 1000000).toLocaleString('vi-VN')} tr
              </div>
            </div>
            <div className="rounded-lg bg-surface-container-low p-2">
              <div className="text-[11px] text-on-surface-variant">Số lần ghé</div>
              <div className="mt-0.5 text-xs font-bold text-on-surface">
                {customer.total_visits} lần
              </div>
            </div>
            <div className="rounded-lg bg-surface-container-low p-2">
              <div className="text-[11px] text-on-surface-variant">Điểm thưởng</div>
              <div className="mt-0.5 text-xs font-bold text-primary">
                {Math.floor(customer.total_spent / 10000).toLocaleString('vi-VN')} đ
              </div>
            </div>
          </div>
        )}
      </SideCard>

      <SideCard icon={BadgeDollarSign} title="Nhân viên & Hoa hồng ca làm" extra={<span className="text-xs text-on-surface-variant">Trực tiếp</span>}>
        <div className="space-y-2.5 text-sm">
          {commissions.map((c) => (
            <div
              key={c.staffId}
              className="flex items-center justify-between rounded-lg bg-surface-container-low p-2.5"
            >
              <div>
                <div className="text-xs font-medium text-on-surface">
                  {c.name} <span className="font-normal text-on-surface-variant">({c.role})</span>
                </div>
                <div className="text-[11px] text-on-surface-variant">{c.rateText}</div>
              </div>
              <div className="text-xs font-semibold text-on-surface">{formatVND(c.amount)}</div>
            </div>
          ))}
          {commissions.length === 0 && (
            <p className="py-2 text-center text-xs text-on-surface-variant">
              Chưa ghi nhận nhân viên thực hiện.
            </p>
          )}
        </div>
        <div className="flex items-center justify-between border-t border-outline pt-2 text-xs">
          <span className="font-medium text-on-surface-variant">Tổng trích quỹ hoa hồng nhân sự:</span>
          <span className="font-bold text-on-surface">{formatVND(totalCommission)}</span>
        </div>
      </SideCard>

      <div className="relative space-y-4 overflow-hidden rounded-xl border-2 border-primary/40 bg-gradient-to-b from-[#F5F8FC] to-surface p-5 shadow-sm">
        <div className="flex items-center justify-between">
          <h3 className="flex items-center gap-2 text-sm font-bold text-on-surface">
            <ChartColumn className="h-[18px] w-[18px] text-primary" />
            Lợi nhuận đơn hàng (P&L)
          </h3>
          <span className="rounded-full bg-primary px-2.5 py-0.5 text-[11px] font-bold tracking-tight text-on-primary shadow-xs">
            BIÊN LÃI {profit.totalMargin.toFixed(1)}%
          </span>
        </div>
        <div className="space-y-2 text-xs">
          <div className="space-y-1.5 rounded-lg border border-outline bg-surface p-3">
            <div className="flex justify-between text-on-surface">
              <span>Doanh thu dịch vụ:</span>
              <span className="font-semibold">+{formatVND(profit.serviceRevenue)}</span>
            </div>
            <div className="flex justify-between text-error">
              <span>Trừ chi phí vật tư:</span>
              <span className="font-medium">-{formatVND(profit.materialCost)}</span>
            </div>
            <div className="flex items-center justify-between border-t border-dashed border-outline pt-1 font-semibold text-on-surface">
              <span>=&gt; Lợi nhuận dịch vụ:</span>
              <span className="font-bold text-primary">
                {formatVND(profit.serviceProfit)}{' '}
                <span className="text-[11px] font-normal text-on-surface-variant">
                  ({profit.serviceMargin.toFixed(1)}%)
                </span>
              </span>
            </div>
          </div>
          <div className="space-y-1.5 rounded-lg border border-outline bg-surface p-3">
            <div className="flex justify-between text-on-surface">
              <span>Doanh thu sản phẩm bán lẻ:</span>
              <span className="font-semibold">+{formatVND(profit.retailRevenue)}</span>
            </div>
            <div className="flex justify-between text-error">
              <span>Trừ vốn sản phẩm:</span>
              <span className="font-medium">-{formatVND(profit.retailCost)}</span>
            </div>
            <div className="flex items-center justify-between border-t border-dashed border-outline pt-1 font-semibold text-on-surface">
              <span>=&gt; Lợi nhuận sản phẩm:</span>
              <span className="font-bold text-primary">
                {formatVND(profit.retailProfit)}{' '}
                <span className="text-[11px] font-normal text-on-surface-variant">
                  ({profit.retailMargin.toFixed(1)}%)
                </span>
              </span>
            </div>
          </div>
        </div>
        <div className="space-y-1.5 border-t border-outline pt-3">
          <div className="flex items-baseline justify-between">
            <span className="text-xs font-bold tracking-wide text-on-surface uppercase">
              Tổng lợi nhuận gộp:
            </span>
            <span className="text-xl font-bold tracking-tight text-primary">
              {formatVND(profit.totalProfit)}
            </span>
          </div>
          <div className="flex items-center justify-between text-xs text-on-surface-variant">
            <span>Tỷ suất biên lợi nhuận ròng:</span>
            <span className="font-semibold text-emerald-600">
              {profit.totalMargin.toFixed(1)}% trên tổng thu {formatVND(profit.totalRevenue)}
            </span>
          </div>
        </div>
      </div>

      <SideCard
        icon={Wallet}
        title="Thông tin Thanh toán & Chứng từ"
        extra={<span className="text-lg text-emerald-600">✓</span>}
      >
        <div className="space-y-2.5 text-xs">
          <div className="flex items-center justify-between py-1">
            <span className="text-on-surface-variant">Phương thức:</span>
            <span className="font-medium text-on-surface">
              {mainPayment ? (METHOD_LABEL[mainPayment.method] ?? mainPayment.method) : '—'}
            </span>
          </div>
          <div className="flex items-center justify-between py-1">
            <span className="text-on-surface-variant">Đã thu đủ:</span>
            <span className="text-base font-bold text-success">{formatVND(paidTotal)}</span>
          </div>
          <div className="flex items-center justify-between py-1">
            <span className="text-on-surface-variant">Thu ngân tiếp nhận:</span>
            <span className="font-medium text-on-surface">{cashierName}</span>
          </div>
          <div className="flex items-center justify-between py-1">
            <span className="text-on-surface-variant">Thời gian hoàn tất:</span>
            <span className="font-medium text-on-surface">
              {formatDateTime(mainPayment?.paid_at ?? order.created_at)}
            </span>
          </div>
          {mainPayment?.provider_txn_id && (
            <div className="flex items-center justify-between border-t border-outline py-1 pt-2">
              <span className="text-on-surface-variant">Mã giao dịch:</span>
              <span className="font-mono font-medium text-primary select-all">
                {mainPayment.provider_txn_id}
              </span>
            </div>
          )}
        </div>
      </SideCard>
    </>
  );
}
