'use client';

import { AlarmClock, Cake, Heart, Mail, MapPin, Phone, Star, TriangleAlert } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { DashDrawer } from '@/components/dashboard/DashDrawer';

const mockCustomer = {
  code: '#KH-0892',
  name: 'Chị Thu Thảo',
  rank: 'Hạng Vàng',
  since: 'Khách hàng thân thiết từ 12/2022',
  avatar:
    'https://lh3.googleusercontent.com/aida/AEtjO1UPHCFbE8KRvrKc0bZdgiuOHrcfFmFE_7MY6ek6RnkH6JT_FnDvpVR2WoIGJb7J0lUBqkd2B5_e0Ev8EXaNcrXyeVOhnSFRj7h_h2uKzjHh1NUMnXCBtJjB8r8uCxxbBZbFJSq8cK6dIlwZC1Xq0g2OX5bu5OItqWX9uUmrRlXSAac88qHKqCLRBB04OGc3SLmslBduhPw_TzoY0AiAcNOQR7AQd-X8BdsSbpkgn-Lc-PiTLDtjKib5T-NC',
  phone: '0903.456.789',
  email: 'thuthao.design@gmail.com',
  birth: '18/08/1994 (30 tuổi)',
  address: 'Phường 6, Quận 3, TP. Hồ Chí Minh',
  stats: { visits: 14, spent: '24.8M', points: '3.420', cancelled: 0 },
  spending: [
    { month: 'T5', h: 'h-10', value: 10 },
    { month: 'T6', h: 'h-12', value: 12 },
    { month: 'T7', h: 'h-8', value: 8 },
    { month: 'T8', h: 'h-14', value: 14 },
    { month: 'T9', h: 'h-11', value: 11 },
    { month: 'T10', h: 'h-16', value: 16, active: true },
  ],
  upcoming: {
    time: '09:30 Hôm nay',
    service: 'Nhuộm Balayage tông khói cao cấp',
    stylist: 'Tuấn (Master Stylist)',
    seat: 'Ghế VIP 02',
  },
  orders: [
    { code: '#DH-1049', service: 'Phục hồi tóc Olaplex', total: '1.850.000 đ' },
    { code: '#DH-0982', service: 'Cắt tóc Layer + Hấp dầu', total: '950.000 đ' },
    { code: '#DH-0844', service: 'Uốn Setting sóng tự nhiên', total: '2.400.000 đ' },
  ],
};

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function CustomerDrawer({ open, onOpenChange }: Props) {
  return (
    <DashDrawer
      open={open}
      onOpenChange={onOpenChange}
      title="Chi tiết khách hàng"
      subtitle={mockCustomer.since}
      badge={
        <span className="text-xs font-mono font-semibold text-primary bg-primary-container/60 px-2 py-0.5 rounded">
          {mockCustomer.code}
        </span>
      }
      footer={
        <>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            className="border-outline text-slate-600"
          >
            Đóng
          </Button>
          <Button className="bg-primary hover:bg-[#436b95] text-on-primary">Chỉnh sửa hồ sơ</Button>
        </>
      }
      widthClass="max-w-[580px] w-[580px]"
    >
      {/* Profile Card - code.html:752-781 */}
      <div className="flex items-start gap-4 p-4 rounded-xl bg-surface-container-low border border-outline">
        <div className="relative shrink-0">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            alt={mockCustomer.name}
            src={mockCustomer.avatar}
            className="w-16 h-16 rounded-xl object-cover ring-2 ring-primary/40 shadow"
          />
          <span className="absolute -bottom-1 -right-1 bg-amber-400 text-slate-900 text-[10px] font-bold px-1.5 py-0.2 rounded-full border border-white">
            VIP
          </span>
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between">
            <h4 className="text-base font-bold text-slate-900 truncate">{mockCustomer.name}</h4>
            <span className="text-[11px] font-semibold text-amber-800 bg-amber-100 px-2 py-0.5 rounded-full border border-amber-200">
              Hạng Vàng
            </span>
          </div>
          <div className="mt-2 space-y-1 text-xs text-on-surface-variant">
            <p className="flex items-center gap-1.5">
              <Phone className="h-3.5 w-3.5 text-slate-500" />
              <span className="font-mono text-slate-800 font-medium">{mockCustomer.phone}</span>
            </p>
            <p className="flex items-center gap-1.5 truncate">
              <Mail className="h-3.5 w-3.5 text-slate-500" />
              <span>{mockCustomer.email}</span>
            </p>
            <p className="flex items-center gap-1.5">
              <Cake className="h-3.5 w-3.5 text-slate-500" />
              <span>{mockCustomer.birth}</span>
            </p>
            <p className="flex items-center gap-1.5 truncate">
              <MapPin className="h-3.5 w-3.5 text-slate-500" />
              <span>{mockCustomer.address}</span>
            </p>
          </div>
        </div>
      </div>

      {/* Stats Matrix */}
      <div className="grid grid-cols-4 gap-2 text-center">
        <div className="p-2.5 rounded-lg bg-surface border border-outline">
          <p className="text-[10px] text-on-surface-variant uppercase font-medium">Lần ghé</p>
          <p className="text-base font-bold text-primary mt-0.5">{mockCustomer.stats.visits}</p>
        </div>
        <div className="p-2.5 rounded-lg bg-surface border border-outline">
          <p className="text-[10px] text-on-surface-variant uppercase font-medium">Tổng chi</p>
          <p className="text-xs font-bold text-slate-900 mt-1">{mockCustomer.stats.spent}</p>
        </div>
        <div className="p-2.5 rounded-lg bg-surface border border-outline">
          <p className="text-[10px] text-on-surface-variant uppercase font-medium">Tích điểm</p>
          <p className="text-base font-bold text-amber-600 mt-0.5">{mockCustomer.stats.points}</p>
        </div>
        <div className="p-2.5 rounded-lg bg-surface border border-outline">
          <p className="text-[10px] text-on-surface-variant uppercase font-medium">Hủy lịch</p>
          <p className="text-base font-bold text-emerald-600 mt-0.5">
            {mockCustomer.stats.cancelled}
          </p>
        </div>
      </div>

      {/* Spending Trend */}
      <div className="bg-surface rounded-xl p-4 border border-outline">
        <div className="flex items-center justify-between mb-3">
          <span className="text-xs font-semibold text-slate-800">
            Mức chi tiêu 6 tháng gần nhất
          </span>
          <span className="text-[11px] text-primary font-medium">Trung bình: 4.1M/tháng</span>
        </div>
        <div className="h-20 flex items-end justify-between gap-3 pt-2">
          {mockCustomer.spending.map((m) => (
            <div key={m.month} className="flex-1 flex flex-col items-center gap-1">
              <div
                className={`w-full rounded-t ${m.active ? 'bg-primary shadow-xs' : 'bg-primary-container'} ${m.h}`}
              />
              <span
                className={`text-[10px] ${m.active ? 'font-bold text-primary' : 'text-on-surface-variant'}`}
              >
                {m.month}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* Upcoming */}
      <div className="space-y-3">
        <h5 className="text-xs font-semibold text-slate-800 uppercase tracking-wider">
          Lịch hẹn sắp tới
        </h5>
        <div className="p-3.5 rounded-xl bg-blue-50/70 border border-blue-200">
          <div className="flex items-center justify-between mb-1.5">
            <span className="inline-flex items-center gap-1 text-xs font-bold text-primary">
              <AlarmClock className="h-3.5 w-3.5" />
              {mockCustomer.upcoming.time}
            </span>
            <span className="text-[10px] font-semibold bg-primary text-white px-2 py-0.5 rounded-full">
              Sắp diễn ra
            </span>
          </div>
          <p className="text-xs font-bold text-slate-900">{mockCustomer.upcoming.service}</p>
          <p className="text-[11px] text-on-surface-variant mt-0.5">
            Stylist thực hiện:{' '}
            <span className="font-medium text-slate-700">{mockCustomer.upcoming.stylist}</span> •{' '}
            {mockCustomer.upcoming.seat}
          </p>
        </div>
      </div>

      {/* Notes */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <h5 className="text-xs font-semibold text-slate-800 uppercase tracking-wider">
            Ghi chú phục vụ & Công thức tóc
          </h5>
          <button className="text-xs text-primary font-medium hover:underline">+ Cập nhật</button>
        </div>
        <div className="p-3.5 rounded-xl bg-amber-50/60 border border-amber-200 text-xs space-y-2 text-slate-800">
          <p className="flex items-start gap-2">
            <TriangleAlert className="text-amber-600 h-3.5 w-3.5 mt-0.5 shrink-0" />
            <span>
              <strong>Lưu ý:</strong> Dị ứng thuốc nhuộm nồng độ amoniac cao. Chỉ dùng dòng Organic
              dưỡng ẩm.
            </span>
          </p>
          <p className="flex items-start gap-2">
            <Heart className="text-amber-600 h-3.5 w-3.5 mt-0.5 shrink-0" />
            <span>
              <strong>Sở thích:</strong> Thích massage da đầu nhẹ nhàng với tinh dầu Argan
              Kerastase; uống trà thảo mộc ít đường.
            </span>
          </p>
        </div>
      </div>

      {/* Orders */}
      <div className="space-y-3">
        <h5 className="text-xs font-semibold text-slate-800 uppercase tracking-wider">
          Lịch sử đơn hàng gần đây
        </h5>
        <div className="border border-outline rounded-xl overflow-hidden text-xs">
          <table className="w-full">
            <thead className="bg-surface-container-low text-on-surface-variant font-medium">
              <tr>
                <th className="py-2 px-3 text-left">Mã đơn</th>
                <th className="py-2 px-3 text-left">Dịch vụ</th>
                <th className="py-2 px-3 text-right">Tổng tiền</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-outline">
              {mockCustomer.orders.map((o) => (
                <tr key={o.code}>
                  <td className="py-2 px-3 font-mono font-medium text-primary">{o.code}</td>
                  <td className="py-2 px-3 text-slate-700">{o.service}</td>
                  <td className="py-2 px-3 text-right font-bold text-slate-900">{o.total}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </DashDrawer>
  );
}
