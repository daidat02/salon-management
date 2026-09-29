export default function AppPage() {
  return (
    <div>
      <h2 className="text-xl font-semibold text-gray-800">Xin chào từ route /app</h2>
      <p className="mt-2 text-gray-600">
        File này là <code>src/app/app/(dashboard)/page.tsx</code> — nội dung render
        bên trong DashboardLayout. Route này đã được middleware bảo vệ.
      </p>
    </div>
  );
}
