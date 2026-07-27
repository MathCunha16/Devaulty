import React from "react";

interface CodeViewerProps {
  code: string;
  language?: string;
  maxHeight?: string;
  height?: string;
  className?: string;
}

export const CodeViewer: React.FC<CodeViewerProps> = ({
  code,
  language: _language,
  maxHeight,
  height = "380px",
  className = "",
}) => {
  void _language;
  const lines = (code || "").split("\n");

  return (
    <div
      className={`relative font-mono text-xs overflow-auto bg-[#0b0b0c] text-foreground rounded border border-border/50 select-text ${className}`}
      style={{ maxHeight, height }}
    >
      <div className="flex min-h-full">
        {/* Line numbers gutter */}
        <div
          className="py-2 px-3 text-right text-muted-foreground/40 bg-zinc-950/60 border-r border-border/30 select-none font-mono text-[11px] leading-5 shrink-0"
          aria-hidden="true"
        >
          {lines.map((_, i) => (
            <div key={i}>{i + 1}</div>
          ))}
        </div>

        {/* Code content */}
        <pre className="py-2 px-3 m-0 font-mono text-[12.5px] leading-5 text-zinc-200 whitespace-pre overflow-x-auto flex-1 font-mono">
          <code>{code}</code>
        </pre>
      </div>
    </div>
  );
};
