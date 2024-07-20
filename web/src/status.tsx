import { useQuery } from "@tanstack/react-query";
import { OpenticketSdk } from "./sdk";
import { PropsWithChildren, ReactNode, createContext, useContext } from "react";
import { StatusResponse } from "./sdk/types.gen";

const StatusContext = createContext<StatusResponse | undefined>(undefined);

type StatusProviderProps = {
  // The StatusProvider component renders this element when the setup is not
  // complete.
  fallback: ReactNode;
};

export function StatusProvider(props: PropsWithChildren<StatusProviderProps>) {
  const sdk = new OpenticketSdk();
  const {
    isLoading,
    isError,
    data: res,
  } = useQuery({
    queryKey: ["status"],
    queryFn: () => sdk.status(),
    staleTime: Infinity,
  });
  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>Something went wrong</div>;
  }

  if (res?.data?.setup) {
    return <>{props.fallback}</>;
  }

  return (
    <StatusContext.Provider value={res}>
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
