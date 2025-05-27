import { useMutation, useQuery } from "@connectrpc/connect-query";
import React from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  createSAMLConnection,
  listSAMLConnections,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function CreateSAMLConnectionButton() {
  const navigate = useNavigate();

  const { data: listSAMLConnectionsResponse } = useQuery(listSAMLConnections);
  const createSAMLConnectionMutation = useMutation(createSAMLConnection);

  async function handleCreateSAMLConnection() {
    const { samlConnection } = await createSAMLConnectionMutation.mutateAsync({
      samlConnection: {
        // if there are no saml connections on the org yet, default to making
        // the first one be primary
        primary: !!listSAMLConnectionsResponse?.samlConnections,
      },
    });

    toast.success("SAML Connection created");
    navigate(`/organization-settings/saml-connections/${samlConnection?.id}`);
  }

  return (
    <Button variant="outline" onClick={handleCreateSAMLConnection}>
      Create SAML Connection
    </Button>
  );
}
