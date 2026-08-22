import React, { useRef, useEffect, useState, useCallback } from "react";
import { gsap } from "gsap";
import type { LucideIcon } from "lucide-react";

export interface BentoCardItem {
  id: string;
  title: string;
  description: string;
  label: string;
  icon?: LucideIcon;
  stat?: string | number;
  statLabel?: string;
  color?: string;
  onClick?: () => void;
}

export interface BentoProps {
  cards?: BentoCardItem[];
  textAutoHide?: boolean;
  enableStars?: boolean;
  enableSpotlight?: boolean;
  enableBorderGlow?: boolean;
  disableAnimations?: boolean;
  spotlightRadius?: number;
  particleCount?: number;
  enableTilt?: boolean;
  glowColor?: string;
  clickEffect?: boolean;
  enableMagnetism?: boolean;
  className?: string;
}

const DEFAULT_PARTICLE_COUNT = 8;
const DEFAULT_SPOTLIGHT_RADIUS = 280;
const DEFAULT_GLOW_COLOR = "16, 185, 129"; // Devaulty Emerald Primary
const MOBILE_BREAKPOINT = 768;

const createParticleElement = (
  x: number,
  y: number,
  color: string = DEFAULT_GLOW_COLOR
): HTMLDivElement => {
  const el = document.createElement("div");
  el.className = "particle";
  el.style.cssText = `
    position: absolute;
    width: 3px;
    height: 3px;
    border-radius: 50%;
    background: rgba(${color}, 1);
    box-shadow: 0 0 6px rgba(${color}, 0.7);
    pointer-events: none;
    z-index: 50;
    left: ${x}px;
    top: ${y}px;
  `;
  return el;
};

const calculateSpotlightValues = (radius: number) => ({
  proximity: radius * 0.5,
  fadeDistance: radius * 0.75,
});

const updateCardGlowProperties = (
  card: HTMLElement,
  mouseX: number,
  mouseY: number,
  glow: number,
  radius: number
) => {
  const rect = card.getBoundingClientRect();
  const relativeX = ((mouseX - rect.left) / rect.width) * 100;
  const relativeY = ((mouseY - rect.top) / rect.height) * 100;

  card.style.setProperty("--glow-x", `${relativeX}%`);
  card.style.setProperty("--glow-y", `${relativeY}%`);
  card.style.setProperty("--glow-intensity", glow.toString());
  card.style.setProperty("--glow-radius", `${radius}px`);
};

const ParticleCard: React.FC<{
  children: React.ReactNode;
  className?: string;
  disableAnimations?: boolean;
  style?: React.CSSProperties;
  particleCount?: number;
  glowColor?: string;
  enableTilt?: boolean;
  clickEffect?: boolean;
  enableMagnetism?: boolean;
  onClick?: () => void;
}> = ({
  children,
  className = "",
  disableAnimations = false,
  style,
  particleCount = DEFAULT_PARTICLE_COUNT,
  glowColor = DEFAULT_GLOW_COLOR,
  enableTilt = true,
  clickEffect = false,
  enableMagnetism = false,
  onClick,
}) => {
  const cardRef = useRef<HTMLDivElement>(null);
  const particlesRef = useRef<HTMLDivElement[]>([]);
  const timeoutsRef = useRef<ReturnType<typeof setTimeout>[]>([]);
  const isHoveredRef = useRef(false);
  const memoizedParticles = useRef<HTMLDivElement[]>([]);
  const particlesInitialized = useRef(false);
  const magnetismAnimationRef = useRef<gsap.core.Tween | null>(null);

  const initializeParticles = useCallback(() => {
    if (particlesInitialized.current || !cardRef.current) return;

    const { width, height } = cardRef.current.getBoundingClientRect();
    memoizedParticles.current = Array.from({ length: particleCount }, () =>
      createParticleElement(
        Math.random() * width,
        Math.random() * height,
        glowColor
      )
    );
    particlesInitialized.current = true;
  }, [particleCount, glowColor]);

  useEffect(() => {
    particlesInitialized.current = false;
    memoizedParticles.current = [];
  }, [particleCount, glowColor]);

  const clearAllParticles = useCallback(() => {
    timeoutsRef.current.forEach(clearTimeout);
    timeoutsRef.current = [];
    magnetismAnimationRef.current?.kill();

    particlesRef.current.forEach((particle) => {
      gsap.killTweensOf(particle);
      particle.remove();
    });
    particlesRef.current = [];
  }, []);

  const animateParticles = useCallback(() => {
    if (!cardRef.current || !isHoveredRef.current) return;

    if (!particlesInitialized.current) {
      initializeParticles();
    }

    memoizedParticles.current.forEach((particle, index) => {
      const timeoutId = setTimeout(() => {
        if (!isHoveredRef.current || !cardRef.current) return;

        const clone = particle.cloneNode(true) as HTMLDivElement;
        cardRef.current.appendChild(clone);
        particlesRef.current.push(clone);

        gsap.fromTo(
          clone,
          { scale: 0, opacity: 0 },
          { scale: 1, opacity: 1, duration: 0.3, ease: "back.out(1.7)" }
        );

        gsap.to(clone, {
          x: (Math.random() - 0.5) * 60,
          y: (Math.random() - 0.5) * 60,
          rotation: Math.random() * 360,
          duration: 1.5 + Math.random() * 1.5,
          ease: "none",
          repeat: -1,
          yoyo: true,
        });

        gsap.to(clone, {
          opacity: 0.2,
          duration: 1.2,
          ease: "power2.inOut",
          repeat: -1,
          yoyo: true,
        });
      }, index * 80);

      timeoutsRef.current.push(timeoutId);
    });
  }, [initializeParticles]);

  useEffect(() => {
    if (disableAnimations || !cardRef.current) return;

    const element = cardRef.current;

    const handleMouseEnter = () => {
      isHoveredRef.current = true;
      animateParticles();

      if (enableTilt) {
        gsap.to(element, {
          rotateX: 3,
          rotateY: 3,
          duration: 0.3,
          ease: "power2.out",
          transformPerspective: 1000,
        });
      }
    };

    const handleMouseLeave = () => {
      isHoveredRef.current = false;
      clearAllParticles();

      if (enableTilt) {
        gsap.to(element, {
          rotateX: 0,
          rotateY: 0,
          duration: 0.3,
          ease: "power2.out",
        });
      }

      if (enableMagnetism) {
        gsap.to(element, {
          x: 0,
          y: 0,
          duration: 0.3,
          ease: "power2.out",
        });
      }
    };

    const handleMouseMove = (e: MouseEvent) => {
      if (!enableTilt && !enableMagnetism) return;

      const rect = element.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;
      const centerX = rect.width / 2;
      const centerY = rect.height / 2;

      if (enableTilt) {
        const rotateX = ((y - centerY) / centerY) * -6;
        const rotateY = ((x - centerX) / centerX) * 6;

        gsap.to(element, {
          rotateX,
          rotateY,
          duration: 0.1,
          ease: "power2.out",
          transformPerspective: 1000,
        });
      }

      if (enableMagnetism) {
        const magnetX = (x - centerX) * 0.03;
        const magnetY = (y - centerY) * 0.03;

        magnetismAnimationRef.current = gsap.to(element, {
          x: magnetX,
          y: magnetY,
          duration: 0.3,
          ease: "power2.out",
        });
      }
    };

    const handleClick = (e: MouseEvent) => {
      if (!clickEffect) return;

      const rect = element.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;

      const maxDistance = Math.max(
        Math.hypot(x, y),
        Math.hypot(x - rect.width, y),
        Math.hypot(x, y - rect.height),
        Math.hypot(x - rect.width, y - rect.height)
      );

      const ripple = document.createElement("div");
      ripple.style.cssText = `
        position: absolute;
        width: ${maxDistance * 2}px;
        height: ${maxDistance * 2}px;
        border-radius: 50%;
        background: radial-gradient(circle, rgba(${glowColor}, 0.35) 0%, rgba(${glowColor}, 0.15) 30%, transparent 70%);
        left: ${x - maxDistance}px;
        top: ${y - maxDistance}px;
        pointer-events: none;
        z-index: 100;
      `;

      element.appendChild(ripple);

      gsap.fromTo(
        ripple,
        { scale: 0, opacity: 1 },
        {
          scale: 1,
          opacity: 0,
          duration: 0.6,
          ease: "power2.out",
          onComplete: () => ripple.remove(),
        }
      );
    };

    element.addEventListener("mouseenter", handleMouseEnter);
    element.addEventListener("mouseleave", handleMouseLeave);
    element.addEventListener("mousemove", handleMouseMove);
    element.addEventListener("click", handleClick);

    return () => {
      isHoveredRef.current = false;
      element.removeEventListener("mouseenter", handleMouseEnter);
      element.removeEventListener("mouseleave", handleMouseLeave);
      element.removeEventListener("mousemove", handleMouseMove);
      element.removeEventListener("click", handleClick);
      clearAllParticles();
      gsap.killTweensOf(element);
    };
  }, [
    animateParticles,
    clearAllParticles,
    disableAnimations,
    enableTilt,
    enableMagnetism,
    clickEffect,
    glowColor,
  ]);

  return (
    <div
      ref={cardRef}
      onClick={onClick}
      className={`${className} relative overflow-hidden cursor-pointer`}
      style={{ ...style, position: "relative", overflow: "hidden" }}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          onClick?.();
        }
      }}
    >
      {children}
    </div>
  );
};

const GlobalSpotlight: React.FC<{
  gridRef: React.RefObject<HTMLDivElement | null>;
  disableAnimations?: boolean;
  enabled?: boolean;
  spotlightRadius?: number;
  glowColor?: string;
}> = ({
  gridRef,
  disableAnimations = false,
  enabled = true,
  spotlightRadius = DEFAULT_SPOTLIGHT_RADIUS,
  glowColor = DEFAULT_GLOW_COLOR,
}) => {
  const spotlightRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (disableAnimations || !gridRef?.current || !enabled) return;

    const spotlight = document.createElement("div");
    spotlight.className = "global-spotlight";
    spotlight.style.cssText = `
      position: fixed;
      width: 600px;
      height: 600px;
      border-radius: 50%;
      pointer-events: none;
      background: radial-gradient(circle,
        rgba(${glowColor}, 0.12) 0%,
        rgba(${glowColor}, 0.06) 20%,
        rgba(${glowColor}, 0.02) 40%,
        transparent 70%
      );
      z-index: 40;
      opacity: 0;
      transform: translate(-50%, -50%);
      mix-blend-mode: screen;
    `;
    document.body.appendChild(spotlight);
    spotlightRef.current = spotlight;

    let rafId: number | null = null;
    let latestMouse: { clientX: number; clientY: number } | null = null;

    const processMouseMove = () => {
      rafId = null;
      if (!latestMouse || !spotlightRef.current || !gridRef.current) return;

      const { clientX, clientY } = latestMouse;
      const section = gridRef.current.closest(".bento-section");
      const rect = section?.getBoundingClientRect();
      const mouseInside =
        rect &&
        clientX >= rect.left &&
        clientX <= rect.right &&
        clientY >= rect.top &&
        clientY <= rect.bottom;

      const cards = gridRef.current.querySelectorAll(".bento-card");

      if (!mouseInside) {
        gsap.to(spotlightRef.current, {
          opacity: 0,
          duration: 0.25,
          ease: "power2.out",
        });
        cards.forEach((card) => {
          (card as HTMLElement).style.setProperty("--glow-intensity", "0");
        });
        return;
      }

      const { proximity, fadeDistance } = calculateSpotlightValues(spotlightRadius);
      let minDistance = Infinity;

      cards.forEach((card) => {
        const cardElement = card as HTMLElement;
        const cardRect = cardElement.getBoundingClientRect();
        const centerX = cardRect.left + cardRect.width / 2;
        const centerY = cardRect.top + cardRect.height / 2;
        const distance =
          Math.hypot(clientX - centerX, clientY - centerY) -
          Math.max(cardRect.width, cardRect.height) / 2;
        const effectiveDistance = Math.max(0, distance);

        minDistance = Math.min(minDistance, effectiveDistance);

        let glowIntensity = 0;
        if (effectiveDistance <= proximity) {
          glowIntensity = 1;
        } else if (effectiveDistance <= fadeDistance) {
          glowIntensity =
            (fadeDistance - effectiveDistance) / (fadeDistance - proximity);
        }

        updateCardGlowProperties(
          cardElement,
          clientX,
          clientY,
          glowIntensity,
          spotlightRadius
        );
      });

      gsap.to(spotlightRef.current, {
        left: clientX,
        top: clientY,
        duration: 0.08,
        ease: "power2.out",
      });

      const targetOpacity =
        minDistance <= proximity
          ? 0.7
          : minDistance <= fadeDistance
          ? ((fadeDistance - minDistance) / (fadeDistance - proximity)) * 0.7
          : 0;

      gsap.to(spotlightRef.current, {
        opacity: targetOpacity,
        duration: 0.2,
        ease: "power2.out",
      });
    };

    const handleMouseMove = (e: MouseEvent) => {
      latestMouse = { clientX: e.clientX, clientY: e.clientY };
      if (rafId === null) {
        rafId = requestAnimationFrame(processMouseMove);
      }
    };

    const handleMouseLeave = () => {
      if (rafId !== null) {
        cancelAnimationFrame(rafId);
        rafId = null;
      }
      gridRef.current?.querySelectorAll(".bento-card").forEach((card) => {
        (card as HTMLElement).style.setProperty("--glow-intensity", "0");
      });
      if (spotlightRef.current) {
        gsap.to(spotlightRef.current, {
          opacity: 0,
          duration: 0.25,
          ease: "power2.out",
        });
      }
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseleave", handleMouseLeave);

    return () => {
      if (rafId !== null) cancelAnimationFrame(rafId);
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseleave", handleMouseLeave);
      if (spotlightRef.current) {
        gsap.killTweensOf(spotlightRef.current);
        spotlightRef.current.remove();
      }
    };
  }, [gridRef, disableAnimations, enabled, spotlightRadius, glowColor]);

  return null;
};

const useMobileDetection = () => {
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth <= MOBILE_BREAKPOINT);
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  return isMobile;
};

export const MagicBento: React.FC<BentoProps> = ({
  cards = [],
  textAutoHide = false,
  enableStars = true,
  enableSpotlight = true,
  enableBorderGlow = true,
  disableAnimations = false,
  spotlightRadius = DEFAULT_SPOTLIGHT_RADIUS,
  particleCount = DEFAULT_PARTICLE_COUNT,
  enableTilt = true,
  glowColor = DEFAULT_GLOW_COLOR,
  clickEffect = true,
  enableMagnetism = true,
  className = "",
}) => {
  const gridRef = useRef<HTMLDivElement>(null);
  const isMobile = useMobileDetection();
  const shouldDisableAnimations = disableAnimations || isMobile;

  return (
    <div className={`bento-section w-full relative ${className}`}>
      <style>
        {`
          .bento-section {
            --glow-x: 50%;
            --glow-y: 50%;
            --glow-intensity: 0;
            --glow-radius: 180px;
            --glow-color: ${glowColor};
          }

          .bento-grid-layout {
            display: grid;
            grid-template-columns: 1fr;
            gap: 0.75rem;
            width: 100%;
          }

          @media (min-width: 640px) {
            .bento-grid-layout {
              grid-template-columns: repeat(2, 1fr);
            }
          }

          @media (min-width: 1024px) {
            .bento-grid-layout {
              grid-template-columns: repeat(3, 1fr);
            }
            
            .bento-grid-layout .bento-card:nth-child(1) {
              grid-column: span 2;
            }
            
            .bento-grid-layout .bento-card:nth-child(2) {
              grid-column: span 1;
            }

            .bento-grid-layout .bento-card:nth-child(3) {
              grid-column: span 1;
            }

            .bento-grid-layout .bento-card:nth-child(4) {
              grid-column: span 2;
            }

            .bento-grid-layout .bento-card:nth-child(5) {
              grid-column: span 1;
            }

            .bento-grid-layout .bento-card:nth-child(6) {
              grid-column: span 2;
            }
          }

          .bento-card--border-glow::after {
            content: '';
            position: absolute;
            inset: 0;
            padding: 1.5px;
            background: radial-gradient(var(--glow-radius) circle at var(--glow-x) var(--glow-y),
                rgba(${glowColor}, calc(var(--glow-intensity) * 0.9)) 0%,
                rgba(${glowColor}, calc(var(--glow-intensity) * 0.4)) 35%,
                transparent 65%);
            border-radius: inherit;
            -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
            -webkit-mask-composite: xor;
            mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
            mask-composite: exclude;
            pointer-events: none;
            opacity: 1;
            transition: opacity 0.2s ease;
            z-index: 10;
          }

          .bento-card:hover {
            box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2), 0 0 20px rgba(${glowColor}, 0.12);
            border-color: rgba(${glowColor}, 0.4);
          }

          .particle::before {
            content: '';
            position: absolute;
            top: -2px;
            left: -2px;
            right: -2px;
            bottom: -2px;
            background: rgba(${glowColor}, 0.3);
            border-radius: 50%;
            z-index: -1;
          }
        `}
      </style>

      {enableSpotlight && (
        <GlobalSpotlight
          gridRef={gridRef}
          disableAnimations={shouldDisableAnimations}
          enabled={enableSpotlight}
          spotlightRadius={spotlightRadius}
          glowColor={glowColor}
        />
      )}

      <div className="bento-grid-layout" ref={gridRef}>
        {cards.map((card) => {
          const Icon = card.icon;
          const baseClassName = `bento-card group flex flex-col justify-between relative min-h-[120px] w-full p-3.5 rounded-xl border border-border bg-card/80 backdrop-blur-sm transition-all duration-300 ease-out hover:-translate-y-0.5 ${
            enableBorderGlow ? "bento-card--border-glow" : ""
          }`;

          const cardStyle = {
            "--glow-x": "50%",
            "--glow-y": "50%",
            "--glow-intensity": "0",
            "--glow-radius": "180px",
          } as React.CSSProperties;

          const cardContent = (
            <>
              {/* Header */}
              <div className="flex items-start justify-between gap-2 relative z-20">
                <div className="flex items-center gap-2">
                  {Icon && (
                    <div
                      className="w-7 h-7 rounded-md flex items-center justify-center bg-secondary/80 border border-border/80 transition-transform group-hover:scale-105 shrink-0"
                      style={{
                        color: card.color || `rgb(${glowColor})`,
                      }}
                    >
                      <Icon size={14} />
                    </div>
                  )}
                  <span className="text-[10px] font-mono uppercase tracking-wider text-muted-foreground group-hover:text-foreground transition-colors font-medium">
                    {card.label}
                  </span>
                </div>

                {card.stat !== undefined && (
                  <div className="flex flex-col items-end">
                    <span
                      className="text-base font-bold font-mono tracking-tight leading-none"
                      style={{ color: card.color || `rgb(${glowColor})` }}
                    >
                      {card.stat}
                    </span>
                    {card.statLabel && (
                      <span className="text-[9px] text-muted-foreground uppercase font-mono mt-0.5">
                        {card.statLabel}
                      </span>
                    )}
                  </div>
                )}
              </div>

              {/* Content / Description */}
              <div className="flex flex-col relative z-20 mt-2">
                <h3
                  className={`font-semibold text-xs text-foreground mb-0.5 group-hover:text-primary transition-colors flex items-center gap-1.5 ${
                    textAutoHide ? "truncate" : ""
                  }`}
                >
                  <span>{card.title}</span>
                </h3>
                <p
                  className={`text-[11px] text-muted-foreground leading-relaxed ${
                    textAutoHide ? "line-clamp-1" : "line-clamp-2"
                  }`}
                >
                  {card.description}
                </p>
              </div>

              {/* Subtle background gradient */}
              <div
                className="absolute inset-0 opacity-0 group-hover:opacity-10 pointer-events-none transition-opacity duration-300 rounded-xl"
                style={{
                  background: `radial-gradient(circle at center, rgba(${glowColor}, 0.4) 0%, transparent 70%)`,
                }}
              />
            </>
          );

          if (enableStars) {
            return (
              <ParticleCard
                key={card.id}
                className={baseClassName}
                style={cardStyle}
                disableAnimations={shouldDisableAnimations}
                particleCount={particleCount}
                glowColor={glowColor}
                enableTilt={enableTilt}
                clickEffect={clickEffect}
                enableMagnetism={enableMagnetism}
                onClick={card.onClick}
              >
                {cardContent}
              </ParticleCard>
            );
          }

          return (
            <div
              key={card.id}
              className={baseClassName}
              style={cardStyle}
              onClick={card.onClick}
              role="button"
              tabIndex={0}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                  e.preventDefault();
                  card.onClick?.();
                }
              }}
            >
              {cardContent}
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default MagicBento;
