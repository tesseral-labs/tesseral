import React, { useState } from "react";
import { Outlet, useLocation, useNavigate, useParams } from "react-router";

import { Step, Steps } from "@/components/core/flows/Steps";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

const steps = ["other", "configure", "metadata", "users"];

export function OtherSamlConnectionFlow() {
  const { samlConnectionId } = useParams();
  const { pathname } = useLocation();
  const navigate = useNavigate();

  const [open, setOpen] = useState(true);

  const currentStep = steps.findIndex((step) => pathname.endsWith(step));

  function handleOpenChange(open: boolean) {
    setOpen(open);
    if (!open) {
      navigate(`/organization/saml-connections/${samlConnectionId}`);
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="w-full">
        <DialogHeader className="text-center">
          <DialogTitle>Set up your SAML Connection</DialogTitle>
        </DialogHeader>

        <div className="space-y-6 w-full">
          <Steps>
            <Step
              label="Create app"
              status={currentStep === 0 ? "active" : "completed"}
            />
            <Step
              label="Configure"
              status={
                currentStep < 1
                  ? "pending"
                  : currentStep === 1
                    ? "active"
                    : "completed"
              }
            />
            <Step
              label="Metadata"
              status={
                currentStep < 2
                  ? "pending"
                  : currentStep === 2
                    ? "active"
                    : "completed"
              }
            />
            <Step
              label="Users"
              status={
                currentStep < 3
                  ? "pending"
                  : currentStep === 3
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
