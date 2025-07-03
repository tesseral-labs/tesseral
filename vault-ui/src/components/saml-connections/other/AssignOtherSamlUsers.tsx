import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

export function AssignOtherSamlUsers() {
  const { samlConnectionId } = useParams();

  return (
    <>
      <div className="space-y-4 text-sm">
        <p className="font-medium">Assign users to the new app</p>

        <p>
          In your Identity Provider, go to your new SAML application's settings
          related to user or group assignments. Assign the appropriate users to
          the application.
        </p>
        <p>
          If you intend to test the application yourself, remember to include
          yourself.
        </p>
      </div>

      <DialogFooter>
        <Link to={`/organization/saml-connections/${samlConnectionId}/test`}>
          <Button type="button" size="sm">
            Continue
          </Button>
        </Link>
      </DialogFooter>
    </>
  );
}
