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
    setIsCollapsed((prev) => {
      const next = !prev;
      try {
        localStorage.setItem(`devaulty:sidebar:${storageKey}:collapsed`, String(next));
      } catch {
        // ignore
      }
      return next;
    });
  }, [storageKey]);

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

  useEffect(() => {
    if (!isDragging) return;

    const onMouseMove = (e: MouseEvent) => {
      const deltaX = e.clientX - startXRef.current;
      const newWidth = Math.min(Math.max(startWidthRef.current + deltaX, minWidth), maxWidth);
      currentWidthRef.current = newWidth;
      setSidebarWidth(newWidth);
    };

    const onMouseUp = () => {
      setIsDragging(false);
      try {
        localStorage.setItem(`devaulty:sidebar:${storageKey}:width`, String(currentWidthRef.current));
      } catch {
        // ignore
      }
    };

    document.addEventListener("mousemove", onMouseMove);
    document.addEventListener("mouseup", onMouseUp);
    document.body.style.cursor = "col-resize";
    document.body.style.userSelect = "none";

    return () => {
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
    toggleCollapse,
    setIsCollapsed,
    handleMouseDown,
    handleDoubleClick,
  };
}
