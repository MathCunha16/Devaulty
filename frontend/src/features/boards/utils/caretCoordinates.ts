// Helper to compute pixel coordinates of the caret in a textarea
// Based on mirror-div technique used by GitHub, Slack, and Discord.

interface CaretCoordinates {
  top: number;
  left: number;
  lineHeight: number;
}

const propertiesToCopy = [
  "direction",
  "boxSizing",
  "width",
  "height",
  "overflowX",
  "overflowY",
  "borderTopWidth",
  "borderRightWidth",
  "borderBottomWidth",
  "borderLeftWidth",
  "borderStyle",
  "paddingTop",
  "paddingRight",
  "paddingBottom",
  "paddingLeft",
  "fontStyle",
  "fontVariant",
  "fontWeight",
  "fontStretch",
  "fontSize",
  "fontSizeAdjust",
  "lineHeight",
  "fontFamily",
  "textAlign",
  "textTransform",
  "textIndent",
  "textDecoration",
  "letterSpacing",
  "wordSpacing",
  "tabSize",
] as const;

export function getCaretCoordinates(
  element: HTMLTextAreaElement,
  position: number
): CaretCoordinates {
  const style = window.getComputedStyle(element);

  const mirror = document.createElement("div");
  mirror.id = "textarea-caret-position-mirror-div";
  document.body.appendChild(mirror);

  // Set mirror base styles
  mirror.style.position = "absolute";
  mirror.style.top = "-9999px";
  mirror.style.left = "-9999px";
  mirror.style.visibility = "hidden";
  mirror.style.whiteSpace = "pre-wrap";
  mirror.style.wordWrap = "break-word";

  // Copy computed styles from textarea
  propertiesToCopy.forEach((prop) => {
    (mirror.style as unknown as Record<string, string>)[prop] = style.getPropertyValue(
      prop.replace(/([A-Z])/g, "-$1").toLowerCase()
    );
  });

  // Text before caret
  mirror.textContent = element.value.substring(0, position);

  const span = document.createElement("span");
  span.textContent = element.value.substring(position) || ".";
  mirror.appendChild(span);

  const coordinates: CaretCoordinates = {
    top: span.offsetTop + parseInt(style.borderTopWidth, 10) - element.scrollTop,
    left: span.offsetLeft + parseInt(style.borderLeftWidth, 10) - element.scrollLeft,
    lineHeight: parseInt(style.lineHeight, 10) || parseInt(style.fontSize, 10) * 1.4,
  };

  document.body.removeChild(mirror);

  return coordinates;
}
