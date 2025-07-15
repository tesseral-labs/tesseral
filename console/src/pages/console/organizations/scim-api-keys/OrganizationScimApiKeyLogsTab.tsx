import React from "react";
import { useParams } from "react-router";

import { ListAuditLogEventsTable } from "@/components/audit-logs/ListAuditLogEventsTable";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ConsoleListAuditLogEventsRequest } from "@/gen/tesseral/backend/v1/backend_pb";

export function OrganizationScimApiKeyLogsTab() {
  const { scimApiKeyId } = useParams();

  return (
    <Card>
      <CardHeader>
        <CardTitle>SCIM API Key Logs</CardTitle>
        <CardDescription>
          View the logs for SCIM API key usage, including requests and
          responses.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ListAuditLogEventsTable
          listParams={
            {
              actorScimApiKeyId: scimApiKeyId as string,
              pageToken: "",
            } as ConsoleListAuditLogEventsRequest
          }
        />
      </CardContent>
    </Card>
  );
}
