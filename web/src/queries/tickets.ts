import { OpenticketSdk } from "@/sdk";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

const TICKETS_QUERY_KEY = ["tickets"];

function ticketQueryKey(ticketId: number) {
  return [...TICKETS_QUERY_KEY, ticketId];
}

export function useTickets() {
  const sdk = new OpenticketSdk();
  return useQuery({
    queryKey: TICKETS_QUERY_KEY,
    queryFn: sdk.tickets,
  });
}

export function useTicket(id: number) {
  const sdk = new OpenticketSdk();
  return useQuery({
    queryKey: ticketQueryKey(id),
    queryFn: () => sdk.ticket(id),
  });
}

export function useCreateTicket() {
  const sdk = new OpenticketSdk();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: sdk.createTicket,
    onSuccess: () => {
      queryClient.refetchQueries({
        queryKey: TICKETS_QUERY_KEY,
      });
    },
  });
}
