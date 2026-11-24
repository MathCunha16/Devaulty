/**
 * Conservative heuristic to detect whether raw text intentionally uses Markdown
 * formatting (including ```mermaid or fenced code blocks), while avoiding false
 * positives on raw stack traces, JSON, SQL, or shell/Python `# comments`.
 */
export function detectMarkdownContent(text: string | undefined | null): boolean {
  if (!text || !text.trim()) return false;

  // 1. Fenced code blocks (e.g. ```mermaid, ```ts, ```sql, ```)
  if (/^```[\w-]*\s*$/m.test(text)) return true;

  // 2. Level 2+ headings (avoids treating bash/python `# comment` as Markdown H1)
  if (/^#{2,6}\s+\S/m.test(text)) return true;

  // 3. Markdown links [label](url)
  if (/\[[^\]]+\]\(https?:\/\/[^)\s]+\)/.test(text)) return true;

  // 4. Task list checkboxes
  if (/^\s*[-*+]\s+\[[ xX]\]\s+/m.test(text)) return true;

  // 5. Markdown tables
  if (/^\|.+\|\s*\n\s*\|[\s:-]+\|/m.test(text)) return true;

  // 6. Paired bold syntax **text**
  if (/\*\*[^*\n]+\*\*/.test(text)) return true;

  return false;
}
