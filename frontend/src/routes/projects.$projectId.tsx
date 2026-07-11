import { createFileRoute } from "@tanstack/react-router";
import { ProjectDetailRouteComponent } from "../components/ProjectDetailView";

export const Route = createFileRoute("/projects/$projectId")({
  component: ProjectDetailRouteComponent,
});
