import React from 'react';
import { useQuery } from '@connectrpc/connect-query';
import { listUsers } from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Link } from 'react-router-dom';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { useParams } from 'react-router';
import { Badge } from '@/components/ui/badge';

export const OrganizationUsersTab = () => {
  const { organizationId } = useParams();
  const { data: listUsersResponse } = useQuery(listUsers, {
    organizationId,
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Users</CardTitle>
        <CardDescription>
          A user is what people using your product log into. Lorem ipsum dolor.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Email</TableHead>
              <TableHead>ID</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead>Updated At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {listUsersResponse?.users?.map((user) => (
              <TableRow key={user.id}>
                <TableCell>
                  <Link
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                    to={`/organizations/${organizationId}/users/${user.id}`}
                  >
                    {user.email}
                  </Link>

                  {user.owner && (
                    <Badge variant="outline" className="ml-2">
                      Owner
                    </Badge>
                  )}
                </TableCell>
                <TableCell className="font-mono">{user.id}</TableCell>
                <TableCell>
                  {user.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(user.createTime),
                    ).toRelative()}
                </TableCell>
                <TableCell>
                  {user.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(user.updateTime),
                    ).toRelative()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};
