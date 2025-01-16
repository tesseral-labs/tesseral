import { useEffect, useState } from 'react'
import { useQuery } from '@connectrpc/connect-query'
import { getProjectUISettings } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { ProjectUISettings } from '@/gen/openauth/intermediate/v1/intermediate_pb'

const useProjectUiSettings = () => {
  const { data: projectUiSettingsRes } = useQuery(getProjectUISettings)

  const [projectUiSettings, setProjectUiSettings] = useState<ProjectUISettings>(
    projectUiSettingsRes?.projectUiSettings || ({} as ProjectUISettings),
  )

  useEffect(() => {
    setProjectUiSettings(
      projectUiSettingsRes?.projectUiSettings || ({} as ProjectUISettings),
    )
  }, [projectUiSettingsRes])

  return projectUiSettings
}

export default useProjectUiSettings
