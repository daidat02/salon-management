import Link from "next/link"
import { ArrowLeft, Home, SearchX, Scissors } from "lucide-react"
import { Button } from "@/components/ui/button"

export default function NotFound() {
  return (
    <div className="theme-app flex min-h-screen flex-col items-center justify-center bg-background px-6 py-12 font-sans text-body antialiased">
      {/* Logo */}
      <Link href="/landing" className="mb-8 flex items-center gap-2.5">
        <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand text-white shadow-sm">
          <Scissors className="h-5 w-5" />
        </div>
        <span className="text-base font-bold tracking-tight text-on-surface">SALON ADMIN</span>
      </Link>

      {/* 404 */}
      <div className="relative flex flex-col items-center text-center">
        <div className="pointer-events-none absolute -top-10 select-none text-[100px] font-extrabold leading-none tracking-tighter text-outline opacity-30 sm:text-[140px]">
          404
        </div>
        <div className="relative mt-8 flex h-20 w-20 items-center justify-center rounded-2xl bg-brand-light text-brand sm:h-24 sm:w-24">
          <SearchX className="h-10 w-10 sm:h-12 sm:w-12" />
        </div>
        <h1 className="mt-6 text-2xl font-bold tracking-tight text-on-surface sm:text-3xl">Không tìm thấy trang</h1>
        <p className="mt-2 max-w-md text-sm leading-relaxed text-on-surface-variant">
          Trang bạn đang tìm không tồn tại, đã được di chuyển hoặc bạn không có quyền truy cập.
          <br className="hidden sm:block" />
          Kiểm tra lại đường dẫn hoặc quay về trang an toàn.
        </p>
        <p className="mt-2 font-mono text-xs text-on-surface-variant/70">Mã lỗi: 404 • NOT_FOUND</p>
      </div>

      {/* Actions */}
      <div className="mt-8 flex flex-col items-center gap-3 sm:flex-row">
        <Button render={<Link href="/landing" />} className="h-10 px-6 text-sm font-semibold shadow-sm">
          <Home className="h-4 w-4" />
          Về trang chủ
        </Button>
        <Button
          variant="outline"
          render={<Link href="/app/overview" />}
          className="h-10 bg-surface px-6 text-sm font-medium border-outline hover:bg-surface-container"
        >
          <ArrowLeft className="h-4 w-4" />
          Về Dashboard
        </Button>
      </div>

      <p className="mt-10 text-center text-xs text-on-surface-variant">
        Cần hỗ trợ? Liên hệ quản lý salon qua{" "}
        <a href="tel:0900000000" className="font-semibold text-primary hover:underline">
          0900 000 000
        </a>
      </p>
    </div>
  )
}
