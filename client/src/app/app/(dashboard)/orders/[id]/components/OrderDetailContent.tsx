'use client';

import * as React from 'react';
import { useParams, useRouter } from 'next/navigation';
import { ArrowLeft, Pencil, Printer, Undo2 } from 'lucide-react';
import PageHeading from '@/components/dashboard/PageHeading';
import { Tag } from '@/components/ui/tag';
import { DashboardSkeleton } from '@/components/dashboard/DashboardSkeletons';
import { getApiErrorMessage } from '@/services/apiClient';
import useOrderDetail from '@/hooks/use-order-detail';
import useCustomer from '@/hooks/use-customer';
import useService from '@/hooks/use-service';
import useProduct from '@/hooks/use-product';
import useStaff from '@/hooks/use-staff';
import type { Customer } from '@/types/customer';
import type { Product } from '@/types/product';
import type { Service } from '@/types/service';
import type { Staff } from '@/types/staff';
import type { OrderStatus } from '@/types/order';
import OrderDetailMain from './OrderDetailMain';
import OrderDetailSide from './OrderDetailSide';
import {
  buildMaterialLines,
  buildRetailLines,
  buildServiceLines,
  buildStaffCommissions,
  calcOrderProfit,
  formatDateTime,
} from './orderDetailUtils';

const statusMeta: Record<OrderStatus, { label: string; variant: 'success' | 'warning' | 'purple' | 'teal' | 'blue' | 'gray' | 'slate' }> = {
  paid: { label: 'Đã thanh toán', variant: 'success' },
  pending_payment: { label: 'Chưa thanh toán', variant: 'warning' },
  draft: { label: 'Nháp', variant: 'slate' },
  serving: { label: 'Đang phục vụ', variant: 'blue' },
  completed: { label: 'Hoàn tất', variant: 'teal' },
  cancelled: { label: 'Đã hủy', variant: 'gray' },
  refunded: { label: 'Đã hoàn tiền', variant: 'purple' },
};

function toMap<T extends { id: string }>(list: (T | null | undefined)[]): Map<string, T> {
  const map = new Map<string, T>();
  for (const item of list ?? []) {
    if (item) map.set(item.id, item);
  }
  return map;
}

export default function OrderDetailContent() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const id = params?.id;

  const { data: detail, isLoading, isError, error } = useOrderDetail(id);
  const { data: customerData } = useCustomer({ page: 1, pageSize: 100 });
  const { data: serviceData } = useService({ page: 1, pageSize: 100 });
  const { data: productData } = useProduct({ page: 1, pageSize: 100 });
  const { data: staffData } = useStaff({ page: 1, pageSize: 100 });

  const services = React.useMemo(
    () => toMap<Service>(serviceData?.data ?? []),
    [serviceData],
  );
  const products = React.useMemo(
    () => toMap<Product>(productData?.data ?? []),
    [productData],
  );
  const staff = React.useMemo(() => toMap<Staff>(staffData?.data ?? []), [staffData]);
  const customers = React.useMemo(
    () => toMap<Customer>(customerData?.data ?? []),
    [customerData],
  );

  const enriched = React.useMemo(() => {
    if (!detail) return null;
    const serviceLines = buildServiceLines(detail.items, services, staff);
    const materialLines = buildMaterialLines(detail.materials, products);
    const retailLines = buildRetailLines(detail.items, products);
    return {
      serviceLines,
      materialLines,
      retailLines,
      profit: calcOrderProfit(detail, retailLines),
      commissions: buildStaffCommissions(detail.items, staff),
      customer: detail.order.customer_id
        ? customers.get(detail.order.customer_id)
        : undefined,
    };
  }, [detail, services, products, staff, customers]);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <DashboardSkeleton />
      </div>
    );
  }

  if (isError || !detail || !enriched) {
    return (
      <div className="space-y-4">
        <div className="rounded-md border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
          Không tải được chi tiết đơn hàng: {getApiErrorMessage(error)} — thử đăng nhập lại hoặc
          kiểm tra kết nối API.
        </div>
      </div>
    );
  }

  const { order } = detail;
  const meta = statusMeta[order.status] ?? { label: order.status, variant: 'slate' as const };

  return (
    <div className="space-y-4">
      <button
        type="button"
        onClick={() => router.back()}
        className="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium text-primary transition-colors hover:bg-surface-container"
      >
        <ArrowLeft className="h-4 w-4" />
        Quay lại danh sách đơn
      </button>

      <PageHeading
        breadcrumbs={[
          { label: 'Trang chủ', href: '#' },
          { label: 'Đơn hàng', href: '/app/orders' },
          { label: order.code },
        ]}
        title={`Chi tiết đơn hàng ${order.code}`}
        badge={meta.label}
        badgeVariant="soft"
        description={
          <p className="mt-0.5 flex items-center gap-2 text-xs text-on-surface-variant">
            <span>Thời gian tạo đơn: {formatDateTime(order.created_at)}</span>
          </p>
        }
        actions={[
          { label: 'In hóa đơn', icon: Printer, variant: 'outline', onClick: () => window.print() },
          { label: 'Hoàn tiền', icon: Undo2, variant: 'outline' },
          { label: 'Chỉnh sửa', icon: Pencil, variant: 'primary' },
        ]}
      />
      <div className="flex items-center gap-2">
        <Tag variant={meta.variant} shape="pill" size="md" dot>
          {meta.label}
        </Tag>
      </div>

      <div className="grid grid-cols-1 items-start gap-6 lg:grid-cols-12">
        <div className="space-y-6 lg:col-span-8">
          <OrderDetailMain
            services={enriched.serviceLines}
            materials={enriched.materialLines}
            retail={enriched.retailLines}
          />
        </div>
        <div className="space-y-6 lg:col-span-4">
          <OrderDetailSide
            customer={enriched.customer}
            commissions={enriched.commissions}
            profit={enriched.profit}
            order={order}
            payments={detail.payments}
            cashierName="Thu ngân"
          />
        </div>
      </div>
    </div>
  );
}
