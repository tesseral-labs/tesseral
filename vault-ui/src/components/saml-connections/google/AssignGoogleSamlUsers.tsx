import React from "react";
import { Link, useParams } from "react-router";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";

export function AssignGoogleSamlUsers() {
  const { samlConnectionId } = useParams();

  return (
    <>
      <div className="space-y-4 text-sm">
        <img
          className="rounded-xl max-w-full border shadow-md"
          src="/videos/saml-setup-wizard/google/users.gif"
        />

        <p className="font-medium">Assign users to the new app.</p>
        <p>
          If you're familiar with Google Workspace organizational units, use
          whatever process you normally use.
        </p>
        <p>
          If you're not familiar with Google Workspace organizational units, or
          you don't normally use them, here's the simplest way to assign users
          to your new Google Workspace SAML application:
        </p>
        <ol className="list-decimal list-inside space-y-2">
          <li>
            To the right of "User access" is a chevron pointing down (). Click
            on it.
          </li>
          <li>Click on "ON for everyone"</li>
          <li>Click "Save".</li>
          <li>
            Allow a minute before users, including yourself, can log in. Google
            Workspace doesn't immediately reflect permissions updates.
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
