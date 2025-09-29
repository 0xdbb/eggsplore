"use client";

import React, { useState } from "react";
import { useAuthStore } from "../../lib/store";
import { toast } from "sonner";
import Link from "next/link";
import { ArrowLeft, Music2, Volume2, User, ShieldCheck } from "lucide-react";

export default function SettingsPage() {
  const { user, setUser, prefs, setPrefs } = useAuthStore();
  const [username, setUsername] = useState(user?.username || "");

  // Change password local state
  const [currentPw, setCurrentPw] = useState("");
  const [newPw, setNewPw] = useState("");
  const [confirmPw, setConfirmPw] = useState("");
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    setSaving(true);
    try {
      // TODO: call backend to update username and prefs as needed
      setUser({ ...(user || {}), username });
      toast.success("Settings saved");
    } catch (e: any) {
      toast.error(e?.message || "Failed to save settings");
    } finally {
      setSaving(false);
    }
  };

  const handleChangePassword = async () => {
    if (!newPw || newPw.length < 8) {
      toast.error("New password must be at least 8 characters");
      return;
    }
    if (newPw !== confirmPw) {
      toast.error("New password and confirmation do not match");
      return;
    }
    try {
      // TODO: call backend /auth/change-password
      setCurrentPw("");
      setNewPw("");
      setConfirmPw("");
      toast.success("Password updated");
    } catch (e: any) {
      toast.error(e?.message || "Failed to update password");
    }
  };

  return (
    <div className="min-h-screen hero-gradient shimmer-bg relative overflow-hidden">
      {/* Background */}
      <div className="pointer-events-none absolute inset-0 opacity-30 z-0">
        <div className="absolute top-20 left-10 w-16 h-16 bg-gradient-egg rounded-full animate-float" />
        <div className="absolute top-40 right-20 w-12 h-12 bg-gradient-adventure rounded-full animate-float" style={{ animationDelay: '1s' }} />
        <div className="absolute bottom-32 left-1/4 w-20 h-20 bg-gradient-nature rounded-full animate-bounce-soft" style={{ animationDelay: '2s' }} />
        <div className="absolute top-32 left-2/3 w-14 h-14 bg-primary/60 rounded-full animate-float" style={{ animationDelay: '1.5s' }} />
      </div>

      <main className="relative z-10 max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="flex items-center justify-between mb-6">
          <Link href="/home" className="inline-flex items-center gap-2 rounded-full px-4 py-2 bg-white/10 border border-white/20 text-white backdrop-blur-sm hover:bg-white/20 text-sm font-semibold">
            <ArrowLeft className="w-4 h-4" /> Home
          </Link>
          <button
            type="button"
            onClick={handleSave}
            disabled={saving}
            className={`rounded-full px-5 py-2 bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-semibold shadow-2xl ${saving ? 'opacity-70' : 'hover:scale-105'} transition-transform`}
          >
            {saving ? 'Saving…' : 'Save'}
          </button>
        </div>

        <div className="space-y-8">
          {/* Profile settings */}
          <section className="bg-card/90 backdrop-blur-sm border border-border rounded-3xl p-6 shadow-game">
            <h2 className="flex items-center gap-2 text-xl font-semibold mb-4"><User className="w-5 h-5" /> Profile</h2>
            <div className="grid sm:grid-cols-2 gap-4">
              <div>
                <label htmlFor="username" className="block text-sm text-muted-foreground mb-1">Username</label>
                <input
                  id="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full rounded-xl border border-border bg-background px-4 py-3 text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40"
                  placeholder="your_username"
                />
              </div>
              <div>
                <label className="block text-sm text-muted-foreground mb-1">Email</label>
                <input
                  disabled
                  value={user?.email || ''}
                  className="w-full rounded-xl border border-border bg-background/70 px-4 py-3 text-foreground/80"
                />
              </div>
            </div>
          </section>

          {/* Preferences */}
          <section className="bg-card/90 backdrop-blur-sm border border-border rounded-3xl p-6 shadow-game">
            <h2 className="flex items-center gap-2 text-xl font-semibold mb-4"><Music2 className="w-5 h-5" /> Preferences</h2>
            <div className="grid sm:grid-cols-2 gap-4">
              <label className="flex items-center justify-between gap-3 rounded-2xl border border-border bg-background px-4 py-3 cursor-pointer">
                <span className="text-sm">Calm background music</span>
                <input
                  type="checkbox"
                  checked={prefs.musicEnabled}
                  onChange={(e) => setPrefs({ musicEnabled: e.target.checked })}
                />
              </label>
              <label className="flex items-center justify-between gap-3 rounded-2xl border border-border bg-background px-4 py-3 cursor-pointer">
                <span className="text-sm">Sound effects</span>
                <input
                  type="checkbox"
                  checked={prefs.sfxEnabled}
                  onChange={(e) => setPrefs({ sfxEnabled: e.target.checked })}
                />
              </label>
            </div>
          </section>

          {/* Security */}
          <section className="bg-card/90 backdrop-blur-sm border border-border rounded-3xl p-6 shadow-game">
            <h2 className="flex items-center gap-2 text-xl font-semibold mb-4"><ShieldCheck className="w-5 h-5" /> Security</h2>
            <div className="grid sm:grid-cols-3 gap-4">
              <div>
                <label htmlFor="current" className="block text-sm text-muted-foreground mb-1">Current password</label>
                <input id="current" type="password" value={currentPw} onChange={(e) => setCurrentPw(e.target.value)} className="w-full rounded-xl border border-border bg-background px-4 py-3" />
              </div>
              <div>
                <label htmlFor="new" className="block text-sm text-muted-foreground mb-1">New password</label>
                <input id="new" type="password" value={newPw} onChange={(e) => setNewPw(e.target.value)} className="w-full rounded-xl border border-border bg-background px-4 py-3" />
              </div>
              <div>
                <label htmlFor="confirm" className="block text-sm text-muted-foreground mb-1">Confirm new password</label>
                <input id="confirm" type="password" value={confirmPw} onChange={(e) => setConfirmPw(e.target.value)} className="w-full rounded-xl border border-border bg-background px-4 py-3" />
              </div>
            </div>
            <div className="mt-4">
              <button
                type="button"
                onClick={handleChangePassword}
                className="rounded-full px-5 py-2 bg-white/10 border border-white/20 text-white hover:bg-white/20"
              >
                Update Password
              </button>
            </div>
          </section>
        </div>
      </main>
    </div>
  );
}
