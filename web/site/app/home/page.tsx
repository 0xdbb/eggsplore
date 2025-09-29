"use client";

import React from "react";
import Link from "next/link";
import Image from "next/image";
import { useAuthStore } from "../../lib/store";
import { Star, Coins, MapPin, ShoppingBag, Users, ClipboardList, Home as HomeIcon } from "lucide-react";

export default function HomePage() {
  const user = useAuthStore((s) => s.user);
  const xp = user?.xp ?? 0;
  const coins = user?.coins ?? 0;

  return (
    <div className="min-h-screen hero-gradient shimmer-bg relative overflow-hidden">
      {/* Floating background shapes */}
      <div className="pointer-events-none absolute inset-0 opacity-30 z-0">
        <div className="absolute top-20 left-10 w-16 h-16 bg-gradient-egg rounded-full animate-float" />
        <div className="absolute top-40 right-20 w-12 h-12 bg-gradient-adventure rounded-full animate-float" style={{ animationDelay: '1s' }} />
        <div className="absolute bottom-32 left-1/4 w-20 h-20 bg-gradient-nature rounded-full animate-bounce-soft" style={{ animationDelay: '2s' }} />
        <div className="absolute top-32 left-2/3 w-14 h-14 bg-primary/60 rounded-full animate-float" style={{ animationDelay: '1.5s' }} />
      </div>

      {/* Header */}
      <header className="relative z-10 flex items-center justify-between px-4 sm:px-6 lg:px-8 pt-5">
        {/* XP and Coins */}
        <div className="flex items-center gap-2">
          <div className="inline-flex items-center gap-2 rounded-full px-3 py-1 bg-card/90 border border-border text-foreground shadow-soft">
            <Star className="w-4 h-4 text-amber-300" />
            <span className="text-xs opacity-80">XP</span>
            <span className="font-semibold">{xp}</span>
          </div>
          <div className="inline-flex items-center gap-2 rounded-full px-3 py-1 bg-card/90 border border-border text-foreground shadow-soft">
            <Coins className="w-4 h-4 text-yellow-300" />
            <span className="text-xs opacity-80">Coins</span>
            <span className="font-semibold">{coins}</span>
          </div>
        </div>

      </header>

      {/* Content */}
      <main className="relative z-10 px-4 sm:px-6 lg:px-8 pt-10 pb-28">
        <h1 className="display-font text-4xl sm:text-5xl font-extrabold bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 bg-clip-text text-transparent">
          Welcome back{user?.username ? `, ${user.username}` : "!"}
        </h1>
        <p className="mt-3 text-muted-foreground max-w-prose">
          Continue your quest. Explore new places, drop eggs, and rack up XP.
        </p>

        <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <Link
            href="/map"
            className="group rounded-2xl bg-card/90 border border-border shadow-soft hover:shadow-game transition-all p-5 flex items-center justify-between"
          >
            <div>
              <div className="text-foreground font-semibold">Open Map</div>
              <div className="text-sm text-muted-foreground">See your eggs and live location</div>
            </div>
            <div className="rounded-xl p-3 bg-gradient-to-br from-rose-300/40 to-amber-300/40 border border-white/10">
              <MapPin className="w-5 h-5 text-rose-400" />
            </div>
          </Link>
          <Link
            href="/quests"
            className="group rounded-2xl bg-card/90 border border-border shadow-soft hover:shadow-game transition-all p-5 flex items-center justify-between"
          >
            <div>
              <div className="text-foreground font-semibold">Quests</div>
              <div className="text-sm text-muted-foreground">Daily and weekly challenges</div>
            </div>
            <div className="rounded-xl p-3 bg-gradient-to-br from-sky-300/40 to-violet-300/40 border border-white/10">
              <ClipboardList className="w-5 h-5 text-sky-400" />
            </div>
          </Link>
          <Link
            href="/shop"
            className="group rounded-2xl bg-card/90 border border-border shadow-soft hover:shadow-game transition-all p-5 flex items-center justify-between"
          >
            <div>
              <div className="text-foreground font-semibold">Shop</div>
              <div className="text-sm text-muted-foreground">Cosmetics and boosts</div>
            </div>
            <div className="rounded-xl p-3 bg-gradient-to-br from-emerald-300/40 to-lime-300/40 border border-white/10">
              <ShoppingBag className="w-5 h-5 text-emerald-500" />
            </div>
          </Link>
        </div>
      </main>

      {/* Bottom navigation (mobile-first) */}
      <nav className="fixed bottom-0 left-0 right-0 z-20 border-t border-white/10 bg-card/90 backdrop-blur-md">
        <div className="mx-auto max-w-3xl grid grid-cols-4">
          <Link href="/home" className="flex flex-col items-center py-3 text-sm text-foreground">
            <HomeIcon className="w-5 h-5" />
            <span className="text-xs mt-1">Home</span>
          </Link>
          <Link href="/shop" className="flex flex-col items-center py-3 text-sm text-muted-foreground hover:text-foreground">
            <ShoppingBag className="w-5 h-5" />
            <span className="text-xs mt-1">Shop</span>
          </Link>
          <Link href="/friends" className="flex flex-col items-center py-3 text-sm text-muted-foreground hover:text-foreground">
            <Users className="w-5 h-5" />
            <span className="text-xs mt-1">Friends</span>
          </Link>
          <Link href="/quests" className="flex flex-col items-center py-3 text-sm text-muted-foreground hover:text-foreground">
            <ClipboardList className="w-5 h-5" />
            <span className="text-xs mt-1">Quests</span>
          </Link>
        </div>
      </nav>
    </div>
  );
}
