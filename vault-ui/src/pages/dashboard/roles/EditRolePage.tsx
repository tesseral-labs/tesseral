import { useMutation, useQuery } from "@connectrpc/connect-query";
import React from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";

import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  getRole,
  updateRole,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { EditRoleForm } from "@/pages/dashboard/roles/EditRoleForm";

export function EditRolePage() {
  const { roleId } = useParams();
  const { mutateAsync: updateRoleAsync } = useMutation(updateRole);
  const navigate = useNavigate();

  const { data: getRoleResponse } = useQuery(getRole, {
    id: roleId,
  });

  async function handleSubmit(role: {
    displayName: string;
    description: string;
    actions: string[];
  }) {
    await updateRoleAsync({
      id: roleId,
      role: {
        displayName: role.displayName,
        description: role.description,
        actions: role.actions,
      },
    });
    toast.success("Role updated");
    navigate(`/roles/${roleId}`);
  }

  return (
    <div>
      <Card>
        <CardHeader>
          <CardTitle>Edit Custom Role</CardTitle>
          <CardDescription>
            Custom roles allow you to define a custom set of permissions you can
            assign to users in your organization.
          </CardDescription>
        </CardHeader>
      </Card>

      <div className="mt-8">
        {getRoleResponse?.role && (
          <EditRoleForm
            role={{
              displayName: getRoleResponse.role.displayName,
              description: getRoleResponse.role.description,
              actions: getRoleResponse.role.actions,
            }}
            onSubmit={handleSubmit}
          />
        )}
      </div>
    </div>
  );
}
