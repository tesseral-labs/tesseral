import { useQuery } from "@connectrpc/connect-query";
import { ExternalLink, Settings2 } from "lucide-react";
import React from "react";
import { Link, Outlet, useLocation } from "react-router";

import { PageContent } from "@/components/page";
import { TabLink, Tabs } from "@/components/page/Tabs";
import { Title } from "@/components/page/Title";
import { Button } from "@/components/ui/button";
import { getProject } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function VaultCustomizationPage() {
  const { pathname } = useLocation();

  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <PageContent>
      <Title title="Vault Customization Settings" />

      <div className="flex justify-between gap-8">
        <div>
          <h1 className="text-2xl font-semibold flex items-center gap-2">
            <Settings2 />
            <span>Vault Customization</span>
          </h1>
          <p className="text-muted-foreground">
            Customize the vault settings to fit your organization's needs.
          </p>
        </div>
        <Link
          to={`https://${getProjectResponse?.project?.vaultDomain}/login`}
          target="_blank"
        >
          <Button variant="outline" size="sm">
            <span>Visit Your Vault</span>
            <ExternalLink className="w-4 h-4" />
          </Button>
        </Link>
      </div>
      <Tabs>
        <TabLink active={pathname === `/settings/vault`} to={`/settings/vault`}>
          Details
        </TabLink>
        <TabLink
          active={pathname === `/settings/vault/domains`}
          to={`/settings/vault/domains`}
        >
          Domains
        </TabLink>
        <TabLink
          active={pathname === `/settings/vault/branding`}
          to={`/settings/vault/branding`}
        >
          Branding
        </TabLink>
      </Tabs>
      <div>
        <Outlet />
      </div>
    </PageContent>
  );
}
