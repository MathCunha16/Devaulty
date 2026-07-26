/**
 * Copy text to clipboard supporting modern Async Clipboard API with
 * automatic fallback to document.execCommand('copy') for Webview/Desktop environments.
 */
export const copyToClipboard = async (text: string): Promise<boolean> => {
  if (!text) return false;

  // Try modern navigator.clipboard API if available in current context
  if (navigator.clipboard && typeof navigator.clipboard.writeText === "function") {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch (err) {
      console.warn("navigator.clipboard.writeText failed, attempting execCommand fallback:", err);
    }
  }

  // Fallback for Webview/Desktop context where navigator.clipboard is blocked or unsupported
  try {
    const textArea = document.createElement("textarea");
    textArea.value = text;
    // Position off-screen without causing scroll or layout shift
    textArea.style.position = "fixed";
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.width = "2em";
    textArea.style.height = "2em";
    textArea.style.padding = "0";
    textArea.style.border = "none";
    textArea.style.outline = "none";
    textArea.style.boxShadow = "none";
    textArea.style.background = "transparent";
    textArea.style.opacity = "0";
    textArea.setAttribute("readonly", "");

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    textArea.setSelectionRange(0, text.length);

    const successful = document.execCommand("copy");
    document.body.removeChild(textArea);
    return successful;
  } catch (fallbackErr) {
    console.error("ExecCommand copy fallback failed:", fallbackErr);
    return false;
  }
};
