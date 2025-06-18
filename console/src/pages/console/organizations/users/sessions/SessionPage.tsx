import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useQuery } from "@connectrpc/connect-query";
import { ArrowLeft } from "lucide-react";
import { DateTime } from "luxon";
import React from "react";
import { Link, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { PageLoading } from "@/components/page/PageLoading";
import { Title } from "@/components/page/Title";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { getSession } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { NotFound } from "@/pages/NotFoundPage";

import { SessionAuditLogsCard } from "./SessionAuditLogsCard";

export function SessionPage() {
  const { organizationId, sessionId, userId } = useParams();

  const {
    data: getSessionResponse,
    isError,
    isLoading,
  } = useQuery(
    getSession,
    {
      id: sessionId,
    },
    {
      retry: 3,
    },
  );

  return (
    <>
      {isLoading ? (
        <PageLoading />
      ) : isError ? (
        <NotFound />
      ) : (
        <PageContent>
          <Title title={`Session ${sessionId}`} />

          <div>
            <Link
              to={`/organizations/${organizationId}/users/${userId}/sessions`}
            >
              <Button variant="ghost" size="sm">
                <ArrowLeft />
                Back to User Sessions
              </Button>
            </Link>
          </div>

          <div>
            <h1 className="text-2xl font-semibold">Session</h1>
            <ValueCopier
              value={getSessionResponse?.session?.id || ""}
              label="Passkey ID"
            />
            <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
              <Badge className="border-0" variant="outline">
                Created{" "}
                {getSessionResponse?.session?.createTime &&
                  DateTime.fromJSDate(
                    timestampDate(getSessionResponse.session.createTime),
                  ).toRelative()}
              </Badge>
              <div>•</div>
              <Badge className="border-0" variant="outline">
                Expires{" "}
                {getSessionResponse?.session?.expireTime &&
                  DateTime.fromJSDate(
                    timestampDate(getSessionResponse.session.expireTime),
                  ).toRelative()}
              </Badge>
              <div>•</div>
              <Badge className="border-0" variant="outline">
                Last Active{" "}
                {getSessionResponse?.session?.lastActiveTime &&
                  DateTime.fromJSDate(
                    timestampDate(getSessionResponse.session.lastActiveTime),
                  ).toRelative()}
              </Badge>
            </div>
          </div>

          <SessionAuditLogsCard />
        </PageContent>
      )}
    </>
  );
}
