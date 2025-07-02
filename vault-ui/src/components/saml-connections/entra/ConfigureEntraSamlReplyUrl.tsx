import { useQuery } from "@connectrpc/connect-query";
import React from "react";
import { Link, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";
import { getSAMLConnection } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function ConfigureEntraSamlReplyUrl() {
  const { samlConnectionId } = useParams();
  const { data: getSamlConnectionResponse } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  });

  const samlConnection = getSamlConnectionResponse?.samlConnection;

  return (
    <>
      <div className="space-y-4 text-sm">
        <p className="font-medium">
          Configure SAML Reply URL (Assertion Consumer Service URL)
        </p>
        <ol className="list-decimal list-inside space-y-2">
          <li>
            Find the "Reply URL (Assertion Consumer Service URL)" section below
            the "Identifier (Entity ID)" section
          </li>
          <li>Click "Add reply URL". An input now appears.</li>
          <li>
            Paste in this value:
            <ValueCopier
              value={samlConnection?.spAcsUrl || ""}
              label="Reply URL"
              maxLength={50}
            />
          </li>
          <li>
            Keep all other settings to their default values. Click "Save" above.
          </li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/entra/metadata`}
        >
          <Button size="sm" type="button">
            Continue
          </Button>
        </Link>
      </DialogFooter>
    </>
  );
}
