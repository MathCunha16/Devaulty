import React, { useRef, useEffect, useCallback } from "react";
import type { LucideIcon } from "lucide-react";

export interface GooeyNavItem {
  id: string;
  label: string;
  icon?: LucideIcon;
  badgeCount?: number;
}

export interface GooeyNavProps {
  items: GooeyNavItem[];
  activeId: string;
  onChange: (id: string) => void;
  projectColor?: string;
  animationTime?: number;
  particleCount?: number;
  particleDistances?: [number, number];
  particleR?: number;
  timeVariance?: number;
  className?: string;
}

const noise = (n = 1) => n / 2 - Math.random() * n;

const getXY = (
  distance: number,
  pointIndex: number,
  totalPoints: number
): [number, number] => {
  const angle =
    ((360 + noise(8)) / totalPoints) * pointIndex * (Math.PI / 180);
  return [distance * Math.cos(angle), distance * Math.sin(angle)];
};

const createParticle = (
  i: number,
  t: number,
  d: [number, number],
  r: number,
  particleCount: number
) => {
  const rotate = noise(r / 10);
  return {
    start: getXY(d[0], particleCount - i, particleCount),
    end: getXY(d[1] + noise(7), particleCount - i, particleCount),
    time: t,
    scale: 1 + noise(0.2),
    rotate: rotate > 0 ? (rotate + r / 20) * 10 : (rotate - r / 20) * 10,
  };
};

export const GooeyNav: React.FC<GooeyNavProps> = ({
  items,
  activeId,
  onChange,
  projectColor = "var(--color-primary)",
  animationTime = 400,
  particleCount = 14,
  particleDistances = [75, 8],
  particleR = 90,
  timeVariance = 200,
  className = "",
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const navRef = useRef<HTMLUListElement>(null);
  const filterRef = useRef<HTMLSpanElement>(null);
  const textRef = useRef<HTMLSpanElement>(null);
  const activeTimeoutsRef = useRef<ReturnType<typeof setTimeout>[]>([]);

  const activeIndex = Math.max(
    0,
    items.findIndex((item) => item.id === activeId)
  );

  const clearAllParticles = useCallback(() => {
    activeTimeoutsRef.current.forEach(clearTimeout);
    activeTimeoutsRef.current = [];
    if (filterRef.current) {
      const particles = filterRef.current.querySelectorAll(".gooey-particle");
      particles.forEach((p) => p.remove());
    }
  }, []);

  const makeParticles = useCallback(
    (element: HTMLElement) => {
      const d: [number, number] = particleDistances;
      const r = particleR;
      const bubbleTime = animationTime * 2 + timeVariance;
      element.style.setProperty("--time", `${bubbleTime}ms`);

      for (let i = 0; i < particleCount; i++) {
        const t = animationTime * 2 + noise(timeVariance * 2);
        const p = createParticle(i, t, d, r, particleCount);
        element.classList.remove("active");

        const timeoutId = setTimeout(() => {
          const particle = document.createElement("span");
          const point = document.createElement("span");
          particle.classList.add("gooey-particle");
          particle.style.setProperty("--start-x", `${p.start[0]}px`);
          particle.style.setProperty("--start-y", `${p.start[1]}px`);
          particle.style.setProperty("--end-x", `${p.end[0]}px`);
          particle.style.setProperty("--end-y", `${p.end[1]}px`);
          particle.style.setProperty("--time", `${p.time}ms`);
          particle.style.setProperty("--scale", `${p.scale}`);
          particle.style.setProperty("--rotate", `${p.rotate}deg`);
          point.classList.add("gooey-point");
          particle.appendChild(point);
          element.appendChild(particle);

          requestAnimationFrame(() => {
            element.classList.add("active");
          });

          const cleanupTimeout = setTimeout(() => {
            particle.remove();
          }, t);
          activeTimeoutsRef.current.push(cleanupTimeout);
        }, 20);

        activeTimeoutsRef.current.push(timeoutId);
      }
    },
    [animationTime, particleCount, particleDistances, particleR, timeVariance]
  );

  const updateEffectPosition = useCallback((element: HTMLElement) => {
    if (!containerRef.current || !filterRef.current || !textRef.current) return;
    const containerRect = containerRef.current.getBoundingClientRect();
    const pos = element.getBoundingClientRect();
    const styles = {
      left: `${pos.x - containerRect.x}px`,
      top: `${pos.y - containerRect.y}px`,
      width: `${pos.width}px`,
      height: `${pos.height}px`,
    };
    Object.assign(filterRef.current.style, styles);
    Object.assign(textRef.current.style, styles);
  }, []);

  const handleItemClick = (
    e: React.MouseEvent<HTMLButtonElement>,
    item: GooeyNavItem
  ) => {
    if (activeId === item.id) return;
    onChange(item.id);

    const liEl = e.currentTarget.parentElement;
    if (liEl) {
      updateEffectPosition(liEl);
      clearAllParticles();
      if (textRef.current) {
        textRef.current.classList.remove("active");
        void textRef.current.offsetWidth;
        textRef.current.classList.add("active");
      }
      if (filterRef.current) {
        makeParticles(filterRef.current);
      }
    }
  };

  useEffect(() => {
    if (!navRef.current || !containerRef.current) return;
    const activeLi = navRef.current.querySelectorAll("li")[
      activeIndex
    ] as HTMLElement;
    if (activeLi) {
      updateEffectPosition(activeLi);
      textRef.current?.classList.add("active");
    }

    const resizeObserver = new ResizeObserver(() => {
      const currentActiveLi = navRef.current?.querySelectorAll("li")[
        activeIndex
      ] as HTMLElement;
      if (currentActiveLi) {
        updateEffectPosition(currentActiveLi);
      }
    });

    resizeObserver.observe(containerRef.current);

    return () => {
      resizeObserver.disconnect();
      clearAllParticles();
    };
  }, [activeIndex, updateEffectPosition, clearAllParticles]);

  return (
    <div
      className={`gooey-nav-wrapper relative inline-flex items-center select-none ${className}`}
      ref={containerRef}
      style={
        {
          "--active-accent": projectColor || "var(--color-primary)",
        } as React.CSSProperties
      }
    >
      <style>
        {`
          .gooey-nav-wrapper {
            position: relative;
          }

          .gooey-nav-container {
            background: color-mix(in srgb, var(--color-card) 70%, transparent);
            backdrop-filter: blur(12px);
            -webkit-backdrop-filter: blur(12px);
            border: 1px solid color-mix(in srgb, var(--color-border) 80%, transparent);
            border-radius: 9999px;
            padding: 3px;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
          }

          .gooey-effect {
            position: absolute;
            opacity: 1;
            pointer-events: none;
            display: grid;
            place-items: center;
            z-index: 1;
            transition: left 0.3s cubic-bezier(0.2, 0.8, 0.2, 1), width 0.3s cubic-bezier(0.2, 0.8, 0.2, 1);
          }

          .gooey-effect.filter {
            filter: blur(5px) contrast(40) blur(0);
            mix-blend-mode: normal;
          }

          .gooey-effect.filter::after {
            content: "";
            position: absolute;
            inset: 0;
            background: var(--active-accent);
            opacity: 0.2;
            border-radius: 9999px;
            transform: scale(0.9);
            transition: transform 0.2s ease, opacity 0.2s ease;
            z-index: -1;
          }

          .gooey-effect.active::after {
            animation: gooey-pill 0.35s cubic-bezier(0.2, 0.8, 0.2, 1) both;
          }

          @keyframes gooey-pill {
            0% {
              transform: scale(0.65);
              opacity: 0.1;
            }
            50% {
              transform: scale(1.15);
              opacity: 0.35;
            }
            100% {
              transform: scale(1);
              opacity: 0.22;
            }
          }

          .gooey-particle,
          .gooey-point {
            display: block;
            opacity: 0;
            width: 16px;
            height: 16px;
            border-radius: 9999px;
            transform-origin: center;
          }

          .gooey-particle {
            --time: 2.5s;
            position: absolute;
            top: calc(50% - 8px);
            left: calc(50% - 8px);
            animation: gooey-particle calc(var(--time)) ease 1 -150ms;
          }

          .gooey-point {
            background: var(--active-accent);
            opacity: 0.9;
            box-shadow: 0 0 10px var(--active-accent);
            animation: gooey-point calc(var(--time)) ease 1 -150ms;
          }

          @keyframes gooey-particle {
            0% {
              transform: rotate(0deg) translate(calc(var(--start-x)), calc(var(--start-y)));
              opacity: 0.9;
              animation-timing-function: cubic-bezier(0.55, 0, 1, 0.45);
            }
            60% {
              transform: rotate(calc(var(--rotate) * 0.6)) translate(calc(var(--end-x) * 1.3), calc(var(--end-y) * 1.3));
              opacity: 0.8;
              animation-timing-function: ease;
            }
            85% {
              transform: rotate(calc(var(--rotate) * 0.85)) translate(calc(var(--end-x)), calc(var(--end-y)));
              opacity: 0.5;
            }
            100% {
              transform: rotate(calc(var(--rotate) * 1.2)) translate(calc(var(--end-x) * 0.4), calc(var(--end-y) * 0.4));
              opacity: 0;
            }
          }

          @keyframes gooey-point {
            0% {
              transform: scale(0);
              opacity: 0;
            }
            30% {
              transform: scale(calc(var(--scale) * 0.4));
              opacity: 1;
            }
            70% {
              transform: scale(var(--scale));
              opacity: 0.8;
            }
            100% {
              transform: scale(0);
              opacity: 0;
            }
          }

          .gooey-nav-item {
            position: relative;
            z-index: 2;
            border-radius: 9999px;
            transition: color 0.2s ease;
          }

          .gooey-nav-item.active-item {
            color: var(--color-foreground);
            font-weight: 600;
          }

          .gooey-nav-item.active-item::after {
            content: "";
            position: absolute;
            inset: 0;
            border-radius: 9999px;
            background-color: var(--color-card);
            border: 1px solid var(--color-border);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12), 0 0 14px color-mix(in srgb, var(--active-accent) 25%, transparent);
            opacity: 1;
            z-index: -1;
            transition: all 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
          }
        `}
      </style>

      <nav className="gooey-nav-container flex relative">
        <ul
          ref={navRef}
          className="flex items-center gap-0.5 list-none p-0 m-0 relative z-[3]"
        >
          {items.map((item) => {
            const isActive = activeId === item.id;
            const Icon = item.icon;

            return (
              <li
                key={item.id}
                className={`gooey-nav-item ${isActive ? "active-item" : ""}`}
              >
                <button
                  type="button"
                  onClick={(e) => handleItemClick(e, item)}
                  className={`flex items-center gap-1.5 py-1 px-2.5 rounded-full text-xs cursor-pointer border-0 bg-transparent transition-all duration-200 outline-none ${
                    isActive
                      ? "text-foreground font-semibold scale-[1.02]"
                      : "text-muted-foreground hover:text-foreground hover:bg-card/30"
                  }`}
                  aria-label={item.label}
                  title={item.label}
                >
                  {Icon && (
                    <Icon
                      size={13}
                      style={{
                        color: isActive
                          ? projectColor || "var(--color-primary)"
                          : "currentColor",
                      }}
                    />
                  )}
                  <span>{item.label}</span>
                  {item.badgeCount !== undefined && item.badgeCount > 0 && (
                    <span className="flex items-center justify-center min-w-[15px] h-3.5 px-1 rounded-full text-[9px] font-mono font-bold bg-destructive/15 text-destructive border border-destructive/30">
                      {item.badgeCount}
                    </span>
                  )}
                </button>
              </li>
            );
          })}
        </ul>
      </nav>
      <span className="gooey-effect filter" ref={filterRef} />
      <span className="gooey-effect" ref={textRef} />
    </div>
  );
};

export default GooeyNav;
