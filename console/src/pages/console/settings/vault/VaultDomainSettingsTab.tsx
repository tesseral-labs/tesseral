import { useQuery } from "@connectrpc/connect-query";
import { ArrowRight, Crown } from "lucide-react";
import React from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  getProjectEntitlements,
  getVaultDomainSettings,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { useHandleUpgrade } from "@/hooks/use-handle-upgrade";

import { VaultCookieDomainCard } from "./domains/VaultCookieDomainCard";
import { VaultDomainRecordsCard } from "./domains/VaultDomainRecordsCard";
import { VaultDomainsCard } from "./domains/VaultDomainsCard";
import { VaultEmailSendFromDomainRecordsCard } from "./domains/VaultEmailSendFromDomainRecordsCard";
import { VaultTrustedDomainsCard } from "./domains/VaultTrustedDomainsCard";

export function VaultDomainSettingsTab() {
  const handleUpgrade = useHandleUpgrade();

  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
  );
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  return (
    <div className="space-y-8">
      {!getProjectEntitlementsResponse?.entitledCustomVaultDomains && (
        <Card className="lg:col-span-1 bg-gradient-to-br from-violet-500 via-purple-500 to-blue-500 border-0 text-white relative overflow-hidden shadow-xl">
          <div className="absolute inset-0 bg-gradient-to-br from-white/10 to-transparent" />
          <CardContent className="relative flex h-full">
            <div className="flex flex-wrap items-center justify-between w-full">
              <div className="space-y-4 md:flex-grow">
                <div className="flex items-center space-x-3">
                  <div className="p-2 rounded-full bg-white/20 backdrop-blur-sm">
                    <Crown className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-white">
                      Upgrade to Growth
                    </h3>
                    <p className="text-xs text-white/80">
                      Unlock advanced features like custom domains.
                    </p>
                  </div>
                </div>
              </div>

              <Button
                className="bg-white text-purple-600 hover:bg-white/90 font-medium cursor-pointer"
                onClick={handleUpgrade}
                size="lg"
              >
                Upgrade Now
                <ArrowRight className="h-4 w-4 ml-2" />
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
      <div className="gap-8 grid grid-cols-1 lg:grid-cols-3">
        <VaultDomainsCard />
        <VaultCookieDomainCard />
        <VaultTrustedDomainsCard />
      </div>
      {getVaultDomainSettingsResponse?.vaultDomainSettings?.pendingDomain && (
        <>
          <VaultDomainRecordsCard />
          <VaultEmailSendFromDomainRecordsCard />
        </>
      )}
    </div>
  );
}
