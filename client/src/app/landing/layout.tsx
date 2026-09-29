// Layout riêng của website landing — header/footer marketing nằm ở đây.
// Không ảnh hưởng dashboard /app.
export default function LandingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="theme-landing min-h-screen bg-background font-sans text-on-surface">
      <header className="border-b border-border bg-white">
        <div className="mx-auto flex max-w-5xl items-center justify-between p-4">
          <span className="font-display text-lg font-bold text-on-surface">Salon</span>
          <nav className="flex gap-4 text-sm font-medium text-on-surface-variant">
            <a href="/landing" className="hover:text-primary">
              Trang chủ
            </a>
            <a href="/app" className="hover:text-primary">
              Dashboard
            </a>
          </nav>
        </div>
      </header>
      {children}
      <footer className="border-t border-border bg-white">
        <div className="mx-auto max-w-5xl p-4 text-center text-sm text-on-surface-variant">
          © 2026 Salon Project
        </div>
      </footer>
    </div>
  );
}
