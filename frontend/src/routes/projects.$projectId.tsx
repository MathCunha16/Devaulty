import { createFileRoute } from "@tanstack/react-router";
import { ProjectDetailRouteComponent } from "../components/ProjectDetailView";

export type ProjectTabType = "overview" | "snippets" | "problems" | "credentials" | "notes" | "links";

interface ProjectSearch {
  tab?: ProjectTabType;
}

export const Route = createFileRoute("/projects/$projectId")({
  validateSearch: (search: Record<string, unknown>): ProjectSearch => {
    const validTabs: ProjectTabType[] = [
      "overview",
      "snippets",
      "problems",
      "credentials",
      "notes",
      "links",
    ];
    const tab = search.tab as ProjectTabType;
    return {
      tab: validTabs.includes(tab) ? tab : "overview",
    };
  },
  component: ProjectDetailRouteComponent,
});
