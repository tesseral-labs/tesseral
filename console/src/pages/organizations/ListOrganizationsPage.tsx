import React from 'react';
import { useInfiniteQuery, useQuery } from '@connectrpc/connect-query';
import { listOrganizations } from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
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
import { Card, CardContent } from '@/components/ui/card';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { PageDescription, PageTitle } from '@/components/page';
import { Button } from '@/components/ui/button';
import { LoaderCircleIcon } from 'lucide-react';

export const ListOrganizationsPage = () => {
  const { data: listOrganizationsResponses, fetchNextPage, hasNextPage, isFetchingNextPage } = useInfiniteQuery(
    listOrganizations,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const organizations = listOrganizationsResponses?.pages?.flatMap(page => page.organizations);

  return (
    <div>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Organizations</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>Organizations</PageTitle>
      <PageDescription>
        An Organization represents one of your business customers.
      </PageDescription>

      <Card className="mt-8 overflow-hidden">
        <CardContent className="-m-6 mt-0">
          <Table>
            <TableHeader className="bg-gray-50">
              <TableRow>
                <TableHead>Display Name</TableHead>
                <TableHead>ID</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {organizations?.map((org) => (
                <TableRow key={org.id}>
                  <TableCell className="font-medium">
                    <Link
                      className="underline underline-offset-2 decoration-muted-foreground/40"
                      to={`/organizations/${org.id}`}
                    >
                      {org.displayName}
                    </Link>
                  </TableCell>
                  <TableCell className="font-mono">{org.id}</TableCell>
                  <TableCell>
                    {org.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(org.createTime),
                      ).toRelative()}
                  </TableCell>
                  <TableCell>
                    {org.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(org.updateTime),
                      ).toRelative()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {hasNextPage && (
        <Button
          className="mt-4"
          variant="outline"
          onClick={() => fetchNextPage()}
        >
          {isFetchingNextPage && <LoaderCircleIcon className="h-4 w-4 animate-spin" />}
          Load more
        </Button>
      )}
    </div>
  );
};
