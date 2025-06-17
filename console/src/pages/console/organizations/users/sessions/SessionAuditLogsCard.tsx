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

export function SessionAuditLogsCard() {
  const { sessionId } = useParams();

  return (
    <Card>
      <CardHeader>
        <CardTitle>Session Logs</CardTitle>
        <CardDescription>
          View Logs associated with this Session
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ListAuditLogEventsTable
          listParams={
            {
              actorSessionId: sessionId,
              pageToken: "",
            } as ConsoleListAuditLogEventsRequest
          }
        />
      </CardContent>
    </Card>
  );
}
