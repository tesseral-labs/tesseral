import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { AlertDialogDescription } from "@radix-ui/react-alert-dialog";
import { Plus } from "lucide-react";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { TableSkeleton } from "@/components/skeletons/TableSkeleton";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  createUserRoleAssignment,
  deleteUserRoleAssignment,
  getRole,
  getUser,
  listRoles,
  listUserRoleAssignments,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { UserRoleAssignment } from "@/gen/tesseral/frontend/v1/models_pb";

export function OrganizationUserRolesTab() {
  const { userId } = useParams();

  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const {
    data: listUserRoleAssignmentsResponse,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listUserRoleAssignments,
    {
      userId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );

  const user = getUserResponse?.user;
  const userRoleAssignments =
    listUserRoleAssignmentsResponse?.pages.flatMap(
      (page) => page.userRoleAssignments,
    ) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Role assignments</CardTitle>
        <CardDescription>
          Roles assigned to <span className="font-semibold">{user?.email}</span>
          .
        </CardDescription>
        <CardAction>
          <AssignRoleButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton />
        ) : (
          <>
            {userRoleAssignments.length === 0 ? (
              <div className="text-muted-foreground pt-8 text-center text-sm">
                No role assignments found for this user.
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Role</TableHead>
                    <TableHead>Actions</TableHead>
                    <TableHead></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {userRoleAssignments?.map((roleAssignment) => (
                    <RoleAssignmentRow
                      key={roleAssignment.id}
                      roleAssignment={roleAssignment}
                    />
                  ))}
                </TableBody>
              </Table>
            )}
          </>
        )}
      </CardContent>
      {hasNextPage && (
        <CardFooter className="justify-center">
          <Button
            variant="outline"
            size="sm"
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
          >
            Load more
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

function RoleAssignmentRow({
  roleAssignment,
}: {
  roleAssignment: UserRoleAssignment;
}) {
  const { data: getRoleResponse } = useQuery(getRole, {
    id: roleAssignment.roleId,
  });
  const role = getRoleResponse?.role;

  return (
    <TableRow>
      <TableCell className="font-medium">
        {role?.displayName || role?.id}
      </TableCell>
      <TableCell className="flex flex-wrap gap-2">
        {role?.actions?.map((action) => (
          <Badge key={action} variant="secondary">
            {action}
          </Badge>
        ))}
      </TableCell>
      <TableCell className="text-right">
        <UnassignRoleButton roleAssignment={roleAssignment} />
      </TableCell>
    </TableRow>
  );
}

const schema = z.object({
  roleId: z.string().min(1, "Role is required"),
  userId: z.string().min(1, "User ID is required"),
});

function AssignRoleButton() {
  const { userId } = useParams();
  const { refetch } = useInfiniteQuery(
    listUserRoleAssignments,
    {
      userId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const createUserRoleAssignmentMutation = useMutation(
    createUserRoleAssignment,
  );
  const { data: listRolesResponse } = useQuery(listRoles);

  const roles = listRolesResponse?.roles || [];
  const user = getUserResponse?.user;

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      roleId: "",
      userId: userId,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await createUserRoleAssignmentMutation.mutateAsync({
        userRoleAssignment: {
          userId: data.userId,
          roleId: data.roleId,
        },
      });
      await refetch();
      form.reset();
      setOpen(false);
      toast.success("Role assigned successfully.");
    } catch {
      toast.error("Failed to assign role. Please try again.");
      setOpen(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus />
          Assign Role
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Assign Role</DialogTitle>
          <DialogDescription>
            Assign a role to {user?.email}. Select the role from the list below.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="roleId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Role</FormLabel>
                  <FormDescription>
                    The rolw to assign to {user?.email}.
                  </FormDescription>
                  <FormMessage />
                  <FormControl>
                    <Select onValueChange={field.onChange} value={field.value}>
                      <SelectTrigger className="w-[180px]">
                        <SelectValue
                          className="max-w-full overflow-x-hidden"
                          placeholder="Select a Role"
                        />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          <SelectLabel>Available Roles</SelectLabel>
                          {roles.map((role) => (
                            <SelectItem
                              className="max-w-full overflow-hidden"
                              key={role.id}
                              value={role.id}
                            >
                              {role.displayName || role.id}
                            </SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </FormControl>
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button variant="outline" onClick={() => setOpen(false)}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={
                  !form.formState.isDirty ||
                  createUserRoleAssignmentMutation.isPending
                }
              >
                Assign Role
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

function UnassignRoleButton({
  roleAssignment,
}: {
  roleAssignment: UserRoleAssignment;
}) {
  const { userId } = useParams();

  const { refetch } = useInfiniteQuery(
    listUserRoleAssignments,
    {
      userId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const deleteUserRoleAssignmentMutation = useMutation(
    deleteUserRoleAssignment,
  );

  const [open, setOpen] = useState(false);

  async function handleUnassign() {
    try {
      await deleteUserRoleAssignmentMutation.mutateAsync({
        id: roleAssignment.id,
      });
      await refetch();
      setOpen(false);
      toast.success("Role unassigned successfully.");
    } catch {
      toast.error("Failed to unassign role. Please try again.");
      setOpen(false);
    }
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button
          className="border-destructive text-destructive hover:bg-destructive hover:text-white"
          variant="outline"
          size="sm"
          onClick={() => setOpen(true)}
        >
          Unassign
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This will remove the role assignment for this user.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleUnassign}>
            Unassign Role
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
