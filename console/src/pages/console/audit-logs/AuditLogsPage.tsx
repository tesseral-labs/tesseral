import { Logs } from "lucide-react";
import React from "react";

import { PageContent } from "@/components/page";

export function AuditLogsPage() {
  return (
    <PageContent>
      <div className="">
        <h1 className="text-2xl font-semibold flex items-center gap-x-2">
          <Logs />
          <span>Audit Logs</span>
        </h1>
      </div>
    </PageContent>
  );
}
