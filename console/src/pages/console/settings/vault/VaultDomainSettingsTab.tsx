import { useQuery } from "@connectrpc/connect-query";
import React from "react";

import { getVaultDomainSettings } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

import { VaultCookieDomainCard } from "./domains/VaultCookieDomainCard";
import { VaultDomainRecordsCard } from "./domains/VaultDomainRecordsCard";
import { VaultDomainsCard } from "./domains/VaultDomainsCard";
import { VaultEmailSendFromDomainRecordsCard } from "./domains/VaultEmailSendFromDomainRecordsCard";
import { VaultTrustedDomainsCard } from "./domains/VaultTrustedDomainsCard";

export function VaultDomainSettingsTab() {
  const { data: getVaultDomainSettingsResponse } = useQuery(
    getVaultDomainSettings,
  );

  return (
    <div className="space-y-8">
      <div className="gap-8 grid grid-cols-3">
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
