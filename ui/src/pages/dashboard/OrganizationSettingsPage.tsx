import { useOrganization } from '@/lib/auth'
import React, { FC } from 'react'

const OrganizationSettingsPage: FC = () => {
  const organization = useOrganization()

  return <div className="text-foreground">{organization?.displayName}</div>
}

export default OrganizationSettingsPage
