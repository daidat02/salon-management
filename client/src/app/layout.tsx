import type { Metadata } from 'next';
import './globals.css';
import { Geist } from 'next/font/google';
import { cn } from '@/lib/utils';
import StoreProvider from './StoreProvider';
import TanStackProvider from './TanstackProvider';
import { Toaster } from '@/components/ui/toast';

const geist = Geist({ subsets: ['latin'], variable: '--font-sans' });

export const metadata: Metadata = {
  title: 'Salon Project',
  description: 'Ứng dụng quản lý salon',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="vi" suppressHydrationWarning className={cn('font-sans', geist.variable)}>
      <body className="bg-background text-body font-sans antialiased" suppressHydrationWarning>
        <StoreProvider>
          <TanStackProvider>
            <Toaster />
            {children}
          </TanStackProvider>
        </StoreProvider>
      </body>
    </html>
  );
}
