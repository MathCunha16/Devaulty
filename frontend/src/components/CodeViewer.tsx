import React, { useMemo } from "react";

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
  
  const lineNumbersText = useMemo(() => {
    const count = (code || "").split("\n").length;
    let result = "";
    for (let i = 1; i <= count; i++) {
      result += (i === 1 ? "" : "\n") + i;
    }
    return result;
  }, [code]);

  return (
    <div
      className={`relative font-mono text-xs overflow-auto bg-[#0b0b0c] text-foreground rounded border border-border/50 select-text ${className}`}
      style={{ maxHeight, height }}
    >
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
    </div>
  );
};
