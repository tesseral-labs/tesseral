import React from "react";
import { useParams } from "react-router";

import { ListAuditLogEventsTable } from "@/components/audit-logs/ListAuditLogEventsTable";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ConsoleListAuditLogEventsRequest } from "@/gen/tesseral/backend/v1/backend_pb";
import { AuditLogEventResourceType } from "@/gen/tesseral/backend/v1/models_pb";

export function OrganizationLogs() {
  const { organizationId } = useParams();

  return (
    <Card>
      <CardHeader>
        <CardTitle>Organization Logs</CardTitle>
        <CardDescription>
          View Logs associated with this Organization
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ListAuditLogEventsTable
          listParams={
            {
              resourceType: AuditLogEventResourceType.ORGANIZATION,
              resourceId: organizationId,
              pageToken: "",
            } as ConsoleListAuditLogEventsRequest
          }
        />
      </CardContent>
    </Card>
  );
}
