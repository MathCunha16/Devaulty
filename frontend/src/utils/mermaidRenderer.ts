type MermaidTheme = "dark" | "default";

const MAX_SVG_CACHE_SIZE = 80;
const svgCache = new Map<string, string>();

let mermaidPromise: Promise<typeof import("mermaid").default> | null = null;
let initializedTheme: MermaidTheme | null = null;
let renderCounter = 0;

/**
 * Lazily loads the Mermaid library only when a rendered Markdown block
 * actually contains a ```mermaid code fence.
 */
async function getMermaid(theme: MermaidTheme) {
  if (!mermaidPromise) {
    mermaidPromise = import("mermaid").then((mod) => mod.default);
  }
  const mermaid = await mermaidPromise;

  if (initializedTheme !== theme) {
    mermaid.initialize({
      startOnLoad: false,
      securityLevel: "strict",
      theme,
      fontFamily: "var(--font-sans, system-ui, sans-serif)",
    });
    initializedTheme = theme;
  }

  return mermaid;
}

function getCachedSvg(key: string): string | undefined {
  const hit = svgCache.get(key);
  if (hit !== undefined) {
    // Refresh LRU order
    svgCache.delete(key);
    svgCache.set(key, hit);
  }
  return hit;
}

function setCachedSvg(key: string, svg: string): void {
  if (svgCache.has(key)) {
    svgCache.delete(key);
  } else if (svgCache.size >= MAX_SVG_CACHE_SIZE) {
    const oldestKey = svgCache.keys().next().value;
    if (oldestKey !== undefined) {
      svgCache.delete(oldestKey);
    }
  }
  svgCache.set(key, svg);
}

function buildDiagramWrapper(svg: string, source: string): HTMLElement {
  const wrapper = document.createElement("div");
  wrapper.className = "devaulty-mermaid-wrapper";

  const toolbar = document.createElement("div");
  toolbar.className = "devaulty-mermaid-toolbar";

  const badge = document.createElement("span");
  badge.className = "devaulty-mermaid-badge";
  badge.textContent = "MERMAID";

  const copyBtn = document.createElement("button");
  copyBtn.type = "button";
  copyBtn.className = "devaulty-mermaid-copy-btn";
  copyBtn.textContent = "Copy source";
  copyBtn.title = "Copy Mermaid source";
  copyBtn.addEventListener("click", (e) => {
    e.stopPropagation();
    void navigator.clipboard?.writeText(source).then(() => {
      copyBtn.textContent = "Copied!";
      setTimeout(() => {
        copyBtn.textContent = "Copy source";
      }, 1500);
    });
  });

  toolbar.appendChild(badge);
  toolbar.appendChild(copyBtn);

  const canvas = document.createElement("div");
  canvas.className = "devaulty-mermaid-canvas";
  canvas.innerHTML = svg;

  wrapper.appendChild(toolbar);
  wrapper.appendChild(canvas);
  return wrapper;
}

/**
 * Hydrates all `pre > code.language-mermaid` blocks inside `container` into
 * themed, cached SVG diagrams. Returns immediately with zero overhead if no
 * Mermaid blocks are present.
 */
export async function hydrateMermaidDiagrams(
  container: HTMLElement,
  themeMode: "dark" | "light",
  signal?: AbortSignal
): Promise<void> {
  const codeNodes = Array.from(
    container.querySelectorAll<HTMLElement>("pre > code.language-mermaid")
  );
  if (codeNodes.length === 0) return;

  const targetTheme: MermaidTheme = themeMode === "dark" ? "dark" : "default";

  // Check if all blocks are already cached to avoid even loading/awaiting when cached
  let mermaidInstance: Awaited<ReturnType<typeof getMermaid>> | null = null;

  for (const codeEl of codeNodes) {
    if (signal?.aborted) return;

    const preEl = codeEl.parentElement;
    if (!preEl || !preEl.parentElement) continue;

    const source = (codeEl.textContent || "").trim();
    if (!source) continue;

    const cacheKey = `${targetTheme}:${source}`;
    let svg = getCachedSvg(cacheKey);

    if (!svg) {
      if (!mermaidInstance) {
        mermaidInstance = await getMermaid(targetTheme);
        if (signal?.aborted) return;
      }

      const uniqueId = `devaulty-mermaid-${++renderCounter}`;
      try {
        const rendered = await mermaidInstance.render(uniqueId, source);
        if (signal?.aborted) return;
        svg = rendered.svg;
        setCachedSvg(cacheKey, svg);
      } catch (err) {
        if (signal?.aborted) return;
        // Clean up any temporary error node Mermaid may have appended to body
        document.getElementById(uniqueId)?.remove();
        document.getElementById(`d${uniqueId}`)?.remove();

        if (!preEl.previousElementSibling?.classList.contains("devaulty-mermaid-error")) {
          const errorBanner = document.createElement("div");
          errorBanner.className = "devaulty-mermaid-error";
          const msg = err instanceof Error ? err.message : "Invalid Mermaid syntax";
          errorBanner.textContent = `Mermaid syntax error: ${msg.split("\n")[0]}`;
          preEl.parentElement.insertBefore(errorBanner, preEl);
        }
        continue;
      }
    }

    if (signal?.aborted) return;
    const wrapper = buildDiagramWrapper(svg, source);
    preEl.replaceWith(wrapper);
  }
}
