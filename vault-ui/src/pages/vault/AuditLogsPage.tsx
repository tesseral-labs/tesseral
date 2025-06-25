import React from "react";
import { Helmet } from "react-helmet";

import { ListAuditLogEventsTable } from "@/components/audit-logs/ListAuditLogEventsTable";
import { PageContent } from "@/components/page";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ListAuditLogEventsRequest } from "@/gen/tesseral/frontend/v1/frontend_pb";

export function AuditLogsPage() {
  return (
    <PageContent>
      <Helmet>
        <title>Audit Logs</title>
      </Helmet>
      <Card>
        <CardHeader>
          <CardTitle>Audit Logs</CardTitle>
          <CardDescription>
            View a history of actions effecting your organization.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ListAuditLogEventsTable
            listParams={{} as ListAuditLogEventsRequest}
          />
        </CardContent>
      </Card>
    </PageContent>
  );
}
