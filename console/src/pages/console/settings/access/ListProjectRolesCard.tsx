import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Edit, Plus, Settings, Trash, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { MouseEvent, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { MultiSelect } from "@/components/core/MuliSelect";
import { ValueCopier } from "@/components/core/ValueCopier";
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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  createRole,
  deleteRole,
  getRBACPolicy,
  getRole,
  listRoles,
  updateRole,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function ListProjectRolesCard() {
  const { data: listRolesResponses } = useInfiniteQuery(
    listRoles,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const roles = listRolesResponses?.pages?.flatMap((page) => page.roles) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Roles</CardTitle>
        <CardDescription>
          Roles are a named collection of Actions, and can be assigned to Users
          and API Keys. These are the Roles available to all Organizations in
          this Project.
        </CardDescription>
        <CardAction>
          <CreateRoleButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Role</TableHead>
              <TableHead>Actions</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {roles.map((role) => (
              <TableRow key={role.id}>
                <TableCell>
                  <div className="space-y-2">
                    <div>{role.displayName}</div>
                    <ValueCopier value={role.id} label="Role ID" />
                  </div>
                </TableCell>
                <TableCell>
                  <div className="flex flex-wrap gap-2">
                    {role.actions.map((action) => (
                      <Badge key={action} variant="secondary">
                        {action}
                      </Badge>
                    ))}
                  </div>
                </TableCell>
                <TableCell>
                  {role.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(role.createTime),
                    ).toRelative()}
                </TableCell>
                <TableCell className="text-right">
                  <ManageRoleButton roleId={role.id} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
  actions: z.array(z.string()),
});

function CreateRoleButton() {
  const { refetch } = useInfiniteQuery(
    listRoles,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getRBACPolicyResponse } = useQuery(getRBACPolicy);
  const createRoleMutation = useMutation(createRole);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
      actions: [],
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await createRoleMutation.mutateAsync({
      role: {
        displayName: data.displayName,
        actions: data.actions,
      },
    });
    await refetch();
    form.reset();
    setOpen(false);
    toast.success("Role created successfully");
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus />
          Create Role
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Role</DialogTitle>
          <DialogDescription>
            Create a new role to assign actions to users.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormDescription>
                      The human-friendly name for this role.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="actions"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Actions</FormLabel>
                    <FormDescription>
                      Select the actions that this role can perform.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <MultiSelect
                        className="block w-full"
                        value={field.value}
                        onValueChange={field.onChange}
                        options={(
                          getRBACPolicyResponse?.rbacPolicy?.actions || []
                        ).map((action) => ({
                          value: action.name,
                          label: action.name,
                        }))}
                        placeholder="Select actions"
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>
            <DialogFooter className="mt-8">
              <Button variant="secondary" onClick={handleCancel}>
                Cancel
              </Button>
              <Button disabled={!form.formState.isDirty} type="submit">
                Create Role
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

function ManageRoleButton({ roleId }: { roleId: string }) {
  const { refetch } = useInfiniteQuery(
    listRoles,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getRBACPolicyResponse } = useQuery(getRBACPolicy);
  const { data: getRoleResponse } = useQuery(getRole, {
    id: roleId,
  });
  const deleteRoleMutation = useMutation(deleteRole);
  const updateRoleMutation = useMutation(updateRole);

  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: getRoleResponse?.role?.displayName || "",
      actions: getRoleResponse?.role?.actions || [],
    },
  });

  async function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setEditOpen(false);
    return false;
  }

  async function handleDelete() {
    await deleteRoleMutation.mutateAsync({ id: roleId });
    await refetch();
    setDeleteOpen(false);
    toast.success("Role deleted successfully");
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updateRoleMutation.mutateAsync({
      id: roleId,
      role: {
        displayName: data.displayName,
        actions: data.actions,
      },
    });
    await refetch();
    form.reset();
    setEditOpen(false);
    toast.success("Role updated successfully");
  }

  useEffect(() => {
    if (getRoleResponse) {
      form.reset({
        displayName: getRoleResponse.role?.displayName || "",
        actions: getRoleResponse.role?.actions || [],
      });
    }
  }, [getRoleResponse, form]);

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" size="sm">
            <Settings />
            Manage
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onClick={() => setEditOpen(true)}>
            <Edit />
            Edit Role
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="group"
            onClick={() => setDeleteOpen(true)}
          >
            <Trash className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Delete Role
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      {/** Edit Dialog */}
      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Role</DialogTitle>
            <DialogDescription>
              Edit the details of the{" "}
              <span className="font-semibold">
                {getRoleResponse?.role?.displayName}
              </span>{" "}
              role.
            </DialogDescription>
          </DialogHeader>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              <div className="space-y-6">
                <FormField
                  control={form.control}
                  name="displayName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Display Name</FormLabel>
                      <FormDescription>
                        The human-friendly name for this role.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="actions"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Actions</FormLabel>
                      <FormDescription>
                        Select the actions that this role can perform.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <MultiSelect
                          className="block w-full"
                          defaultValue={field.value}
                          value={field.value}
                          onValueChange={field.onChange}
                          options={(
                            getRBACPolicyResponse?.rbacPolicy?.actions || []
                          ).map((action) => ({
                            value: action.name,
                            label: action.name,
                          }))}
                          placeholder="Select actions"
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>
              <DialogFooter className="mt-8">
                <Button variant="outline" onClick={handleCancel}>
                  Cancel
                </Button>
                <Button disabled={!form.formState.isDirty} type="submit">
                  Update Role
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/** Delete Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center space-x-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a role cannot be undone. This will permanently remove the{" "}
              <span className="font-semibold">
                {getRoleResponse?.role?.displayName}
              </span>{" "}
              role.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete Role
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
