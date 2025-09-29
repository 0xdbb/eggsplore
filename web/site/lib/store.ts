"use client";

import { create } from "zustand";
import type { Egg, Session, User } from "./types";

type StoreState = {
  user?: Partial<User> | null;
  session?: Session | null;
  eggs: Egg[];
  prefs: {
    musicEnabled: boolean;
    sfxEnabled: boolean;
  };
  setUser: (user: Partial<User> | null) => void;
  setSession: (session: Session | null) => void;
  setXP: (xp: number) => void;
  addXP: (delta: number) => void;
  addEgg: (egg: Egg) => void;
  setPrefs: (prefs: Partial<StoreState["prefs"]>) => void;
  clear: () => void;
};

export const useAuthStore = create<StoreState>((set) => ({
  user: null,
  session: null,
  eggs: [],
  prefs: { musicEnabled: false, sfxEnabled: true },
  setUser: (user) => set({ user }),
  setSession: (session) => set({ session }),
  setXP: (xp) => set((s) => ({ user: { ...(s.user || {}), xp } })),
  addXP: (delta) => set((s) => ({ user: { ...(s.user || {}), xp: Math.max(0, (s.user?.xp || 0) + delta) } })),
  addEgg: (egg) => set((s) => ({ eggs: [egg, ...s.eggs] })),
  setPrefs: (prefs) => set((s) => ({ prefs: { ...s.prefs, ...prefs } })),
  clear: () => set({ user: null, session: null, eggs: [], prefs: { musicEnabled: false, sfxEnabled: true } }),
}));
