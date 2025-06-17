import { useQuery } from "@connectrpc/connect-query";
import { useEffect, useState } from "react";

import { getSettings } from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { Settings } from "@/gen/tesseral/intermediate/v1/intermediate_pb";

export function useSettings() {
  const { data: settingsRes } = useQuery(getSettings);

  const [settings, setSettings] = useState<Settings | undefined>(
    settingsRes?.settings,
  );

  useEffect(() => {
    setSettings(settingsRes?.settings);
  }, [settingsRes]);

  return settings;
}
