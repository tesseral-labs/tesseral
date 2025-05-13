import React from 'react';
import {
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import { EditRoleForm } from '@/pages/roles/EditRoleForm';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getOrganization,
  getRole,
  updateRole,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useNavigate, useParams } from 'react-router';
import { toast } from 'sonner';

export function EditRolePage() {
  const { roleId } = useParams();
  const { mutateAsync: updateRoleAsync } = useMutation(updateRole);
  const navigate = useNavigate();

  const { data: getRoleResponse } = useQuery(getRole, {
    id: roleId,
  });

  const { data: getOrganizationResponse } = useQuery(
    getOrganization,
    {
      id: getRoleResponse?.role?.organizationId,
    },
    {
      enabled: !!getRoleResponse?.role?.organizationId,
    },
  );

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
    toast.success('Role updated');
    navigate(`/roles/${roleId}`);
  }

  return (
    <>
      <PageHeader>
        <PageTitle>Edit Role</PageTitle>
        <PageDescription>
          Roles are a named collection of Actions, and can be assigned to Users.
        </PageDescription>
      </PageHeader>

      <PageContent>
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
      </PageContent>
    </>
  );
}
