import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight, ListCheck, Plus } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
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
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { UserRoleAssignment } from "@/gen/tesseral/backend/v1/models_pb";

export function UserRolesTab() {
  const { userId } = useParams();
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const {
    data: listUserRoleAssignmentsResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery(
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

  const userRoleAssignments =
    listUserRoleAssignmentsResponses?.pages.flatMap(
      (page) => page.userRoleAssignments,
    ) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>User Roles</CardTitle>
        <CardDescription>
          Roles assigned to{" "}
          <span className="font-semibold">
            {getUserResponse?.user?.displayName || getUserResponse?.user?.email}
          </span>
          .
        </CardDescription>
        <CardAction>
          <AssignRoleButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {userRoleAssignments.length === 0 ? (
          <div className="col-span-3 text-center text-sm text-muted-foreground">
            No roles assigned to this user.
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
              {userRoleAssignments.map((assignment) => (
                <UserRoleAssignmentRow
                  key={assignment.roleId}
                  userRoleAssignment={assignment}
                />
              ))}
            </TableBody>
          </Table>
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
            Load More
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

function UserRoleAssignmentRow({
  userRoleAssignment,
}: {
  userRoleAssignment: UserRoleAssignment;
}) {
  const { roleId } = userRoleAssignment;
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
  const { data: getRoleResponse } = useQuery(getRole, {
    id: roleId,
  });
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });

  const deleteUserRoleAssignmentMutation = useMutation(
    deleteUserRoleAssignment,
  );

  const [unassignOpen, setUnassignOpen] = useState(false);

  async function handleUnassign() {
    await deleteUserRoleAssignmentMutation.mutateAsync({
      id: userRoleAssignment.id,
    });
    await refetch();
    setUnassignOpen(false);
    toast.success("Role unassigned successfully.");
  }

  return (
    <>
      <TableRow>
        <TableCell className="font-medium">
          {getRoleResponse?.role?.displayName || userRoleAssignment.roleId}
        </TableCell>
        <TableCell>
          {getRoleResponse?.role?.actions?.map((action) => (
            <Badge key={action} variant="secondary">
              {action}
            </Badge>
          ))}
        </TableCell>
        <TableCell className="text-right">
          <Button
            variant="outline"
            size="sm"
            className="border-destructive text-destructive hover:bg-destructive hover:text-white"
            onClick={() => setUnassignOpen(true)}
          >
            Unassign
          </Button>
        </TableCell>
      </TableRow>

      {/* Unassign Confirmation Dialog */}
      <AlertDialog open={unassignOpen} onOpenChange={setUnassignOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Unassign Role</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to unassign the role{" "}
              <span className="font-semibold">
                {getRoleResponse?.role?.displayName ||
                  getRoleResponse?.role?.id}
              </span>{" "}
              from this{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.email}
              </span>
              ?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="flex justify-end items-center gap-2 mt-8">
            <Button variant="outline" onClick={() => setUnassignOpen(false)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleUnassign}
              disabled={!userRoleAssignment.id}
            >
              Unassign Role
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

const schema = z.object({
  roleId: z.string(),
  userId: z.string(),
});

function AssignRoleButton() {
  const { organizationId, userId } = useParams();

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
  const { data: listProjectRolesResponse } = useQuery(listRoles, {});
  const { data: listOrganizationRolesResponse } = useQuery(
    listRoles,
    {
      organizationId,
    },
    {
      enabled: !!getUserResponse?.user?.organizationId,
    },
  );
  const createUserRoleAssignmentMutation = useMutation(
    createUserRoleAssignment,
  );

  const roles = [
    ...(listProjectRolesResponse?.roles || []),
    ...(listOrganizationRolesResponse?.roles || []),
  ];

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      roleId: "",
      userId: userId || "",
    },
  });

  function handleCancel(e: React.MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    form.reset();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await createUserRoleAssignmentMutation.mutateAsync({
      userRoleAssignment: {
        roleId: data.roleId,
        userId: data.userId,
      },
    });
    await refetch();
    form.reset();
    setOpen(false);
    toast.success("Role assigned successfully.");
  }

  useEffect(() => {
    if (getUserResponse?.user) {
      form.reset({
        userId: getUserResponse.user.id,
      });
    }
  }, [getUserResponse, form]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus />
          Assign Role
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Asign Role</DialogTitle>
          <DialogDescription>
            Select a role to assign to{" "}
            <span className="font-semibold">
              {getUserResponse?.user?.email}
            </span>
            .
          </DialogDescription>
        </DialogHeader>
        {!roles.length ? (
          <div className="text-sm text-muted-foreground">
            No roles available to assign. Please create a role first.
            <Link to="/settings/access">
              <Button variant="outline" className="mt-4">
                <ListCheck />
                Manage Roles
                <ArrowRight />
              </Button>
            </Link>
          </div>
        ) : (
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              <div className="space-y-4">
                <FormField
                  control={form.control}
                  name="roleId"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Role</FormLabel>
                      <FormDescription>
                        The Role to assign to{" "}
                        <span className="font-semibold">
                          {getUserResponse?.user?.email}
                        </span>
                        .
                      </FormDescription>
                      <FormMessage />
                      <FormControl className="w-full">
                        <Select
                          onValueChange={field.onChange}
                          value={field.value}
                        >
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
                                  {!!role.displayName && (
                                    <>
                                      <span>{role.displayName}</span>
                                      <span> - </span>
                                    </>
                                  )}
                                  {role.id}
                                </SelectItem>
                              ))}
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>
              <CardFooter className="flex justify-end items-center gap-2 mt-8">
                <Button variant="outline" onClick={handleCancel}>
                  Cancel
                </Button>
                <Button type="submit" disabled={!form.formState.isDirty}>
                  Assign Role
                </Button>
              </CardFooter>
            </form>
          </Form>
        )}
      </DialogContent>
    </Dialog>
  );
}
