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
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { OpenticketSdk } from "@/sdk";
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

const TICKETS_QUERY_KEY = ["tickets"];

export function IndexDashboardPage() {
  const usePanelsResult = usePanels();
  const { panels, openPanel, createPanel } = usePanelsResult;
  const sdk = new OpenticketSdk();
  const ticketsQuery = useQuery({
    queryKey: TICKETS_QUERY_KEY,
    queryFn: sdk.tickets,
  });
  const tickets = ticketsQuery.data?.data;

  return (
    <>
      <Helmet>
        <title>Openticket - Tickets</title>
      </Helmet>

      <header>
        <div className="px-6 border-b py-4 flex items-center justify-between">
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
                <TableCell className="font-medium pl-6">{t.title}</TableCell>
                <TableCell>
                  {t.labels.length > 0 && (
                    <div className="space-x-1">
                      {t.labels.map((l) => (
                        <Badge key={l}>{l}</Badge>
                      ))}
                    </div>
                  )}
                </TableCell>
                <TableCell>{t.created_by.name}</TableCell>
                <TableCell className="capitalize pr-6">{t.status}</TableCell>
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
  const queryClient = useQueryClient();
  const sdk = new OpenticketSdk();

  const createTicketMutation = useMutation({
    mutationFn: sdk.createTicket,
    onSuccess: () => {
      queryClient.refetchQueries({
        queryKey: TICKETS_QUERY_KEY,
      });
    },
  });
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
