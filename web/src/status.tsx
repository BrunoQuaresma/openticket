import { useQuery, useQueryClient } from "@tanstack/react-query";
import { OpenticketSdk } from "./sdk";
import { PropsWithChildren, ReactNode, createContext, useContext } from "react";
import { StatusResponse, User } from "./sdk/types.gen";

type StatusContextValue = {
  data: NonNullable<StatusResponse["data"]>;
  authenticate: (user: User) => void;
  finishSetup: () => void;
};

const StatusContext = createContext<StatusContextValue | undefined>(undefined);

type StatusProviderProps = {
  fallback: ReactNode;
};

const statusQueryKey = ["status"];

export function StatusProvider(props: PropsWithChildren<StatusProviderProps>) {
  const queryClient = useQueryClient();
  const sdk = new OpenticketSdk();
  const { isError, data: res } = useQuery({
    queryKey: statusQueryKey,
    queryFn: () => sdk.status(),
    staleTime: Infinity,
  });

  function finishSetup() {
    const newStatusData: StatusResponse = {
      data: {
        setup: true,
        user: undefined,
      },
    };

    queryClient.setQueryData(statusQueryKey, newStatusData);
  }

  function authenticate(user: User) {
    const newStatusData: StatusResponse = {
      data: {
        setup: true,
        user,
      },
    };

    queryClient.setQueryData(statusQueryKey, newStatusData);
  }

  if (!res) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>Something went wrong</div>;
  }

  return (
    <StatusContext.Provider
      value={{ data: res.data, authenticate, finishSetup }}
    >
      {res.data.setup ? props.children : props.fallback}
    </StatusContext.Provider>
  );
}

export function useStatus() {
  const context = useContext(StatusContext);

  if (context === undefined) {
    throw new Error("useStatus must be used within a StatusProvider");
  }

  return context;
}
