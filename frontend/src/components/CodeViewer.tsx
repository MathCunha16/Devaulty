import React, { useMemo } from "react";

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

  return (
    <div
      className={`relative font-mono text-xs overflow-auto bg-[#0b0b0c] text-foreground rounded border border-border/50 select-text ${className}`}
      style={{ maxHeight, height }}
    >
      {wrapLines ? (
        <div className="w-full font-mono text-[12.5px] leading-5 py-2">
          {lines.map((line, index) => (
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
          ))}
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
