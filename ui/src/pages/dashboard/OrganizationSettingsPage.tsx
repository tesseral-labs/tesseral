import React, { FC } from 'react'
import { DateTime } from 'luxon'
import { useOrganization, useProject } from '@/lib/auth'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { timestampDate } from '@bufbuild/protobuf/wkt'

const OrganizationSettingsPage: FC = () => {
  const organization = useOrganization()
  const project = useProject()

  return (
    <div className="dark:text-foreground">
      <div className="mb-4">
        <h1 className="text-2xl font-bold mb-2">{organization?.displayName}</h1>
        <span className="text-xs border px-2 py-1 rounded text-gray-400 dark:text-gray-700 bg-gray-200 dark:bg-gray-900 dark:border-gray-800">
          {organization?.id}
        </span>
      </div>
      <Card className="my-8">
        <CardHeader className="py-4">
          <CardTitle className="text-xl">General configuration</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-x-2 text-sm">
            <div className="border-r border-gray-200 pr-8 dark:border-gray-700">
              <div className="font-semibold mb-2">Display Name</div>
              <div className="text-sm text-gray-500">
                {organization?.displayName}
              </div>
            </div>
            <div className="border-r border-gray-200 pl-8 pr-8 dark:border-gray-700">
              <div className="font-semibold mb-2">Created</div>
              <div className="text-sm text-gray-500">
                {organization?.createTime &&
                  DateTime.fromJSDate(
                    new Date(organization.updateTime),
                  ).toRelative()}
              </div>
            </div>
            <div className="px-8">
              <div className="font-semibold mb-2">Last updated</div>
              <div className="text-sm text-gray-500">
                {organization?.updateTime &&
                  DateTime.fromJSDate(
                    new Date(organization.updateTime),
                  ).toRelative()}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default OrganizationSettingsPage
