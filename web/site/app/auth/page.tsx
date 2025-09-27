"use client";

import React, { useState } from "react";
import Link from "next/link";
import { z } from "zod";
import { toast } from "sonner";
import { ArrowLeft } from "lucide-react";
import dynamic from "next/dynamic";
import { useRouter } from "next/navigation";

// Dynamically import confetti (client-only)
const ConfettiBoom = dynamic(() => import("react-confetti-boom"), { ssr: false });

const authSchema = z.object({
  email: z
    .string()
    .min(1, "Email is required")
    .email("Enter a valid email"),
  password: z
    .string()
    .min(8, "Password must be at least 8 characters")
    .max(72, "Password is too long")
    .regex(/\d/, "Include at least one number")
    .regex(/[A-Za-z]/, "Include at least one letter"),
  mode: z.enum(["signin", "register"]).default("signin"),
});

export default function AuthPage() {
  const router = useRouter();
  const [mode, setMode] = useState<"signin" | "register">("signin");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({});
  const [boom, setBoom] = useState(false);

  const validate = () => {
    const result = authSchema.safeParse({ email, password, mode });
    if (!result.success) {
      const fieldErrors: { email?: string; password?: string } = {};
      for (const issue of result.error.issues) {
        if (issue.path[0] === "email") fieldErrors.email = issue.message;
        if (issue.path[0] === "password") fieldErrors.password = issue.message;
      }
      setErrors(fieldErrors);
      return false;
    }
    setErrors({});
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    setSubmitting(true);
    try {
      // TODO: integrate with your auth provider (e.g., Supabase)
      await new Promise((r) => setTimeout(r, 600));
      if (mode === "signin") {
        toast.success("Signed in successfully. Welcome back!");
        router.push("/map");
      } else {
        toast.success("Account created! You can start your quest.");
        setBoom(true);
        // Show confetti briefly, then navigate to map
        setTimeout(() => {
          setBoom(false);
          router.push("/map");
        }, 2000);
      }
    } catch (err) {
      toast.error("Something went wrong. Please try again.");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-violet-900 via-indigo-900 to-sky-900 shimmer-bg">
      <section className="relative min-h-screen flex items-center justify-center px-4 sm:px-6 lg:px-8 hero-gradient overflow-hidden">
        {/* Fancy back button */}
        <div className="absolute top-20 left-4 z-20 sm:top-24">
          <Link
            href="/"
            className="inline-flex items-center gap-2 rounded-full px-4 py-2 bg-white/10 border border-white/20 text-white backdrop-blur-sm hover:bg-white/20 transition-colors text-sm font-semibold shadow-soft"
          >
            <ArrowLeft className="w-4 h-4" /> Back Home
          </Link>
        </div>

        {/* Background bubbles reused */}
        <div className="pointer-events-none absolute inset-0 opacity-20 z-0">
          <div className="absolute top-20 left-10 w-16 h-16 bg-gradient-egg rounded-full animate-float" style={{ animationDelay: '0s' }} />
          <div className="absolute top-40 right-20 w-12 h-12 bg-gradient-adventure rounded-full animate-float" style={{ animationDelay: '1s' }} />
          <div className="absolute bottom-32 left-1/4 w-20 h-20 bg-gradient-nature rounded-full animate-bounce-soft" style={{ animationDelay: '2s' }} />
          <div className="absolute top-32 left-2/3 w-14 h-14 bg-primary/60 rounded-full animate-float" style={{ animationDelay: '1.5s' }} />
          <div className="absolute bottom-20 right-1/3 w-[4.5rem] h-[4.5rem] bg-secondary/70 rounded-full animate-bounce-soft" style={{ animationDelay: '0.5s' }} />
        </div>

        <div className="relative z-10 w-full max-w-md mx-auto bg-card/90 backdrop-blur-sm rounded-3xl border border-border shadow-game p-6 sm:p-8">
          {boom && (
            <div className="pointer-events-none fixed top-6 left-1/2 -translate-x-1/2 z-[9999]">
              {/* Confetti tuned to fall from above */}
              {/* Some props may be ignored depending on library version, but are safe */}
              <ConfettiBoom
                // @ts-ignore
                mode="fall"
                particleCount={220}
                shapeSize={10}
                colors={["#fda4af", "#f9a8d4", "#a78bfa", "#93c5fd", "#fde68a"]}
                fadeOutHeight={160}
              />
            </div>
          )}
          <div className="text-center mb-6">
            <h1 className="display-font text-4xl font-extrabold bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 bg-clip-text text-transparent">
              {mode === "signin" ? "Sign In" : "Create Account"}
            </h1>
            <p className="mt-2 text-sm text-muted-foreground">
              {mode === "signin" ? "Welcome back!" : "Join the adventure."}
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-muted-foreground">Email</label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                onBlur={validate}
                className="mt-1 w-full rounded-xl border border-border bg-background px-4 py-3 text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40"
                placeholder="you@example.com"
                autoComplete="email"
              />
              {errors.email && <p className="mt-1 text-sm text-rose-400">{errors.email}</p>}
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-muted-foreground">Password</label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                onBlur={validate}
                className="mt-1 w-full rounded-xl border border-border bg-background px-4 py-3 text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40"
                placeholder="••••••••"
                autoComplete={mode === "signin" ? "current-password" : "new-password"}
              />
              {errors.password && <p className="mt-1 text-sm text-rose-400">{errors.password}</p>}
              <p className="mt-2 text-xs text-muted-foreground">Must be 8+ chars, include a letter and a number.</p>
            </div>

            <button
              type="submit"
              disabled={submitting}
              className={`w-full inline-flex items-center justify-center gap-2 px-6 py-3 rounded-full bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-semibold shadow-2xl transition-transform ${submitting ? "opacity-70" : "hover:scale-105"}`}
            >
              {submitting ? "Please wait..." : mode === "signin" ? "Sign In" : "Create Account"}
            </button>
          </form>

          <div className="mt-6 text-center text-sm">
            {mode === "signin" ? (
              <button className="text-primary hover:underline" onClick={() => setMode("register")}>
                Don't have an account? Register
              </button>
            ) : (
              <button className="text-primary hover:underline" onClick={() => setMode("signin")}>
                Already have an account? Sign in
              </button>
            )}
          </div>

          {/* Demo buttons removed per request */}
        </div>

        {/* Footer */}
        <footer className="absolute bottom-2 left-0 right-0 text-center text-xs text-gray-400">
          © 2025 Eggsplore.quest - Hatch your adventure today
        </footer>
      </section>
    </div>
  );
}
