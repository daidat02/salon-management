"use client";

import * as React from "react";
import { CalendarClock, Clock3, MapPin } from "lucide-react";
import { cn } from "cn";
import {
  WEEK_DAYS,
  WEEK_EVENTS,
  WEEK_HOURS,
  type WeekEvent,
  type WeekEventStatus,
} from "./weekData";

const ROW_H = 80; // h-20 — 1 dòng giờ

const STATUS_STYLE: Record<WeekEventStatus, { card: string; time: string; badge: string; badgeLabel: string }> = {
  done: {
    card: "bg-emerald-50 hover:bg-emerald-100/90",
    time: "text-emerald-700",
    badge: "bg-white text-emerald-700",
    badgeLabel: "Xong",
  },
  serving: {
    card: "bg-primary hover:bg-primary/90 text-white",
    time: "text-blue-100",
    badge: "bg-white/20 text-white",
    badgeLabel: "Đang làm",
  },
  pending: {
    card: "bg-amber-50 hover:bg-amber-100/90",
    time: "text-amber-700",
    badge: "bg-amber-200/70 text-amber-800",
    badgeLabel: "Chờ",
  },
  confirmed: {
    card: "bg-indigo-50 hover:bg-indigo-100",
    time: "text-indigo-700",
    badge: "bg-indigo-200/70 text-indigo-800",
    badgeLabel: "Đã hẹn",
  },
};

const LEGEND: { dot: string; label: string }[] = [
  { dot: "bg-emerald-500", label: "Đã hoàn thành" },
  { dot: "bg-primary", label: "Đang phục vụ" },
  { dot: "bg-amber-500", label: "Chờ phục vụ" },
  { dot: "bg-indigo-400", label: "Đã xác nhận" },
];

function EventCard({ event, tall }: { event: WeekEvent; tall?: boolean }) {
  const s = STATUS_STYLE[event.status];
  const serving = event.status === "serving";
  return (
    <div
      className={cn(
        "flex h-full cursor-pointer flex-col justify-between rounded-lg p-2 shadow-xs transition-all",
        s.card,
      )}
    >
      <div className="flex min-w-0 items-center justify-between gap-1">
        <span className={cn("min-w-0 flex-1 truncate text-[11px] font-bold", serving ? "text-white" : "text-on-surface")}>
          {event.customer}
        </span>
        <span className="flex shrink-0 items-center gap-1">
          {event.vip && (
            <span className="rounded bg-amber-400 px-1 py-px text-[9px] font-extrabold text-slate-900">
              VIP
            </span>
          )}
          <span className={cn("rounded px-1 py-px text-[9px] font-semibold", s.badge)}>
            {serving && <span className="mr-0.5 inline-block h-1.5 w-1.5 animate-pulse rounded-full bg-emerald-300" />}
            {s.badgeLabel}
          </span>
        </span>
      </div>
      <div className={cn("truncate text-[10px]", serving ? "text-blue-100" : "text-on-surface-variant")}>
        {event.service}
      </div>
      <div className={cn("flex items-center gap-1 text-[9px] font-medium", s.time)}>
        <Clock3 className="h-3 w-3" />
        <span className="truncate font-mono">{event.timeRange}</span>
      </div>
      {tall && (
        <div className={cn("truncate text-[10px] font-semibold", serving ? "text-white" : "text-on-surface")}>
          {event.staff}
        </div>
      )}
    </div>
  );
}

export default function WeekSchedule() {
  const nowRef = React.useRef<HTMLDivElement>(null);

  // Ô (day, hour) đã bị event span phủ thì không render cell trống đè lên
  const covered = React.useMemo(() => {
    const set = new Set<string>();
    for (const e of WEEK_EVENTS) {
      const span = e.span ?? 1;
      for (let h = e.startHour + 1; h < e.startHour + span; h++) set.add(`${e.day}-${h}`);
    }
    return set;
  }, []);

  const eventAt = (day: number, hour: number) =>
    WEEK_EVENTS.find((e) => e.day === day && e.startHour === hour);

  const scrollToNow = () => {
    nowRef.current?.scrollIntoView({ behavior: "smooth", block: "center" });
  };

  return (
    <div className="space-y-3">
      {/* Legend + jump-to-now — bám preview viewLichTuan */}
      <div className="flex flex-wrap items-center justify-between gap-3 px-1">
        <div className="flex flex-wrap items-center gap-4 text-xs">
          <span className="text-[11px] font-semibold tracking-wider text-on-surface-variant uppercase">
            Trạng thái:
          </span>
          {LEGEND.map((l) => (
            <div key={l.label} className="flex items-center gap-1.5">
              <span className={cn("h-2.5 w-2.5 rounded-full", l.dot)} />
              <span className="text-on-surface-variant">{l.label}</span>
            </div>
          ))}
        </div>
        <button
          type="button"
          onClick={scrollToNow}
          className="flex items-center gap-1 text-xs font-semibold text-primary hover:underline"
        >
          <MapPin className="h-3.5 w-3.5" />
          <span>Đến khung giờ hiện tại (10:15)</span>
        </button>
      </div>

      {/* Card lịch tuần — cùng width với bảng (overflow-x-auto bên trong, không tràn) */}
      <div className="overflow-hidden rounded-md border border-outline bg-surface shadow-sm">
        <div className="overflow-x-auto">
          <div className="min-w-[1024px]">
            {/* Header 7 ngày */}
            <div className="grid grid-cols-[64px_repeat(7,minmax(0,1fr))] bg-surface-container select-none">
              <div className="flex flex-col items-center justify-center bg-surface-container p-3">
                <Clock3 className="h-[18px] w-[18px] text-on-surface-variant" />
                <span className="mt-0.5 text-[10px] font-bold tracking-wider text-on-surface-variant uppercase">
                  Giờ
                </span>
              </div>
              {WEEK_DAYS.map((d) => (
                <div
                  key={d.date}
                  className={cn(
                    "relative p-3 text-center",
                    d.isToday && "bg-brand-light/60",
                  )}
                >
                  {d.isToday && (
                    <span className="absolute top-1.5 right-1.5 h-1.5 w-1.5 animate-pulse rounded-full bg-brand" />
                  )}
                  <div
                    className={cn(
                      "text-[11px] font-semibold uppercase",
                      d.isToday ? "font-bold text-brand" : d.isWeekend ? "font-bold text-amber-800" : "text-on-surface-variant",
                    )}
                  >
                    {d.weekday}
                  </div>
                  {d.isToday ? (
                    <div className="mx-auto mt-0.5 flex h-7 w-7 items-center justify-center rounded-full bg-brand text-base font-bold text-white shadow-sm">
                      {d.date}
                    </div>
                  ) : (
                    <div className="mt-0.5 text-base font-bold text-on-surface">{d.date}</div>
                  )}
                  <span
                    className={cn(
                      "mt-1 inline-block rounded-full px-2 py-0.5 text-[10px] font-medium",
                      d.isToday
                        ? "bg-white font-bold text-brand shadow-2xs"
                        : d.isWeekend
                          ? "bg-amber-100 font-semibold text-amber-700"
                          : "bg-surface-container text-on-surface-variant",
                    )}
                  >
                    {d.countLabel}
                  </span>
                </div>
              ))}
            </div>

            {/* Lưới giờ — relative để ghim vạch "hiện tại" */}
            <div className="relative grid grid-cols-[64px_repeat(7,minmax(0,1fr))]">
              {/* Vạch giờ hiện tại: Thứ 5 lúc 10:15 */}
              <div
                ref={nowRef}
                className="pointer-events-none absolute right-0 left-0 z-20 flex items-center"
                style={{ top: 2 * ROW_H + 20 }}
              >
                <div className="w-[64px] pr-1 text-right">
                  <span className="rounded bg-error px-1 py-0.5 font-mono text-[10px] font-bold text-white shadow-sm">
                    10:15
                  </span>
                </div>
                <div className="relative flex flex-1 items-center">
                  <div className="h-0.5 w-full bg-error/80" />
                  <div className="absolute left-[46%] h-3 w-3 rounded-full bg-error shadow ring-4 ring-red-100" />
                </div>
              </div>

              {WEEK_HOURS.map((hour) => (
                <React.Fragment key={hour}>
                  <div className="flex flex-col justify-start bg-surface-container/60 p-2 text-right font-mono text-[11px] font-semibold text-on-surface-variant" style={{ height: ROW_H }}>
                    {String(hour).padStart(2, "0")}:00
                  </div>
                  {WEEK_DAYS.map((d, dayIdx) => {
                    const isToday = !!d.isToday;
                    if (covered.has(`${dayIdx}-${hour}`)) {
                      return (
                        <div
                          key={dayIdx}
                          style={{ height: ROW_H }}
                          className={cn(isToday && "bg-brand-light/20")}
                        />
                      );
                    }
                    const event = eventAt(dayIdx, hour);
                    const span = event?.span ?? 1;
                    return (
                      <div
                        key={dayIdx}
                        style={{ height: ROW_H }}
                        className={cn(
                          "relative p-1 transition-colors",
                          isToday
                            ? "bg-brand-light/20 hover:bg-brand-light/40"
                            : d.isWeekend
                              ? "bg-amber-50/20 hover:bg-amber-50/50"
                              : "bg-surface hover:bg-surface-container-low",
                        )}
                      >
                        {event &&
                          (span > 1 ? (
                            <div
                              className="absolute top-1 right-1 left-1 z-10"
                              style={{ height: span * ROW_H - 8 }}
                            >
                              <EventCard event={event} tall />
                            </div>
                          ) : (
                            <EventCard event={event} />
                          ))}
                      </div>
                    );
                  })}
                </React.Fragment>
              ))}
            </div>
          </div>
        </div>

        {/* Footer mini giống footer bảng */}
        <div className="flex items-center gap-1.5 border-t border-outline bg-surface px-6 py-3 text-xs text-on-surface-variant">
          <CalendarClock className="h-3.5 w-3.5" />
          <span>
            Tuần 21/10 – 27/10 • <span className="font-semibold text-on-surface">98</span> lịch hẹn
          </span>
        </div>
      </div>
    </div>
  );
}
