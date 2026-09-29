"use client"

import { Minus, Plus, User, X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Tag } from "@/components/ui/tag"
import type { CatalogItem } from "./ServiceCatalog"

export type InvoiceItem = CatalogItem & { quantity: number; dotColor: string }

type Props = {
  items: InvoiceItem[]
  onQuantity: (id: string, delta: number) => void
  onRemove: (id: string) => void
  onClear: () => void
  onCheckout: () => void
}

function formatPrice(n: number) {
  return n.toLocaleString("vi-VN") + " đ"
}

export default function InvoicePanel({ items, onQuantity, onRemove, onClear, onCheckout }: Props) {
  const customer = {
    name: "Chị Thu Thảo",
    points: "3.420 pts",
    vip: true,
  } as const

  const subtotal = items.reduce((sum, it) => sum + it.price * it.quantity, 0)

  return (
    <section className="bg-surface border-outline flex w-full shrink-0 flex-col overflow-hidden rounded-md border shadow-sm lg:sticky lg:top-4 lg:h-[calc(100dvh-64px)] lg:max-h-[calc(100dvh-64px)] lg:w-[440px] xl:w-[480px]">
      {/* Header + Khách hàng gộp ngang hàng với mã đơn - chỉ tên + điểm */}
      <div className="bg-surface-container-low border-outline flex shrink-0 flex-col gap-1 border-b px-2.5 py-2">
        <div className="flex items-center justify-between gap-2">
          <div className="flex flex-1 items-center gap-2 min-w-0">
            <div className="flex items-center gap-1.5 shrink-0">
              <h2 className="text-[13px] font-bold leading-none text-on-surface">Hóa đơn #INV-2026-001</h2>
              <Tag variant="warning" size="sm" shape="rounded" className="px-1.5 py-0 text-[10px] leading-none">
                Đang tạo
              </Tag>
            </div>
            <div className="hidden sm:flex items-center gap-1.5 border-l border-outline pl-2 ml-1">
              <User className="text-primary h-3 w-3 shrink-0" />
              <span className="text-[11px] font-bold leading-none text-on-surface truncate">{customer.name}</span>
              {customer.vip && <span className="rounded bg-amber-500 px-1 py-0 text-[9px] font-bold leading-none text-white">VIP</span>}
              <span className="text-on-surface-variant text-[11px] leading-none whitespace-nowrap">• {customer.points}</span>
            </div>
          </div>
          <div className="flex items-center gap-1 shrink-0">
            {/* Khi có khách thì ẩn nút Khách mới - hiện đã có khách nên ẩn */}
          </div>
        </div>
        <div className="flex sm:hidden items-center gap-1.5">
          <User className="text-primary h-3 w-3" />
          <span className="text-[11px] font-bold text-on-surface">{customer.name}</span>
          {customer.vip && <span className="rounded bg-amber-500 px-1 py-0 text-[9px] font-bold text-white">VIP</span>}
          <span className="text-on-surface-variant text-[11px]">• {customer.points}</span>
        </div>
        <p className="text-on-surface-variant text-[10px] leading-none">Hôm nay, 24/10/2026 - Máy POS 01</p>
      </div>

      {/* Items - chiếm hết không gian còn lại */}
      <div className="custom-scroll flex-1 min-h-[180px] overflow-y-auto px-2.5 py-2">
        {items.length === 0 ? (
          <div className="py-10 text-center text-xs text-on-surface-variant">Chưa có sản phẩm — chọn bên trái để thêm</div>
        ) : (
          <div className="divide-y divide-slate-100">
            {items.map((it) => (
              <div key={it.id} className="flex flex-col gap-1 py-2 first:pt-0 last:pb-0">
                <div className="flex items-start justify-between">
                  <div className="flex-1 pr-2">
                    <div className="flex items-center gap-1">
                      <span className={`h-1.5 w-1.5 shrink-0 rounded-full ${it.dotColor}`} />
                      <h4 className="text-[11px] font-bold leading-none text-on-surface line-clamp-1">{it.name}</h4>
                    </div>
                    <p className="text-on-surface-variant pl-2.5 text-[10px] leading-none">{it.desc}</p>
                  </div>
                  <button onClick={() => onRemove(it.id)} className="shrink-0 p-0.5 text-slate-400 transition-colors hover:text-error">
                    <X className="h-3.5 w-3.5" />
                  </button>
                </div>
                <div className="flex items-center justify-between pl-2.5">
                  <span className="text-primary text-[10px] font-semibold">{formatPrice(it.price)}</span>
                  <div className="flex items-center gap-1">
                    <div className="border-outline flex items-center rounded border bg-white">
                      <button onClick={() => onQuantity(it.id, -1)} className="hover:bg-slate-100 flex h-5 w-5 items-center justify-center rounded-l text-on-surface-variant">
                        <Minus className="h-2.5 w-2.5" />
                      </button>
                      <span className="w-6 text-center text-[11px] font-bold text-on-surface">{it.quantity}</span>
                      <button onClick={() => onQuantity(it.id, 1)} className="hover:bg-slate-100 flex h-5 w-5 items-center justify-center rounded-r text-on-surface-variant">
                        <Plus className="h-2.5 w-2.5" />
                      </button>
                    </div>
                    <span className="w-20 text-right text-[11px] font-bold text-on-surface">{formatPrice(it.price * it.quantity)}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Summary - chỉ giá tạm tính */}
      <div className="bg-surface-container-low border-outline flex shrink-0 flex-col gap-2 border-t px-3 py-3">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium text-on-surface">Tạm tính ({items.length} mục)</span>
          <span className="text-primary font-mono text-lg font-extrabold tracking-tight">{formatPrice(subtotal)}</span>
        </div>
        <p className="text-on-surface-variant text-[11px] leading-none">Chưa bao gồm khuyến mãi - thanh toán tại Đơn hàng</p>
      </div>

      {/* Actions - cao hơn, đổi thành Tạo đơn */}
      <div className="bg-surface-container-low border-outline flex shrink-0 gap-2 border-t px-3 py-3">
        <Button
          variant="outline"
          onClick={onClear}
          className="bg-white border-outline text-on-surface-variant hover:text-error hover:border-red-200 hover:bg-red-50 h-10 px-4 text-sm font-medium"
        >
          <X className="h-4 w-4" /> Hủy đơn
        </Button>
        <Button onClick={onCheckout} className="flex h-10 flex-1 items-center justify-center gap-2 text-sm font-bold shadow-sm">
          <Plus className="h-4 w-4" /> Tạo đơn
        </Button>
      </div>
    </section>
  )
}
