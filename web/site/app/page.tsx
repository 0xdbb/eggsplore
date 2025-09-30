"use client";

import React, { useState } from 'react';
import { MapPin, Zap, Trophy, Play, Users, Star, ArrowRight, Egg, Sparkles } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import HowItWorks from '../components/HowItWorks';
import { useAuthStore } from '../lib/store';

export default function HomePage() {
  const [isSignUpHovered, setIsSignUpHovered] = useState(false);
  const [isSignInHovered, setIsSignInHovered] = useState(false);
  const [isFinalCtaHovered, setIsFinalCtaHovered] = useState(false);
  const user = useAuthStore((s) => s.user);
  const session = useAuthStore((s) => s.session);
  const cookieAuthed = typeof document !== 'undefined' && document.cookie.includes('access_token=');
  const isAuthed = !!(user?.id || session?.accessToken || cookieAuthed);

  return (
    <div className="min-h-screen hero-gradient shimmer-bg overflow-x-hidden">

      {/* Hero Section */}
      <section className="relative min-h-screen flex items-center justify-center px-4 sm:px-6 lg:px-8 overflow-hidden">
        {/* Login button (top-right) — hidden if authenticated */}
        {!isAuthed && (
          <div className="absolute top-4 right-4 z-20">
            <Link
              href="/auth"
              className="inline-flex items-center gap-2 rounded-full px-5 py-2 bg-white/10 border border-white/20 text-white backdrop-blur-sm hover:bg-white/20 transition-colors text-sm font-semibold shadow-soft"
            >
              Login
            </Link>
          </div>
        )}
        {/* Floating background eggs inside the hero (behind content) */}
        <div className="pointer-events-none absolute inset-0 opacity-30 z-0">
          <div className="absolute top-20 left-10 w-16 h-16 bg-gradient-egg rounded-full animate-float" style={{ animationDelay: '0s' }}></div>
          <div className="absolute top-40 right-20 w-12 h-12 bg-gradient-adventure rounded-full animate-float" style={{ animationDelay: '1s' }}></div>
          <div className="absolute bottom-32 left-1/4 w-20 h-20 bg-gradient-nature rounded-full animate-bounce-soft" style={{ animationDelay: '2s' }}></div>
          <div className="absolute top-32 left-2/3 w-14 h-14 bg-primary/60 rounded-full animate-float" style={{ animationDelay: '1.5s' }}></div>
          <div className="absolute bottom-20 right-1/3 w-[4.5rem] h-[4.5rem] bg-secondary/70 rounded-full animate-bounce-soft" style={{ animationDelay: '0.5s' }}></div>
        </div>
        <div className="max-w-6xl mx-auto text-center relative z-10">
          {/* Logo Area */}
          <div className="mb-8">
            <div className="inline-flex items-center gap-3 mb-4">
              <div className="relative w-16 h-16 sm:w-20 sm:h-20 rounded-xl overflow-hidden border border-white/10 bg-white/5 shadow-soft">
                <Image src="/logo.png" alt="Eggsplore Logo" fill className="object-cover" priority />
              </div>
              <h1 className="display-font text-5xl sm:text-7xl font-black text-transparent bg-clip-text bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 drop-shadow-[0_10px_20px_rgba(244,114,182,0.25)]">
                Eggsplore
              </h1>
            </div>
          </div>

          {/* Tagline */}
          {/* <h2 className="display-font text-3xl sm:text-4xl lg:text-5xl font-extrabold text-white mb-6 leading-tight text-glow">
            <span className="text-yellow-400">Hatch</span>, <span className="text-green-400">Explore</span>, <span className="text-purple-400">Conquer</span>
          </h2> */}

          {/* Subtitle */}
          <p className="text-lg sm:text-xl text-gray-200 mb-12 max-w-3xl mx-auto leading-relaxed">
            Drop virtual eggs anywhere and hatch amazing creatures by revisiting them.
          </p>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
          <Link
            href={isAuthed ? "/home" : "/auth"}
            onMouseEnter={() => setIsFinalCtaHovered(true)}
            onMouseLeave={() => setIsFinalCtaHovered(false)}
            className={`inline-flex items-center gap-3 px-5 py-3 bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-bold text-xl rounded-full shadow-2xl transform transition-all duration-300 hover:scale-105 hover:shadow-rose-400/25 ${isFinalCtaHovered ? 'scale-105' : ''}`}
          >
            Join the Adventure
            <ArrowRight className={`w-6 h-6 transition-transform duration-300 ${isFinalCtaHovered ? 'translate-x-1' : ''}`} />
          </Link>
          </div>

          {/* Easter preview */}
          <div className="mt-16 flex justify-center items-center gap-8 text-6xl animate-pulse"> 
            <div className="transform rotate-12 hover:scale-110 transition-transform duration-300">🐰</div>
            <div className="transform -rotate-6 hover:scale-110 transition-transform duration-300">🥚</div>
            <div className="transform rotate-6 hover:scale-110 transition-transform duration-300">🌷</div>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <HowItWorks isAuthed={isAuthed} />

  

   
      {/* Footer */}
      <footer className="py-2 px-4 sm:px-6 lg:px-8 bg-black/20 backdrop-blur-sm border-t border-white/10">
        <div className="max-w-4xl mx-auto text-center">
          <p className="text-gray-400 text-xs">
            © 2025 Eggsplore.quest - Hatch your adventure today
          </p>
        </div>
      </footer>
    </div>
  );
}
