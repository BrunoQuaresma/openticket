import { OpenticketSdk } from "@/sdk";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ticketQueryKey } from "./tickets";

export function useLabels(name?: string) {
  const sdk = new OpenticketSdk();
  return useQuery({
    queryKey: ["labels", name],
    queryFn: () => sdk.labels(name),
  });
}

export function useCreateLabel() {
  const queryClient = useQueryClient();
  const sdk = new OpenticketSdk();

  return useMutation({
    mutationFn: sdk.createLabel,
    onSettled: async (_data, _err, req) => {
      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: ["labels"],
        }),
        queryClient.invalidateQueries({
          queryKey: ticketQueryKey(req.ticket_id),
        }),
      ]);
    },
  });
}
