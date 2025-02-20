import { useEffect, useState } from 'react';
import { useQuery } from '@connectrpc/connect-query';
import { getSettings } from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { Settings } from '@/gen/tesseral/intermediate/v1/intermediate_pb';

const useSettings = () => {
  const { data: settingsRes } = useQuery(getSettings);

  const [settings, setSettings] = useState<Settings | undefined>(
    settingsRes?.settings,
  );

  useEffect(() => {
    setSettings(settingsRes?.settings);
  }, [settingsRes]);

  return settings;
};

export default useSettings;
