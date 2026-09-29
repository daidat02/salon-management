import { NextRequest, NextResponse } from "next/server";

// Bảo vệ toàn bộ /app* trừ /app/login.
// Middleware chỉ đọc được cookie (refresh_token httpOnly do Go set),
// không đọc được localStorage/sessionStorage — nên đây là chốt chặn vòng ngoài.
// Verify access_token thật thì làm ở Server Component / route handler gọi API Go.
const PUBLIC_APP_PATHS = ["/app/login", "/app/register", "/app/forgot-password"];

export function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  if (PUBLIC_APP_PATHS.some((p) => pathname === p || pathname.startsWith(p + "/"))) {
    return NextResponse.next();
  }

  const refreshToken = req.cookies.get("refresh_token")?.value;
  if (!refreshToken) {
    const loginUrl = new URL("/app/login", req.url);
    loginUrl.searchParams.set("next", pathname);
    return NextResponse.redirect(loginUrl);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/app/:path*"],
};
