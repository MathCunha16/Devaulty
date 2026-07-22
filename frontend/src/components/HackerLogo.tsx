import React, { useEffect, useState, useRef } from "react";
import { LogoDevaulty } from "./LogoDevaulty";
import styles from "./HackerLogo.module.css";

const BINARY_DIGITS = ["00", "01"];

interface Particle {
  id: number;
  digit: string;
  x: number;
  y: number;
  size: number;
  speedY: number;
  speedX: number;
  opacity: number;
}

interface HackerLogoProps {
  height?: number | string;
  width?: number | string;
  className?: string;
}

export const HackerLogo: React.FC<HackerLogoProps> = ({
  height = "100%",
  width = "auto",
  className,
}) => {
  const [isVaultActive, setIsVaultActive] = useState(
    () => typeof document !== "undefined" && document.documentElement.dataset.vaultActive === "true"
  );
  const [particles, setParticles] = useState<Particle[]>([]);
  const nextIdRef = useRef(0);

  useEffect(() => {
    const checkActive = () => {
      const active = document.documentElement.dataset.vaultActive === "true";
      setIsVaultActive(active);
    };

    checkActive();

    const observer = new MutationObserver(checkActive);
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["data-vault-active"],
    });

    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (!isVaultActive) return;

    const interval = setInterval(() => {
      setParticles((prev) => {
        const newParticle: Particle = {
          id: nextIdRef.current++,
          digit: BINARY_DIGITS[Math.floor(Math.random() * BINARY_DIGITS.length)],
          x: Math.random() * 85 + 5,
          y: Math.random() * 20 + 45,
          size: Math.floor(Math.random() * 5 + 10), // 10px to 14px font
          speedY: Math.random() * 0.7 + 0.5,
          speedX: (Math.random() - 0.5) * 0.6,
          opacity: 0.9,
        };

        const updated = prev
          .map((p) => ({
            ...p,
            y: p.y - p.speedY * 1.8,
            x: p.x + p.speedX * 1.2,
            opacity: p.opacity - 0.025,
          }))
          .filter((p) => p.opacity > 0 && p.y > -35);

        if (updated.length < 14) {
          updated.push(newParticle);
        }

        return updated;
      });
    }, 70);

    return () => {
      clearInterval(interval);
      setParticles([]);
    };
  }, [isVaultActive]);

  return (
    <div className={`${styles.container} ${className || ""}`}>
      {isVaultActive && <div className={styles.organicAura} aria-hidden="true" />}

      {isVaultActive && (
        <div className={styles.particleCanvas} aria-hidden="true">
          {particles.map((p) => (
            <span
              key={p.id}
              className={styles.particle}
              style={{
                left: `${p.x}%`,
                top: `${p.y}%`,
                fontSize: `${p.size}px`,
                fontFamily: "var(--font-mono)",
                fontWeight: 750,
                opacity: p.opacity,
              }}
            >
              {p.digit}
            </span>
          ))}
        </div>
      )}

      <LogoDevaulty
        height={height}
        width={width}
        className={`${styles.logoSvg} ${isVaultActive ? styles.greenLogo : ""}`}
      />
    </div>
  );
};
