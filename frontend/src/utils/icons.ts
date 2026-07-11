import React from "react";
import * as Icons from "lucide-react";

export const ICON_MAPPING: Record<
  string,
  React.ComponentType<{ size?: number; className?: string; style?: React.CSSProperties }>
> = {
  Folder: Icons.Folder,
  Terminal: Icons.Terminal,
  Database: Icons.Database,
  Globe: Icons.Globe,
  Cpu: Icons.Cpu,
  Activity: Icons.Activity,
  BookOpen: Icons.BookOpen,
  Code: Icons.Code,
};

/**
 * Secures and resolves a Lucide icon component via allowlist mapping lookup.
 * Defaults to `Icons.Folder` if not found.
 */
export const getIconComponent = (iconName?: string) => {
  if (iconName && iconName in ICON_MAPPING) {
    return ICON_MAPPING[iconName];
  }
  return Icons.Folder;
};
