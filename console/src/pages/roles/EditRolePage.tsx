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
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createRole,
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
    <div>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          {getRoleResponse?.role?.organizationId ? (
            <>
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link to="/organizations">Organizations</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link
                    to={`/organizations/${getRoleResponse?.role?.organizationId}`}
                  >
                    {getOrganizationResponse?.organization?.displayName}
                  </Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link
                    to={`/organizations/${getRoleResponse?.role?.organizationId}/roles`}
                  >
                    Roles
                  </Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage>
                  {getRoleResponse?.role?.displayName}
                </BreadcrumbPage>
              </BreadcrumbItem>
            </>
          ) : (
            <>
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link to="/project-settings">Project settings</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink asChild>
                  <Link to="/project-settings/rbac-settings">Roles</Link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage>
                  {getRoleResponse?.role?.displayName}
                </BreadcrumbPage>
              </BreadcrumbItem>
            </>
          )}
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>Edit Role</PageTitle>
      <PageDescription>
        Roles are a named collection of Actions, and can be assigned to Users.
      </PageDescription>

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
