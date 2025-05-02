import { z } from 'zod';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createUserRoleAssignment,
  deleteUserRoleAssignment,
  getRBACPolicy,
  getRole,
  getUser,
  listRoles,
  listUserRoleAssignments,
  updateRole,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useParams } from 'react-router';
import { useForm } from 'react-hook-form';
import React, { useEffect, useState } from 'react';
import { toast } from 'sonner';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent, AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Checkbox } from '@/components/ui/checkbox';

const schema = z.object({
  roles: z.array(z.string()),
});

export function AssignUserRolesButton() {
  const { userId } = useParams();
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });

  const { data: listUserRoleAssignmentsResponse, refetch } = useQuery(listUserRoleAssignments, {
    userId,
  });

  const { data: listProjectRolesResponse } = useQuery(listRoles, {});
  const { data: listOrganizationRolesResponse } = useQuery(
    listRoles,
    {
      organizationId: getUserResponse?.user?.organizationId,
    },
    {
      enabled: !!getUserResponse?.user?.organizationId,
    },
  );

  const roles = [
    ...(listProjectRolesResponse?.roles || []),
    ...(listOrganizationRolesResponse?.roles || []),
  ];

  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      roles: [],
    },
  });

  const { mutateAsync: createUserRoleAssignmentAsync } = useMutation(
    createUserRoleAssignment,
  );
  const { mutateAsync: deleteUserRoleAssignmentAsync } = useMutation(
    deleteUserRoleAssignment,
  );

  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (listUserRoleAssignmentsResponse?.userRoleAssignments) {
      form.reset({
        roles: listUserRoleAssignmentsResponse?.userRoleAssignments?.map(
          (userRoleAssignment) => userRoleAssignment.roleId,
        ),
      });
    }
  }, [listUserRoleAssignmentsResponse]);

  const handleSubmit = async (data: z.infer<typeof schema>) => {
    // create new assignments for new roles, and delete assignments for old
    // roles; no-op otherwise

    // appease typescript in the unlikely event of hitting "Save" before having
    // loaded in user role assignments
    if (!listUserRoleAssignmentsResponse) {
      return;
    }

    const oldRoles = listUserRoleAssignmentsResponse.userRoleAssignments.map(
      (userRoleAssignment) => userRoleAssignment.roleId,
    );
    const newRoles = data.roles;

    const createAssignments = newRoles.filter((role) => !oldRoles?.includes(role));
    const deleteAssignments = oldRoles.filter((role) => !newRoles?.includes(role));

    await Promise.all(
      createAssignments.map((roleId) =>
        createUserRoleAssignmentAsync({
          userRoleAssignment: {
            userId,
            roleId,
          }
        }),
      ),
    );

    await Promise.all(
      deleteAssignments.map((roleId) =>
        deleteUserRoleAssignmentAsync({
          id: listUserRoleAssignmentsResponse.userRoleAssignments.find(
            (userRoleAssignment) => userRoleAssignment.roleId === roleId,
          )!.id,
        }),
      ),
    );

    await refetch();
    setOpen(false);
    toast.success('User Role assignments updated successfully');
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">Assign Roles</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Assign User to Roles</AlertDialogTitle>
          <AlertDialogDescription>
            Users can be assigned to one or more Roles.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-8"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <div className="space-y-4">
              {roles.map((role) => (
                <FormField
                  key={role.id}
                  control={form.control}
                  name="roles"
                  render={({ field }) => {
                    return (
                      <FormItem
                        key={role.id}
                        className="flex flex-row items-start space-x-3 space-y-0"
                      >
                        <FormControl>
                          <Checkbox
                            checked={field.value?.includes(role.id)}
                            onCheckedChange={(checked) => {
                              return checked
                                ? field.onChange([...field.value, role.id])
                                : field.onChange(
                                    field.value?.filter(
                                      (value) => value !== role.id,
                                    ),
                                  );
                            }}
                          />
                        </FormControl>
                        <div className="grid gap-1.5 leading-none">
                          <FormLabel className="font-normal">
                            {role.displayName}
                          </FormLabel>
                          <FormDescription>
                            {role.description}
                          </FormDescription>
                        </div>
                      </FormItem>
                    );
                  }}
                />
              ))}
            </div>

            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
