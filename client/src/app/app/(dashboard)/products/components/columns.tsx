import { Tag } from '@/components/ui/tag';
import type { DashColumn } from '@/components/dashboard/DashTable';
import { Product } from '@/types/product';

export const productColumns: DashColumn<Product>[] = [
  {
    id: 'sku',
    header: 'Mã SKU',
    className:
      'font-mono text-primary font-semibold text-[11px] min-w-[120px] hover:underline hover:text-primary/80 transition-colors cursor-pointer',
    accessorKey: 'sku',
  },
  {
    id: 'name',
    header: 'Tên sản phẩm',
    className: 'min-w-[240px]',
    cell: (row) => (
      <div>
        <div className="font-semibold text-on-surface">{row.name}</div>
        <div className="text-[11px] text-on-surface-variant line-clamp-1">Đơn vị: {row.unit}</div>
      </div>
    ),
  },
  {
    id: 'category',
    header: 'Danh mục',
    className: 'min-w-[120px] ',
    cell: (row) => 'Chưa phân loại',
  },
  {
    id: 'product_type',
    header: 'Loại sản phẩm',
    className: 'min-w-[120px]',
    cell: (row) => {
      switch (row.product_type) {
        case 'retail':
          return 'Bán lẻ';
        case 'material':
          return 'Vật tư tiêu hao';
        case 'both':
          return 'Cả hai';
        default:
          return 'Chưa xác định';
      }
    },
  },
  {
    id: 'cost_price',
    header: 'Đơn giá nhập',
    className: 'text-right font-medium',
    cell: (row) =>
      row.cost_price == 0 ? 'Chưa nhập' : row.cost_price.toLocaleString('vi-VN') + ' đ', // Format tiền tệ VNĐ cho đẹp
  },
  {
    id: 'sell_price',
    header: 'Đơn giá bán',
    className: 'text-right font-medium',
    cell: (row) =>
      row.sell_price == 0 ? 'Chưa nhập' : row.sell_price.toLocaleString('vi-VN') + ' đ', // Format tiền tệ VNĐ cho đẹp
  },
  {
    id: 'revenue_mounth',
    header: 'Doanh thu tháng',
    className: 'text-center',
    cell: (row) => 'Chưa Thống Kê',
  },
  {
    id: 'stock_quantity',
    header: 'Tồn kho',
    className: 'text-center',
    cell: (row) => (
      <Tag
        variant={
          row.stock_quantity === 0
            ? 'danger'
            : row.stock_quantity <= row.min_stock
              ? 'warning'
              : 'slate'
        }
        size="sm"
        shape="rounded"
      >
        {row.stock_quantity
          ? row.net_amount && row.net_amount > 0
            ? `${row.stock_quantity / row.net_amount} ${row.unit ?? ''}`
            : `${row.stock_quantity} ${row.unit ?? ''}`
          : 'Hết hàng'}
      </Tag>
    ),
  },
  {
    id: 'is_active',
    header: 'Trạng thái',
    className: 'text-center',
    cell: (row) => (
      <Tag variant={row.is_active ? 'success' : 'gray'} shape="pill" size="md" dot>
        {row.is_active ? 'Đang kinh doanh' : 'Ngừng bán'}
      </Tag>
    ),
  },
  {
    id: 'updated_at',
    header: 'Cập nhật lần cuối',
    className: 'text-center',
    cell: (row) =>
      new Date(row.updated_at).toLocaleString('vi-VN', { dateStyle: 'short', timeStyle: 'short' }),
  },
];
