import React from "react";

import { ListAuditLogEventsTable } from "@/components/audit-logs/ListAuditLogEventsTable";
import { Title } from "@/components/core/Title";
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
      <Title title="Audit Logs" />
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
