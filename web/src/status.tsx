import { useQuery } from "@tanstack/react-query";
import { OpenticketSdk } from "./sdk";
import { StatusResponse } from "./sdk/types";
import { PropsWithChildren, ReactNode, createContext } from "react";

const StatusContext = createContext<StatusResponse | undefined>(undefined);

type StatusProviderProps = {
  // The StatusProvider component renders this element when the setup is not
  // complete.
  setup: ReactNode;
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

  if (!res?.data.setup) {
    return <>{props.setup}</>;
  }

  return (
    <StatusContext.Provider value={res}>
      {props.children}
    </StatusContext.Provider>
  );
}