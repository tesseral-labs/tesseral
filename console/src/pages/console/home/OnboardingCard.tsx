import { useQuery } from "@connectrpc/connect-query";
import { ArrowRight, LogIn, Settings2, Shield } from "lucide-react";
import React from "react";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { getProjectEntitlements } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { cn } from "@/lib/utils";

export function OnboardingCard() {
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
  );

  return (
    <Card
      className={cn(
        getProjectEntitlementsResponse?.entitledBackendApiKeys
          ? "lg:col-span-3"
          : "lg:col-span-2",
      )}
    >
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Getting Started</CardTitle>
            <CardDescription>
              Complete these steps to set up your authentication platform
            </CardDescription>
          </div>
          <div className="text-right">
            <div className="text-2xl font-bold">1/3</div>
            <p className="text-xs text-muted-foreground">Completed</p>
          </div>
        </div>
        <Progress value={(1 / 3) * 100} className="w-full mt-4" />
      </CardHeader>
      <CardContent className="space-y-3 grid grid-cols-1 lg:grid-cols-3 gap-x-4 gap-y-0 lg:gap-y-4">
        <div className="flex shadow items-center space-x-3 p-4 rounded-lg cursor-pointer transition-all text-muted-foreground bg-muted text-sm">
          <div className="flex-1 min-w-0">
            <div className="flex items-center space-x-2 mb-1">
              <Shield className="h-4 w-4" />
              <p className="font-medium line-through">
                Configure Authentication
              </p>
            </div>
            <p className="text-sm line-through">
              Set up SAML, OAuth, and Multi-factor authentication.
            </p>
          </div>
        </div>
        <div className="group border shadow flex items-center space-x-3 p-4 rounded-lg cursor-pointer transition-all text-sm">
          <div className="flex-1 min-w-0">
            <div className="flex items-center space-x-2 mb-1">
              <Settings2 className="h-4 w-4" />
              <p className="font-medium">Customize your Vault</p>
            </div>
            <p className="text-sm">
              Customize the look and feel of your Vault pages and authentication
              flows.
            </p>
          </div>
          <ArrowRight className="h-4 w-4 flex-shrink-0 opacity-0 -translate-x-8 group-hover:translate-x-0 group-hover:opacity-100 transition-all" />
        </div>
        <div className="group border border-border/50 shadow flex items-center space-x-3 p-4 rounded-lg cursor-pointer text-muted-foreground text-sm">
          <div className="flex-1 min-w-0">
            <div className="flex items-center space-x-2 mb-1">
              <LogIn className="h-4 w-4" />
              <p className="font-medium">Log in to your Vault</p>
            </div>
            <p className="text-sm">
              Test your authentication setup with a live login via your Vault.
            </p>
          </div>
          <ArrowRight className="h-4 w-4 flex-shrink-0" />
        </div>
      </CardContent>
    </Card>
  );
}
