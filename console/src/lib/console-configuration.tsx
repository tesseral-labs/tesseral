import { useQuery } from "@connectrpc/connect-query";
import React, { PropsWithChildren, createContext, useContext } from "react";

import { consoleGetConfiguration } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { ConsoleConfiguration } from "@/gen/tesseral/backend/v1/models_pb";

const Context = createContext<ConsoleConfiguration | undefined>(undefined);

export function ConsoleConfigurationProvider({ children }: PropsWithChildren) {
  const { data: consoleGetConfigurationResponse } = useQuery(
    consoleGetConfiguration,
  );
  if (!consoleGetConfigurationResponse?.configuration) {
    return null;
  }

  return (
    <Context.Provider value={consoleGetConfigurationResponse.configuration}>
      {children}
    </Context.Provider>
  );
}

export function useConsoleConfiguration(): ConsoleConfiguration {
  const configuration = useContext(Context);
  if (!configuration) {
    throw new Error(
      "useConsoleConfiguration must be used within a ConsoleConfigurationProvider",
    );
  }
  return configuration;
}
