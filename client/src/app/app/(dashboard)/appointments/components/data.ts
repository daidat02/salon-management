export type AppointmentItem = {
  id: string
  time: string
  customer: string
  phone: string
  service: string
  duration: string
  staff: string
  status: "pending" | "confirmed" | "serving" | "done" | "cancelled"
}

export const appointmentData: AppointmentItem[] = [
  { id: "LH-001", time: "09:00 24/10", customer: "Chị Thu Thảo", phone: "0903 456 789", service: "Nhuộm Balayage tông khói", duration: "120p • Olaplex", staff: "Tuấn Stylist", status: "serving" },
  { id: "LH-002", time: "10:00 24/10", customer: "Chị Kim Oanh", phone: "0912 345 222", service: "Gội dưỡng sinh thảo dược", duration: "60p", staff: "Ngọc Spa", status: "confirmed" },
  { id: "LH-003", time: "11:00 24/10", customer: "Anh Minh Hoàng", phone: "0988 111 333", service: "Cắt nam + Uốn phồng", duration: "90p", staff: "Linh Barber", status: "pending" },
  { id: "LH-004", time: "13:30 24/10", customer: "Chị Lan Hương", phone: "0901 222 333", service: "Phục hồi Keratin", duration: "90p", staff: "Hương Stylist", status: "done" },
]
