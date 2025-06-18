import { useQuery } from "@connectrpc/connect-query";
import {
  ArrowRight,
  Building2,
  Key,
  Settings2,
  Shield,
  Webhook,
} from "lucide-react";
import React from "react";
import { Link } from "react-router-dom";

import { PageContent } from "@/components/page";
import { Title } from "@/components/page/Title";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  getProjectEntitlements,
  getProjectWebhookManagementURL,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

import { UpgradeCard } from "./UpgradeCard";
import { VisitVaultCard } from "./VisitVaultCard";
import { WelcomeCard } from "./WelcomeCard";

export function HomePage() {
  const {
    data: getProjectEntitlementsResponse,
    isLoading: isLoadingEntitlements,
  } = useQuery(getProjectEntitlements);
  const { data: getProjectWebhookManagementUrlResponse } = useQuery(
    getProjectWebhookManagementURL,
  );

  return (
    <PageContent>
      <Title title="Home" />

      <div className={"grid grid-cols-1 lg:grid-cols-3 gap-8"}>
        <WelcomeCard />
        <VisitVaultCard />
        {!isLoadingEntitlements &&
          !getProjectEntitlementsResponse?.entitledBackendApiKeys && (
            <UpgradeCard />
          )}
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 items-stretch">
        <Link to="/organizations">
          <Card className="hover:shadow-md transition-all ease-in-out">
            <CardHeader>
              <CardTitle className="flex items-center gap-x-2">
                <Building2 />
                Organizations
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground">
                Manage your customers, invite users, and configure
                organization-specific settings.
              </div>
            </CardContent>
            <CardFooter className="justify-end text-sm">
              View Organizations <ArrowRight className="ml-2 h-4 w-4" />
            </CardFooter>
          </Card>
        </Link>
        <Link to="/settings/authentication">
          <Card className="hover:shadow-md transition-all">
            <CardHeader>
              <CardTitle className="flex items-center gap-x-2">
                <Shield />
                Authentication
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground">
                Configure SAML, SCIM, OAuth, and Multi-factor authentication
                settings for your project.
              </div>
            </CardContent>
            <CardFooter className="justify-end text-sm">
              Configure Authentication <ArrowRight className="ml-2 h-4 w-4" />
            </CardFooter>
          </Card>
        </Link>
        <Link to="/settings/api-keys">
          <Card className="hover:shadow-md transition-all">
            <CardHeader>
              <CardTitle className="flex items-center gap-x-2">
                <Key />
                API Keys
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground">
                Generate and manage API keys for secure programmatic access to
                your authentication platform.
              </div>
            </CardContent>
            <CardFooter className="justify-end text-sm">
              Manage API Keys <ArrowRight className="ml-2 h-4 w-4" />
            </CardFooter>
          </Card>
        </Link>
        <Link to="/settings/vault">
          <Card className="hover:shadow-md transition-all">
            <CardHeader>
              <CardTitle className="flex items-center gap-x-2">
                <Settings2 />
                Vault Customization
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground">
                Customize the look and feel of your Vault pages to match your
                brand identity.
              </div>
            </CardContent>
            <CardFooter className="justify-end text-sm">
              Customize your Vault <ArrowRight className="ml-2 h-4 w-4" />
            </CardFooter>
          </Card>
        </Link>
        <Link to={getProjectWebhookManagementUrlResponse?.url || "#"}>
          <Card className="hover:shadow-md transition-all">
            <CardHeader>
              <CardTitle className="flex items-center gap-x-2">
                <Webhook />
                Webhooks
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground">
                Set up webhooks to receive real-time notifications about
                authentication events and user activities.
              </div>
            </CardContent>
            <CardFooter className="justify-end text-sm">
              Configure Webhooks <ArrowRight className="ml-2 h-4 w-4" />
            </CardFooter>
          </Card>
        </Link>
      </div>
    </PageContent>
  );
}
