import { PanelStatus } from "@/ui/panels";
import { useState } from "react";

type PanelState = {
  id: string;
  status: PanelStatus;
  createdAt: Date;
};

export function usePanels() {
  const [panels, setPanels] = useState<Record<string, PanelState>>({});

  function createPanel(id: string = crypto.randomUUID()) {
    setPanels((prev) => {
      const next = { ...prev };
      Object.entries(next).forEach(([key]) => {
        next[key] = { ...prev[key], status: "minimized" };
      });
      next[id] = { id, status: "open", createdAt: new Date() };
      return next;
    });
  }

  function openPanel(id: string) {
    setPanels((prev) => {
      const next = { ...prev };
      next[id] = { ...prev[id], status: "open" };
      return next;
    });
  }

  function minimizePanel(id: string) {
    setPanels((prev) => ({
      ...prev,
      [id]: { ...prev[id], status: "minimized" },
    }));
  }

  function closePanel(id: string) {
    setPanels((prev) => {
      const next = { ...prev };
      delete next[id];
      return next;
    });
  }

  return {
    createPanel,
    openPanel,
    minimizePanel,
    closePanel,
    panels,
  };
}

export type UsePanelsResult = ReturnType<typeof usePanels>;
