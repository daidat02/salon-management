"use client"

import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "cn"

const tagVariants = cva(
  "inline-flex items-center gap-1 font-medium border whitespace-nowrap",
  {
    variants: {
      variant: {
        default: "bg-slate-100 text-slate-700 border-slate-200",
        primary: "bg-primary-container text-primary border-transparent",
        blue: "bg-blue-50 text-blue-700 border-blue-100",
        purple: "bg-purple-50 text-purple-700 border-purple-100",
        amber: "bg-amber-50 text-amber-800 border-amber-200",
        teal: "bg-teal-50 text-teal-700 border-teal-100",
        indigo: "bg-indigo-50 text-indigo-700 border-indigo-100",
        emerald: "bg-emerald-50 text-emerald-800 border-emerald-200",
        success: "bg-emerald-50 text-success border-emerald-200",
        warning: "bg-amber-100 text-amber-800 border-amber-200",
        error: "bg-red-100 text-error border-red-200",
        danger: "bg-red-100 text-error border-transparent",
        gray: "bg-gray-200 text-gray-700 border-transparent",
        slate: "bg-slate-100 text-slate-700 border-slate-200",
        neutral: "bg-gray-200 text-gray-700 border-transparent",
      },
      size: {
        sm: "px-2 py-0.5 text-[11px]",
        md: "px-2.5 py-0.5 text-[11px]",
        lg: "px-2.5 py-0.5 text-xs",
      },
      shape: {
        pill: "rounded-full",
        rounded: "rounded",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "md",
      shape: "rounded",
    },
  }
)

export interface TagProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof tagVariants> {
  dot?: boolean
}

function Tag({ className, variant, size, shape, dot, children, ...props }: TagProps) {
  return (
    <span className={cn(tagVariants({ variant, size, shape }), className)} {...props}>
      {dot && <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-current" />}
      {children}
    </span>
  )
}

export { Tag, tagVariants }
