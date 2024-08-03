import { OpenticketSdk } from "@/sdk";
import { CreateCommentRequest } from "@/sdk/types.gen";
import { useMutation, useQuery } from "@tanstack/react-query";

export function useCreateComment(ticketId: number) {
  const sdk = new OpenticketSdk();

  return useMutation({
    mutationFn: (req: CreateCommentRequest) => sdk.createComment(ticketId, req),
  });
}

export function useComments(ticketId: number) {
  const sdk = new OpenticketSdk();

  return useQuery({
    queryKey: ["comments", ticketId],
    queryFn: () => sdk.comments(ticketId),
  });
}
