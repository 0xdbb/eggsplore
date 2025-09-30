import { NextResponse, NextRequest } from 'next/server';

// Public routes that do not require auth
const PUBLIC_PATHS = new Set<string>([
  '/',
  '/auth',
]);

// Paths that are always allowed (static/assets)
function isAlwaysAllowed(pathname: string) {
  const isRootImage = /\.(png|jpg|jpeg|webp|gif|svg)$/i.test(pathname) && pathname.split('/').length <= 2;
  return (
    pathname.startsWith('/_next/') ||
    pathname.startsWith('/static/') ||
    pathname.startsWith('/public/') ||
    pathname.startsWith('/api/') || // allow Next API routes if any
    pathname.startsWith('/favicon') ||
    pathname.startsWith('/icons') ||
    pathname.startsWith('/apple-touch-icon') ||
    pathname.startsWith('/manifest') ||
    pathname.startsWith('/sw.js') ||
    pathname === '/logo.png' ||
    pathname.startsWith('/icon-') ||
    pathname.startsWith('/maskable-') ||
    isRootImage
  );
}

export function middleware(req: NextRequest) {
  const { pathname, search } = req.nextUrl;

  if (isAlwaysAllowed(pathname) || PUBLIC_PATHS.has(pathname)) {
    return NextResponse.next();
  }

  // Check auth cookie
  const access = req.cookies.get('access_token')?.value;

  if (!access) {
    const url = req.nextUrl.clone();
    url.pathname = '/auth';
    const nextParam = `next=${encodeURIComponent(req.nextUrl.pathname + (search || ''))}`;
    url.search = url.search ? `${url.search}&${nextParam}` : `?${nextParam}`;
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    // Protect all app routes except public ones (landing/auth) and assets
    '/((?!_next/|static/|public/|favicon|icons|apple-touch-icon|manifest|sw\.js).*)',
  ],
};
