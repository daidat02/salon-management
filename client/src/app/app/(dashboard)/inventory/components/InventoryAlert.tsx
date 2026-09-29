import { TriangleAlert } from 'lucide-react';

export default function InventoryAlert() {
  return (
    <section className="flex flex-col gap-4 rounded-md border border-warning bg-[#FFFBEB] p-4 shadow-sm sm:flex-row sm:items-center sm:justify-between">
      <div className="flex items-start gap-3">
        <div className="mt-0.5 shrink-0 rounded-md bg-[#FEF3C7] p-1.5 text-[#D97706]">
          <TriangleAlert className="h-[22px] w-[22px]" />
        </div>
        <div className="text-xs">
          <span className="mb-0.5 block text-sm font-bold text-[#92400E]">
            Cảnh báo tồn kho an toàn: Có 7 sản phẩm đang dưới ngưỡng tối thiểu
          </span>
          <p className="leading-relaxed text-[#78350F]">
            Các mặt hàng như <em>Dầu gội Keratin Kérastase</em>,{' '}
            <em>Thuốc nhuộm Majirel #09</em> đang thiếu hụt. Vui lòng tạo phiếu nhập kho kịp thời.
          </p>
        </div>
      </div>
      <div className="flex shrink-0 items-center gap-2 self-end sm:self-center">
        <span className="rounded-md bg-warning px-3 py-1.5 text-xs font-semibold text-white shadow-sm">
          Tạo phiếu nhập đề xuất
        </span>
        <span className="rounded-md border border-warning/50 bg-white px-3 py-1.5 text-xs font-medium text-[#92400E]">
          Xem chi tiết 7 SP
        </span>
      </div>
    </section>
  );
}
