import { useMutation } from "@connectrpc/connect-query";
import { ArrowRight, CheckCircle2, Crown } from "lucide-react";
import React from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { createStripeCheckoutLink } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function UpgradeCard() {
  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  async function handleUpgrade() {
    const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
    window.location.href = url;
  }
  return (
    <Card className="lg:col-span-1 bg-gradient-to-br from-violet-500 via-purple-500 to-blue-500 border-0 text-white relative overflow-hidden shadow-xl">
      <div className="absolute inset-0 bg-gradient-to-br from-white/10 to-transparent" />
      <CardContent className="p-6 relative flex h-full">
        <div className="flex flex-wrap w-full">
          <div className="w-full space-y-4 md:flex-grow">
            <div className="flex items-center space-x-3">
              <div className="p-2 rounded-full bg-white/20 backdrop-blur-sm">
                <Crown className="h-6 w-6 text-white" />
              </div>
              <div>
                <h3 className="font-semibold text-white">Upgrade to Growth</h3>
                <p className="text-xs text-white/80">
                  Unlock advanced features
                </p>
              </div>
            </div>

            <div className="space-y-4 flex-grow">
              <p className="font-semibold">
                Upgrade to the Growth tier to unlock more features.
              </p>
              <div className="text-sm space-y-2">
                <div className="flex items-center">
                  <CheckCircle2 className="inline h-4 w-4 mr-2" />
                  Run Tesseral on your own domain
                </div>
                <div className="flex items-center">
                  <CheckCircle2 className="inline h-4 w-4 mr-2" />
                  Manage Tesseral resources over an API
                </div>
                <div className="flex items-center">
                  <CheckCircle2 className="inline h-4 w-4 mr-2" />
                  Use Tesseral to manage API keys
                </div>
                <div className="flex items-center">
                  <CheckCircle2 className="inline h-4 w-4 mr-2" />
                  Real-time sync with webhooks
                </div>
                <div className="flex items-center">
                  <CheckCircle2 className="inline h-4 w-4 mr-2" />
                  Dedicated support over email
                </div>
              </div>
            </div>
          </div>

          <div className="md:mt-auto w-full">
            <Button
              className="w-full mt-8 bg-white text-purple-600 hover:bg-white/90 font-medium cursor-pointer"
              onClick={handleUpgrade}
              size="lg"
            >
              Upgrade Now
              <ArrowRight className="h-4 w-4 ml-2" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
