"use client";

import React from "react";
import Image from "next/image";
import Link from "next/link";

type Step = {
  number: string;
  title: string;
  description: string;
  color: string; // tailwind class e.g. bg-gradient-egg
  imageSrc: string; // path under public/
};

const steps: Step[] = [
  {
    number: "01",
    title: "Plant Your Egg",
    description:
      "Drop a virtual egg at any real location",
    color: "bg-gradient-egg",
    imageSrc: "/step-1-plant.jpg",
  },
  {
    number: "02",
    title: "Come Back & Hatch",
    description:
      "Return to your egg's location before it decays",
    color: "bg-gradient-nature",
    imageSrc: "/step-2-hatch.jpg",
  },
  {
    number: "03",
    title: "Watch It Grow",
    description:
      "Discover what adorable creature emerges!",
    color: "bg-gradient-adventure",
    imageSrc: "/step-3-grow.jpg",
  },
];

export default function HowItWorks({ onStart }: { onStart?: () => void }) {
  return (
    <section className="relative overflow-hidden py-24">
      {/* Floating background eggs to match hero */}
      <div className="pointer-events-none absolute inset-0 opacity-30 z-0">
        <div className="absolute top-20 left-10 w-16 h-16 bg-gradient-egg rounded-full animate-float" style={{ animationDelay: '0s' }} />
        <div className="absolute top-40 right-20 w-12 h-12 bg-gradient-adventure rounded-full animate-float" style={{ animationDelay: '1s' }} />
        <div className="absolute bottom-32 left-1/4 w-20 h-20 bg-gradient-nature rounded-full animate-bounce-soft" style={{ animationDelay: '2s' }} />
        <div className="absolute top-32 left-2/3 w-14 h-14 bg-primary/60 rounded-full animate-float" style={{ animationDelay: '1.5s' }} />
        <div className="absolute bottom-20 right-1/3 w-[4.5rem] h-[4.5rem] bg-secondary/70 rounded-full animate-bounce-soft" style={{ animationDelay: '0.5s' }} />
      </div>

      <div className="container mx-auto px-4 relative z-10">
        <div className="text-center mb-16">
          <h2 className="text-4xl lg:text-5xl font-bold text-foreground mb-6">
            How It <span className="bg-gradient-adventure bg-clip-text text-transparent">Works</span>
          </h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Ready to start your egg-citing adventure? Follow these simple steps to begin hatching!
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8 mb-16">
          {steps.map((step, index) => (
            <div key={step.number} className="group relative">
              <div className="bg-card rounded-3xl p-8 shadow-soft hover:shadow-game transition-all duration-500 transform hover:-translate-y-2 hover:scale-105">
                {/* Step number */}
                <div
                  className={`inline-flex items-center justify-center w-16 h-16 ${step.color} rounded-2xl text-2xl font-bold text-foreground mb-6 shadow-game`}
                >
                  {step.number}
                </div>

                {/* Step image */}
                <div className="mb-6 relative overflow-hidden rounded-2xl h-48">
                  <Image
                    src={step.imageSrc}
                    alt={step.title}
                    fill
                    sizes="(min-width: 1024px) 28vw, (min-width: 768px) 45vw, 100vw"
                    className="object-cover transform group-hover:scale-110 transition-transform duration-500"
                    priority={false}
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent" />
                </div>

                {/* Step content */}
                <h3 className="text-2xl font-bold text-foreground mb-4 group-hover:text-primary transition-colors">
                  {step.title}
                </h3>
                <p className="text-muted-foreground leading-relaxed">{step.description}</p>
              </div>

              {/* Connecting line (hidden on mobile) */}
              {index < steps.length - 1 && (
                <div className="hidden md:block absolute top-1/2 -right-4 w-8 h-0.5 bg-gradient-to-r from-primary to-accent transform -translate-y-1/2 z-10">
                  <div className="absolute right-0 top-1/2 transform -translate-y-1/2 w-2 h-2 bg-accent rounded-full animate-pulse" />
                </div>
              )}
            </div>
          ))}
        </div>

        <div className="text-center">
          <Link
            href="/auth"
            className={`inline-flex items-center gap-3 px-5 py-3 bg-gradient-to-r from-rose-400 via-pink-400 to-amber-300 text-white font-bold text-xl rounded-full shadow-2xl transform transition-all duration-300 hover:scale-105 hover:shadow-rose-400/25 `}
          >
            Start Your Quest
          </Link>
        </div>
      </div>
    </section>
  );
}
