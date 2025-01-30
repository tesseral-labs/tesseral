import { useEffect, useState } from 'react'
import { useQuery } from '@connectrpc/connect-query'
import { getSettings } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { Settings } from '@/gen/openauth/intermediate/v1/intermediate_pb'

const useSettings = () => {
  const { data: settingsRes } = useQuery(getSettings)

  const [projectUiSettings, setProjectUiSettings] = useState<Settings>(
    settingsRes?.settings || ({} as Settings),
  )

  useEffect(() => {
    setProjectUiSettings(settingsRes?.settings || ({} as Settings))
  }, [settingsRes])

  return projectUiSettings
}

export default useSettings
