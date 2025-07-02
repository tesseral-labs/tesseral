import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

export function CreateOtherSamlApplication() {
  const { samlConnectionId } = useParams();

  return (
    <>
      <div className="space-y-4 text-sm">
        <p className="font-medium">Create a new SAML application:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>Go to your SAML provider's admin console.</li>
          <li>Create a new SAML application.</li>
        </ol>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/other/configure`}
        >
          <Button type="button" size="sm">
            Continue
          </Button>
        </Link>
      </DialogFooter>
    </>
  );
}
