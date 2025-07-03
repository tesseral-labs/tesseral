import React, { useState } from "react";
import { Outlet, useLocation, useNavigate, useParams } from "react-router";

import { Step, Steps } from "@/components/core/flows/Steps";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

const steps = ["entra", "identifier", "reply-url", "metadata", "users"];

export function EntraSamlConnectionFlow() {
  const { samlConnectionId } = useParams();
  const navigate = useNavigate();
  const { pathname } = useLocation();

  const currentStep = steps.findIndex((step) => pathname.endsWith(step));

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
          <DialogTitle>Set up your Microsoft Entra SAML Connection</DialogTitle>
        </DialogHeader>

        <div className="space-y-6 w-full">
          <Steps>
            <Step
              label="Create app"
              status={currentStep === 0 ? "active" : "completed"}
            />
            <Step
              label="Identifier"
              status={
                currentStep < 1
                  ? "pending"
                  : currentStep === 1
                    ? "active"
                    : "completed"
              }
            />
            <Step
              label="Reply URL"
              status={
                currentStep < 2
                  ? "pending"
                  : currentStep === 2
                    ? "active"
                    : "completed"
              }
            />
            <Step
              label="Metadata"
              status={
                currentStep < 3
                  ? "pending"
                  : currentStep === 3
                    ? "active"
                    : "completed"
              }
            />
            <Step
              label="Users"
              status={
                currentStep < 4
                  ? "pending"
                  : currentStep === 4
                    ? "active"
                    : "completed"
              }
            />
          </Steps>

          <Outlet />
        </div>
      </DialogContent>
    </Dialog>
  );
}
