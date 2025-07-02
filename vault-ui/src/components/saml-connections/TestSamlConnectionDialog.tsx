import { useQuery } from "@connectrpc/connect-query";
import React, { useState } from "react";
import { Link, useNavigate, useParams } from "react-router";

import { getProject } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog";

export function TestSamlConnectionDialog() {
  const { samlConnectionId } = useParams();
  const navigate = useNavigate();

  const { data: getProjectResponse } = useQuery(getProject);

  const project = getProjectResponse?.project;

  const [open, setOpen] = useState(true);

  function handleOpenChange(open: boolean) {
    setOpen(open);
    navigate(`/organization/saml-connections/${samlConnectionId}`);
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Test your SAML Connection</DialogTitle>
          <DialogDescription className="space-y-2">
            <p>You've finished configuring your SAML connection.</p>
            <p>You can now test it to ensure everything is set up correctly.</p>
          </DialogDescription>
        </DialogHeader>

        <DialogFooter className="mt-4">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => handleOpenChange(false)}
          >
            Skip testing
          </Button>
          <Link
            to={`https://${project?.vaultDomain}/api/saml/v1/${samlConnectionId}/init`}
            target="_blank"
          >
            <Button type="button" disabled={!project?.vaultDomain} size="sm">
              Test SAML Connection
            </Button>
          </Link>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
