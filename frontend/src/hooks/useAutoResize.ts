import { useEffect, useRef } from "react";

/**
 * Hook that automatically resizes a textarea to fit its content.
 * Returns a ref to attach to the textarea element.
 *
 * @param value The current value of the textarea (used to trigger resize on change)
 * @param minHeight Minimum height in pixels (default: 80)
 * @param trigger Optional extra dependency to re-run resize (e.g. tab switch remount)
 */
export function useAutoResize(value: string, minHeight = 80, trigger?: unknown) {
  const ref = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    el.style.height = `${minHeight}px`;
    const scrollHeight = el.scrollHeight;
    el.style.height = `${Math.max(scrollHeight, minHeight)}px`;
  }, [value, minHeight, trigger]);

  return ref;
}
