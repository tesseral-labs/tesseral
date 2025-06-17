import React, { useMemo } from "react";
import { Link, useLocation } from "react-router";

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
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
          label: titleCaseSlug(part, index === parts.length - 1),
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
                index={index}
                last={index === breadcrumbs.length - 1}
              />
            ))}
          </BreadcrumbList>
        </Breadcrumb>
      </div>
    </div>
  );
}

function BreadcrumbSlug({
  breadcrumb,
  index,
  last = false,
}: {
  breadcrumb: { label: string; path: string };
  index: number;
  last?: boolean;
}) {
  return (
    <>
      <BreadcrumbSeparator />
      {last ? (
        <BreadcrumbPage>{breadcrumb.label}</BreadcrumbPage>
      ) : (
        <BreadcrumbLink href={breadcrumb.path}>
          {breadcrumb.label}
        </BreadcrumbLink>
      )}
    </>
  );
}
