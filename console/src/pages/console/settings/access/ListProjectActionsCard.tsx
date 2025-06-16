import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Edit, Plus, Settings, Trash, TriangleAlert } from "lucide-react";
import React, { MouseEvent, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { ValueCopier } from "@/components/core/ValueCopier";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
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
  getRBACPolicy,
  listRoles,
  updateRBACPolicy,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { Action } from "@/gen/tesseral/backend/v1/models_pb";

export function ListProjectActionsCard() {
  const { data: getRBACPolicyResponse } = useQuery(getRBACPolicy);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Role-based Action Control Policy</CardTitle>
        <CardDescription>
          A Role-Based Access Control Policy is the set of fine-grained Actions
          in a Project.
        </CardDescription>
        <CardAction>
          <CreateProjectActionButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Action</TableHead>
              <TableHead>Description</TableHead>
              <TableHead></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {getRBACPolicyResponse?.rbacPolicy?.actions.map((action) => (
              <TableRow key={action.name}>
                <TableCell>
                  <ValueCopier value={action.name} label="Action" />
                </TableCell>
                <TableCell>
                  {action.description || (
                    <span className="text-muted-foreground">—</span>
                  )}
                </TableCell>
                <TableCell className="text-right">
                  <ManageProjectActionButtion action={action} />
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
  name: z.string().regex(/^[a-z0-9_]+\.[a-z0-9_]+\.[a-z0-9_]+$/i, {
    message:
      "Action name must contain only lowercase letters, numbers, and underscores, and must be of the form 'x.y.z'.",
  }),
  description: z.string(),
});

function CreateProjectActionButton() {
  const { data: getRBACPolicyResponse, refetch } = useQuery(getRBACPolicy);
  const createActionMutation = useMutation(updateRBACPolicy);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: "",
      description: "",
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    // Get a fresh copy of the RBAC policy before updating actions
    await refetch();

    await createActionMutation.mutateAsync({
      rbacPolicy: {
        actions: [
          ...(getRBACPolicyResponse?.rbacPolicy?.actions || []),
          {
            name: data.name,
            description: data.description,
          },
        ],
      },
    });
    await refetch();
    form.reset();
    setOpen(false);
    toast.success("Action created successfully");
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus />
          Create Action
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Project Action</DialogTitle>
          <DialogDescription>
            Define a new action for your project.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Action Name</FormLabel>
                    <FormDescription>
                      The name of the action, in the format <code>x.y.z</code>.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="x.y.z" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormDescription>
                      Optional description of the action.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="Optional" {...field} />
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
                Create Action
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

function ManageProjectActionButtion({ action }: { action: Action }) {
  const { refetch: refetchRoles } = useInfiniteQuery(
    listRoles,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getRBACPolicyResponse, refetch } = useQuery(getRBACPolicy);
  const updateActionMutation = useMutation(updateRBACPolicy);

  const [actionToDelete, setActionToDelete] = useState<string | null>(null);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: getRBACPolicyResponse?.rbacPolicy?.actions[0]?.name || "",
      description:
        getRBACPolicyResponse?.rbacPolicy?.actions[0]?.description || "",
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setEditOpen(false);
    return false;
  }

  async function handleDelete() {
    if (!actionToDelete) {
      return;
    }

    // Get a fresh copy of the RBAC policy before updating actions
    await refetch();

    const updatedActions = (
      getRBACPolicyResponse?.rbacPolicy?.actions || []
    ).filter((action) => action.name !== actionToDelete);

    await updateActionMutation.mutateAsync({
      rbacPolicy: {
        actions: updatedActions,
      },
    });
    await refetch();
    await refetchRoles();
    setDeleteOpen(false);
    toast.success("Action deleted successfully");
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    // Get a fresh copy of the RBAC policy before updating actions
    await refetch();

    console.log("data.name", data.name);

    const updatedActions = (
      getRBACPolicyResponse?.rbacPolicy?.actions || []
    ).map((qAction) => (qAction.name === action.name ? data : qAction));

    await updateActionMutation.mutateAsync({
      rbacPolicy: {
        actions: updatedActions,
      },
    });
    await refetch();
    await refetchRoles();
    form.reset(data);
    setEditOpen(false);
    toast.success("Action updated successfully");
  }

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
            Edit Action
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="group"
            onClick={() => {
              setActionToDelete(action.name);
              setDeleteOpen(true);
            }}
          >
            <Trash className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Delete Action
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      {/* Edit Action Dialog */}
      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Project Action</DialogTitle>
            <DialogDescription>
              Update the details of the action.
            </DialogDescription>
          </DialogHeader>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              <div className="space-y-6">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Action Name</FormLabel>
                      <FormDescription>
                        The name of the action, in the format <code>x.y.z</code>
                        .
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input placeholder="x.y.z" {...field} />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Description</FormLabel>
                      <FormDescription>
                        Optional description of the action.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input placeholder="Optional" {...field} />
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
                  Update Action
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Delete Action Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent className="max-w-sm">
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This cannot be undone. This will permanently delete the{" "}
              <span className="font-semibold">{actionToDelete}</span> action.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="mt-8">
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button
              className="ml-2"
              variant="destructive"
              onClick={handleDelete}
            >
              Delete Action
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
