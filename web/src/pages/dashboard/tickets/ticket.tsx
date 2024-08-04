import { humanTimeAgo } from "@/pages/utils/time";
import { useComments, useCreateComment } from "@/queries/comments";
import { useTicket } from "@/queries/tickets";
import { Button } from "@/ui/button";
import { Form, FormField } from "@/ui/form";
import { Skeleton } from "@/ui/skeleton";
import { Textarea } from "@/ui/textarea";
import { UserAvatar } from "@/ui/user-avatar";
import { zodResolver } from "@hookform/resolvers/zod";
import { PlusIcon } from "lucide-react";
import { Helmet } from "react-helmet-async";
import { useForm } from "react-hook-form";
import { Link, useParams } from "react-router-dom";
import { z } from "zod";

export function TicketPage() {
  const params = useParams() as { ticketId: string };
  const ticketId = parseInt(params.ticketId, 10);
  const ticketQuery = useTicket(ticketId);

  return (
    <>
      <Helmet>
        <title>{ticketQuery.data?.title ?? "Loading..."} - Openticket</title>
      </Helmet>

      <header className="border-b">
        <div className="px-6 py-6 flex items-center justify-between max-w-screen-xl mx-auto">
          <hgroup className="space-y-1">
            {ticketQuery.data ? (
              <>
                <h1 className="text-2xl font-bold">{ticketQuery.data.title}</h1>
                <span className="text-sm text-muted-foreground">
                  Created by {ticketQuery.data.created_by.name}
                </span>
              </>
            ) : (
              <>
                <Skeleton className="h-[32px] w-[240px] rounded" />
                <Skeleton className="h-[16.5px] w-[160px] rounded" />
              </>
            )}
          </hgroup>

          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm">
              Edit
            </Button>
            <Button size="sm" asChild>
              <Link to="/dashboard/tickets/new">
                <PlusIcon className="w-3 h-3 mr-2" />
                New ticket
              </Link>
            </Button>
          </div>
        </div>
      </header>

      <div className="px-6 py-6 max-w-screen-xl mx-auto">
        <div className="space-y-4">
          <Comments ticketId={ticketId} />
          <CommentForm ticketId={ticketId} />
        </div>
      </div>
    </>
  );
}

type CommentsProps = { ticketId: number };

function Comments({ ticketId }: CommentsProps) {
  const commentsQuery = useComments(ticketId);

  if (!commentsQuery.data) {
    return "Loading comments....";
  }

  return (
    <div className="space-y-4">
      {commentsQuery.data.map((c) => (
        <div key={c.id} className="flex gap-4 w-full">
          <UserAvatar name={c.created_by.name} />

          <div className="flex-1 space-y-1 rounded-lg border p-4">
            <header>
              <div className="text-sm">
                <span className="text-sm font-medium">{c.created_by.name}</span>{" "}
                <span className="text-muted-foreground text-xs">
                  commented {humanTimeAgo(new Date(c.created_at))}
                </span>
              </div>
            </header>
            <div>
              <p className="text-sm">{c.content}</p>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

const commentFormSchema = z.object({
  content: z.string(),
});

type CommentFormProps = { ticketId: number };

function CommentForm({ ticketId }: CommentFormProps) {
  const createCommentMutation = useCreateComment(ticketId);
  const form = useForm({
    resolver: zodResolver(commentFormSchema),
    defaultValues: { content: "" },
  });

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(async (v) => {
          await createCommentMutation.mutateAsync(v);
          form.reset();
        })}
      >
        <fieldset
          className="space-y-2"
          disabled={createCommentMutation.isPending}
        >
          <FormField
            control={form.control}
            name="content"
            render={({ field }) => (
              <Textarea
                {...field}
                rows={5}
                aria-label="Comment"
                placeholder="Type your comment here..."
              />
            )}
          />
          <div className="flex justify-end">
            <Button type="submit">Comment</Button>
          </div>
        </fieldset>
      </form>
    </Form>
  );
}
