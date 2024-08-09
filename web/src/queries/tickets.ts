import { isSuccess, OpenticketSdk } from "@/sdk";
import { PatchTicketRequest, Ticket } from "@/sdk/types.gen";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

const TICKETS_QUERY_KEY = ["tickets"];

export function ticketQueryKey(ticketId: number) {
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
    onSuccess: async (res) => {
      if (!isSuccess(res)) {
        return;
      }
      queryClient.setQueryData(ticketQueryKey(res.data.id), res.data);
      await queryClient.invalidateQueries({
        queryKey: TICKETS_QUERY_KEY,
      });
    },
  });
}

export function usePatchTicket(ticketId: number) {
  const sdk = new OpenticketSdk();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (req: PatchTicketRequest) => sdk.patchTicket(ticketId, req),
    onMutate: async (req) => {
      const prevData = queryClient.getQueryData<Ticket>(
        ticketQueryKey(ticketId)
      );
      queryClient.setQueryData(ticketQueryKey(ticketId), {
        ...prevData,
        ...req,
      });
      console.log({
        ...prevData,
        ...req,
      });
      await queryClient.invalidateQueries({
        queryKey: TICKETS_QUERY_KEY,
        exact: true,
      });
    },
  });
}
