import { useQuery } from "@connectrpc/connect-query";
import React, { useEffect, useMemo, useState } from "react";
import { Link, useLocation } from "react-router";

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import {
  getOrganization,
  getUser,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { titleCaseSlug } from "@/lib/utils";

export function BreadcrumbBar() {
  const { pathname } = useLocation();

  const breadcrumbs: {
    label: string;
    path: string;
  }[] = useMemo(() => {
    const parts = pathname.split("/").filter(Boolean);
    if (parts.length === 0) {
      return [];
    } else {
      return parts.map((part, index) => {
        return {
          label: part,
          path: `/${parts.slice(0, index + 1).join("/")}`,
        };
      });
    }
  }, [pathname]);

  return (
    <div className="hidden lg:flex items-center space-x-2 text-sm border-t border-b p-2 bg-muted/90 backdrop-blur-lg supports-[backdrop-filter]:bg-muted/80 relative -z-10">
      <div className="container px-4 m-auto">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              {breadcrumbs.length === 0 ? (
                <BreadcrumbPage>Home</BreadcrumbPage>
              ) : (
                <BreadcrumbLink asChild>
                  <Link to="/">Home</Link>
                </BreadcrumbLink>
              )}
            </BreadcrumbItem>
            {breadcrumbs.map((breadcrumb, index) => (
              <BreadcrumbSlug
                key={index}
                breadcrumb={breadcrumb}
                last={index === breadcrumbs.length - 1}
              />
            ))}
          </BreadcrumbList>
        </Breadcrumb>
      </div>
    </div>
  );
}

const organizationRegex = /org_([a-z0-9-]+)/;
const userRegex = /user_([a-z0-9-]+)/;

function BreadcrumbSlug({
  breadcrumb,
  last = false,
}: {
  breadcrumb: { label: string; path: string };
  last?: boolean;
}) {
  const [organizationId, setOrganizationId] = useState<string>();
  const [userId, setUserId] = useState<string>();

  const { data: getOrganizationResponse } = useQuery(
    getOrganization,
    {
      id: organizationId,
    },
    {
      enabled: !!organizationId,
    },
  );
  const { data: getUserResponse } = useQuery(
    getUser,
    {
      id: userId,
    },
    {
      enabled: !!userId,
    },
  );

  const label = useMemo(
    () =>
      getOrganizationResponse?.organization?.displayName ||
      getUserResponse?.user?.email ||
      titleCaseSlug(breadcrumb.label),
    [getOrganizationResponse, getUserResponse, breadcrumb.label],
  );

  useEffect(() => {
    if (organizationRegex.test(breadcrumb.label)) {
      setOrganizationId(breadcrumb.label);
    }
    if (userRegex.test(breadcrumb.label)) {
      setUserId(breadcrumb.label);
    }
  }, [breadcrumb, setOrganizationId, setUserId]);

  return (
    <>
      <BreadcrumbSeparator />
      {last ? (
        <BreadcrumbPage>{label}</BreadcrumbPage>
      ) : (
        <BreadcrumbLink href={breadcrumb.path}>{label}</BreadcrumbLink>
      )}
    </>
  );
}
