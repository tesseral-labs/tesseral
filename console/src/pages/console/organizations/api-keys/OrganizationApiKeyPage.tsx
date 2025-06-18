import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useQuery } from "@connectrpc/connect-query";
import { ArrowLeft, ChevronDown } from "lucide-react";
import { DateTime } from "luxon";
import React from "react";
import { Link, Outlet, useLocation, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { PageLoading } from "@/components/page/PageLoading";
import { Tab, Tabs } from "@/components/page/Tabs";
import { Title } from "@/components/page/Title";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { getAPIKey } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { NotFound } from "@/pages/NotFoundPage";

export function OrganizationApiKeyPage() {
  const { pathname } = useLocation();
  const { apiKeyId, organizationId } = useParams();

  const {
    data: getApiKeyResponse,
    isError,
    isLoading,
  } = useQuery(
    getAPIKey,
    {
      id: apiKeyId,
    },
    {
      retry: false,
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
          <Title title={`API Key ${apiKeyId}`} />

          <div>
            <Link to={`/organizations/${organizationId}/api-keys`}>
              <Button variant="ghost" size="sm">
                <ArrowLeft />
                Back to API Keys
              </Button>
            </Link>
          </div>

          <div className="">
            <h1 className="text-2xl font-semibold">
              {getApiKeyResponse?.apiKey?.displayName}
            </h1>
            <ValueCopier
              value={getApiKeyResponse?.apiKey?.id || ""}
              label="API Key ID"
            />
            <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
              <Badge className="border-0" variant="outline">
                Created{" "}
                {getApiKeyResponse?.apiKey?.createTime &&
                  DateTime.fromJSDate(
                    timestampDate(getApiKeyResponse.apiKey.createTime),
                  ).toRelative()}
              </Badge>
              <div>•</div>
              <Badge className="border-0" variant="outline">
                Updated{" "}
                {getApiKeyResponse?.apiKey?.updateTime &&
                  DateTime.fromJSDate(
                    timestampDate(getApiKeyResponse.apiKey.updateTime),
                  ).toRelative()}
              </Badge>
            </div>
          </div>
          <Tabs className="hidden lg:inline-block">
            <Tab
              active={
                pathname ===
                `/organizations/${organizationId}/api-keys/${apiKeyId}`
              }
            >
              <Link
                to={`/organizations/${organizationId}/api-keys/${apiKeyId}`}
              >
                Details
              </Link>
            </Tab>
            <Tab
              active={pathname.startsWith(
                `/organizations/${organizationId}/api-keys/${apiKeyId}/roles`,
              )}
            >
              <Link
                to={`/organizations/${organizationId}/api-keys/${apiKeyId}/roles`}
              >
                Roles
              </Link>
            </Tab>
            <Tab
              active={
                pathname ===
                `/organizations/${organizationId}/api-keys/${apiKeyId}/logs`
              }
            >
              <Link
                to={`/organizations/${organizationId}/api-keys/${apiKeyId}/logs`}
              >
                Audit Logs
              </Link>
            </Tab>
          </Tabs>

          <div className="block lg:hidden space-y-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  className="flex items-center gap-2"
                  variant="outline"
                  size="sm"
                >
                  <span>
                    {pathname ===
                      `/organizations/${organizationId}/api-keys/${apiKeyId}` &&
                      "Details"}
                    {pathname ===
                      `/organizations/${organizationId}/api-keys/${apiKeyId}/roles` &&
                      "Roles"}
                    {pathname ===
                      `/organizations/${organizationId}/api-keys/${apiKeyId}/logs` &&
                      "Audit Logs"}
                  </span>
                  <ChevronDown className="w-4 h-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuItem asChild>
                  <Link
                    to={`/organizations/${organizationId}/api-keys/${apiKeyId}`}
                  >
                    Details
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link
                    to={`/organizations/${organizationId}/api-keys/${apiKeyId}/roles`}
                  >
                    Roles
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link
                    to={`/organizations/${organizationId}/api-keys/${apiKeyId}/logs`}
                  >
                    Audit Logs
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div>
            <Outlet />
          </div>
        </PageContent>
      )}
    </>
  );
}
