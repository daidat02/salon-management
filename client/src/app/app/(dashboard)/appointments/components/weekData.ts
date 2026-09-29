// Mock data cho lịch tuần — bám preview/admin_preview (content lich-hen)
// + reuse nội dung mock của appointments/components/data.ts cho cột "hôm nay".
// KHÔNG kết nối API.

export type WeekEventStatus = "done" | "serving" | "pending" | "confirmed";

export type WeekDay = {
  weekday: string;
  date: number;
  countLabel: string;
  isToday?: boolean;
  isWeekend?: boolean;
};

export type WeekEvent = {
  id: string;
  /** 0 = Thứ 2 ... 6 = Chủ nhật */
  day: number;
  /** giờ bắt đầu (khớp 1 dòng giờ trong lưới) */
  startHour: number;
  /** chiếm bao nhiêu dòng giờ (mặc định 1) */
  span?: number;
  customer: string;
  vip?: boolean;
  service: string;
  timeRange: string;
  staff: string;
  status: WeekEventStatus;
  /** card nổi bật toàn chiều rộng (kiểu hero "Đang phục vụ" trong preview) */
  featured?: boolean;
};

export const WEEK_HOURS = [8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18];

export const WEEK_DAYS: WeekDay[] = [
  { weekday: "Thứ 2", date: 21, countLabel: "9 lịch" },
  { weekday: "Thứ 3", date: 22, countLabel: "11 lịch" },
  { weekday: "Thứ 4", date: 23, countLabel: "10 lịch" },
  { weekday: "Thứ 5", date: 24, countLabel: "Hôm nay • 14 lịch", isToday: true },
  { weekday: "Thứ 6", date: 25, countLabel: "13 lịch" },
  { weekday: "Thứ 7", date: 26, countLabel: "18 lịch (Đông)", isWeekend: true },
  { weekday: "Chủ nhật", date: 27, countLabel: "13 lịch (Kín)", isWeekend: true },
];

export const WEEK_EVENTS: WeekEvent[] = [
  // ---- Thứ 5 (hôm nay): reuse appointmentData ----
  { id: "LH-001", day: 3, startHour: 8, span: 2, customer: "Chị Thu Thảo", vip: true, service: "Nhuộm Balayage tông khói", timeRange: "08:30 - 10:30", staff: "Tuấn Stylist", status: "serving", featured: true },
  { id: "LH-002", day: 3, startHour: 10, customer: "Chị Kim Oanh", service: "Gội dưỡng sinh thảo dược", timeRange: "10:00 - 11:00", staff: "Ngọc Spa", status: "pending" },
  { id: "LH-003", day: 3, startHour: 11, customer: "Anh Minh Hoàng", service: "Cắt nam + Uốn phồng", timeRange: "11:00 - 12:30", staff: "Linh Barber", status: "pending" },
  { id: "LH-004", day: 3, startHour: 13, span: 2, customer: "Chị Lan Hương", service: "Phục hồi Keratin", timeRange: "13:30 - 15:00", staff: "Hương Stylist", status: "confirmed" },
  { id: "LH-005", day: 3, startHour: 15, customer: "Chị Hà My", vip: true, service: "Chăm sóc Kérastase phục hồi", timeRange: "15:45 - 17:00", staff: "Hương Stylist", status: "confirmed" },
  { id: "LH-006", day: 3, startHour: 17, customer: "Chị Bích Ngọc", service: "Cắt uốn setting sóng lơi", timeRange: "17:15 - 18:15", staff: "Hương Stylist", status: "confirmed" },
  // ---- Các ngày còn lại: lấy mẫu từ preview ----
  { id: "LH-MON01", day: 0, startHour: 9, span: 2, customer: "Chị Thu Thủy", service: "Nhuộm phủ bóng Nano", timeRange: "09:00 - 11:00", staff: "Hương Stylist", status: "done" },
  { id: "LH-MON02", day: 0, startHour: 16, customer: "Anh Việt Anh", service: "Cắt tạo kiểu + Gội massage", timeRange: "16:00 - 17:00", staff: "Linh Barber", status: "done" },
  { id: "LH-TUE01", day: 1, startHour: 8, customer: "Anh Hoàng Nam", service: "Cắt Barber + Gội", timeRange: "08:00 - 09:00", staff: "Linh Barber", status: "done" },
  { id: "LH-TUE02", day: 1, startHour: 14, span: 2, customer: "Chị Kim Tuyến", service: "Nhuộm nâu chocolate", timeRange: "14:00 - 16:00", staff: "Tuấn Stylist", status: "done" },
  { id: "LH-WED01", day: 2, startHour: 9, span: 2, customer: "Chị Quỳnh Anh", service: "Uốn sóng lơi Hàn Quốc", timeRange: "09:15 - 11:30", staff: "Hương Stylist", status: "done" },
  { id: "LH-WED02", day: 2, startHour: 15, customer: "Anh Hoàng Bách", service: "Cắt tóc nam Classic", timeRange: "11:00 - 11:45", staff: "Linh Barber", status: "confirmed" },
  { id: "LH-FRI01", day: 4, startHour: 9, customer: "Chị Bảo Châu", service: "Nối mi Volume thiết kế", timeRange: "09:00 - 10:30", staff: "Trang Eyelash", status: "confirmed" },
  { id: "LH-FRI02", day: 4, startHour: 13, customer: "Chị Thùy Dung", service: "Cắt Layer nữ + Sấy kiểu", timeRange: "13:00 - 14:00", staff: "Hương Stylist", status: "confirmed" },
  { id: "LH-SAT01", day: 5, startHour: 8, customer: "Chị Thanh Vân", service: "Phục hồi Kérastase", timeRange: "08:00 - 09:30", staff: "Hương Stylist", status: "confirmed" },
  { id: "LH-SAT02", day: 5, startHour: 10, span: 2, customer: "Anh Minh Đức", service: "Cắt Fade + Uốn con sâu", timeRange: "10:00 - 11:15", staff: "Linh Barber", status: "pending" },
  { id: "LH-SUN01", day: 6, startHour: 9, span: 2, customer: "Chị Hồng Hạnh", service: "Gội dưỡng sinh 90p", timeRange: "09:30 - 11:00", staff: "Ngọc Spa", status: "confirmed" },
  { id: "LH-SUN02", day: 6, startHour: 16, customer: "Chị Cẩm Tú", service: "Nối mi Volume thiết kế", timeRange: "16:00 - 17:30", staff: "Trang Eyelash", status: "pending" },
];
