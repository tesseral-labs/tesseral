import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useQuery } from "@connectrpc/connect-query";
import { ArrowLeft, ChevronDown } from "lucide-react";
import { DateTime } from "luxon";
import React from "react";
import { Link, Outlet, useLocation, useParams } from "react-router";

import { Title } from "@/components/core/Title";
import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { TabLink, Tabs } from "@/components/page/Tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { getAPIKey } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function ApiKeyPage() {
  const { apiKeyId } = useParams();

  const { data: getApiKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const apiKey = getApiKeyResponse?.apiKey;

  return (
    <PageContent>
      <Title title={`${apiKey?.displayName || "API Key"} Details`} />

      <div>
        <Link to={`/organization/api-keys`}>
          <Button variant="ghost" size="sm">
            <ArrowLeft />
            Back to API Keys
          </Button>
        </Link>
      </div>

      <div>
        <div>
          <h1 className="text-2xl font-semibold">API Key</h1>
          <ValueCopier value={apiKey?.id || ""} label="API Key ID" />
          <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
            {apiKey?.revoked ? (
              <Badge variant="secondary">Revoked</Badge>
            ) : (
              <Badge>Active</Badge>
            )}
            <div>•</div>
            <Badge className="border-0" variant="outline">
              Created{" "}
              {apiKey?.createTime &&
                DateTime.fromJSDate(
                  timestampDate(apiKey.createTime),
                ).toRelative()}
            </Badge>
            <div>•</div>
            <Badge className="border-0" variant="outline">
              Updated{" "}
              {apiKey?.updateTime &&
                DateTime.fromJSDate(
                  timestampDate(apiKey.updateTime),
                ).toRelative()}
            </Badge>
          </div>
        </div>
      </div>

      <ApiKeyTabs />

      <div>
        <Outlet />
      </div>
    </PageContent>
  );
}

function ApiKeyTabs() {
  const { pathname } = useLocation();
  const { apiKeyId } = useParams();

  return (
    <>
      <Tabs className="hidden lg:inline-flex">
        <TabLink
          active={pathname === `/organization/api-keys/${apiKeyId}`}
          to={`/organization/api-keys/${apiKeyId}`}
        >
          Details
        </TabLink>
        <TabLink
          active={pathname === `/organization/api-keys/${apiKeyId}/roles`}
          to={`/organization/api-keys/${apiKeyId}/roles`}
        >
          Roles
        </TabLink>
      </Tabs>

      <div className="lg:hidden">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="outline"
              size="sm"
              className="w-full justify-between"
            >
              {pathname === `/organization/api-keys/${apiKeyId}` && "Details"}
              {pathname === `/organization/api-keys/${apiKeyId}/roles` &&
                "Roles"}

              <ChevronDown className="ml-2 h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start">
            <DropdownMenuItem asChild>
              <Link to={`/organization/api-keys/${apiKeyId}`}>Details</Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link to={`/organization/api-keys/${apiKeyId}/roles`}>Roles</Link>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </>
  );
}
