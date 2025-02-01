import { useEffect, useState } from 'react'
import { useQuery } from '@connectrpc/connect-query'
import { getSettings } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { Settings } from '@/gen/openauth/intermediate/v1/intermediate_pb'
import { LoginLayouts } from './views'

export const useLayout = () => {
  const { data: settingsRes } = useQuery(getSettings)

  const [layout, setLayout] = useState<LoginLayouts>()

  useEffect(() => {
    setLayout(
      ((settingsRes?.settings as any) || {}).layout || LoginLayouts.Centered,
    )
  }, [settingsRes])

  return layout
}

const useSettings = () => {
  const { data: settingsRes } = useQuery(getSettings)

  const [settings, setSettings] = useState<Settings | undefined>(
    settingsRes?.settings,
  )

  useEffect(() => {
    setSettings(settingsRes?.settings)
  }, [settingsRes])

  return settings
}

export default useSettings
