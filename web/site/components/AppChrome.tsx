"use client";

import React, { useEffect, useRef, useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { Bell, ChevronDown, LogOut, Settings, User2 } from "lucide-react";
import { useAuthStore } from "../lib/store";

export default function AppChrome() {
  const user = useAuthStore((s) => s.user);
  const [open, setOpen] = useState(false);
  const [hovered, setHovered] = useState(false);
  const menuRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const onDocClick = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("click", onDocClick);
    return () => document.removeEventListener("click", onDocClick);
  }, []);

  return (
    <>
      {/* Top-right chrome */}
      <div className="fixed top-4 right-4 z-50 flex items-center gap-2">
        {/* Notifications bell */}
        <Link
          href="/notifications"
          className="inline-flex items-center gap-2 rounded-full px-3 py-2 bg-card/90 border border-border text-foreground shadow-soft hover:bg-white/10"
          aria-label="Notifications"
        >
          <Bell className="w-4 h-4" />
          <span className="text-xs hidden sm:inline">Alerts</span>
        </Link>

        {/* Avatar with hover animation and dropdown */}
        <div ref={menuRef} className="relative">
          <button
            type="button"
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
            onClick={() => setOpen((v) => !v)}
            className={`inline-flex items-center gap-2 rounded-full pl-2 pr-3 py-1.5 bg-card/90 border border-border text-foreground shadow-soft hover:bg-white/10 transition-transform ${hovered ? "scale-105" : ""}`}
            aria-haspopup="menu"
            aria-expanded={open}
          >
            <Image
              src="/logo.png"
              alt={user?.username || user?.email || "Profile"}
              width={28}
              height={28}
              className="rounded-xl border border-white/10"
            />
            <ChevronDown className={`w-4 h-4 transition-transform ${open ? "rotate-180" : ""}`} />
          </button>
          {open && (
            <div
              role="menu"
              className="absolute right-0 mt-2 w-44 rounded-xl bg-card border border-border shadow-game overflow-hidden"
            >
              <Link
                href="/profile"
                className="flex items-center gap-2 px-3 py-2 text-sm hover:bg-white/10"
                role="menuitem"
                onClick={() => setOpen(false)}
              >
                <User2 className="w-4 h-4" /> View Profile
              </Link>
              <Link
                href="/settings"
                className="flex items-center gap-2 px-3 py-2 text-sm hover:bg-white/10"
                role="menuitem"
                onClick={() => setOpen(false)}
              >
                <Settings className="w-4 h-4" /> Settings
              </Link>
              <button
                type="button"
                className="w-full text-left flex items-center gap-2 px-3 py-2 text-sm hover:bg-white/10"
                role="menuitem"
                onClick={() => {
                  // TODO: call logout endpoint then reset store
                  // useAuthStore.getState().clear();
                  setOpen(false);
                }}
              >
                <LogOut className="w-4 h-4" /> Logout
              </button>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
