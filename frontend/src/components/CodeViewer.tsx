import React, { useEffect, useMemo, useRef, useState } from "react";

interface CodeViewerProps {
  code: string;
  language?: string;
  maxHeight?: string;
  height?: string;
  className?: string;
  wrapLines?: boolean;
}

export const CodeViewer: React.FC<CodeViewerProps> = ({
  code,
  language: _language,
  maxHeight,
  height = "380px",
  className = "",
  wrapLines = false,
}) => {
  void _language;

  const lines = useMemo(() => (code || "").split("\n"), [code]);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [scrollState, setScrollState] = useState({ top: 0, height: 380, width: 0 });

  const lineNumbersText = useMemo(() => {
    let result = "";
    for (let i = 1; i <= lines.length; i++) {
      result += (i === 1 ? "" : "\n") + i;
    }
    return result;
  }, [lines.length]);

  const gutterWidth = useMemo(() => {
    return Math.max(2.5, String(lines.length).length * 0.65 + 1.2);
  }, [lines.length]);

  const wrappedLineMetrics = useMemo(() => {
    const gutterPixels = gutterWidth * 16;
    const contentWidth = Math.max(1, scrollState.width - gutterPixels - 24);
    const charactersPerLine = Math.max(1, Math.floor(contentWidth / 7.5));
    const rowHeights = lines.map((line) =>
      Math.max(1, Math.ceil(Math.max(line.length, 1) / charactersPerLine)) * 20
    );
    const offsets = [0];

    for (const rowHeight of rowHeights) {
      offsets.push(offsets[offsets.length - 1] + rowHeight);
    }

    return { rowHeights, offsets, totalHeight: offsets[offsets.length - 1] };
  }, [gutterWidth, lines, scrollState.width]);

  const visibleWrappedLines = useMemo(() => {
    if (!wrapLines || lines.length === 0) return { start: 0, end: 0 };

    const { offsets } = wrappedLineMetrics;
    const overscan = 400;
    const startOffset = Math.max(0, scrollState.top - overscan);
    const endOffset = scrollState.top + scrollState.height + overscan;
    let start = 0;

    while (start < lines.length && offsets[start + 1] < startOffset) start++;
    let end = start;
    while (end < lines.length && offsets[end] < endOffset) end++;

    return { start, end };
  }, [lines.length, scrollState.height, scrollState.top, wrapLines, wrappedLineMetrics]);

  useEffect(() => {
    const element = scrollContainerRef.current;
    if (!element) return;

    const updateSize = () => {
      setScrollState((previous) => ({
        ...previous,
        height: element.clientHeight,
        width: element.clientWidth,
      }));
    };

    updateSize();
    const observer = new ResizeObserver(updateSize);
    observer.observe(element);
    return () => observer.disconnect();
  }, []);

  return (
    <div
      ref={scrollContainerRef}
      className={`relative font-mono text-xs overflow-auto bg-[#0b0b0c] text-foreground rounded border border-border/50 select-text ${className}`}
      style={{ maxHeight, height }}
      onScroll={(event) => {
        const top = event.currentTarget.scrollTop;
        setScrollState((previous) => ({ ...previous, top }));
      }}
    >
      {wrapLines ? (
        <div
          className="relative w-full font-mono text-[12.5px] leading-5"
          style={{ height: wrappedLineMetrics.totalHeight + 16 }}
        >
          <div
            className="absolute inset-x-0 top-2"
            style={{ transform: `translateY(${wrappedLineMetrics.offsets[visibleWrappedLines.start]}px)` }}
          >
            {lines.slice(visibleWrappedLines.start, visibleWrappedLines.end).map((line, offset) => {
              const index = visibleWrappedLines.start + offset;
              return (
                <div key={index} className="flex min-w-0 hover:bg-zinc-900/40">
                  <span
                    className="text-right text-muted-foreground/40 bg-zinc-950/60 border-r border-border/30 select-none px-3 shrink-0 font-mono text-[12.5px] leading-5 tabular-nums"
                    style={{ minWidth: `${gutterWidth}rem` }}
                    aria-hidden="true"
                  >
                    {index + 1}
                  </span>
                  <span className="px-3 flex-1 min-w-0 font-mono text-[12.5px] leading-5 text-zinc-200 whitespace-pre-wrap break-words [overflow-wrap:anywhere]">
                    {line || "\u00A0"}
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      ) : (
        <div className="flex min-h-full">
          {/* Line numbers gutter — single DOM element for arbitrarily large content */}
          <pre
            className="py-2 px-3 m-0 text-right text-muted-foreground/40 bg-zinc-950/60 border-r border-border/30 select-none font-mono text-[12.5px] leading-5 shrink-0"
            aria-hidden="true"
          >
            {lineNumbersText}
          </pre>

          {/* Code content */}
          <pre className="py-2 px-3 m-0 font-mono text-[12.5px] leading-5 text-zinc-200 whitespace-pre overflow-x-auto flex-1 font-mono">
            <code>{code}</code>
          </pre>
        </div>
      )}
    </div>
  );
};
