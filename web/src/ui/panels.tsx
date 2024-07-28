import { ComponentProps, HTMLProps } from "react";
import { Card, CardContent } from "./card";
import { Button, ButtonProps } from "./button";
import {
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "./tooltip";

export type PanelStatus = "open" | "minimized";

export function Panels(props: HTMLProps<HTMLDivElement>) {
  return (
    <div
      className="fixed bottom-0 w-full h-full max-h-[680px] px-6 pt-6 flex justify-end items-end pointer-events-none space-x-2"
      {...props}
    />
  );
}

export function Panel({
  status,
  ...props
}: ComponentProps<typeof Card> & { status: "open" | "minimized" }) {
  return (
    <Card
      data-status={status}
      className="group w-full max-w-xl data-[status=minimized]:max-w-xs h-[inherit] data-[status=minimized]:h-[fit-content]  flex flex-col rounded-md rounded-b-none overflow-hidden pointer-events-auto"
      {...props}
    />
  );
}

export function PanelHeader(props: HTMLProps<HTMLElement>) {
  return (
    <header
      className="bg-muted border-b group-data-[status=minimized]:border-0 p-3 py-2 text-sm flex-row flex items-center"
      {...props}
    />
  );
}

export function PanelHeaderAction({
  title,
  ...props
}: ButtonProps & { title: string }) {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            size="icon"
            variant="ghost"
            className="hover:bg-primary/5 w-7 h-7"
            {...props}
          />
        </TooltipTrigger>
        <TooltipContent>
          <span>{title}</span>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
}

export function PanelHeaderActions(props: HTMLProps<HTMLDivElement>) {
  return <div className="inline-flex ml-auto" {...props} />;
}

export function PanelContent(props: ComponentProps<typeof CardContent>) {
  return (
    <CardContent
      className="flex-1 group-data-[status=minimized]:hidden"
      {...props}
    />
  );
}
