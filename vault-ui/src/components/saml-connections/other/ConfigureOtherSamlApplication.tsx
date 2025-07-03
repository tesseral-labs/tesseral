import { useQuery } from "@connectrpc/connect-query";
import React from "react";
import { Link, useParams } from "react-router";

import { ValueCopier } from "@/components/core/ValueCopier";
import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";
import { getSAMLConnection } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function ConfigureOtherSamlApplication() {
  const { samlConnectionId } = useParams();

  const { data: getSamlConnectionResponse } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  });

  const samlConnection = getSamlConnectionResponse?.samlConnection;

  return (
    <>
      <div className="space-y-4 text-sm">
        <p className="font-medium">Configure your SAML application</p>
        <p>
          Inside your Identity Provider, configure your application's service
          provider settings.
        </p>

        <p>
          Your service provider will ask for a Service Provider ACS URL, or some
          variation of one of these names:
        </p>

        <ul className="list-disc list-inside">
          <li>Assertion Consumer Service URL</li>
          <li>SAML Start URL</li>
          <li>SAML Sign-On URL</li>
        </ul>

        <p>These terms all refer to the same concept. Input the following:</p>
        <p>
          <ValueCopier
            value={samlConnection?.spAcsUrl || ""}
            maxLength={50}
            label="ACS URL"
          />
        </p>

        <p>
          Your service provider will also ask for a Service Provider Entity ID,
          or some variation of:
        </p>

        <ul className="list-disc list-inside">
          <li>Service Provider Entity ID</li>
          <li>Audience URI</li>
          <li>Relying Party Identifier / ID</li>
        </ul>

        <p>These terms all refer to the same concept. Input the following:</p>
        <p>
          <ValueCopier
            value={samlConnection?.spEntityId || ""}
            maxLength={50}
            label="Entity ID"
          />
        </p>
      </div>

      <DialogFooter>
        <Link
          to={`/organization/saml-connections/${samlConnectionId}/setup/other/metadata`}
        >
          <Button type="button" size="sm">
            Continue
          </Button>
        </Link>
      </DialogFooter>
    </>
  );
}
