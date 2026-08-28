import { createFileRoute } from "@tanstack/react-router";
import { ProjectDetailRouteComponent } from "../components/ProjectDetailView";

export type ProjectTabType = "overview" | "boards" | "snippets" | "problems" | "credentials" | "notes" | "links";

export interface ProjectSearch {
  tab?: ProjectTabType;
  itemId?: string;
}

export const Route = createFileRoute("/projects/$projectId")({
  validateSearch: (search: Record<string, unknown>): ProjectSearch => {
    const validTabs: ProjectTabType[] = [
      "overview",
      "boards",
      "snippets",
      "problems",
      "credentials",
      "notes",
      "links",
    ];
    const tab = search.tab as ProjectTabType;
    return {
      tab: validTabs.includes(tab) ? tab : "overview",
      itemId: typeof search.itemId === "string" ? search.itemId : undefined,
    };
  },
  component: ProjectDetailRouteComponent,
});
