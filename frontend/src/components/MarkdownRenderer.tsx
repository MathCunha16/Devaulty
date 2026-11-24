import React, { useLayoutEffect, useMemo, useRef } from "react";
import { marked } from "marked";
import DOMPurify, { type Config as DOMPurifyConfig } from "dompurify";
import { useTheme } from "../hooks/useTheme";
import { hydrateMermaidDiagrams } from "../utils/mermaidRenderer";

export interface MarkdownRendererProps {
  /** Raw markdown text to parse, sanitize, and render. */
  content: string;
  /** Optional CSS class name(s) applied to the root container. */
  className?: string;
  /** Optional HTML transform applied after marked.parse and before DOMPurify.sanitize. */
  transformHtml?: (rawHtml: string) => string;
  /** Optional HTML transform applied after DOMPurify.sanitize (e.g. enabling task checkboxes). */
  postSanitizeHtml?: (sanitizedHtml: string) => string;
  /** Optional DOMPurify configuration overrides. */
  sanitizeConfig?: DOMPurifyConfig;
  /** Whether to enable GFM breaks (default: true). */
  breaks?: boolean;
  /** Optional click handler on the rendered container. */
  onClick?: React.MouseEventHandler<HTMLDivElement>;
  /** Optional ref callback or RefObject to access the underlying container element. */
  containerRef?: React.RefObject<HTMLDivElement | null>;
}

/**
 * Shared, performance-first Markdown renderer with DOMPurify sanitization and
 * lazy, theme-aware, LRU-cached Mermaid diagram hydration.
 */
export const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({
  content,
  className = "",
  transformHtml,
  postSanitizeHtml,
  sanitizeConfig,
  breaks = true,
  onClick,
  containerRef: externalRef,
}) => {
  const { theme } = useTheme();
  const internalRef = useRef<HTMLDivElement>(null);

  const sanitizedHtml = useMemo(() => {
    if (!content) return "";
    try {
      const parsed = marked.parse(content, { breaks, gfm: true, async: false }) as string;
      const transformed = transformHtml ? transformHtml(parsed) : parsed;
      const sanitized = sanitizeConfig
        ? (DOMPurify.sanitize(transformed, sanitizeConfig) as string)
        : DOMPurify.sanitize(transformed);
      return postSanitizeHtml ? postSanitizeHtml(sanitized) : sanitized;
    } catch {
      return DOMPurify.sanitize(content);
    }
  }, [content, breaks, transformHtml, sanitizeConfig, postSanitizeHtml]);

  useLayoutEffect(() => {
    const el = internalRef.current;
    if (!el) return;

    el.innerHTML = sanitizedHtml;

    const controller = new AbortController();
    void hydrateMermaidDiagrams(el, theme === "dark" ? "dark" : "light", controller.signal);

    return () => {
      controller.abort();
    };
  }, [sanitizedHtml, theme]);

  const setRefs = (node: HTMLDivElement | null) => {
    internalRef.current = node;
    if (externalRef) {
      externalRef.current = node;
    }
  };

  return <div ref={setRefs} className={className} onClick={onClick} />;
};
