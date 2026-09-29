import type { LucideIcon } from "lucide-react";
import { cn } from "cn";

type StatCardProps = {
  label: string;
  value: React.ReactNode;
  sub?: React.ReactNode;
  icon: LucideIcon;
  /** Màu nền icon — mặc định primary tint */
  iconWrapClassName?: string;
  iconClassName?: string;
  className?: string;
};

export default function StatCard({
  label,
  value,
  sub,
  icon: Icon,
  iconWrapClassName = "bg-primary-container text-primary",
  iconClassName,
  className,
}: StatCardProps) {
  return (
    <div
      className={cn(
        "flex flex-col justify-between rounded-xl border border-outline bg-surface px-4 py-3.5 shadow-sm hover:border-primary/40 transition-colors",
        className
      )}
    >
      <div className="flex items-center justify-between gap-3">
        <span className="text-sm font-semibold leading-none text-on-surface-variant">{label}</span>
        <div className={cn("flex h-7 w-7 shrink-0 items-center justify-center rounded-lg", iconWrapClassName)}>
          <Icon className={cn("h-4 w-4", iconClassName)} />
        </div>
      </div>
      <div className="mt-2">
        <div className="text-xl font-bold leading-none tracking-tight text-on-surface">{value}</div>
        {sub && <div className="mt-1.5 text-[11px] leading-none text-on-surface-variant">{sub}</div>}
      </div>
    </div>
  );
}
