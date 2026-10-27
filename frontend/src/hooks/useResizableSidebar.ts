import { useState, useEffect, useCallback, useRef } from "react";

interface UseResizableSidebarOptions {
  storageKey: string;
  defaultWidth?: number;
  minWidth?: number;
  maxWidth?: number;
}

export function useResizableSidebar({
  storageKey,
  defaultWidth = 340,
  minWidth = 240,
  maxWidth = 560,
}: UseResizableSidebarOptions) {
  const [sidebarWidth, setSidebarWidth] = useState<number>(() => {
    try {
      const saved = localStorage.getItem(`devaulty:sidebar:${storageKey}:width`);
      if (saved) {
        const parsed = parseInt(saved, 10);
        if (!isNaN(parsed) && parsed >= minWidth && parsed <= maxWidth) {
          return parsed;
        }
      }
    } catch {
      // ignore localStorage error
    }
    return defaultWidth;
  });

  const [isCollapsed, setIsCollapsed] = useState<boolean>(() => {
    try {
      return localStorage.getItem(`devaulty:sidebar:${storageKey}:collapsed`) === "true";
    } catch {
      return false;
    }
  });

  const [isDragging, setIsDragging] = useState(false);
  const startXRef = useRef<number>(0);
  const startWidthRef = useRef<number>(sidebarWidth);
  const currentWidthRef = useRef<number>(sidebarWidth);

  const toggleCollapse = useCallback(() => {
    const next = !isCollapsed;
    setIsCollapsed(next);
    try {
      localStorage.setItem(`devaulty:sidebar:${storageKey}:collapsed`, String(next));
    } catch {
      // ignore
    }
  }, [isCollapsed, storageKey]);

  const handleMouseDown = useCallback(
    (e: React.MouseEvent) => {
      e.preventDefault();
      setIsDragging(true);
      startXRef.current = e.clientX;
      startWidthRef.current = currentWidthRef.current;
    },
    []
  );

  const handleDoubleClick = useCallback(() => {
    toggleCollapse();
  }, [toggleCollapse]);

  const updateWidth = useCallback(
    (newWidth: number) => {
      const clamped = Math.min(Math.max(newWidth, minWidth), maxWidth);
      setSidebarWidth(clamped);
      currentWidthRef.current = clamped;
      try {
        localStorage.setItem(`devaulty:sidebar:${storageKey}:width`, String(clamped));
      } catch {
        // ignore
      }
    },
    [minWidth, maxWidth, storageKey]
  );

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      const step = 20;
      if (e.key === "ArrowLeft") {
        e.preventDefault();
        updateWidth(sidebarWidth - step);
      } else if (e.key === "ArrowRight") {
        e.preventDefault();
        updateWidth(sidebarWidth + step);
      } else if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        toggleCollapse();
      }
    },
    [sidebarWidth, updateWidth, toggleCollapse]
  );

  useEffect(() => {
    if (!isDragging) return;

    const onMouseUp = () => {
      setIsDragging(false);
      try {
        localStorage.setItem(`devaulty:sidebar:${storageKey}:width`, String(currentWidthRef.current));
      } catch {
        // ignore
      }
    };

    const onMouseMove = (e: MouseEvent) => {
      if (e.buttons === 0) {
        onMouseUp();
        return;
      }
      const deltaX = e.clientX - startXRef.current;
      const newWidth = Math.min(Math.max(startWidthRef.current + deltaX, minWidth), maxWidth);
      currentWidthRef.current = newWidth;
      setSidebarWidth(newWidth);
    };

    window.addEventListener("blur", onMouseUp);
    document.addEventListener("mousemove", onMouseMove);
    document.addEventListener("mouseup", onMouseUp);
    document.body.style.cursor = "col-resize";
    document.body.style.userSelect = "none";

    return () => {
      window.removeEventListener("blur", onMouseUp);
      document.removeEventListener("mousemove", onMouseMove);
      document.removeEventListener("mouseup", onMouseUp);
      document.body.style.cursor = "";
      document.body.style.userSelect = "";
    };
  }, [isDragging, minWidth, maxWidth, storageKey]);

  return {
    sidebarWidth,
    isCollapsed,
    isDragging,
    minWidth,
    maxWidth,
    toggleCollapse,
    setIsCollapsed,
    handleMouseDown,
    handleDoubleClick,
    handleKeyDown,
  };
}
