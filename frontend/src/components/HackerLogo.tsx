import React, { useEffect, useState } from "react";
import { LogoDevaulty } from "./LogoDevaulty";
import styles from "./HackerLogo.module.css";

const STATIC_PARTICLES = [
  { id: 1, digit: "01", left: "12%", top: "50%", size: "11px", delay: "0s", duration: "2.8s" },
  { id: 2, digit: "00", left: "38%", top: "55%", size: "13px", delay: "0.9s", duration: "3.2s" },
  { id: 3, digit: "01", left: "64%", top: "48%", size: "10px", delay: "1.7s", duration: "2.6s" },
  { id: 4, digit: "00", left: "82%", top: "52%", size: "12px", delay: "0.4s", duration: "3.0s" },
  { id: 5, digit: "01", left: "25%", top: "58%", size: "12px", delay: "2.1s", duration: "3.4s" },
  { id: 6, digit: "00", left: "52%", top: "45%", size: "11px", delay: "1.3s", duration: "2.9s" },
];

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

  return (
    <div className={`${styles.container} ${className || ""}`}>
      {isVaultActive && <div className={styles.organicAura} aria-hidden="true" />}

      {isVaultActive && (
        <div className={styles.particleCanvas} aria-hidden="true">
          {STATIC_PARTICLES.map((p) => (
            <span
              key={p.id}
              className={styles.cssParticle}
              style={{
                left: p.left,
                top: p.top,
                fontSize: p.size,
                animationDelay: p.delay,
                animationDuration: p.duration,
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
