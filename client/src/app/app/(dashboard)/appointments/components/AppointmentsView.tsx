"use client";

import * as React from "react";
import { CalendarRange, List } from "lucide-react";
import { cn } from "cn";
import AppointmentsTable from "./AppointmentsTable";
import WeekSchedule from "./WeekSchedule";

type View = "list" | "week";

export default function AppointmentsView() {
  const [view, setView] = React.useState<View>("list");

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3 text-xs">
        <span className="rounded-lg bg-surface-container px-3 py-1.5">Hôm nay, 24/10/2024</span>
        <div className="flex items-center gap-1 rounded-lg bg-surface-container p-1">
          <button
            type="button"
            onClick={() => setView("list")}
            className={cn(
              "flex items-center gap-1.5 rounded-md px-3 py-1.5 font-medium transition-all",
              view === "list"
                ? "bg-primary font-semibold text-white shadow-xs"
                : "text-on-surface-variant hover:text-on-surface",
            )}
          >
            <List className="h-3.5 w-3.5" />
            <span>Danh sách</span>
          </button>
          <button
            type="button"
            onClick={() => setView("week")}
            className={cn(
              "flex items-center gap-1.5 rounded-md px-3 py-1.5 font-medium transition-all",
              view === "week"
                ? "bg-primary font-semibold text-white shadow-xs"
                : "text-on-surface-variant hover:text-on-surface",
            )}
          >
            <CalendarRange className="h-3.5 w-3.5" />
            <span>Lịch tuần</span>
          </button>
        </div>
      </div>

      {view === "list" ? <AppointmentsTable /> : <WeekSchedule />}
    </div>
  );
}
