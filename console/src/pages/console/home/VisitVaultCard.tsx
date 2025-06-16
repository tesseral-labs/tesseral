import { useQuery } from "@connectrpc/connect-query";
import { ExternalLink, Vault } from "lucide-react";
import React from "react";
import { Link } from "react-router";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getProject } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function VisitVaultCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-x-2">
          <Vault />
          <span className="text-lg">Your Tesseral Vault</span>
        </CardTitle>
        <CardDescription>
          Visit your Tesseral Vault to see how your users will log in and manage
          their Organization.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow space-y-6 text-sm text-muted-foreground">
        <p>
          The Tesseral Vault is the collection of pages your customers interract
          with to log in, manage users, manage organization settings, and
          configure their own user settings.
        </p>
        <p>
          Your Vault is created – and automatically configured based on your
          signup settings – when you create your Project.
        </p>
      </CardContent>
      <CardFooter className="mt-8">
        <Link
          className="w-full"
          to={`https://${getProjectResponse?.project?.vaultDomain}/login`}
          target="_blank"
        >
          <Button className="w-full" variant="outline">
            <ExternalLink />
            Log in to Your Vault
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
