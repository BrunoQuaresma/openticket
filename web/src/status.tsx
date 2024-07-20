import { useQuery, useQueryClient } from "@tanstack/react-query";
import { OpenticketSdk } from "./sdk";
import { PropsWithChildren, ReactNode, createContext, useContext } from "react";
import { StatusResponse, User } from "./sdk/types.gen";

type StatusContextValue = {
  data: NonNullable<StatusResponse["data"]>;
  authenticate: (user: User) => void;
};

const StatusContext = createContext<StatusContextValue | undefined>(undefined);

type StatusProviderProps = {
  fallback: ReactNode;
};

const statusQueryKey = ["status"];

export function StatusProvider(props: PropsWithChildren<StatusProviderProps>) {
  const queryClient = useQueryClient();
  const sdk = new OpenticketSdk();
  const {
    isLoading,
    isError,
    data: res,
  } = useQuery({
    queryKey: statusQueryKey,
    queryFn: () => sdk.status(),
    staleTime: Infinity,
  });

  function authenticate(user: User) {
    const newStatusData: StatusResponse = {
      data: {
        setup: true,
        user,
      },
    };

    queryClient.setQueryData(statusQueryKey, newStatusData);
  }

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>Something went wrong</div>;
  }

  if (!res?.data?.setup) {
    return <>{props.fallback}</>;
  }

  return (
    <StatusContext.Provider value={{ data: res.data, authenticate }}>
      {props.children}
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
