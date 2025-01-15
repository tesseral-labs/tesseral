import { useQuery } from '@connectrpc/connect-query'
import { listProjectAPIKeys } from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Link } from 'react-router-dom'
import React from 'react'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'

export function ListProjectAPIKeysPage() {
  const { data: listProjectAPIKeysResponse } = useQuery(listProjectAPIKeys, {})

  return (
    <div>
      <h1 className="font-semibold text-xl">Project API Keys</h1>
      <Table className="mt-4">
        <TableHeader>
          <TableRow>
            <TableCell>Display Name</TableCell>
            <TableHead>ID</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Created At</TableHead>
            <TableHead>Updated At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {listProjectAPIKeysResponse?.projectApiKeys?.map((projectAPIKey) => (
            <TableRow key={projectAPIKey.id}>
              <TableCell className="font-medium">
                <Link to={`/project-api-keys/${projectAPIKey.id}`}>
                  {projectAPIKey.displayName}
                </Link>
              </TableCell>
              <TableCell className="font-mono">{projectAPIKey.id}</TableCell>
              <TableCell>
                {projectAPIKey?.revoked ? 'Revoked' : 'Active'}
              </TableCell>
              <TableCell>
                {DateTime.fromJSDate(
                  timestampDate(projectAPIKey.createTime!),
                ).toRelative()}
              </TableCell>
              <TableCell>
                {DateTime.fromJSDate(
                  timestampDate(projectAPIKey.updateTime!),
                ).toRelative()}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
