import { useQuery } from "@connectrpc/connect-query";
import React, { createContext } from "react";

import { getSettings } from "../gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { Settings } from "../gen/tesseral/intermediate/v1/intermediate_pb";

const Context = createContext<Settings | undefined>(undefined);

export function ProjectSettingsProvider({
  children,
}: {
  children?: React.ReactNode;
}) {
  const { data: getSettingsResponse } = useQuery(getSettings);
  if (!getSettingsResponse?.settings) {
    return null;
  }

  return (
    <Context.Provider value={getSettingsResponse.settings}>
      {children}
    </Context.Provider>
  );
}

export function useProjectSettings(): Settings {
  const settings = React.useContext(Context);
  if (!settings) {
    throw new Error(
      "useProjectSettings must be used within a ProjectSettingsProvider"
    );
  }
  return settings;
}
