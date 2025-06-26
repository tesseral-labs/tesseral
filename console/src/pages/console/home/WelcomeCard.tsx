import { useQuery } from "@connectrpc/connect-query";
import { Users } from "lucide-react";
import React from "react";
import { Link } from "react-router";

import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getProjectEntitlements } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { cn } from "@/lib/utils";

export function WelcomeCard() {
  const {
    data: getProjectEntitlementsResponse,
    isLoading: isLoadingEntitlements,
  } = useQuery(getProjectEntitlements);

  return (
    <Card
      className={cn(
        "col-span-1",
        isLoadingEntitlements ||
          getProjectEntitlementsResponse?.entitledBackendApiKeys
          ? "lg:col-span-2"
          : "",
      )}
    >
      <CardHeader>
        <CardTitle className="text-xl flex flex-wrap items-center gap-x-2">
          <span>Welcome to Tesseral</span>
          <Badge className="bg-gray-50 border-gray-200 text-muted-foreground">
            Beta
          </Badge>
        </CardTitle>
        <CardDescription>
          <Badge className="bg-indigo-50 border-indigo-200 text-indigo-800">
            <Users />
            From the Founders
          </Badge>
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4 flex-grow text-sm text-muted-foreground ">
        <p>We consider it a privilege to support auth for your app.</p>
        <p>
          We've designed Tesseral to be powerful and easy to use. If you
          encounter issues with your implementation, please let us know
          immediately. We want to get this right.
        </p>
        <p>
          This project is still in its early days. We're always grateful to
          receive feedback from the community as we continue to advance the
          project.
        </p>
        <p>
          &mdash;{" "}
          <Link
            to="https://www.linkedin.com/in/ucarion/"
            className="underline underline-offset-2 decoration-muted-foreground/40 hover:text-primary hover:decoration-primary transition-colors"
          >
            Ulysse
          </Link>{" "}
          and{" "}
          <Link
            to="https://www.linkedin.com/in/nedoleary/"
            className="underline underline-offset-2 decoration-muted-foreground/40 hover:text-primary hover:decoration-primary transition-colors"
          >
            Ned
          </Link>
        </p>
      </CardContent>
      <CardFooter className="mt-8 z-index-0 relative">
        <img
          src="/images/tesseral-logo-black.svg"
          alt="Tesseral Logo"
          className="w-32 mt-at opacity-20"
        />
      </CardFooter>
    </Card>
  );
}
