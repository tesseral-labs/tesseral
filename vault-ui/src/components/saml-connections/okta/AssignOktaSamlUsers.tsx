import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

export function AssignOktaSamlUsers() {
  const { samlConnectionId } = useParams();

  return (
    <>
      <div className="space-y-4 text-sm">
        <img
          className="rounded-xl max-w-full border shadow-md"
          src="/videos/saml-setup-wizard/okta/users.gif"
        />

        <p className="font-medium">Assign users to the new app.</p>
        <p>
          If you intend to test the connection yourself, remember to assign
          yourself too.
        </p>

        <ol className="list-decimal list-inside space-y-2">
          <li>Click on the "Assignments" tab</li>
          <li>Click "Assign"</li>
          <li>
            Click "Assign to People" or "Assign to Groups", whichever you
            usually use
          </li>
        </ol>
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
