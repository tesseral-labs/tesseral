import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Copy,
  Edit,
  LoaderCircle,
  Plus,
  Settings,
  Trash,
  TriangleAlert,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { MouseEvent, useState } from "react";
import { useForm } from "react-hook-form";
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
import { Switch } from "@/components/ui/switch";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  createPublishableKey,
  deletePublishableKey,
  listPublishableKeys,
  updatePublishableKey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { PublishableKey } from "@/gen/tesseral/backend/v1/models_pb";

export function ListPublishableKeysCard() {
  const {
    data: listPublishableKeysResponse,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listPublishableKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const publishableKeys =
    listPublishableKeysResponse?.pages?.flatMap(
      (page) => page.publishableKeys,
    ) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Publishable Keys</CardTitle>
        <CardDescription>
          List of publishable keys for your project. Publishable keys are used
          to identify your Project and are safe to expose publicly.
        </CardDescription>
        <CardAction>
          <CreatePublishableKeyButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton columns={6} />
        ) : (
          <>
            {publishableKeys.length === 0 ? (
              <div className="text-center text-muted-foreground text-sm py-6">
                No Publishable Keys found
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Key</TableHead>
                    <TableHead>Mode</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead>Updated</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {publishableKeys.map((publishableKey) => (
                    <TableRow key={publishableKey.id}>
                      <TableCell className="font-medium">
                        {publishableKey.displayName}
                      </TableCell>
                      <TableCell>
                        <div
                          className="px-2 py-1 bg-muted text-muted-foreground rounded inline hover:text-foreground cursor-pointer font-mono text-xs"
                          onClick={() => {
                            navigator.clipboard.writeText(publishableKey.id);
                            toast.success(
                              "Publishable Key copied to clipboard",
                            );
                          }}
                        >
                          {publishableKey.id}
                          <Copy className="inline h-4 w-4 ml-2" />
                        </div>
                      </TableCell>
                      <TableCell>
                        {publishableKey.devMode ? (
                          <Badge variant="secondary">Development</Badge>
                        ) : (
                          <Badge>Production</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        {publishableKey.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(publishableKey.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell>
                        {publishableKey.updateTime &&
                          DateTime.fromJSDate(
                            timestampDate(publishableKey.updateTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <ManagePublishableKeyButton
                          publishableKey={publishableKey}
                        />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </>
        )}
      </CardContent>
      {hasNextPage && (
        <CardFooter className="flex justify-center">
          <Button
            disabled={isFetchingNextPage}
            variant="outline"
            size="sm"
            onClick={() => fetchNextPage()}
          >
            Load more
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

function ManagePublishableKeyButton({
  publishableKey,
}: {
  publishableKey: PublishableKey;
}) {
  const { refetch } = useInfiniteQuery(
    listPublishableKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const deletePublishableKeyMutation = useMutation(deletePublishableKey);
  const updatePublishableKeyMutation = useMutation(updatePublishableKey);

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: publishableKey.displayName,
      devMode: publishableKey.devMode,
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    form.reset({
      displayName: publishableKey.displayName,
      devMode: publishableKey.devMode,
    });
    setEditOpen(false);
  }

  async function handleDelete() {
    await deletePublishableKeyMutation.mutateAsync({
      id: publishableKey.id,
    });
    await refetch();
    setDeleteOpen(false);
    toast.success("Publishable Key deleted successfully");
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    await updatePublishableKeyMutation.mutateAsync({
      id: publishableKey.id,
      publishableKey: {
        displayName: data.displayName,
        devMode: data.devMode,
      },
    });
    await refetch();
    setEditOpen(false);
    toast.success("Publishable Key updated successfully");
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
            Edit Publishable Key
          </DropdownMenuItem>
          <DropdownMenuSeparator />

          <DropdownMenuItem
            className="group"
            onClick={() => setDeleteOpen(true)}
          >
            <Trash className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Delete API Key
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Edit Publishable Key Dialog */}
      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Publishable Key</DialogTitle>
            <DialogDescription>
              Edit the details of this Publishable Key.
            </DialogDescription>
          </DialogHeader>

          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              <div className="mt-8 space-y-4">
                <FormField
                  control={form.control}
                  name="displayName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Display Name</FormLabel>
                      <FormDescription>
                        The human-friendly name for this Publishable Key
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input placeholder="My Publishable Key" {...field} />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="devMode"
                  render={({ field }) => (
                    <FormItem className="flex items-center space-x-2">
                      <div className="space-y-2">
                        <FormLabel>Development Mode</FormLabel>
                        <FormDescription>
                          Enable this if you want to use this key in development
                          environments.
                        </FormDescription>
                        <FormMessage />
                      </div>
                      <FormControl>
                        <Switch
                          checked={field.value}
                          onCheckedChange={field.onChange}
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
                <Button
                  disabled={
                    !form.formState.isDirty ||
                    updatePublishableKeyMutation.isPending
                  }
                  type="submit"
                >
                  {updatePublishableKeyMutation.isPending && (
                    <LoaderCircle className="animate-spin" />
                  )}
                  {updatePublishableKeyMutation.isPending
                    ? "Saving changes"
                    : "Save changes"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation AlertDialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are your sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete this Backend API Key. This action
              cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
  devMode: z.boolean(),
});

function CreatePublishableKeyButton() {
  const { refetch } = useInfiniteQuery(
    listPublishableKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
      onSuccess: () => {
        refetch();
      },
    },
  );
  const createPublishableKeyMutation = useMutation(createPublishableKey);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
      devMode: false,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await createPublishableKeyMutation.mutateAsync({
      publishableKey: {
        displayName: data.displayName,
        devMode: data.devMode,
      },
    });
    form.reset(data);
    await refetch();
    setOpen(false);
    toast.success("Publishable key created successfully");
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus />
          Create Publishable Key
        </Button>
      </DialogTrigger>
      <DialogContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <DialogHeader>
              <DialogTitle>Create Publishable Key</DialogTitle>
              <DialogDescription>
                Create a new Publishable Key for your project. Publishable Keys
                are used to identify your Project and are safe to expose
                publicly.
              </DialogDescription>
            </DialogHeader>

            <div className="mt-8 space-y-4">
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormDescription>
                      The human-friendly name for this Publishable Key
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="My Publishable Key" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="devMode"
                render={({ field }) => (
                  <FormItem className="flex items-center space-x-2">
                    <div className="space-y-2">
                      <FormLabel>Development Mode</FormLabel>
                      <FormDescription>
                        Enable this if you want to use this key in development
                        environments.
                      </FormDescription>
                      <FormMessage />
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>

            <DialogFooter className="mt-4 justify-end space-x-2">
              <Button onClick={() => setOpen(false)} variant="outline">
                Cancel
              </Button>
              <Button
                disabled={
                  !form.formState.isDirty ||
                  createPublishableKeyMutation.isPending
                }
                type="submit"
              >
                {createPublishableKeyMutation.isPending && (
                  <LoaderCircle className="animate-spin" />
                )}
                {createPublishableKeyMutation.isPending
                  ? "Creating Publishable Key"
                  : "Create Publishable Key"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
