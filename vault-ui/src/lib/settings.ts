import { useQuery } from "@connectrpc/connect-query";
import { useEffect, useState } from "react";

import { getSettings } from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { Settings } from "@/gen/tesseral/intermediate/v1/intermediate_pb";

import { LoginLayouts } from "./views";

export function useLayout() {
  const { data: settingsRes } = useQuery(getSettings);

  const [layout, setLayout] = useState<LoginLayouts>();

  useEffect(() => {
    setLayout(
      (settingsRes?.settings || ({} as Settings)).logInLayout as LoginLayouts,
    );
  }, [settingsRes]);

  return layout;
}

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
