import React from "react";
import * as Icons from "lucide-react";

export const DrawingCanvasSkeleton: React.FC = () => {
  return (
    <div className="flex-1 w-full h-full flex flex-col items-center justify-center bg-card/60 relative overflow-hidden select-none">
      {/* Background grid shimmer */}
      <div
        className="absolute inset-0 opacity-[0.03] pointer-events-none"
        style={{
          backgroundImage:
            "radial-gradient(circle, currentColor 1px, transparent 1px)",
          backgroundSize: "24px 24px",
        }}
      />
      <div className="flex flex-col items-center gap-3 z-10">
        <div className="w-10 h-10 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary animate-pulse">
          <Icons.Shapes size={22} className="animate-spin text-primary" style={{ animationDuration: '3s' }} />
        </div>
        <div className="flex flex-col items-center gap-1 text-center">
          <span className="text-xs font-medium text-foreground/80">
            Loading Excalidraw canvas...
          </span>
          <span className="text-[11px] text-muted-foreground">
            Preparing hand-drawn graphics engine
          </span>
        </div>
      </div>
    </div>
  );
};
