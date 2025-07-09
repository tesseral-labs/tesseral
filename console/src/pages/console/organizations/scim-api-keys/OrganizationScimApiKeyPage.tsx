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
import { getSCIMAPIKey } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { NotFound } from "@/pages/NotFoundPage";

export function OrganizationScimApiKeyPage() {
  const { organizationId, scimApiKeyId } = useParams();

  const {
    data: getScimApiKeyResponse,
    isError,
    isLoading,
  } = useQuery(
    getSCIMAPIKey,
    {
      id: scimApiKeyId,
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
          <Title title={`SCIM API Key ${scimApiKeyId}`} />

          <div>
            <Link to={`/organizations/${organizationId}/authentication`}>
              <Button variant="ghost" size="sm">
                <ArrowLeft />
                Back to Authentication
              </Button>
            </Link>
          </div>

          <div className="">
            <h1 className="text-2xl font-semibold">
              {getScimApiKeyResponse?.scimApiKey?.displayName}
            </h1>
            <ValueCopier
              value={getScimApiKeyResponse?.scimApiKey?.id || ""}
              label="SCIM API Key ID"
            />
            <div className="flex flex-wrap mt-2 gap-2 text-muted-foreground/50">
              <Badge className="border-0" variant="outline">
                Created{" "}
                {getScimApiKeyResponse?.scimApiKey?.createTime &&
                  DateTime.fromJSDate(
                    timestampDate(getScimApiKeyResponse.scimApiKey.createTime),
                  ).toRelative()}
              </Badge>
              <div>•</div>
              <Badge className="border-0" variant="outline">
                Updated{" "}
                {getScimApiKeyResponse?.scimApiKey?.updateTime &&
                  DateTime.fromJSDate(
                    timestampDate(getScimApiKeyResponse.scimApiKey.updateTime),
                  ).toRelative()}
              </Badge>
            </div>
          </div>

          <OrganizationScimApiKeysPageTabs />

          <div>
            <Outlet />
          </div>
        </PageContent>
      )}
    </>
  );
}

function OrganizationScimApiKeysPageTabs() {
  const { pathname } = useLocation();
  const { organizationId, scimApiKeyId } = useParams();

  return (
    <>
      {/* Desktop tabs */}
      <Tabs className="hidden lg:inline-block">
        <Tab
          active={
            pathname ===
            `/organizations/${organizationId}/scim-api-keys/${scimApiKeyId}`
          }
        >
          <Link
            to={`/organizations/${organizationId}/scim-api-keys/${scimApiKeyId}`}
          >
            Details
          </Link>
        </Tab>
        <Tab
          active={
            pathname ===
            `/organizations/${organizationId}/scim-api-keys/${scimApiKeyId}/logs`
          }
        >
          <Link
            to={`/organizations/${organizationId}/scim-api-keys/${scimApiKeyId}/logs`}
          >
            Audit Logs
          </Link>
        </Tab>
      </Tabs>

      {/* Mobile tabs */}
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
                  `/organizations/${organizationId}/scim-api-keys/${scimApiKeyId}` &&
                  "Details"}
                {pathname ===
                  `/organizations/${organizationId}/scim-api-keys/${scimApiKeyId}/logs` &&
                  "Audit Logs"}
              </span>
              <ChevronDown className="w-4 h-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem asChild>
              <Link
                to={`/organizations/${organizationId}/scim-api-keys/${scimApiKeyId}`}
              >
                Details
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link
                to={`/organizations/${organizationId}/scim-api-keys/${scimApiKeyId}/logs`}
              >
                Audit Logs
              </Link>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </>
  );
}
