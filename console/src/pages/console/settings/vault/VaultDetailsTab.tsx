import React from "react";

import { VaultBehaviorSettingsCard } from "./details/VaultBehaviorSettingsCard";

export function VaultDetailsTab() {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <VaultBehaviorSettingsCard />
    </div>
  );
}
