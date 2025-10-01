export const metadata = {
  title: 'Eggsplore',
  description: 'Eggsplore and Touch Grass',
};

import './globals.css';
import React from 'react';
import { Chewy, Nunito } from 'next/font/google';
import { Toaster } from 'sonner';
import { ReactQueryProvider } from '../lib/queryClient';
import ClientShell from '../components/ClientShell';

const display = Chewy({ subsets: ['latin'], weight: '400', variable: '--font-display' });
const body = Nunito({ subsets: ['latin'], variable: '--font-body' });

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={`${display.variable} ${body.variable}`}>
      <head>
        <link rel="icon" href="/favicon.ico" />
        <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
        <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
        <link rel="apple-touch-icon" href="/apple-touch-icon.png" />
        <link rel="manifest" href="/manifest.webmanifest" />
        <meta name="theme-color" content="#0f172a" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
      </head>
      <body className="body-font bg-background text-foreground">
        <ClientShell>
          <ReactQueryProvider>
            {children}
          </ReactQueryProvider>
          {/* Global toaster for notifications */}
          <Toaster
            theme="dark"
            position="top-center"
            toastOptions={{
              classNames: {
                toast: 'bg-card text-foreground border border-border shadow-game',
                success: 'bg-card text-foreground border border-border',
                error: 'bg-card text-foreground border border-border',
                warning: 'bg-card text-foreground border border-border',
                info: 'bg-card text-foreground border border-border',
              },
            }}
          />
        </ClientShell>
      </body>
    </html>
  );
}
