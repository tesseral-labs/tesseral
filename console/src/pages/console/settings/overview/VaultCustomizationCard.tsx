import { useQuery } from "@connectrpc/connect-query";
import { ExternalLink, Settings2 } from "lucide-react";
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

export function VaultCustomizationCard() {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Settings2 />
          Vault Customization
        </CardTitle>
        <CardDescription>
          Customize the appearance of your Vault, including colors, logos, and
          more.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="space-y-4">
          <div className="space-y-2">
            <div className="font-semibold">Vault Domain</div>
            <Link
              className="inline-flex items-center gap-1 text-xs font-mono bg-muted text-muted-foreground"
              to={`https://${getProjectResponse?.project?.vaultDomain}`}
              target="_blank"
            >
              <span>{getProjectResponse?.project?.vaultDomain}</span>
              <ExternalLink className="h-3 w-3" />
            </Link>
          </div>
          <div className="space-y-2">
            <div className="font-semibold">Default Redirect URL</div>
            <span className="inline-flex items-center gap-1 text-xs font-mono bg-muted text-muted-foreground">
              <span>{getProjectResponse?.project?.redirectUri}</span>
            </span>
          </div>

          <div className="space-y-2">
            <div className="font-semibold">Cookie Domain</div>
            <span className="inline-flex items-center gap-1 text-xs font-mono bg-muted text-muted-foreground">
              <span>{getProjectResponse?.project?.cookieDomain}</span>
            </span>
          </div>
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <Link className="w-full" to="/settings/vault">
          <Button className="w-full" variant="outline">
            Customize Vault
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
