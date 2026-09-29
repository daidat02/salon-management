import Link from "next/link";

// Home của website landing — route `/landing`.
export default function LandingPage() {
  return (
    <main className="mx-auto flex min-h-[70vh] max-w-5xl flex-col items-center justify-center gap-6 p-8 text-center">
      <h1 className="text-3xl font-bold text-gray-800">
        Salon Project — Landing Page
      </h1>
      <p className="text-gray-600">
        Đây là trang chủ của website landing tại <code>/landing</code>.
      </p>
      <Link
        href="/app"
        className="rounded-lg bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700"
      >
        Vào /app (Dashboard)
      </Link>
    </main>
  );
}
