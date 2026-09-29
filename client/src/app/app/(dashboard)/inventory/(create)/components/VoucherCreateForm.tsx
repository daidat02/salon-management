'use client';

import * as React from 'react';
import { useRouter } from 'next/navigation';
import { Download, Info, Printer } from 'lucide-react';
import { toast } from '@/components/ui/toast';
import { getApiErrorMessage } from '@/services/apiClient';
import useCategory from '@/hooks/use-category';
import useProduct from '@/hooks/use-product';
import { useCreateDocument } from '@/hooks/use-create-document';
import type { Product } from '@/types/product';
import VoucherHeading from './VoucherHeading';
import VoucherTypeBanner from './VoucherTypeBanner';
import VoucherInfoCard from './VoucherInfoCard';
import VoucherItemsCard from './VoucherItemsCard';
import VoucherFinanceCard from './VoucherFinanceCard';
import VoucherPaymentCard from './VoucherPaymentCard';
import VoucherAutomationCard from './VoucherAutomationCard';
import {
  VOUCHER_META,
  VOUCHER_TYPE_MAP,
  makeInitialFormState,
  type VoucherFormState,
  type VoucherItem,
  type VoucherType,
} from './voucherData';
import { buildLines, calcTotals, findOverStockLines } from './voucherUtils';

/** Container giữ state phiếu + layout 2 cột, đấu API category / product / tạo phiếu. */
export default function VoucherCreateForm({ initialType }: { initialType: VoucherType }) {
  const router = useRouter();
  const [form, setForm] = React.useState<VoucherFormState>(() => makeInitialFormState(initialType));
  const [productQuery, setProductQuery] = React.useState('');
  const [debouncedQuery, setDebouncedQuery] = React.useState('');
  const [categoryFilter, setCategoryFilter] = React.useState('');
  const [submitError, setSubmitError] = React.useState<string | null>(null);

  const { mutate: createDocument, isPending: isSubmitting } = useCreateDocument();

  // Danh mục loại product để lọc gợi ý thêm dòng.
  const { data: categoryData } = useCategory({ page: 1, pageSize: 100 });
  const productCategories = React.useMemo(
    () => (categoryData?.data ?? []).filter((c) => c && c.type === 'product' && c.is_active),
    [categoryData],
  );

  // Tìm kiếm sản phẩm qua API (debounce).
  React.useEffect(() => {
    const t = setTimeout(() => setDebouncedQuery(productQuery.trim()), 400);
    return () => clearTimeout(t);
  }, [productQuery]);
  const { data: productData, isLoading: isProductsLoading } = useProduct({
    page: 1,
    pageSize: 20,
    search: debouncedQuery,
  });
  const products = React.useMemo(
    () => (productData?.data ?? []).filter((p): p is NonNullable<typeof p> => p != null),
    [productData],
  );

  // Cache sản phẩm đã thêm vào phiếu (dialog picker có list riêng,
  // list tìm kiếm của form có thể không chứa các SP này).
  const [productCache, setProductCache] = React.useState<Product[]>([]);
  const allProducts = React.useMemo(() => {
    const map = new Map<string, Product>();
    for (const p of [...products, ...productCache]) map.set(p.id, p);
    return Array.from(map.values());
  }, [products, productCache]);

  const patch = React.useCallback((p: Partial<VoucherFormState>) => {
    setForm((prev) => ({ ...prev, ...p }));
  }, []);

  const handleUpdateItem = React.useCallback((productId: string, p: Partial<VoucherItem>) => {
    setForm((prev) => ({
      ...prev,
      items: prev.items.map((it) => (it.productId === productId ? { ...it, ...p } : it)),
    }));
  }, []);

  const handlePriceChange = React.useCallback((productId: string, unitPrice: number) => {
    setForm((prev) => ({
      ...prev,
      items: prev.items.map((it) => (it.productId === productId ? { ...it, unitPrice } : it)),
    }));
  }, []);

  const handleRemoveItem = React.useCallback((productId: string) => {
    setForm((prev) => ({ ...prev, items: prev.items.filter((it) => it.productId !== productId) }));
  }, []);

  const handleAddItem = React.useCallback((productId: string) => {
    setForm((prev) =>
      prev.items.some((it) => it.productId === productId)
        ? prev
        : { ...prev, items: [...prev.items, { productId, qty: 1, discount: 0 }] },
    );
  }, []);

  const handleAddItems = React.useCallback((picked: Product[]) => {
    setProductCache((prev) => {
      const known = new Set(prev.map((p) => p.id));
      const fresh = picked.filter((p) => !known.has(p.id));
      return fresh.length > 0 ? [...prev, ...fresh] : prev;
    });
    setForm((prev) => {
      const existing = new Set(prev.items.map((it) => it.productId));
      const fresh = picked
        .filter((p) => !existing.has(p.id))
        .map((p) => ({ productId: p.id, qty: 1, discount: 0 }));
      return fresh.length > 0 ? { ...prev, items: [...prev.items, ...fresh] } : prev;
    });
  }, []);

  const lines = React.useMemo(() => buildLines(form.items, allProducts), [form.items, allProducts]);
  const totals = React.useMemo(
    () => calcTotals(lines, form.shippingFee),
    [lines, form.shippingFee],
  );
  const overStockLines = React.useMemo(
    () => (form.type === 'outbound' ? findOverStockLines(lines) : []),
    [form.type, lines],
  );

  const handleSaveDraft = React.useCallback(() => {
    toast.add({
      title: 'Đã lưu nháp (tạm)',
      description: 'Chức năng lưu nháp server chưa hỗ trợ, phiếu vẫn giữ trên màn hình.',
      type: 'info',
    });
  }, []);

  const handleComplete = React.useCallback(() => {
    setSubmitError(null);
    if (form.items.length === 0) {
      const msg = 'Vui lòng thêm ít nhất 1 sản phẩm vào phiếu';
      setSubmitError(msg);
      toast.add({ title: 'Thiếu sản phẩm', description: msg, type: 'error' });
      return;
    }
    if (!form.code.trim()) {
      const msg = 'Vui lòng nhập mã chứng từ';
      setSubmitError(msg);
      toast.add({ title: 'Thiếu mã phiếu', description: msg, type: 'error' });
      return;
    }
    if (lines.some((l) => l.qty <= 0)) {
      const msg = 'Số lượng mỗi dòng phải lớn hơn 0';
      setSubmitError(msg);
      toast.add({ title: 'Số lượng không hợp lệ', description: msg, type: 'error' });
      return;
    }
    if (overStockLines.length > 0) {
      const names = overStockLines
        .map((l) => `${l.name} (tồn ${l.stock.toLocaleString('vi-VN')} ${l.baseUnit})`)
        .join(', ');
      const msg = `Phiếu xuất vượt tồn kho: ${names}`;
      setSubmitError(msg);
      toast.add({ title: 'Vượt tồn kho', description: msg, type: 'error' });
      return;
    }

    createDocument(
      {
        document_code: form.code.trim(),
        type: VOUCHER_TYPE_MAP[form.type],
        note: form.note.trim() || null,
        reason: form.reason || null,
        items: lines.map((l) => ({
          product_id: l.productId,
          // Gửi số lượng đã quy đổi ra đơn vị nhỏ nhất (khớp tồn kho backend),
          // kèm đơn giá đã quy về 1 đơn vị nhỏ nhất để total = quantity × unit_price đúng
          quantity: l.convertedQty,
          unit_price: l.pricePerNet,
        })),
      },
      {
        onSuccess: (doc) => {
          toast.add({
            title: 'Tạo phiếu thành công',
            description: `Đã lưu ${doc.document_code} với ${lines.length} dòng hàng`,
            type: 'success',
          });
          router.push('/app/inventory');
        },
        onError: (err) => {
          const msg = getApiErrorMessage(err);
          setSubmitError(msg);
          toast.add({ title: 'Tạo phiếu thất bại', description: msg, type: 'error' });
        },
      },
    );
  }, [form, lines, overStockLines, createDocument, router]);

  return (
    <div className="space-y-4">
      <VoucherHeading type={form.type} onSaveDraft={handleSaveDraft} onComplete={handleComplete} />
      {submitError && (
        <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-700">
          {submitError}
        </div>
      )}
      <VoucherTypeBanner type={form.type} />

      <div className="grid grid-cols-1 items-start gap-6 lg:grid-cols-12">
        <div className="space-y-6 lg:col-span-8">
          <VoucherInfoCard form={form} onChange={patch} />
          <VoucherItemsCard
            lines={lines}
            voucherType={form.type}
            products={products}
            isProductsLoading={isProductsLoading}
            query={productQuery}
            onQueryChange={setProductQuery}
            categories={productCategories}
            categoryFilter={categoryFilter}
            onCategoryFilterChange={setCategoryFilter}
            onUpdateItem={handleUpdateItem}
            onRemoveItem={handleRemoveItem}
            onAddItem={handleAddItem}
            onAddItems={handleAddItems}
            onPriceChange={handlePriceChange}
          />

          <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-surface p-4 text-xs shadow-sm">
            <div className="flex items-center gap-2">
              <Info className="h-5 w-5 text-primary" />
              <span className="text-on-surface-variant">
                Lưu ý: Sau khi hoàn tất phiếu, số lượng tồn thực tế của các chi nhánh sẽ cập nhật
                ngay lập tức.
              </span>
            </div>
            <div className="flex items-center gap-2">
              <button
                type="button"
                className="flex items-center gap-1 rounded bg-surface-container px-2.5 py-1.5 font-medium text-on-surface transition-colors hover:bg-outline-variant"
              >
                <Printer className="h-[15px] w-[15px]" /> In phiếu kho
              </button>
              <button
                type="button"
                className="flex items-center gap-1 rounded bg-surface-container px-2.5 py-1.5 font-medium text-on-surface transition-colors hover:bg-outline-variant"
              >
                <Download className="h-[15px] w-[15px]" /> Xuất Excel
              </button>
            </div>
          </div>
        </div>

        <div className="space-y-6 lg:col-span-4">
          <VoucherFinanceCard
            totals={totals}
            shippingFee={form.shippingFee}
            onShippingFeeChange={(fee) => patch({ shippingFee: fee })}
          />
          {form.type === 'inbound' && <VoucherPaymentCard form={form} onChange={patch} />}
          <VoucherAutomationCard
            form={form}
            onChange={patch}
            onComplete={handleComplete}
            isSubmitting={isSubmitting}
            submitLabel={VOUCHER_META[form.type].submitLabel}
          />
        </div>
      </div>
    </div>
  );
}
