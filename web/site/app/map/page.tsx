"use client";

import React, { useEffect, useMemo, useRef, useState } from "react";
import dynamic from "next/dynamic";
import Link from "next/link";
import { ArrowLeft, MapPin, X, Egg as EggIcon, Hammer, Wrench, Shield, Crosshair } from "lucide-react";
import Image from "next/image";
import { useAuthStore } from "../../lib/store";
// MapLibre styles
import "maplibre-gl/dist/maplibre-gl.css";

// Deck.gl
import DeckGL from "@deck.gl/react";
import { ScatterplotLayer } from "@deck.gl/layers";

// React Map GL using MapLibre entrypoint (no mapbox-gl)
import Map from "react-map-gl/maplibre";
import maplibregl from "maplibre-gl";

export default function MapPage() {
  const [drawerOpen, setDrawerOpen] = useState(true);
  const [mobileHeight, setMobileHeight] = useState<number>(420); // px
  const dragStartY = useRef<number | null>(null);
  const dragStartHeight = useRef<number>(0);
  const [userPos, setUserPos] = useState<{ latitude: number; longitude: number } | null>(null);
  const xp = useAuthStore((s) => s.user?.xp ?? 0);
  const addEgg = useAuthStore((s) => s.addEgg);
  const eggs = useAuthStore((s) => s.eggs);
  const user = useAuthStore((s) => s.user);
  const [showDropModal, setShowDropModal] = useState(false);
  const [eggType, setEggType] = useState<"bunny" | "dragon" | null>(null);
  const [activePanel, setActivePanel] = useState<"eggs" | "inventory">("eggs");

  // Hide global header logo on map page
  useEffect(() => {
    const el = document.getElementById("global-logo");
    if (el) el.style.display = "none";
    return () => {
      if (el) el.style.display = "";
    };
  }, []);

  const initialViewState = {
    longitude: -122.4194,
    latitude: 37.7749,
    zoom: 11,
    pitch: 0,
    bearing: 0,
  } as const;
  const [viewState, setViewState] = useState<any>(initialViewState);

  // Get initial position from onboarding (sessionStorage), then watch live location
  useEffect(() => {
    try {
      const raw = sessionStorage.getItem("user_location");
      if (raw) {
        const loc = JSON.parse(raw);
        if (loc?.latitude && loc?.longitude) {
          setUserPos({ latitude: loc.latitude, longitude: loc.longitude });
          setViewState((vs: any) => ({ ...vs, latitude: loc.latitude, longitude: loc.longitude, zoom: 14 }));
        }
      }
    } catch {}

    if ("geolocation" in navigator) {
      const id = navigator.geolocation.watchPosition(
        (p) => {
          const { latitude, longitude } = p.coords;
          setUserPos({ latitude, longitude });
        },
        () => {},
        { enableHighAccuracy: true, maximumAge: 5000 }
      );
      return () => navigator.geolocation.clearWatch(id);
    }
  }, []);

  // Demo deck.gl layer (pastel dot) matching theme
  const layers = useMemo(() => {
    const base: any[] = [
      new ScatterplotLayer({
        id: "demo-eggs",
        data: [
          { position: [-122.4194, 37.7749], size: 80 },
          { position: [-122.3894, 37.7849], size: 60 },
        ],
        getPosition: (d: any) => d.position,
        getRadius: (d: any) => d.size,
        radiusUnits: "meters",
        getFillColor: [244, 114, 182, 160],
        getLineColor: [99, 102, 241, 180],
        lineWidthUnits: "pixels",
        lineWidthMinPixels: 1.5,
        pickable: true,
      }),
    ];
    if (userPos) {
      base.push(
        new ScatterplotLayer({
          id: "user-location",
          data: [{ position: [userPos.longitude, userPos.latitude], size: 50 }],
          getPosition: (d: any) => d.position,
          getRadius: (d: any) => d.size,
          radiusUnits: "meters",
          getFillColor: [147, 197, 253, 200], // sky-300
          getLineColor: [59, 130, 246, 220], // blue-500
          lineWidthUnits: "pixels",
          lineWidthMinPixels: 2,
          pickable: false,
          updateTriggers: {
            getPosition: [userPos.longitude, userPos.latitude],
            getRadius: [userPos.longitude, userPos.latitude],
          },
        })
      );
    }
    return base;
  }, [userPos]);

  return (
    <div className="h-screen w-screen bg-background text-foreground relative">
      {/* Back Home */}
      <div className="absolute top-4 left-4 z-30">
        <Link
          href="/"
          className="inline-flex items-center gap-2 rounded-full px-4 py-2 bg-white/40 border border-white/40 text-white backdrop-blur-sm hover:bg-white/50 transition-colors text-sm font-semibold shadow-soft"
        >
          <ArrowLeft className="w-4 h-4" /> Home
        </Link>
      </div>

      {/* Removed mobile drawer toggle per design */}

      {/* DeckGL + Map */}
      <DeckGL
        initialViewState={initialViewState}
        viewState={viewState}
        onViewStateChange={({ viewState: vs }: any) => setViewState(vs)}
        controller={true}
        layers={layers}
        style={{ position: "absolute"}}
      >
        <Map
          reuseMaps
          mapLib={maplibregl as any}
          mapStyle="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
        />
      </DeckGL>

      {/* XP Badge (top-right) */}
      <div className="absolute top-4 right-4 z-30 inline-flex items-center gap-2 rounded-full px-4 py-2 bg-card border border-border text-foreground shadow-soft">
        <span className="text-xs opacity-80">XP</span>
        <span className="font-semibold">{xp}</span>
      </div>

      {/* Floating Add/Drop Egg button with hover-expand */
      }
      <button
        className="group absolute right-4 bottom-24 md:bottom-8 z-30 inline-flex items-center rounded-full px-3 py-1 bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-semibold shadow-2xl transition-[padding,opacity] hover:pr-5 hover:pl-4"
        onClick={() => {
          // Placeholder: open a creation flow; for now, center on user
          setShowDropModal(true);
        }}
      >
        <MapPin className="w-6 h-6" />
        <span className="ml-2 overflow-hidden max-w-0 opacity-0 transition-all duration-200 ease-out group-hover:max-w-[120px] group-hover:opacity-100">
          Drop Egg
        </span>
      </button>

      {/* Center on user button */}
      <button
        className="absolute left-4 bottom-24 md:bottom-8 z-30 inline-flex items-center rounded-full p-3 bg-card border border-border text-foreground shadow-soft hover:bg-white/10"
        onClick={() => {
          if (userPos) {
            setViewState((vs: any) => ({ ...vs, latitude: userPos.latitude, longitude: userPos.longitude, zoom: 15 }));
          }
        }}
        aria-label="Center on me"
      >
        <Crosshair className="w-4 h-4" />
      </button>

      {/* Drop Egg Modal */}
      {showDropModal && (
        <div
          className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm flex items-center justify-center px-4"
          onClick={() => setShowDropModal(false)}
        >
          <div
            className="relative w-full max-w-md bg-card border border-border rounded-2xl shadow-game p-6"
            onClick={(e) => e.stopPropagation()}
          >
            <button
              className="absolute top-3 right-3 rounded-md p-1.5 bg-white/10 border border-white/20 text-white hover:bg-white/20"
              onClick={() => setShowDropModal(false)}
              aria-label="Close"
            >
              <X className="w-4 h-4" />
            </button>
            <h3 className="text-xl font-semibold mb-4">Select Egg Type</h3>
            <div className="grid grid-cols-2 gap-4">
              <button
                type="button"
                onClick={() => setEggType("bunny")}
                className={`flex flex-col items-center gap-2 rounded-2xl border p-4 transition-colors ${eggType === "bunny" ? "border-pink-300 bg-white/10" : "border-border bg-background"}`}
              >
                <EggIcon className="w-6 h-6 text-rose-300" />
                <span className="text-sm">Bunny Egg</span>
              </button>
              <button
                type="button"
                onClick={() => setEggType("dragon")}
                className={`flex flex-col items-center gap-2 rounded-2xl border p-4 transition-colors ${eggType === "dragon" ? "border-amber-300 bg-white/10" : "border-border bg-background"}`}
              >
                <EggIcon className="w-6 h-6 text-amber-300" />
                <span className="text-sm">Dragon Egg</span>
              </button>
            </div>
            <div className="mt-6 text-right">
              <button
                type="button"
                disabled={!eggType || !userPos}
                onClick={() => {
                  if (!eggType || !userPos) return;
                  const id = String(Date.now());
                  addEgg({
                    id,
                    ownerId: (user?.id as string) || "me",
                    title: eggType === "bunny" ? "Bunny Egg" : "Dragon Egg",
                    description: undefined,
                    coords: { latitude: userPos.latitude, longitude: userPos.longitude },
                    createdAt: new Date().toISOString(),
                    color: eggType === "bunny" ? "rose" : "amber",
                    rarity: "common",
                  });
                  setViewState((vs: any) => ({ ...vs, latitude: userPos.latitude, longitude: userPos.longitude, zoom: 16 }));
                  setShowDropModal(false);
                }}
                className={`rounded-full px-5 py-2 bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-semibold shadow-2xl ${!eggType || !userPos ? "opacity-60" : "hover:scale-105"} transition-transform`}
              >
                Drop Here
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Desktop side panel removed per design */}

      {/* Mobile/Tablet drawer (bottom) - draggable height with panels */}
      <aside
        className={`md:hidden absolute left-0 right-0 bottom-0 px-4 pt-2 pb-4 bg-card border-t border-border z-20 transition-transform`}
        style={{ height: drawerOpen ? mobileHeight : 48 }}
      >
        {/* Drag handle */}
        <div
          className="mx-auto h-1.5 w-12 rounded-full bg-white/40 mb-3"
          onMouseDown={(e) => {
            dragStartY.current = e.clientY;
            dragStartHeight.current = mobileHeight;
            const onMove = (ev: MouseEvent) => {
              if (dragStartY.current !== null) {
                const dy = dragStartY.current - ev.clientY;
                const next = Math.max(120, Math.min(window.innerHeight * 0.85, dragStartHeight.current + dy));
                setMobileHeight(next);
              }
            };
            const onUp = () => {
              dragStartY.current = null;
              window.removeEventListener('mousemove', onMove);
              window.removeEventListener('mouseup', onUp);
            };
            window.addEventListener('mousemove', onMove);
            window.addEventListener('mouseup', onUp);
          }}
          onTouchStart={(e) => {
            const t = e.touches[0];
            dragStartY.current = t.clientY;
            dragStartHeight.current = mobileHeight;
            const onMove = (ev: TouchEvent) => {
              if (dragStartY.current !== null) {
                const dy = dragStartY.current - ev.touches[0].clientY;
                const next = Math.max(120, Math.min(window.innerHeight * 0.85, dragStartHeight.current + dy));
                setMobileHeight(next);
              }
            };
            const onUp = () => {
              dragStartY.current = null;
              window.removeEventListener('touchmove', onMove);
              window.removeEventListener('touchend', onUp);
            };
            window.addEventListener('touchmove', onMove);
            window.addEventListener('touchend', onUp);
          }}
        />
        <div className="flex items-center justify-between mb-2">
          <div className="inline-flex items-center gap-2 text-sm text-foreground">
            <Image src="/logo.png" alt="Eggsplore" width={16} height={16} className="rounded-sm border border-border" />
            <span className="text-xs">Nearby Eggs</span>
          </div>
        
        </div>
        {/* Segmented control for panels */}
        <div className="mt-1 relative rounded-full bg-background border border-border p-1">
          <div
            className={`absolute top-1 bottom-1 w-1/2 rounded-full bg-white/10 transition-transform ${activePanel === "inventory" ? "translate-x-full" : "translate-x-0"}`}
          />
          <div className="relative grid grid-cols-2 text-sm">
            <button
              className={`z-10 py-1.5 rounded-full ${activePanel === "eggs" ? "text-foreground" : "text-muted-foreground"}`}
              onClick={() => setActivePanel("eggs")}
            >
              My Eggs
            </button>
            <button
              className={`z-10 py-1.5 rounded-full ${activePanel === "inventory" ? "text-foreground" : "text-muted-foreground"}`}
              onClick={() => setActivePanel("inventory")}
            >
              Inventory
            </button>
          </div>
        </div>

        {/* Panel content */}
        {activePanel === "eggs" ? (
          <div className="mt-3 space-y-2 text-sm text-foreground overflow-y-auto">
            {eggs.length === 0 && (
              <div className="p-3 rounded-xl bg-background border border-border text-muted-foreground">
                No eggs yet — drop one to get started!
              </div>
            )}
            {eggs.map((e) => (
              <button
                key={e.id}
                className="w-full text-left p-3 rounded-xl bg-background border border-border hover:bg-white/10"
                onClick={() => setViewState((vs: any) => ({ ...vs, latitude: e.coords.latitude, longitude: e.coords.longitude, zoom: 16 }))}
              >
                <div className="flex items-center justify-between">
                  <span className="font-semibold">{e.title || "Egg"}</span>
                  <EggIcon className="w-4 h-4 text-rose-300" />
                </div>
                <div className="text-xs text-muted-foreground mt-1">{e.coords.latitude.toFixed(5)}, {e.coords.longitude.toFixed(5)}</div>
              </button>
            ))}
          </div>
        ) : (
          <div className="mt-3 grid grid-cols-2 gap-2 text-sm text-foreground">
            <div className="p-3 rounded-xl bg-background border border-border flex items-center gap-2"><Hammer className="w-4 h-4" /> Hammer</div>
            <div className="p-3 rounded-xl bg-background border border-border flex items-center gap-2"><Wrench className="w-4 h-4" /> Wrench</div>
            <div className="p-3 rounded-xl bg-background border border-border flex items-center gap-2"><Shield className="w-4 h-4" /> Shield</div>
            <div className="p-3 rounded-xl bg-background border border-border flex items-center gap-2"><EggIcon className="w-4 h-4" /> Nest</div>
          </div>
        )}
      </aside>
    </div>
  );
}
