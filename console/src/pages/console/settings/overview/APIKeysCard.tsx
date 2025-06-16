import { useMutation, useQuery } from "@connectrpc/connect-query";
import { Crown, Key } from "lucide-react";
import React from "react";
import { Link } from "react-router";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  createStripeCheckoutLink,
  getProjectEntitlements,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function ApiKeysCard() {
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
  );
  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  async function handleUpgrade() {
    const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
    window.location.href = url;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Key />
          API Keys
        </CardTitle>
        <CardDescription>
          Manage Publishable and Backend API Keys for your Project. Publishable
          Keys can be used to identify your application and are safe to be
          publicly expose. Backend API Keys are used to authenticate to the
          Tesseral API.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="space-y-4">
          <div className="space-y-2">
            <div>
              <div className="font-semibold">Publishable API Keys</div>
              <div className="text-xs text-muted-foreground">
                Publishable API Keys are used to identify your application and
                can be safely exposed.
              </div>
            </div>
            <div className="flex flex-wrap gap-2">
              <Badge>Always Enabled</Badge>
            </div>
          </div>
          <div className="space-y-2">
            <div>
              <div className="font-semibold">Backend API Keys</div>
              <div className="text-xs text-muted-foreground">
                Backend API Keys are used to authenticate to the Tesseral API.
              </div>
            </div>
            <div className="flex flex-wrap gap-2">
              {getProjectEntitlementsResponse?.entitledBackendApiKeys ? (
                <Badge>Enabled</Badge>
              ) : (
                <div
                  className="bg-gradient-to-br from-violet-500 via-purple-500 to-blue-500 border-0 text-white w-full mt-4 p-4 rounded-md"
                  onClick={handleUpgrade}
                >
                  <div className="absolute inset-0 bg-gradient-to-br from-white/10 to-transparent" />

                  <div className="flex flex-wrap w-full gap-4">
                    <div className="w-full space-y-4 md:flex-grow">
                      <div className="flex items-center space-x-3">
                        <div className="p-2 rounded-full bg-white/20 backdrop-blur-sm">
                          <Crown className="h-6 w-6 text-white" />
                        </div>
                        <div>
                          <h3 className="font-semibold text-white">
                            Upgrade to Growth
                          </h3>
                          <p className="text-xs text-white/80">
                            Access Backend API Keys and more
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter className="mt-4">
        <Link className="w-full" to="/settings/api-keys">
          <Button className="w-full" variant="outline">
            Manage API Keys
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
