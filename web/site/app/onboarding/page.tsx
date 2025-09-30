"use client";

import React, { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { MapPin, ArrowLeft } from "lucide-react";

export default function OnboardingPage() {
  const router = useRouter();
  const [locating, setLocating] = useState(false);
  const [step, setStep] = useState<1 | 2 | 3>(1);

  const handleDropEgg = () => {
    if (!("geolocation" in navigator)) {
      toast.error("Geolocation is not supported on this device.");
      return;
    }
    setLocating(true);
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude, longitude } = pos.coords;
        // store for map to use on first load
        try {
          sessionStorage.setItem("user_location", JSON.stringify({ latitude, longitude }));
        } catch {}
        toast.success("Location acquired! Let the adventure begin!.");
        router.push("/map");
      },
      (err) => {
        toast.error(err.message || "Unable to get location");
        setLocating(false);
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 }
    );
  };

  return (
    <div className="min-h-screen hero-gradient shimmer-bg relative overflow-hidden">
      {/* Background floating eggs */}
      <div className="pointer-events-none absolute inset-0 opacity-30 z-0">
        <div className="absolute top-20 left-10 w-16 h-16 bg-gradient-egg rounded-full animate-float" />
        <div className="absolute top-40 right-20 w-12 h-12 bg-gradient-adventure rounded-full animate-float" style={{ animationDelay: '1s' }} />
        <div className="absolute bottom-32 left-1/4 w-20 h-20 bg-gradient-nature rounded-full animate-bounce-soft" style={{ animationDelay: '2s' }} />
        <div className="absolute top-32 left-2/3 w-14 h-14 bg-primary/60 rounded-full animate-float" style={{ animationDelay: '1.5s' }} />
      </div>
      <div className="absolute top-16 left-6">
            <Link href="/" className="inline-flex items-center gap-2 rounded-full px-4 py-2 bg-white/10 border border-white/20 text-white backdrop-blur-sm hover:bg-white/20 text-sm font-semibold">
              <ArrowLeft className="w-4 h-4" /> Home
            </Link>
          </div>

      <main className="relative z-10 min-h-screen flex items-center justify-center px-4 mt-8 sm:px-6 lg:px-8">
        <div className="w-full max-w-xl mx-auto bg-card/90 backdrop-blur-sm rounded-3xl border border-border shadow-game p-6 sm:p-10 text-center">
          {step === 1 && (
            <>
              <h1 className="display-font text-5xl font-extrabold bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 bg-clip-text text-transparent">
                Welcome, fellow Eggsplorer!
              </h1>
              <p className="mt-4 text-muted-foreground">
                Ready to touch grass and make your adventures count? Let’s get you started.
              </p>
              <div className="mt-10">
                <button
                  type="button"
                  onClick={() => setStep(2)}
                  className="inline-flex items-center gap-2 px-6 py-2 rounded-full bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-semibold shadow-2xl hover:scale-105 transition-transform"
                >
                  Next
                </button>
              </div>
            </>
          )}

          {step === 2 && (
            <>
              <h2 className="display-font text-4xl font-extrabold bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 bg-clip-text text-transparent">
                Plant eggs where life happens
              </h2>
              <p className="mt-4 text-muted-foreground">
                Drop eggs at the places you want to visit regularly in the real world — parks, gyms, cafés, trails.
                Every visit can earn you XP and unlock sweet rewards.
              </p>
              <div className="mt-10 flex items-center justify-center gap-3">
                <button
                  type="button"
                  onClick={() => setStep(1)}
                  className="inline-flex items-center gap-2 px-5 py-2 rounded-full bg-white/10 border border-white/20 text-white hover:bg-white/20"
                >
                  Back
                </button>
                <button
                  type="button"
                  onClick={() => setStep(3)}
                  className="inline-flex items-center gap-2 px-6 py-2 rounded-full bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-semibold shadow-2xl hover:scale-105 transition-transform"
                >
                  Next
                </button>
              </div>
            </>
          )}

          {step === 3 && (
            <>
              <h1 className="display-font text-5xl font-extrabold bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 bg-clip-text text-transparent">
                Drop your first egg
              </h1>
              <p className="mt-3 text-muted-foreground">
                We’ll use your current location to place your first egg on the map. You can add more later!
              </p>

              <div className="mt-10">
                <button
                  type="button"
                  onClick={handleDropEgg}
                  disabled={locating}
                  className={`inline-flex items-center gap-2 px-6 py-2 rounded-full bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-semibold shadow-2xl transition-transform ${locating ? "opacity-70" : "hover:scale-105"}`}
                >
                  <MapPin className="w-5 h-5" /> {locating ? "Locating..." : "Drop Egg Here"}
                </button>
              </div>

              <p className="mt-6 text-xs text-muted-foreground">
                We respect your privacy. Your location is only used to place your egg.
              </p>
              <div className="mt-6">
                <button
                  type="button"
                  onClick={() => setStep(2)}
                  className="inline-flex items-center gap-2 px-5 py-2 rounded-full bg-white/10 border border-white/20 text-white hover:bg-white/20"
                >
                  Back
                </button>
              </div>
            </>
          )}
        </div>
      </main>
    </div>
  );
}
