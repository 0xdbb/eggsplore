"use client";

import React, { useEffect } from "react";
import { usePathname } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import AppChrome from "./AppChrome";
import { useAuthStore } from "../lib/store";

export default function ClientShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const showGlobalLogo = pathname === "/" || pathname?.startsWith("/auth");
  const musicEnabled = useAuthStore((s) => s.prefs.musicEnabled);

  useEffect(() => {
    if (typeof window !== "undefined" && "serviceWorker" in navigator) {
      navigator.serviceWorker
        .register("/sw.js", { scope: "/" })
        .catch(() => {
          // no-op; registration may fail in dev or unsupported contexts
        });
    }
  }, []);

  return (
    <>
      {showGlobalLogo ? (
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
            <span className="hidden sm:inline display-font text-lg bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 bg-clip-text text-transparent">
              Eggsplore
            </span>
          </Link>
        </div>
      ) : pathname?.startsWith("/home") ? (
        <AppChrome />
      ) : null}
      {/* Background calm music element; controlled from Settings */}
      {!showGlobalLogo && (
        <audio src="/audio/ambient.mp3" loop autoPlay={musicEnabled} muted={!musicEnabled} preload="auto" />
      )}
      {children}
    </>
  );
}
