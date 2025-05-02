import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import React from 'react';
import { Button, buttonVariants } from '@/components/ui/button';
import { useInfiniteQuery, useQuery } from '@connectrpc/connect-query';
import {
  getRBACPolicy,
  listRoles,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { Link } from 'react-router-dom';

export function RBACSettingsTab() {
  return (
    <div className="space-y-8">
      <RBACPolicyCard />
      <RolesCard />
    </div>
  );
}

function RBACPolicyCard() {
  const { data: getRBACPolicyResponse } = useQuery(getRBACPolicy, {});

  return (
    <Card>
      <CardHeader className="flex-row justify-between items-center gap-x-2">
        <div className="flex flex-col space-y-1.5">
          <CardTitle>Role-Based Access Control Policy</CardTitle>
          <CardDescription>
            A Role-Based Access Control Policy is the set of fine-grained
            Actions in a Project.
          </CardDescription>
        </div>

        <div className="shrink-0 space-x-4">
          <Link
            className={buttonVariants({ variant: 'outline' })}
            to="/project-settings/rbac-settings/rbac-policy/edit"
          >
            Edit
          </Link>
        </div>
      </CardHeader>

      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Action Name</TableHead>
              <TableHead>Description</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {getRBACPolicyResponse?.rbacPolicy?.actions?.map((action) => (
              <TableRow key={action.name}>
                <TableCell className="font-medium font-mono">
                  {action.name}
                </TableCell>
                <TableCell>{action.description}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

function RolesCard() {
  const {
    data: listRolesResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    refetch,
  } = useInfiniteQuery(
    listRoles,
    {
      pageToken: '',
    },
    {
      pageParamKey: 'pageToken',
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const roles = listRolesResponses?.pages?.flatMap((page) => page.roles);

  return (
    <Card>
      <CardHeader className="flex-row justify-between items-center gap-x-2">
        <div className="flex flex-col space-y-1.5">
          <CardTitle>Roles</CardTitle>
          <CardDescription>
            Roles are a named collection of Actions, and can be assigned to
            Users. These are the Roles available to all Organizations in this
            Project.
          </CardDescription>
        </div>

        <div className="shrink-0 space-x-4">
          <Link
            to="/roles/new"
            className={buttonVariants({ variant: 'outline' })}
          >
            Create Role
          </Link>
        </div>
      </CardHeader>

      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Role Display Name</TableHead>
              <TableHead>Actions</TableHead>
              <TableHead>Created</TableHead>
              <TableHead>Updated</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {roles?.map((role) => (
              <TableRow key={role.id}>
                <TableCell>
                  <Link
                    to={`/roles/${role.id}`}
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                  >
                    {role.displayName}
                  </Link>
                </TableCell>
                <TableCell className="font-mono">
                  {role.actions.join(' ')}
                </TableCell>
                <TableCell>
                  {role?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(role.createTime),
                    ).toRelative()}
                </TableCell>
                <TableCell>
                  {role?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(role.updateTime),
                    ).toRelative()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
