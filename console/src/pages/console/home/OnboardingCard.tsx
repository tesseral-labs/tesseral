import { timestampNow } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { Building2, Check, LogIn, Shield, Star, XIcon } from "lucide-react";
import React, { useState } from "react";
import { useNavigate } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  getProjectOnboardingProgress,
  listPublishableKeys,
  updateProjectOnboardingProgress,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { cn } from "@/lib/utils";

export function OnboardingCard() {
  const navigate = useNavigate();

  const [exampleAppOpen, setExampleAppOpen] = useState(false);

  const { data: getProjectOnboardingProgressResponse, refetch } = useQuery(
    getProjectOnboardingProgress,
  );
  const { data: listPublishableKeysResponse } = useQuery(listPublishableKeys);
  const updateProgressMutation = useMutation(updateProjectOnboardingProgress);

  const publishableKey = listPublishableKeysResponse?.publishableKeys?.find(
    (key) => key.crossDomainMode === true,
  );

  // Onboarding status
  const configureAuthenticationCompleted =
    !!getProjectOnboardingProgressResponse?.progress
      ?.configureAuthenticationTime;
  const logInToVaultCompleted =
    !!getProjectOnboardingProgressResponse?.progress?.logInToVaultTime;
  const manageOrganizationsCompleted =
    !!getProjectOnboardingProgressResponse?.progress?.manageOrganizationsTime;

  const stepsCompleted =
    Number(configureAuthenticationCompleted) +
    Number(logInToVaultCompleted) +
    Number(manageOrganizationsCompleted);

  async function handleConfigureAuthentication() {
    await updateProgressMutation.mutateAsync({
      progress: {
        configureAuthenticationTime: timestampNow(),
      },
    });

    navigate("/settings/authentication");
  }

  async function handleLogInToVault() {
    setExampleAppOpen(true);
  }

  async function handleRunExampleApp() {
    setExampleAppOpen(false);

    await updateProgressMutation.mutateAsync({
      progress: {
        logInToVaultTime: timestampNow(),
      },
    });

    await refetch();
  }

  async function handleManageOrganizations() {
    await updateProgressMutation.mutateAsync({
      progress: {
        manageOrganizationsTime: timestampNow(),
      },
    });

    await refetch();

    navigate("/organizations");
  }

  async function handleSkipOnboarding() {
    await updateProgressMutation.mutateAsync({
      progress: {
        onboardingSkipped: true,
      },
    });

    await refetch();
  }

  return (
    <>
      <Card className="bg-gradient-to-br from-gray-100 via-white to-gray-100 shadow-lg relative">
        <div className="absolute top-4 right-4">
          <XIcon
            className="w-5 h-5 text-muted-foreground cursor-pointer"
            onClick={handleSkipOnboarding}
          />
        </div>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center space-x-1">
                <Star className="w-5 h-5" />
                <span>Getting Started</span>
              </CardTitle>
              <CardDescription>
                Complete these steps to set up your Tesseral Project.
              </CardDescription>
            </div>
          </div>
          <div className="flex items-end gap-4">
            <Progress
              value={(stepsCompleted / 3) * 100}
              className="w-full mt-4"
            />
            <div className="text-right">
              <div className="text-2xl font-bold">{stepsCompleted}/3</div>
              <p className="text-xs text-muted-foreground">Completed</p>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
            <div
              className={cn(
                "flex shadow items-center space-x-3 p-4 rounded-lg transition-all text-sm flex-wrap gap-4 flex-grow text-muted-foreground relative col-span-1 bg-white",
                configureAuthenticationCompleted
                  ? "bg-muted/50"
                  : "border border-border/50",
              )}
            >
              {configureAuthenticationCompleted && (
                <div className="absolute bg-primary/10 text-muted rounded-full w-5 h-5 top-4 right-0 flex items-center justify-center">
                  <Check className="w-3 h-3 m-auto" />
                </div>
              )}
              <div
                className={cn(configureAuthenticationCompleted && "opacity-50")}
              >
                <div className="flex items-center space-x-2 mb-1">
                  <Shield className="h-4 w-4" />
                  <p className="font-medium">Configure Authentication</p>
                </div>
                <p className="text-sm">
                  Set up SAML, OAuth, and Multi-factor authentication.
                </p>
              </div>
              <Button
                disabled={configureAuthenticationCompleted}
                className={cn(
                  "w-full cursor-pointer",
                  configureAuthenticationCompleted &&
                    "bg-primary/10 text-primary/40",
                )}
                onClick={handleConfigureAuthentication}
              >
                Configure Authentication
              </Button>
            </div>
            <div
              className={cn(
                "flex shadow items-center space-x-3 p-4 rounded-lg transition-all text-sm flex-wrap gap-4 flex-grow text-muted-foreground relative col-span-1 bg-white",
                logInToVaultCompleted
                  ? "bg-muted/50"
                  : "border border-border/50",
                !configureAuthenticationCompleted &&
                  !logInToVaultCompleted &&
                  "text-muted-foreground/50 bg-white/50",
              )}
            >
              {logInToVaultCompleted && (
                <div className="absolute bg-primary/10 text-muted rounded-full w-5 h-5 top-4 right-0 flex items-center justify-center">
                  <Check className="w-3 h-3 m-auto" />
                </div>
              )}
              <div className={cn(logInToVaultCompleted && "opacity-50")}>
                <div className="flex items-center space-x-2 mb-1">
                  <LogIn className="h-4 w-4" />
                  <p className="font-medium">Log in to your Vault</p>
                </div>
                <p className="text-sm">
                  Test your authentication setup with a live login via your
                  Vault.
                </p>
              </div>
              <Button
                disabled={
                  !configureAuthenticationCompleted || logInToVaultCompleted
                }
                variant={
                  !configureAuthenticationCompleted ? "outline" : "default"
                }
                onClick={handleLogInToVault}
                className={cn(
                  "w-full cursor-pointer",
                  logInToVaultCompleted && "bg-primary/10 text-primary/40",
                )}
              >
                Run the Example App
              </Button>
            </div>
            <div
              className={cn(
                "flex shadow items-center space-x-3 p-4 rounded-lg transition-all text-sm flex-wrap gap-4 flex-grow text-muted-foreground relative col-span-1 bg-white",
                manageOrganizationsCompleted
                  ? "bg-muted/50"
                  : "border border-border/50",
                !logInToVaultCompleted &&
                  !manageOrganizationsCompleted &&
                  "text-muted-foreground/50 bg-white/50",
              )}
            >
              {manageOrganizationsCompleted && (
                <div className="absolute bg-primary/10 text-muted rounded-full w-5 h-5 top-4 right-0 flex items-center justify-center">
                  <Check className="w-3 h-3 m-auto" />
                </div>
              )}
              <div className={cn(manageOrganizationsCompleted && "opacity-50")}>
                <div className="flex items-center space-x-2 mb-1">
                  <Building2 className="h-4 w-4" />
                  <p className="font-medium">Manage Organizations</p>
                </div>
                <p className="text-sm">
                  View and manage Organizations within your Tesseral Project.
                </p>
              </div>
              <Button
                disabled={
                  !logInToVaultCompleted || manageOrganizationsCompleted
                }
                variant={!logInToVaultCompleted ? "outline" : "default"}
                onClick={handleManageOrganizations}
                className={cn(
                  "w-full cursor-pointer",
                  manageOrganizationsCompleted &&
                    "bg-primary/10 text-primary/40",
                )}
              >
                Manage Organizations
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Dialog open={exampleAppOpen} onOpenChange={setExampleAppOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Run the Example App</DialogTitle>
            <DialogDescription>
              Execute the code below in your terminal to run the example app.
            </DialogDescription>
          </DialogHeader>

          <Tabs>
            <TabsList>
              <TabsTrigger value="mac-linux">Mac or Linux</TabsTrigger>
              <TabsTrigger value="windows">Windows</TabsTrigger>
            </TabsList>
            <TabsContent value="mac-linux">
              <ValueCopier
                maxLength={50}
                value={`curl -fsSL https://github.com/tesseral-labs/tesseral-example/releases/download/0.0.1/install-and-run.sh | bash -s ${publishableKey?.id}`}
              />
            </TabsContent>
            <TabsContent value="windows">
              <ValueCopier
                maxLength={50}
                value={`$env:TESSERAL_PUBLISHABLE_KEY="${publishableKey?.id}"; iwr -useb https://github.com/tesseral-labs/tesseral-example/releases/download/0.0.1/install-and-run.ps1 | iex
`}
              />
            </TabsContent>
          </Tabs>

          <DialogFooter>
            <Button variant="outline" onClick={handleRunExampleApp}>
              Done
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
