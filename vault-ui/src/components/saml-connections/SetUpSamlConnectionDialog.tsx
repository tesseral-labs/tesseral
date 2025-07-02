import { ChevronRight, ShieldEllipsis } from "lucide-react";
import React, { useState } from "react";
import { Link, useNavigate, useParams } from "react-router";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog";
import { Label } from "../ui/label";
import { Separator } from "../ui/separator";

export function SetUpSamlConnectionDialog() {
  const { samlConnectionId } = useParams();
  const navigate = useNavigate();

  const [open, setOpen] = useState(true);

  function handleOpenChange(open: boolean) {
    setOpen(open);
    if (!open) {
      navigate(`/organization/saml-connections/${samlConnectionId}`);
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Choose your Identity Provider.</DialogTitle>
          <DialogDescription>
            Select the Identify Provider you use to start setting up your SAML
            Connection.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <Link
            className="flex items-center gap-4"
            to={`/organization/saml-connections/${samlConnectionId}/setup/okta`}
          >
            <img
              className="w-8 h-8"
              src="/images/saml-connections/logo-okta.svg"
              alt="Okta"
            />

            <div className="space-y-1">
              <Label className="cursor-pointer">Okta</Label>
              <div className="text-xs text-muted-foreground">
                Set up your SAML connection with your corportate Okta.
              </div>
            </div>

            <div className="flex items-center ml-auto text-muted-foreground">
              <ChevronRight className="w-4 h-4" />
            </div>
          </Link>

          <Separator />

          <Link
            className="flex items-center gap-4"
            to={`/organization/saml-connections/${samlConnectionId}/setup/google`}
          >
            <img
              className="w-8 h-8"
              src="/images/saml-connections/logo-google.svg"
              alt="Okta"
            />

            <div className="space-y-1">
              <Label className="cursor-pointer">Google</Label>
              <div className="text-xs text-muted-foreground">
                Set up your SAML connection with your Google Workspace.
              </div>
            </div>

            <div className="flex items-center ml-auto text-muted-foreground">
              <ChevronRight className="w-4 h-4" />
            </div>
          </Link>

          <Separator />

          <Link
            className="flex items-center gap-4"
            to={`/organization/saml-connections/${samlConnectionId}/setup/entra`}
          >
            <img
              className="w-8 h-8"
              src="/images/saml-connections/logo-entra.svg"
              alt="Okta"
            />

            <div className="space-y-1">
              <Label className="cursor-pointer">Microsoft Entra</Label>
              <div className="text-xs text-muted-foreground">
                Set up your SAML connection with your Microsoft Entra.
              </div>
            </div>

            <div className="flex items-center ml-auto text-muted-foreground">
              <ChevronRight className="w-4 h-4" />
            </div>
          </Link>

          <Separator />

          <Link
            className="flex items-center gap-4"
            to={`/organization/saml-connections/${samlConnectionId}/setup/other`}
          >
            <ShieldEllipsis className="h-8" />

            <div className="space-y-1">
              <Label className="cursor-pointer">Other</Label>
              <div className="text-xs text-muted-foreground">
                Set up your SAML connection with another identity provider.
              </div>
            </div>

            <div className="flex items-center ml-auto text-muted-foreground">
              <ChevronRight className="w-4 h-4" />
            </div>
          </Link>
        </div>
      </DialogContent>
    </Dialog>
  );
}
