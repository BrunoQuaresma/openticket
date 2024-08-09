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
import { Maximize2Icon, MinusIcon, PlusIcon, XIcon } from "lucide-react";
import { usePanels, UsePanelsResult } from "./use-panels";
import { Form, FormField } from "@/ui/form";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/ui/table";
import { Badge } from "@/ui/badge";
import { Link } from "react-router-dom";
import { useCreateTicket, useTickets } from "@/queries/tickets";
import { UserAvatar } from "@/ui/user-avatar";
import { TicketLabels } from "@/ui/ticket-labels";

export function IndexDashboardPage() {
  const usePanelsResult = usePanels();
  const { panels, openPanel, createPanel } = usePanelsResult;
  const ticketsQuery = useTickets();
  const tickets = ticketsQuery.data?.data;

  return (
    <>
      <Helmet>
        <title>Tickets - Openticket</title>
      </Helmet>

      <div className="h-full flex flex-col">
        <header className="px-6 border-b py-4 flex items-center justify-between">
          <h1 className="text-lg font-bold">Tickets</h1>
          <div>
            <Button
              onClick={() => {
                const id = "new-ticket";
                panels[id] ? openPanel(id) : createPanel(id);
              }}
            >
              <PlusIcon className="w-4 h-4 mr-2" />
              Open a ticket
            </Button>
          </div>
        </header>

        {!tickets ? (
          <div className="w-full h-full flex items-center justify-center text-center">
            <span>Loading tickets...</span>
          </div>
        ) : tickets.length > 0 ? (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="pl-6">Title</TableHead>
                <TableHead>Labels</TableHead>
                <TableHead>Created By</TableHead>
                <TableHead className="pr-6">Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {tickets.map((t) => (
                <TableRow key={t.id}>
                  <TableCell className="pl-6 space-x-2">
                    <Link
                      to={`/tickets/${t.id}`}
                      className="hover:text-link hover:underline underline-offset-2"
                    >
                      <span className="font-medium">{t.title}</span>
                    </Link>
                    <span className="text-muted-foreground">#{t.id}</span>
                  </TableCell>
                  <TableCell>
                    {t.labels.length > 0 ? (
                      <TicketLabels labels={t.labels} />
                    ) : (
                      <span className="text-muted-foreground">No labels</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2 items-center">
                      <UserAvatar size="sm" name={t.created_by.name} />
                      {t.created_by.name}
                    </div>
                  </TableCell>
                  <TableCell className="capitalize pr-6">
                    <Badge
                      variant="outline"
                      className="bg-emerald-50 border-emerald-300 text-emerald-900 rounded-full"
                    >
                      {t.status}
                    </Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : (
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
        )}
      </div>
      <TicketPanels {...usePanelsResult} />
    </>
  );
}

const createTicketFormSchema = z.object({
  title: z.string(),
  description: z.string(),
});

type CreateTicketFormValues = z.infer<typeof createTicketFormSchema>;

function TicketPanels(props: UsePanelsResult) {
  const { panels, closePanel, minimizePanel, openPanel } = props;
  const createTicketMutation = useCreateTicket();
  const form = useForm<CreateTicketFormValues>({
    defaultValues: {
      title: "",
      description: "",
    },
    resolver: zodResolver(createTicketFormSchema),
  });

  return (
    <Panels>
      {Object.values(panels).map((p) => (
        <Panel status={p.status} key={p.id}>
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
                  form.reset();
                }}
              >
                <XIcon className="w-3.5 h-3.5" />
              </PanelHeaderAction>
            </PanelHeaderActions>
          </PanelHeader>
          <PanelContent className="p-0 h-full">
            <Form {...form}>
              <form
                className="h-full"
                onSubmit={form.handleSubmit(async (values) => {
                  await createTicketMutation.mutateAsync(values);
                  closePanel(p.id);
                  form.reset();
                })}
              >
                <fieldset
                  className="flex flex-col h-full"
                  disabled={createTicketMutation.isPending}
                >
                  <FormField
                    name="title"
                    control={form.control}
                    render={({ field }) => (
                      <input
                        autoFocus
                        placeholder="Title"
                        type="text"
                        className="text-sm p-3 w-full outline-none border-b"
                        {...field}
                      />
                    )}
                  />
                  <FormField
                    name="description"
                    control={form.control}
                    render={({ field }) => (
                      <textarea
                        className="flex-1 text-sm p-3 block outline-none"
                        {...field}
                      />
                    )}
                  />

                  <footer className="p-3 flex justify-end">
                    <Button type="submit">
                      {createTicketMutation.isPending
                        ? "Creating..."
                        : "Create"}
                    </Button>
                  </footer>
                </fieldset>
              </form>
            </Form>
          </PanelContent>
        </Panel>
      ))}
    </Panels>
  );
}
