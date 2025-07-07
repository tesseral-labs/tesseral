import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  ArrowRightFromLine,
  ListCheck,
  Plus,
  TriangleAlert,
} from "lucide-react";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { TableSkeleton } from "@/components/skeletons/TableSkeleton";
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
  createAPIKeyRoleAssignment,
  deleteAPIKeyRoleAssignment,
  getAPIKey,
  getRole,
  listAPIKeyRoleAssignments,
  listRoles,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { APIKeyRoleAssignment } from "@/gen/tesseral/backend/v1/models_pb";

export function OrganizationApiKeyRolesTab() {
  const { apiKeyId } = useParams();

  const {
    data: listApiKeyRoleAssignmentsResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId: apiKeyId as string,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const apiKeyRoleAssignments =
    listApiKeyRoleAssignmentsResponses?.pages.flatMap(
      (page) => page.apiKeyRoleAssignments || [],
    ) || [];
  return (
    <Card>
      <CardHeader>
        <CardTitle>API Key Roles</CardTitle>
        <CardDescription>
          Manage roles associated with this API key.
        </CardDescription>
        <CardAction>
          <AssignRoleButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton columns={3} />
        ) : (
          <>
            {apiKeyRoleAssignments.length === 0 ? (
              <div className="text-sm text-muted-foreground text-center py-6">
                No roles assigned to this API Key.
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
                  {apiKeyRoleAssignments.map((assignment) => (
                    <ApiKeyRoleAssignmentRow
                      key={assignment.id}
                      apiKeyRoleAssignment={assignment}
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
            Load More
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

function ApiKeyRoleAssignmentRow({
  apiKeyRoleAssignment,
}: {
  apiKeyRoleAssignment: APIKeyRoleAssignment;
}) {
  const { refetch } = useInfiniteQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId: apiKeyRoleAssignment.apiKeyId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getRoleResponse } = useQuery(getRole, {
    id: apiKeyRoleAssignment.roleId,
  });
  const deleteApiKeyRoleAssignmentMutation = useMutation(
    deleteAPIKeyRoleAssignment,
  );

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteApiKeyRoleAssignmentMutation.mutateAsync({
      id: apiKeyRoleAssignment.id,
    });
    await refetch();
    toast.success("API Key Role Assignment deleted successfully.");
    setDeleteOpen(false);
  }

  return (
    <>
      <TableRow>
        <TableCell className="font-medium">
          {getRoleResponse?.role?.displayName || apiKeyRoleAssignment.roleId}
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
            onClick={() => setDeleteOpen(true)}
          >
            Unassign
          </Button>
        </TableCell>
      </TableRow>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              Are you sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will remove the{" "}
              <span className="font-semibold">
                {getRoleResponse?.role?.displayName ||
                  getRoleResponse?.role?.id}
              </span>{" "}
              Role from the API Key.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
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
  apiKeyId: z.string(),
});

function AssignRoleButton() {
  const { apiKeyId, organizationId } = useParams();

  const { data: listAPIKeyRoleAssignmentsResponse, refetch } = useInfiniteQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId: apiKeyId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getApiKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const { data: listProjectRolesResponse } = useQuery(listRoles, {});
  const { data: listOrganizationRolesResponse } = useQuery(
    listRoles,
    {
      organizationId,
    },
    {
      enabled: !!getApiKeyResponse?.apiKey?.organizationId,
    },
  );
  const createApiKeyRoleAssignmentMutation = useMutation(
    createAPIKeyRoleAssignment,
  );

  const currentRoleAssignments =
    listAPIKeyRoleAssignmentsResponse?.pages.flatMap(
      (page) => page.apiKeyRoleAssignments || [],
    ) || [];
  const roles = [
    ...(listProjectRolesResponse?.roles || []),
    ...(listOrganizationRolesResponse?.roles || []),
  ].filter((role) => {
    return !currentRoleAssignments.some(
      (assignment) => assignment.roleId === role.id,
    );
  });

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      roleId: "",
      apiKeyId: apiKeyId,
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
    try {
      await createApiKeyRoleAssignmentMutation.mutateAsync({
        apiKeyRoleAssignment: {
          roleId: data.roleId,
          apiKeyId: data.apiKeyId,
        },
      });
      await refetch();
      form.reset();
      setOpen(false);
      toast.success("Role assigned successfully.");
    } catch (error) {
      toast.error("Failed to assign role.");
    }
  }

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
            Select a role to assign to the{" "}
            <span className="font-semibold">
              {getApiKeyResponse?.apiKey?.displayName ||
                getApiKeyResponse?.apiKey?.id}
            </span>{" "}
            API Key.
          </DialogDescription>
        </DialogHeader>
        {!roles.length ? (
          <div className="text-sm text-muted-foreground">
            No roles available to assign. Please create a role first.
            <Link to="/settings/access">
              <Button variant="outline" className="mt-4">
                <ListCheck />
                Manage Roles
                <ArrowRightFromLine />
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
                        The Role to assign to the{" "}
                        <span className="font-semibold">
                          {getApiKeyResponse?.apiKey?.displayName ||
                            getApiKeyResponse?.apiKey?.id}
                        </span>{" "}
                        API Key.
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
