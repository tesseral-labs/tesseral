import { useMutation } from "@connectrpc/connect-query";
import React from "react";
import { useNavigate } from "react-router";
import { Link, useSearchParams } from "react-router-dom";
import { toast } from "sonner";

import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { createRole } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { EditRoleForm } from "@/pages/dashboard/roles/EditRoleForm";

export function CreateRolePage() {
  const [searchParams] = useSearchParams();
  const { mutateAsync: createRoleAsync } = useMutation(createRole);
  const navigate = useNavigate();

  async function handleSubmit(role: {
    displayName: string;
    description: string;
    actions: string[];
  }) {
    const createRoleResponse = await createRoleAsync({
      role: {
        organizationId: searchParams.get("organization-id") ?? "",
        displayName: role.displayName,
        description: role.description,
        actions: role.actions,
      },
    });
    toast.success("Role created");
    navigate(`/organization-settings/advanced`);
  }

  return (
    <div>
      <Card>
        <CardHeader>
          <CardTitle>Create Custom Role</CardTitle>
          <CardDescription>
            Custom roles allow you to define a custom set of permissions you can
            assign to users in your organization.
          </CardDescription>
        </CardHeader>
      </Card>

      <div className="mt-8">
        <EditRoleForm
          role={{ displayName: "", description: "", actions: [] }}
          onSubmit={handleSubmit}
        />
      </div>
    </div>
  );
}
