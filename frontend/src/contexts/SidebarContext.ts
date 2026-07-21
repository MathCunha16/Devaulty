import { createContext } from "react";

export interface SidebarContextValue {
  isOpen: boolean;
  open: () => void;
  close: () => void;
  toggle: () => void;
}

export const SidebarContext = createContext<SidebarContextValue | null>(null);
