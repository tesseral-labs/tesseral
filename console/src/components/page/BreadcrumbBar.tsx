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
  }[] = useMemo(() => {
    const parts = pathname.split("/").filter(Boolean);
    if (parts.length === 0) {
      return [];
    } else {
      return parts.map((part, index) => {
        return {
          label: titleCaseSlug(part, index === parts.length - 1),
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
              <>
                <BreadcrumbSeparator />
                {index === breadcrumbs.length - 1 ? (
                  <BreadcrumbPage key={index}>
                    {breadcrumb.label}
                  </BreadcrumbPage>
                ) : (
                  <BreadcrumbItem key={index}>
                    {breadcrumb.label}
                  </BreadcrumbItem>
                )}
              </>
            ))}
          </BreadcrumbList>
        </Breadcrumb>
      </div>
    </div>
  );
}
