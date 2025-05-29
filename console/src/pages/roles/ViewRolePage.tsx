import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  deleteRole,
  deleteUserRoleAssignment,
  getOrganization,
  getRole,
  getUser,
  listUserRoleAssignments,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import React, { useState } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Link } from 'react-router-dom';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import {
  ConsoleCard,
  ConsoleCardDetails,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardHeader,
  ConsoleCardTitle,
  ConsoleCardTableContent,
} from '@/components/ui/console-card';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { Button, buttonVariants } from '@/components/ui/button';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { toast } from 'sonner';
import { UserRoleAssignment } from '@/gen/tesseral/backend/v1/models_pb';

export function ViewRolePage() {
  const { roleId } = useParams();
  const { data: getRoleResponse } = useQuery(getRole, {
    id: roleId,
  });
  const { data: listUserRoleAssignmentsResponse } = useQuery(
    listUserRoleAssignments,
    {
      roleId,
    },
  );
  const { data: getOrganizationResponse } = useQuery(
    getOrganization,
    {
      id: getRoleResponse?.role?.organizationId,
    },
    {
      enabled: !!getRoleResponse?.role?.organizationId,
    },
  );

  return (
    <>
      <PageHeader>
        <PageTitle>{getRoleResponse?.role?.displayName}</PageTitle>
        <PageCodeSubtitle>{roleId}</PageCodeSubtitle>
        <PageDescription>
          Roles are a named collection of Actions, and can be assigned to Users.
        </PageDescription>
      </PageHeader>

      <PageContent>
        <div className="space-y-8">
          <ConsoleCard>
            <ConsoleCardHeader className="flex-row justify-between items-center gap-x-2">
              <ConsoleCardDetails>
                <ConsoleCardTitle>General settings</ConsoleCardTitle>
                <ConsoleCardDescription>
                  Basic settings for this Role.
                </ConsoleCardDescription>
              </ConsoleCardDetails>
              <Link to={`/roles/${roleId}/edit`}>
                <Button variant="outline">Edit</Button>
              </Link>
            </ConsoleCardHeader>
            <ConsoleCardContent>
              <DetailsGrid>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Display Name</DetailsGridKey>
                    <DetailsGridValue>
                      {getRoleResponse?.role?.displayName}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                  <DetailsGridEntry>
                    <DetailsGridKey>Description</DetailsGridKey>
                    <DetailsGridValue>
                      {getRoleResponse?.role?.description}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                  <DetailsGridEntry>
                    <DetailsGridKey>Role Type</DetailsGridKey>
                    <DetailsGridValue>
                      {getRoleResponse?.role?.organizationId
                        ? 'Organization-Specific Role'
                        : 'Project-Level Role'}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                  {getRoleResponse?.role?.organizationId && (
                    <DetailsGridEntry>
                      <DetailsGridKey>Organization</DetailsGridKey>
                      <DetailsGridValue>
                        <Link
                          className="underline underline-offset-2 decoration-muted-foreground/40"
                          to={`/organizations/${getRoleResponse?.role?.organizationId}`}
                        >
                          {getOrganizationResponse?.organization?.displayName}
                        </Link>
                      </DetailsGridValue>
                    </DetailsGridEntry>
                  )}
                </DetailsGridColumn>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Role Actions</DetailsGridKey>
                    <DetailsGridValue>
                      {getRoleResponse?.role?.actions.map((action) => (
                        <div>{action}</div>
                      ))}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
                <DetailsGridColumn>
                  <DetailsGridEntry>
                    <DetailsGridKey>Created</DetailsGridKey>
                    <DetailsGridValue>
                      {getRoleResponse?.role?.createTime &&
                        DateTime.fromJSDate(
                          timestampDate(getRoleResponse?.role?.createTime),
                        ).toRelative()}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                  <DetailsGridEntry>
                    <DetailsGridKey>Updated</DetailsGridKey>
                    <DetailsGridValue>
                      {getRoleResponse?.role?.updateTime &&
                        DateTime.fromJSDate(
                          timestampDate(getRoleResponse?.role?.updateTime),
                        ).toRelative()}
                    </DetailsGridValue>
                  </DetailsGridEntry>
                </DetailsGridColumn>
              </DetailsGrid>
            </ConsoleCardContent>
          </ConsoleCard>

          <ConsoleCard>
            <ConsoleCardHeader>
              <ConsoleCardTitle>Assigned Users</ConsoleCardTitle>
              <ConsoleCardDescription>
                Users assigned to this Role.
              </ConsoleCardDescription>
            </ConsoleCardHeader>
            <ConsoleCardTableContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>User Email</TableHead>
                    <TableHead>User Organization</TableHead>
                    <TableHead></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {listUserRoleAssignmentsResponse?.userRoleAssignments?.map(
                    (userRoleAssignment) => (
                      <UserRoleAssignmentRow
                        key={userRoleAssignment.roleId}
                        userRoleAssignment={userRoleAssignment}
                      />
                    ),
                  )}
                </TableBody>
              </Table>
            </ConsoleCardTableContent>
          </ConsoleCard>
        </div>
        <DangerZoneCard />
      </PageContent>
    </>
  );
}

function UserRoleAssignmentRow({
  userRoleAssignment,
}: {
  userRoleAssignment: UserRoleAssignment;
}) {
  const { refetch } = useQuery(listUserRoleAssignments, {
    roleId: userRoleAssignment.roleId,
  });
  const { data: getRoleResponse } = useQuery(getRole, {
    id: userRoleAssignment.roleId,
  });
  const { data: getUserResponse } = useQuery(getUser, {
    id: userRoleAssignment.userId,
  });
  const { data: getOrganizationResponse } = useQuery(
    getOrganization,
    {
      id: getUserResponse?.user?.organizationId,
    },
    {
      enabled: !!getUserResponse?.user?.organizationId,
    },
  );

  const { mutateAsync: deleteUserRoleAssignmentAsync } = useMutation(
    deleteUserRoleAssignment,
  );

  async function handleUnassign() {
    await deleteUserRoleAssignmentAsync({ id: userRoleAssignment.id });
    await refetch();
    toast.success('User unassigned');
  }

  const [open, setOpen] = useState(false);

  return (
    <TableRow>
      <TableCell>
        <Link
          className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
          to={`/organizations/${getUserResponse?.user?.organizationId}/${userRoleAssignment.userId}`}
        >
          {getUserResponse?.user?.email}
        </Link>
      </TableCell>
      <TableCell>
        <Link
          className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
          to={`/organizations/${getUserResponse?.user?.organizationId}`}
        >
          {getOrganizationResponse?.organization?.displayName}
        </Link>
      </TableCell>
      <TableCell className="text-right">
        <AlertDialog open={open} onOpenChange={setOpen}>
          <AlertDialogTrigger asChild>
            <Button size="sm" variant="link">
              Unassign
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Unassign Role</AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogDescription>
              Are you sure you want to unassign{' '}
              <span className="font-medium">
                {getUserResponse?.user?.email}
              </span>{' '}
              from the Role{' '}
              <span className="font-medium">
                {getRoleResponse?.role?.displayName}
              </span>
              ?
            </AlertDialogDescription>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction onClick={handleUnassign}>
                Unassign
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </TableCell>
    </TableRow>
  );
}

function DangerZoneCard() {
  return (
    <ConsoleCard className="mt-8 border-destructive">
      <ConsoleCardHeader>
        <ConsoleCardTitle>Danger Zone</ConsoleCardTitle>
      </ConsoleCardHeader>

      <ConsoleCardContent>
        <div className="flex justify-between items-center">
          <div>
            <div className="text-sm font-semibold">Delete Role</div>
            <p className="text-sm">
              Unassign all Users from this Role and delete this Role. This
              cannot be undone.
            </p>
          </div>

          <DeleteRoleButton />
        </div>
      </ConsoleCardContent>
    </ConsoleCard>
  );
}

function DeleteRoleButton() {
  const { roleId } = useParams();
  const { mutateAsync: deleteRoleAsync } = useMutation(deleteRole);

  const navigate = useNavigate();

  async function handleDelete() {
    await deleteRoleAsync({ id: roleId });
    navigate('/project-settings/rbac-settings');
    toast.success('Role deleted');
  }

  const [open, setOpen] = useState(false);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="destructive">Delete</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete Role</AlertDialogTitle>
        </AlertDialogHeader>
        <AlertDialogDescription>
          Are you sure you want to delete this Role? This cannot be undone.
        </AlertDialogDescription>

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            className={buttonVariants({ variant: 'destructive' })}
            onClick={handleDelete}
          >
            Delete Role
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
