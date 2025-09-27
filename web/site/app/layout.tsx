export const metadata = {
  title: 'Eggsplore.quest',
  description: 'Hatch, Explore, Conquer',
};

import './globals.css';
import React from 'react';
import { Chewy, Nunito } from 'next/font/google';
import { Toaster } from 'sonner';
import Link from 'next/link';
import Image from 'next/image';

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
      </head>
      <body className="body-font bg-background text-foreground">
        {/* Global top-left logo */}
        <div id="global-logo" className="fixed top-4 left-4 z-50">
          <Link href="/" className="inline-flex items-center gap-2">
            <Image
              src="/logo.png"
              alt="Eggsplore Logo"
              width={40}
              height={40}
              className="rounded-xl shadow-soft border border-white/10 bg-white/5"
              priority
            />
            <span className="hidden sm:inline display-font text-lg bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 bg-clip-text text-transparent">Eggsplore</span>
          </Link>
        </div>
        {children}
        {/* Global toaster for notifications */}
        <Toaster
          theme="dark"
          richColors
          position="top-center"
          toastOptions={{
            classNames: {
              toast: 'bg-card text-foreground border border-border shadow-soft',
              title: 'text-foreground',
              description: 'text-muted-foreground',
              actionButton: 'bg-primary text-primary-foreground',
              cancelButton: 'bg-white/10 text-foreground border border-white/20',
              closeButton: 'text-muted-foreground hover:text-foreground',
              loader: 'bg-primary',
            },
          }}
        />
      </body>
    </html>
  );
}
