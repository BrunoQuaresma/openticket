import { Button } from "@/ui/button";
import { Helmet } from "react-helmet-async";
import { CardTitle } from "@/ui/card";
import {
  Panel,
  PanelContent,
  PanelHeader,
  PanelHeaderAction,
  PanelHeaderActions,
  Panels,
} from "@/ui/panels";
import { Maximize2Icon, MinusIcon, XIcon } from "lucide-react";
import { usePanels, UsePanelsResult } from "./use-panels";

export function IndexDashboardPage() {
  const usePanelsResult = usePanels();
  const { panels, openPanel, createPanel } = usePanelsResult;

  return (
    <>
      <Helmet>
        <title>Openticket - Tickets</title>
      </Helmet>
      <div className="w-full h-full flex items-center justify-center text-center">
        <div className="space-y-4">
          <hgroup className="space-y-1">
            <h3 className="text-2xl font-semibold tracking-tight">
              No tickets open
            </h3>
            <p className="text-muted-foreground">
              New open tickets will show up automatically
            </p>
          </hgroup>
          <Button
            variant="secondary"
            onClick={() => {
              const id = "new-ticket";
              panels[id] ? openPanel(id) : createPanel(id);
            }}
          >
            Open a ticket
          </Button>
        </div>
      </div>
      <TicketPanels {...usePanelsResult} />
    </>
  );
}

function TicketPanels(props: UsePanelsResult) {
  const { panels, closePanel, minimizePanel, openPanel } = props;
  return (
    <Panels>
      {Object.values(panels).map((p) => (
        <Panel status={p.status} id={p.id}>
          <PanelHeader>
            <CardTitle>New Ticket</CardTitle>
            <PanelHeaderActions>
              <PanelHeaderAction
                title="Minimize"
                onClick={() => {
                  minimizePanel(p.id);
                }}
              >
                <MinusIcon className="w-3.5 h-3.5" />
              </PanelHeaderAction>
              <PanelHeaderAction
                title="Maximize"
                onClick={() => {
                  openPanel(p.id);
                }}
              >
                <Maximize2Icon className="w-3.5 h-3.5" />
              </PanelHeaderAction>
              <PanelHeaderAction
                title="Close"
                onClick={() => {
                  closePanel(p.id);
                }}
              >
                <XIcon className="w-3.5 h-3.5" />
              </PanelHeaderAction>
            </PanelHeaderActions>
          </PanelHeader>
          <PanelContent className="p-0 flex flex-col h-full">
            <input
              autoFocus
              placeholder="Title"
              type="text"
              className="text-sm p-3 w-full outline-none border-b"
            />
            <textarea className="flex-1 text-sm p-3 block outline-none" />
            <footer className="p-3 flex justify-end">
              <Button>Create</Button>
            </footer>
          </PanelContent>
        </Panel>
      ))}
    </Panels>
  );
}
