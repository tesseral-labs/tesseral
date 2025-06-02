import React from 'react';
import { useSearchParams } from 'react-router-dom';
import {
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import { EditRoleForm } from '@/pages/roles/EditRoleForm';
import { useMutation } from '@connectrpc/connect-query';
import { createRole } from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useNavigate } from 'react-router';
import { toast } from 'sonner';

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
        organizationId: searchParams.get('organization-id') ?? '',
        displayName: role.displayName,
        description: role.description,
        actions: role.actions,
      },
    });
    toast.success('Role created');
    navigate(`/roles/${createRoleResponse.role!.id}`);
  }

  return (
    <>
      <PageHeader>
        <PageTitle>Create Role</PageTitle>
        <PageDescription>
          Roles are a named collection of Actions, and can be assigned to Users.
        </PageDescription>
      </PageHeader>

      <PageContent>
        <EditRoleForm
          role={{ displayName: '', description: '', actions: [] }}
          onSubmit={handleSubmit}
        />
      </PageContent>
    </>
  );
}
