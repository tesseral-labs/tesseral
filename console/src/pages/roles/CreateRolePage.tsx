import React from 'react';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { Link, useSearchParams } from 'react-router-dom';
import { PageDescription, PageTitle } from '@/components/page';
import { EditRoleForm } from '@/pages/roles/EditRoleForm';
import { useMutation } from '@connectrpc/connect-query';
import { createRole } from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useNavigate, useParams } from 'react-router';
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
    <div>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Create Role</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>Create Role</PageTitle>
      <PageDescription>
        Roles are a named collection of Actions, and can be assigned to Users.
      </PageDescription>

      <div className="mt-8">
        <EditRoleForm
          role={{ displayName: '', description: '', actions: [] }}
          onSubmit={handleSubmit}
        />
      </div>
    </div>
  );
}
