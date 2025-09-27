"use client";

import React, { useEffect, useState } from "react";
import { z } from "zod";
import { toast } from "sonner";

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

export type AuthModalProps = {
  open: boolean;
  onClose: () => void;
};

export default function AuthModal({ open, onClose }: AuthModalProps) {
  const [mode, setMode] = useState<"signin" | "register">("signin");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({});

  useEffect(() => {
    if (!open) {
      setEmail("");
      setPassword("");
      setErrors({});
      setSubmitting(false);
      setMode("signin");
    }
  }, [open]);

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
      // Placeholder: integrate your auth (e.g., Supabase) here.
      await new Promise((r) => setTimeout(r, 600));
      if (mode === "signin") {
        toast.success("Signed in successfully. Welcome back!", { className: "bg-card text-foreground" });
      } else {
        toast.success("Account created! You can start your quest.", { className: "bg-card text-foreground" });
      }
      onClose();
    } catch (err) {
      toast.error("Something went wrong. Please try again.", { className: "bg-card text-foreground" });
    } finally {
      setSubmitting(false);
    }
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <button
        aria-label="Close auth"
        onClick={onClose}
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
      />

      {/* Modal */}
      <div className="relative z-10 w-full max-w-md mx-4 rounded-3xl border border-border bg-card text-foreground shadow-game">
        <div className="p-6 sm:p-8">
          <div className="mb-6 text-center">
            <h2 className="display-font text-3xl font-extrabold bg-gradient-to-r from-pink-300 via-rose-300 to-sky-300 bg-clip-text text-transparent">
              {mode === "signin" ? "Sign In" : "Create Account"}
            </h2>
            <p className="mt-2 text-sm text-muted-foreground">
              {mode === "signin" ? "Welcome back!" : "Join the adventure."}
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-muted-foreground">
                Email
              </label>
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
              <label htmlFor="password" className="block text-sm font-medium text-muted-foreground">
                Password
              </label>
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
              <button className="text-primary hover:underline" onClick={() => setMode("register")}>Don\'t have an account? Register</button>
            ) : (
              <button className="text-primary hover:underline" onClick={() => setMode("signin")}>Already have an account? Sign in</button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
