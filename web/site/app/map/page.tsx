"use client";

import React, { useEffect, useMemo, useRef, useState } from "react";
import dynamic from "next/dynamic";
import Link from "next/link";
import { ArrowLeft, MapPin, Menu, X } from "lucide-react";
import Image from "next/image";
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

  // Demo deck.gl layer (pastel dot) matching theme
  const layers = useMemo(() => {
    return [
      new ScatterplotLayer({
        id: "demo-eggs",
        data: [
          { position: [-122.4194, 37.7749], size: 80 },
          { position: [-122.3894, 37.7849], size: 60 },
        ],
        getPosition: (d: any) => d.position,
        getRadius: (d: any) => d.size,
        radiusUnits: "meters",
        getFillColor: [244, 114, 182, 160], // rose-400 pastel
        getLineColor: [99, 102, 241, 180], // indigo-500 outline
        lineWidthUnits: "pixels",
        lineWidthMinPixels: 1.5,
        pickable: true,
      }),
    ];
  }, []);

  return (
    <div className="h-screen w-screen bg-background text-foreground relative">
      {/* Back Home */}
      <div className="absolute top-4 left-4 z-30">
        <Link
          href="/"
          className="inline-flex items-center gap-2 rounded-full px-4 py-2 bg-white/10 border border-white/20 text-white backdrop-blur-sm hover:bg-white/20 transition-colors text-sm font-semibold shadow-soft"
        >
          <ArrowLeft className="w-4 h-4" /> Home
        </Link>
      </div>

      {/* Drawer toggle (mobile) */}
      <button
        onClick={() => setDrawerOpen((v) => !v)}
        className="md:hidden absolute bottom-6 right-4 z-30 inline-flex items-center gap-2 rounded-full px-4 py-2 bg-card border border-border text-foreground hover:bg-white/10 transition-colors text-sm font-semibold shadow-soft"
      >
        <Menu className="w-4 h-4" /> Panel
      </button>

      {/* DeckGL + Map */}
      <DeckGL
        initialViewState={initialViewState}
        controller={true}
        layers={layers}
        style={{ position: "absolute", inset: 0 }}
      >
        <Map
          reuseMaps
          mapLib={maplibregl as any}
          mapStyle="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
        />
      </DeckGL>

      {/* Desktop drawer (left side) */}
      <aside
        className={`hidden md:flex flex-col gap-4 absolute top-0 left-0 h-full w-80 p-4 bg-card border-r border-border z-20 transition-transform ${
          drawerOpen ? "translate-x-0" : "-translate-x-80"
        }`}
      >
        <div className="flex items-center justify-between">
          <div className="inline-flex items-center gap-2 text-sm text-foreground">
            <Image src="/logo.png" alt="Eggsplore" width={18} height={18} className="rounded-sm border border-border" />
            <span className="text-xs">Panel</span>
          </div>
          <button
            onClick={() => setDrawerOpen(false)}
            className="hidden md:inline-flex rounded-md p-1.5 bg-white/10 border border-white/20 text-white hover:bg-white/20"
            aria-label="Close panel"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
        <div className="space-y-3 text-sm text-foreground overflow-y-auto pr-2">
          <div className="p-3 rounded-xl bg-background border border-border">
            <div className="font-semibold mb-1">Golden Egg</div>
            <div className="text-xs text-muted-foreground">San Francisco • 120m away</div>
          </div>
          <div className="p-3 rounded-xl bg-background border border-border">
            <div className="font-semibold mb-1">Mystic Bunny Nest</div>
            <div className="text-xs text-muted-foreground">Waterfront • 350m away</div>
          </div>
        </div>
      </aside>

      {/* Mobile/Tablet drawer (bottom) - draggable height */}
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
          <button
            onClick={() => setDrawerOpen(false)}
            className="rounded-md p-1.5 bg-white/10 border border-white/20 text-white hover:bg-white/20"
            aria-label="Close panel"
          >
            <X className="w-4 h-4" />
          </button>
        </div>
        <div className="grid grid-cols-2 gap-3 text-sm text-foreground">
          <div className="p-3 rounded-xl bg-background border border-border">
            <div className="font-semibold mb-1">Golden Egg</div>
            <div className="text-xs text-muted-foreground">120m away</div>
          </div>
          <div className="p-3 rounded-xl bg-background border border-border">
            <div className="font-semibold mb-1">Mystic Bunny Nest</div>
            <div className="text-xs text-muted-foreground">350m away</div>
          </div>
        </div>
      </aside>
    </div>
  );
}
