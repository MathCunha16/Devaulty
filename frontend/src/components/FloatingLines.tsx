import { useEffect, useRef } from 'react';
import {
  Clock,
  Mesh,
  OrthographicCamera,
  PlaneGeometry,
  Scene,
  ShaderMaterial,
  Vector2,
  Vector3,
  WebGLRenderer,
} from 'three';

import './FloatingLines.css';

// ── Shaders ──────────────────────────────────────────────────

const vertexShader = /* glsl */ `
precision highp float;
void main() {
  gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
}
`;

const fragmentShader = /* glsl */ `
precision highp float;

uniform float iTime;
uniform vec3  iResolution;
uniform float animationSpeed;
uniform bool  isDark;

uniform bool enableTop;
uniform bool enableMiddle;
uniform bool enableBottom;

uniform int topLineCount;
uniform int middleLineCount;
uniform int bottomLineCount;

uniform float topLineDistance;
uniform float middleLineDistance;
uniform float bottomLineDistance;

uniform vec3 topWavePosition;
uniform vec3 middleWavePosition;
uniform vec3 bottomWavePosition;

uniform vec2 iMouse;
uniform bool interactive;
uniform float bendRadius;
uniform float bendStrength;
uniform float bendInfluence;

uniform bool parallax;
uniform float parallaxStrength;
uniform vec2 parallaxOffset;

uniform vec3 lineGradient[8];
uniform int lineGradientCount;

const vec3 BLACK = vec3(0.0);
const vec3 EMERALD_DARK = vec3(6.0, 78.0, 59.0) / 255.0;
const vec3 EMERALD_LIGHT = vec3(16.0, 185.0, 129.0) / 255.0;

mat2 rotate(float r) {
  return mat2(cos(r), sin(r), -sin(r), cos(r));
}

vec3 background_color(vec2 uv) {
  vec3 col = vec3(0.0);
  float y = sin(uv.x - 0.2) * 0.3 - 0.1;
  float m = uv.y - y;
  col += mix(EMERALD_DARK, BLACK, smoothstep(0.0, 1.0, abs(m)));
  col += mix(EMERALD_LIGHT, BLACK, smoothstep(0.0, 1.0, abs(m - 0.8)));
  return col * 0.5;
}

vec3 getGradientStop(int idx) {
  if (idx == 0) return lineGradient[0];
  if (idx == 1) return lineGradient[1];
  if (idx == 2) return lineGradient[2];
  if (idx == 3) return lineGradient[3];
  if (idx == 4) return lineGradient[4];
  if (idx == 5) return lineGradient[5];
  if (idx == 6) return lineGradient[6];
  if (idx == 7) return lineGradient[7];
  return lineGradient[0];
}

vec3 getLineColor(float t, vec3 baseColor) {
  if (lineGradientCount <= 0) {
    return baseColor;
  }
  vec3 gradientColor;
  if (lineGradientCount == 1) {
    gradientColor = lineGradient[0];
  } else {
    float clampedT = clamp(t, 0.0, 0.9999);
    float scaled = clampedT * float(lineGradientCount - 1);
    int idx = int(floor(scaled));
    float f = fract(scaled);
    int idx2 = min(idx + 1, lineGradientCount - 1);
    vec3 c1 = getGradientStop(idx);
    vec3 c2 = getGradientStop(idx2);
    gradientColor = mix(c1, c2, f);
  }
  return gradientColor;
}

float wave(vec2 uv, float offset, vec2 screenUv, vec2 mouseUv, bool shouldBend) {
  float time = iTime * animationSpeed;
  float x_offset   = offset;
  float x_movement = time * 0.1;
  float amp        = sin(offset + time * 0.2) * 0.3;
  float y          = sin(uv.x + x_offset + x_movement) * amp;

  if (shouldBend) {
    vec2 d = screenUv - mouseUv;
    float influence = exp(-dot(d, d) * bendRadius);
    float bendOffset = (mouseUv.y - screenUv.y) * influence * bendStrength * bendInfluence;
    y += bendOffset;
  }

  float m = uv.y - y;
  return 0.0175 / max(abs(m) + 0.01, 1e-3) + 0.01;
}

void mainImage(out vec4 fragColor, in vec2 fragCoord) {
  vec2 baseUv = (2.0 * fragCoord - iResolution.xy) / iResolution.y;
  baseUv.y *= -1.0;

  if (parallax) {
    baseUv += parallaxOffset;
  }

  vec3 b = lineGradientCount > 0 ? vec3(0.0) : background_color(baseUv);

  vec2 mouseUv = vec2(0.0);
  if (interactive) {
    mouseUv = (2.0 * iMouse - iResolution.xy) / iResolution.y;
    mouseUv.y *= -1.0;
  }

  if (isDark) {
    // ── DARK MODE: Vibrant cyber-emerald luminous lines ──
    vec3 col = vec3(0.0);

    if (enableBottom) {
      for (int i = 0; i < bottomLineCount; ++i) {
        float fi = float(i);
        float t = fi / max(float(bottomLineCount - 1), 1.0);
        vec3 lineCol = getLineColor(t, b);
        float angle = bottomWavePosition.z * log(length(baseUv) + 1.0);
        vec2 ruv = baseUv * rotate(angle);
        col += lineCol * wave(
          ruv + vec2(bottomLineDistance * fi + bottomWavePosition.x, bottomWavePosition.y),
          1.5 + 0.2 * fi, baseUv, mouseUv, interactive
        ) * 0.22;
      }
    }

    if (enableMiddle) {
      for (int i = 0; i < middleLineCount; ++i) {
        float fi = float(i);
        float t = fi / max(float(middleLineCount - 1), 1.0);
        vec3 lineCol = getLineColor(t, b);
        float angle = middleWavePosition.z * log(length(baseUv) + 1.0);
        vec2 ruv = baseUv * rotate(angle);
        col += lineCol * wave(
          ruv + vec2(middleLineDistance * fi + middleWavePosition.x, middleWavePosition.y),
          2.0 + 0.15 * fi, baseUv, mouseUv, interactive
        ) * 0.85;
      }
    }

    if (enableTop) {
      for (int i = 0; i < topLineCount; ++i) {
        float fi = float(i);
        float t = fi / max(float(topLineCount - 1), 1.0);
        vec3 lineCol = getLineColor(t, b);
        float angle = topWavePosition.z * log(length(baseUv) + 1.0);
        vec2 ruv = baseUv * rotate(angle);
        ruv.x *= -1.0;
        col += lineCol * wave(
          ruv + vec2(topLineDistance * fi + topWavePosition.x, topWavePosition.y),
          1.0 + 0.2 * fi, baseUv, mouseUv, interactive
        ) * 0.12;
      }
    }

    float maxVal = max(col.r, max(col.g, col.b));
    float alpha = clamp(maxVal * 2.0, 0.0, 0.95);
    fragColor = vec4(col, alpha);
  } else {
    // ── LIGHT MODE: Silky, multi-layered jade & emerald ribbon waves ──
    vec3 waveColor = vec3(0.0);
    float totalAlpha = 0.0;

    if (enableBottom) {
      for (int i = 0; i < bottomLineCount; ++i) {
        float fi = float(i);
        float t = fi / max(float(bottomLineCount - 1), 1.0);
        vec3 lineCol = getLineColor(t, b);
        float angle = bottomWavePosition.z * log(length(baseUv) + 1.0);
        vec2 ruv = baseUv * rotate(angle);
        float w = wave(
          ruv + vec2(bottomLineDistance * fi + bottomWavePosition.x, bottomWavePosition.y),
          1.5 + 0.2 * fi, baseUv, mouseUv, interactive
        );
        vec3 strandColor = mix(lineCol, vec3(0.03, 0.40, 0.26), 0.3);
        float lineAlpha = clamp(w * 0.035, 0.0, 0.12);
        waveColor += strandColor * lineAlpha;
        totalAlpha += lineAlpha;
      }
    }

    if (enableMiddle) {
      for (int i = 0; i < middleLineCount; ++i) {
        float fi = float(i);
        float t = fi / max(float(middleLineCount - 1), 1.0);
        vec3 lineCol = getLineColor(t, b);
        float angle = middleWavePosition.z * log(length(baseUv) + 1.0);
        vec2 ruv = baseUv * rotate(angle);
        float w = wave(
          ruv + vec2(middleLineDistance * fi + middleWavePosition.x, middleWavePosition.y),
          2.0 + 0.15 * fi, baseUv, mouseUv, interactive
        );
        vec3 strandColor = mix(lineCol, vec3(0.02, 0.46, 0.30), 0.2);
        float lineAlpha = clamp(w * 0.075, 0.0, 0.22);
        waveColor += strandColor * lineAlpha;
        totalAlpha += lineAlpha;
      }
    }

    float finalAlpha = clamp(totalAlpha, 0.0, 0.42);
    vec3 finalColor = totalAlpha > 0.001 ? (waveColor / totalAlpha) : vec3(0.02, 0.45, 0.30);
    fragColor = vec4(finalColor, finalAlpha);
  }
}

void main() {
  vec4 color = vec4(0.0);
  mainImage(color, gl_FragCoord.xy);
  gl_FragColor = color;
}
`;

// ── Constants ────────────────────────────────────────────────

const MAX_GRADIENT_STOPS = 8;
const MAX_PIXEL_RATIO = 1.5;

const DEFAULT_TOP_POS = { x: 10.0, y: 0.5, rotate: -0.4 };
const DEFAULT_MID_POS = { x: 5.0, y: 0.0, rotate: 0.2 };
const DEFAULT_BOT_POS = { x: 2.0, y: -0.7, rotate: 0.4 };
const DEFAULT_ENABLED_WAVES: Array<'top' | 'middle' | 'bottom'> = ['top', 'middle', 'bottom'];
const DEFAULT_LINE_COUNT = [6];
const DEFAULT_LINE_DISTANCE = [5];

// ── Types ────────────────────────────────────────────────────

export type WavePosition = {
  x: number;
  y: number;
  rotate: number;
};

export type FloatingLinesProps = {
  theme?: 'light' | 'dark';
  linesGradient?: string[];
  enabledWaves?: Array<'top' | 'middle' | 'bottom'>;
  lineCount?: number | number[];
  lineDistance?: number | number[];
  topWavePosition?: WavePosition;
  middleWavePosition?: WavePosition;
  bottomWavePosition?: WavePosition;
  animationSpeed?: number;
  interactive?: boolean;
  bendRadius?: number;
  bendStrength?: number;
  parallax?: boolean;
  parallaxStrength?: number;
  mixBlendMode?: React.CSSProperties['mixBlendMode'];
};

// ── Helpers ──────────────────────────────────────────────────

function hexToVec3(hex: string): Vector3 {
  let value = hex.trim();
  if (value.startsWith('#')) value = value.slice(1);

  let r = 255, g = 255, b = 255;

  if (value.length === 3) {
    r = parseInt(value[0] + value[0], 16);
    g = parseInt(value[1] + value[1], 16);
    b = parseInt(value[2] + value[2], 16);
  } else if (value.length === 6) {
    r = parseInt(value.slice(0, 2), 16);
    g = parseInt(value.slice(2, 4), 16);
    b = parseInt(value.slice(4, 6), 16);
  }

  return new Vector3(r / 255, g / 255, b / 255);
}

// ── Component ────────────────────────────────────────────────

export default function FloatingLines({
  theme = 'dark',
  linesGradient,
  enabledWaves = DEFAULT_ENABLED_WAVES,
  lineCount = DEFAULT_LINE_COUNT,
  lineDistance = DEFAULT_LINE_DISTANCE,
  topWavePosition = DEFAULT_TOP_POS,
  middleWavePosition = DEFAULT_MID_POS,
  bottomWavePosition = DEFAULT_BOT_POS,
  animationSpeed = 0.6,
  interactive = false,
  bendRadius = 5.0,
  bendStrength = -0.5,
  parallax = false,
  parallaxStrength = 0.2,
  mixBlendMode = 'normal',
}: FloatingLinesProps) {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const uniformsRef = useRef<Record<string, { value: unknown }> | null>(null);

  // Helper resolvers
  const resolveLineCount = (waveType: 'top' | 'middle' | 'bottom'): number => {
    if (typeof lineCount === 'number') return lineCount;
    if (!enabledWaves.includes(waveType)) return 0;
    const index = enabledWaves.indexOf(waveType);
    return lineCount[index] ?? 6;
  };

  const resolveLineDistance = (waveType: 'top' | 'middle' | 'bottom'): number => {
    if (typeof lineDistance === 'number') return lineDistance;
    if (!enabledWaves.includes(waveType)) return 0.05;
    const index = enabledWaves.indexOf(waveType);
    return (lineDistance[index] ?? 5) * 0.01;
  };

  // ── Main Effect: Mount WebGL Scene ONCE ─────────────────────
  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    let active = true;
    let raf = 0;

    const scene = new Scene();
    const camera = new OrthographicCamera(-1, 1, 1, -1, 0, 1);
    camera.position.z = 1;

    const renderer = new WebGLRenderer({
      antialias: true,
      alpha: true,
      powerPreference: 'low-power',
    });
    renderer.setClearColor(0x000000, 0.0);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, MAX_PIXEL_RATIO));
    renderer.domElement.style.width = '100%';
    renderer.domElement.style.height = '100%';
    renderer.domElement.style.display = 'block';
    container.appendChild(renderer.domElement);

    const initialUniforms = {
      iTime: { value: 0 },
      iResolution: { value: new Vector3(1, 1, 1) },
      animationSpeed: { value: animationSpeed },
      isDark: { value: theme === 'dark' },

      enableTop: { value: enabledWaves.includes('top') },
      enableMiddle: { value: enabledWaves.includes('middle') },
      enableBottom: { value: enabledWaves.includes('bottom') },

      topLineCount: { value: enabledWaves.includes('top') ? resolveLineCount('top') : 0 },
      middleLineCount: { value: enabledWaves.includes('middle') ? resolveLineCount('middle') : 0 },
      bottomLineCount: { value: enabledWaves.includes('bottom') ? resolveLineCount('bottom') : 0 },

      topLineDistance: { value: resolveLineDistance('top') },
      middleLineDistance: { value: resolveLineDistance('middle') },
      bottomLineDistance: { value: resolveLineDistance('bottom') },

      topWavePosition: {
        value: new Vector3(topWavePosition.x, topWavePosition.y, topWavePosition.rotate),
      },
      middleWavePosition: {
        value: new Vector3(middleWavePosition.x, middleWavePosition.y, middleWavePosition.rotate),
      },
      bottomWavePosition: {
        value: new Vector3(bottomWavePosition.x, bottomWavePosition.y, bottomWavePosition.rotate),
      },

      iMouse: { value: new Vector2(-1000, -1000) },
      interactive: { value: interactive },
      bendRadius: { value: bendRadius },
      bendStrength: { value: bendStrength },
      bendInfluence: { value: 0 },

      parallax: { value: parallax },
      parallaxStrength: { value: parallaxStrength },
      parallaxOffset: { value: new Vector2(0, 0) },

      lineGradient: {
        value: Array.from({ length: MAX_GRADIENT_STOPS }, () => new Vector3(1, 1, 1)),
      },
      lineGradientCount: { value: 0 },
    };

    if (linesGradient && linesGradient.length > 0) {
      const stops = linesGradient.slice(0, MAX_GRADIENT_STOPS);
      initialUniforms.lineGradientCount.value = stops.length;
      stops.forEach((hex, i) => {
        const color = hexToVec3(hex);
        initialUniforms.lineGradient.value[i].set(color.x, color.y, color.z);
      });
    }

    uniformsRef.current = initialUniforms;

    const material = new ShaderMaterial({
      uniforms: initialUniforms,
      vertexShader,
      fragmentShader,
      transparent: true,
    });
    const geometry = new PlaneGeometry(2, 2);
    const mesh = new Mesh(geometry, material);
    scene.add(mesh);

    const clock = new Clock();

    const setSize = () => {
      if (!container || !active) return;
      const width = container.clientWidth || 1;
      const height = container.clientHeight || 1;
      renderer.setSize(width, height, false);
      const canvasWidth = renderer.domElement.width;
      const canvasHeight = renderer.domElement.height;
      initialUniforms.iResolution.value.set(canvasWidth, canvasHeight, 1);
    };
    setSize();

    const ro =
      typeof ResizeObserver !== 'undefined'
        ? new ResizeObserver(() => {
            if (!active) return;
            setSize();
          })
        : null;
    if (ro) ro.observe(container);

    const handleVisibility = () => {
      if (document.hidden) {
        cancelAnimationFrame(raf);
        clock.stop();
      } else if (active) {
        clock.start();
        raf = requestAnimationFrame(renderLoop);
      }
    };
    document.addEventListener('visibilitychange', handleVisibility);

    const renderLoop = () => {
      if (!active) return;
      initialUniforms.iTime.value = clock.getElapsedTime();
      renderer.render(scene, camera);
      raf = requestAnimationFrame(renderLoop);
    };
    renderLoop();

    return () => {
      active = false;
      cancelAnimationFrame(raf);
      document.removeEventListener('visibilitychange', handleVisibility);
      if (ro) ro.disconnect();

      uniformsRef.current = null;
      geometry.dispose();
      material.dispose();
      renderer.dispose();
      renderer.forceContextLoss();
      if (renderer.domElement.parentElement) {
        renderer.domElement.parentElement.removeChild(renderer.domElement);
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Mount ONCE to guarantee continuous, flicker-free background

  // ── Props Update Effect: Update Uniforms without remounting WebGL ──
  useEffect(() => {
    const uniforms = uniformsRef.current;
    if (!uniforms) return;

    uniforms.isDark.value = theme === 'dark';
    uniforms.animationSpeed.value = animationSpeed;
    uniforms.enableTop.value = enabledWaves.includes('top');
    uniforms.enableMiddle.value = enabledWaves.includes('middle');
    uniforms.enableBottom.value = enabledWaves.includes('bottom');

    uniforms.topLineCount.value = enabledWaves.includes('top') ? resolveLineCount('top') : 0;
    uniforms.middleLineCount.value = enabledWaves.includes('middle') ? resolveLineCount('middle') : 0;
    uniforms.bottomLineCount.value = enabledWaves.includes('bottom') ? resolveLineCount('bottom') : 0;

    uniforms.topLineDistance.value = resolveLineDistance('top');
    uniforms.middleLineDistance.value = resolveLineDistance('middle');
    uniforms.bottomLineDistance.value = resolveLineDistance('bottom');

    uniforms.interactive.value = interactive;
    uniforms.parallax.value = parallax;
    uniforms.bendRadius.value = bendRadius;
    uniforms.bendStrength.value = bendStrength;
    uniforms.parallaxStrength.value = parallaxStrength;

    if (topWavePosition && uniforms.topWavePosition.value instanceof Vector3) {
      uniforms.topWavePosition.value.set(topWavePosition.x, topWavePosition.y, topWavePosition.rotate);
    }
    if (middleWavePosition && uniforms.middleWavePosition.value instanceof Vector3) {
      uniforms.middleWavePosition.value.set(middleWavePosition.x, middleWavePosition.y, middleWavePosition.rotate);
    }
    if (bottomWavePosition && uniforms.bottomWavePosition.value instanceof Vector3) {
      uniforms.bottomWavePosition.value.set(bottomWavePosition.x, bottomWavePosition.y, bottomWavePosition.rotate);
    }

    if (linesGradient && linesGradient.length > 0) {
      const stops = linesGradient.slice(0, MAX_GRADIENT_STOPS);
      uniforms.lineGradientCount.value = stops.length;
      const arr = uniforms.lineGradient.value as Vector3[];
      stops.forEach((hex, i) => {
        const color = hexToVec3(hex);
        if (arr[i]) arr[i].set(color.x, color.y, color.z);
      });
    } else {
      uniforms.lineGradientCount.value = 0;
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    theme,
    linesGradient,
    enabledWaves,
    lineCount,
    lineDistance,
    animationSpeed,
    interactive,
    parallax,
    bendRadius,
    bendStrength,
    parallaxStrength,
    topWavePosition,
    middleWavePosition,
    bottomWavePosition,
  ]);

  return (
    <div
      ref={containerRef}
      className="floating-lines-container"
      style={{ mixBlendMode }}
    />
  );
}
