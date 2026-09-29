"use client"

import * as React from "react"
import { Toast as ToastPrimitive } from "@base-ui/react/toast"
import { cn } from "cn"

import { Button } from "@/components/ui/button"
import { XIcon, CircleCheckIcon, InfoIcon, TriangleAlertIcon, OctagonXIcon, Loader2Icon } from "lucide-react"

const toast = ToastPrimitive.createToastManager()

function ToastProvider({ ...props }: ToastPrimitive.Provider.Props) {
  return <ToastPrimitive.Provider {...props} />
}

function ToastPortal({ ...props }: ToastPrimitive.Portal.Props) {
  return <ToastPrimitive.Portal data-slot="toast-portal" {...props} />
}

function ToastViewport({ className, ...props }: ToastPrimitive.Viewport.Props) {
  return (
    <ToastPrimitive.Viewport
      data-slot="toast-viewport"
      className={cn(
        "pointer-events-none fixed left-1/2 top-6 z-50 flex w-full max-w-sm -translate-x-1/2 flex-col items-center gap-2 outline-none",
        className
      )}
      {...props}
    />
  )
}

function Toast({ className, ...props }: ToastPrimitive.Root.Props) {
  return (
    <ToastPrimitive.Root
      data-slot="toast"
      className={cn(
        "group/toast pointer-events-auto relative flex w-full items-center rounded-full border shadow-xl will-change-transform outline-none select-none",
        // stack centered
        "data-[swipe-direction=up]:animate-out data-[swipe-direction=up]:fade-out-0 data-[swipe-direction=up]:slide-out-to-top-2",
        "data-[swipe-direction=down]:animate-out data-[swipe-direction=down]:fade-out-0",
        className
      )}
      {...props}
    />
  )
}

function ToastContent({ className, ...props }: ToastPrimitive.Content.Props) {
  return (
    <ToastPrimitive.Content
      data-slot="toast-content"
      className={cn(
        "flex w-full items-center gap-1.5 rounded-full px-2.5 py-1",
        className
      )}
      {...props}
    />
  )
}

function ToastTitle({ className, ...props }: ToastPrimitive.Title.Props) {
  return (
    <ToastPrimitive.Title
      data-slot="toast-title"
      className={cn("text-xs font-semibold", className)}
      {...props}
    />
  )
}

function ToastDescription({
  className,
  ...props
}: ToastPrimitive.Description.Props) {
  return (
    <ToastPrimitive.Description
      data-slot="toast-description"
      className={cn("text-[11px] leading-none opacity-80", className)}
      {...props}
    />
  )
}

function ToastAction({
  className,
  render = <Button variant="outline" size="sm" />,
  ...props
}: ToastPrimitive.Action.Props) {
  return (
    <ToastPrimitive.Action
      data-slot="toast-action"
      render={render}
      className={cn("shrink-0", className)}
      {...props}
    />
  )
}

function ToastClose({
  className,
  children,
  render,
  ...props
}: ToastPrimitive.Close.Props) {
  const closeRender =
    render ?? (
      <button className="flex h-4 w-4 shrink-0 items-center justify-center rounded-full text-current/50 hover:text-current hover:bg-transparent">
        <XIcon className="h-2.5 w-2.5" aria-hidden="true" />
      </button>
    )
  return (
    <ToastPrimitive.Close
      data-slot="toast-close"
      aria-label="Close toast"
      render={closeRender}
      className={cn("shrink-0", className)}
      {...props}
    >
      {children}
    </ToastPrimitive.Close>
  )
}

function ToastIcon({ type }: { type: string | undefined }) {
  const base = "flex h-4 w-4 shrink-0 items-center justify-center rounded-full [&_svg]:size-3"
  if (type === "success") {
    return (
      <span data-slot="toast-icon" className={cn(base, "bg-emerald-100 text-emerald-600")}>
        <CircleCheckIcon aria-hidden="true" />
      </span>
    )
  }
  if (type === "info") {
    return (
      <span data-slot="toast-icon" className={cn(base, "bg-blue-100 text-blue-600")}>
        <InfoIcon aria-hidden="true" />
      </span>
    )
  }
  if (type === "warning") {
    return (
      <span data-slot="toast-icon" className={cn(base, "bg-amber-100 text-amber-600")}>
        <TriangleAlertIcon aria-hidden="true" />
      </span>
    )
  }
  if (type === "error") {
    return (
      <span data-slot="toast-icon" className={cn(base, "bg-red-100 text-red-600")}>
        <OctagonXIcon aria-hidden="true" />
      </span>
    )
  }
  if (type === "loading") {
    return (
      <span data-slot="toast-icon" className={cn(base, "bg-slate-100 text-slate-500")}>
        <Loader2Icon className="animate-spin" aria-hidden="true" />
      </span>
    )
  }
  return null
}

function getToastVariant(type: string | undefined) {
  switch (type) {
    case "success":
      return "bg-emerald-50 border-emerald-200 text-emerald-900"
    case "error":
      return "bg-red-50 border-red-200 text-red-900"
    case "warning":
      return "bg-amber-50 border-amber-200 text-amber-900"
    case "info":
      return "bg-blue-50 border-blue-200 text-blue-900"
    default:
      return "bg-white border-slate-200 text-slate-900"
  }
}

function ToastList() {
  const { toasts } = ToastPrimitive.useToastManager()

  // Chỉ hiện toast success (theo yêu cầu: chỉ cần message thành công)
  const visibleToasts = toasts.filter((t) => !t.type || t.type === "success")

  return visibleToasts.map((toastItem) => (
    <Toast
      key={toastItem.id}
      toast={toastItem}
      className={cn("min-w-[240px] max-w-[320px] justify-between py-0", getToastVariant(toastItem.type))}
    >
      <ToastContent className="flex-1 py-1">
        <ToastIcon type={toastItem.type} />
        {/* Chỉ một dòng message */}
        <ToastTitle className="flex-1 truncate text-[10px] font-medium leading-none" />
        <ToastClose className="ml-1" />
      </ToastContent>
    </Toast>
  ))
}

function Toaster({
  children,
  toastManager = toast,
  ...props
}: ToastPrimitive.Provider.Props) {
  return (
    <ToastProvider toastManager={toastManager} {...props}>
      {children}
      <ToastPortal>
        <ToastViewport>
          <ToastList />
        </ToastViewport>
      </ToastPortal>
    </ToastProvider>
  )
}

const createToastManager = ToastPrimitive.createToastManager
const useToastManager = ToastPrimitive.useToastManager

export {
  Toaster,
  Toast,
  ToastAction,
  ToastClose,
  ToastContent,
  ToastDescription,
  ToastPortal,
  ToastProvider,
  ToastTitle,
  ToastViewport,
  createToastManager,
  toast,
  useToastManager,
}
