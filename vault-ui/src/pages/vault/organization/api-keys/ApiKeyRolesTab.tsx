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
  createAPIKeyRoleAssignment,
  deleteAPIKeyRoleAssignment,
  getAPIKey,
  getRole,
  listAPIKeyRoleAssignments,
  listRoles,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { APIKeyRoleAssignment } from "@/gen/tesseral/frontend/v1/models_pb";

export function ApiKeyRolesTab() {
  const { apiKeyId } = useParams();

  const { data: getApiKeyResponse } = useQuery(getAPIKey, {
    id: apiKeyId,
  });
  const {
    data: listAPIKeyRoleAssignmentsResponse,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );

  const apiKey = getApiKeyResponse?.apiKey;
  const apiKeyRoleAssignments =
    listAPIKeyRoleAssignmentsResponse?.pages.flatMap(
      (page) => page.apiKeyRoleAssignments,
    ) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Role assignments</CardTitle>
        <CardDescription>
          Roles assigned to{" "}
          <span className="font-semibold">
            {apiKey?.displayName || apiKey?.id}
          </span>
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
            {apiKeyRoleAssignments.length === 0 ? (
              <div className="text-muted-foreground pt-8 text-center text-sm">
                No role assignments found for this apiKey.
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
                  {apiKeyRoleAssignments?.map((roleAssignment) => (
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
  roleAssignment: APIKeyRoleAssignment;
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
  apiKeyId: z.string().min(1, "User ID is required"),
});

function AssignRoleButton() {
  const { apiKeyId } = useParams();
  const { refetch } = useInfiniteQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId,
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
  const createApiKeyRoleAssignmentMutation = useMutation(
    createAPIKeyRoleAssignment,
  );
  const { data: listRolesResponse } = useQuery(listRoles);

  const roles = listRolesResponse?.roles || [];
  const apiKey = getApiKeyResponse?.apiKey;

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      roleId: "",
      apiKeyId: apiKeyId,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await createApiKeyRoleAssignmentMutation.mutateAsync({
        apiKeyRoleAssignment: {
          apiKeyId: data.apiKeyId,
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
            Assign a role to {apiKey?.displayName || apiKey?.id}. Select the
            role from the list below.
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
                    The rolw to assign to {apiKey?.displayName || apiKey?.id}.
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
                  createApiKeyRoleAssignmentMutation.isPending
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
  roleAssignment: APIKeyRoleAssignment;
}) {
  const { apiKeyId } = useParams();

  const { refetch } = useInfiniteQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const deleteApiKeyRoleAssignmentMutation = useMutation(
    deleteAPIKeyRoleAssignment,
  );

  const [open, setOpen] = useState(false);

  async function handleUnassign() {
    try {
      await deleteApiKeyRoleAssignmentMutation.mutateAsync({
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
            This will remove the role assignment for this apiKey.
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
