import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlignLeft,
  KeyRound,
  LoaderCircle,
  Logs,
  Plus,
  Settings,
  ShieldPlus,
  Trash,
  UserLock,
  VenetianMask,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { MouseEvent, useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { ValueCopier } from "@/components/core/ValueCopier";
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
  createUser,
  createUserImpersonationToken,
  deleteUser,
  getOrganization,
  getProject,
  listUsers,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { User } from "@/gen/tesseral/backend/v1/models_pb";

export function OrganizationUsersTab() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const {
    data: listUsersResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listUsers,
    {
      organizationId: organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const users = listUsersResponses?.pages?.flatMap((page) => page.users) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Users</CardTitle>
        <CardDescription>
          Manage users for{" "}
          <span className="font-semibold">
            {getOrganizationResponse?.organization?.displayName}
          </span>
          .
        </CardDescription>
        <CardAction>
          <CreateUserButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton />
        ) : (
          <>
            {users.length === 0 ? (
              <div className="text-center text-muted-foreground text-sm py-6">
                No Users found in this Organization.
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>User</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead>Role</TableHead>
                    <TableHead>Auth Methods</TableHead>
                    <TableHead>Last Updated</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {users?.map((user) => (
                    <TableRow key={user.id}>
                      <TableCell>
                        <Link
                          to={`/organizations/${organizationId}/users/${user.id}`}
                        >
                          <div className="flex flex-col items-start gap-2">
                            {user.displayName && (
                              <span className="font-medium text-sm">
                                {user.displayName}
                              </span>
                            )}
                            <span className="text-muted-foreground text-sm">
                              {user.email}
                            </span>
                          </div>
                        </Link>
                      </TableCell>
                      <TableCell>
                        <ValueCopier value={user.id} label="User ID" />
                      </TableCell>
                      <TableCell>
                        {user.owner ? (
                          <Badge>Owner</Badge>
                        ) : (
                          <Badge variant="secondary">Member</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center flex-wrap gap-2">
                          <Badge variant="outline">Email</Badge>
                          {user.googleUserId && (
                            <Badge variant="outline">Google</Badge>
                          )}
                          {user.microsoftUserId && (
                            <Badge variant="outline">Microsoft</Badge>
                          )}
                          {user.githubUserId && (
                            <Badge variant="outline">GitHub</Badge>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        {user.updateTime &&
                          DateTime.fromJSDate(
                            timestampDate(user.updateTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <ManageUserButton user={user} />
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

function ManageUserButton({ user }: { user: User }) {
  const { organizationId } = useParams();
  const { refetch } = useInfiniteQuery(
    listUsers,
    {
      organizationId: organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const createUserImpersonationTokenMutation = useMutation(
    createUserImpersonationToken,
  );
  const deleteUserMutation = useMutation(deleteUser);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteUserMutation.mutateAsync({ id: user.id });
    await refetch();
    toast.success("User deleted successfully");
  }

  async function handleImpersonate() {
    const { userImpersonationToken } =
      await createUserImpersonationTokenMutation.mutateAsync({
        userImpersonationToken: {
          impersonatedId: user.id,
        },
      });

    window.location.href = `https://${getProjectResponse?.project?.vaultDomain}/impersonate?secret-user-impersonation-token=${userImpersonationToken?.secretToken}`;
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
          <DropdownMenuItem>
            <Link
              className="w-full"
              to={`/organizations/${organizationId}/users/${user.id}`}
            >
              <div className="w-full flex items-center">
                <AlignLeft className="inline mr-2" />
                Details
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link
              className="w-full"
              to={`/organizations/${organizationId}/users/${user.id}/sessions`}
            >
              <div className="w-full flex items-center">
                <ShieldPlus className="inline mr-2" />
                Sessions
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link
              className="w-full"
              to={`/organizations/${organizationId}/users/${user.id}/roles`}
            >
              <div className="w-full flex items-center">
                <UserLock className="inline mr-2" />
                Roles
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link
              className="w-full"
              to={`/organizations/${organizationId}/users/${user.id}/passkeys`}
            >
              <div className="w-full flex items-center">
                <KeyRound className="inline mr-2" />
                Passkeys
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link
              className="w-full"
              to={`/organizations/${organizationId}/users/${user.id}/logs`}
            >
              <div className="w-full flex items-center">
                <Logs className="inline mr-2" />
                Audit Logs
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem className="group" onClick={handleImpersonate}>
            <VenetianMask className="text-destructive group-hover:text:destructive" />
            <span className="text-destructive group-hover:text:destructive">
              Impersonate User
            </span>
          </DropdownMenuItem>
          <DropdownMenuItem
            className="group"
            onClick={() => setDeleteOpen(true)}
          >
            <Trash className="text-destructive group-hover:text:destructive" />
            <span className="text-destructive group-hover:text:destructive">
              Delete User
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription className="space-y-4">
              <p>
                You are about to delete{" "}
                <span className="font-semibold">{user.email}</span> from the{" "}
                <span className="font-semibold">
                  {getOrganizationResponse?.organization?.displayName}
                </span>{" "}
                Organization.
              </p>
              <p className="font-semibold">This action cannot be undone.</p>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="space-x-2 justify-end">
            <Button onClick={() => setDeleteOpen(false)} variant="outline">
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Delete User
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

const schema = z.object({
  email: z.string().email("Invalid email address"),
  displayName: z.string().optional(),
  googleUserId: z.string().optional(),
  microsoftUserId: z.string().optional(),
  githubUserId: z.string().optional(),
  owner: z.boolean(),
});

function CreateUserButton() {
  const { organizationId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const { refetch } = useInfiniteQuery(
    listUsers,
    {
      organizationId: organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const createUserMutation = useMutation(createUser);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
      displayName: "",
      googleUserId: "",
      microsoftUserId: "",
      githubUserId: "",
      owner: false,
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    const newUser: Partial<User> = {
      organizationId: organizationId as string,
      email: data.email,
      owner: data.owner,
    };

    if (data.googleUserId && data.googleUserId != "") {
      newUser.googleUserId = data.googleUserId;
    }

    if (data.microsoftUserId && data.microsoftUserId != "") {
      newUser.microsoftUserId = data.microsoftUserId;
    }

    if (data.githubUserId && data.githubUserId != "") {
      newUser.githubUserId = data.githubUserId;
    }

    await createUserMutation.mutateAsync({
      user: newUser as User,
    });
    await refetch();
    form.reset();
    setOpen(false);
    toast.success("User created successfully.");
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus />
          Create User
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create User</DialogTitle>
          <DialogDescription>
            Create a new user in{" "}
            <span className="font-semibold">
              {getOrganizationResponse?.organization?.displayName}
            </span>
            .
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-4">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormDescription>
                      The email address of the user login.
                    </FormDescription>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder="email@example.com"
                        {...field}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormDescription>
                      The display name of the user
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="John Doe" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="googleUserId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Google User ID</FormLabel>
                    <FormDescription>
                      The Google User ID of the user, if applicable
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="Google User ID" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="microsoftUserId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Microsoft User ID</FormLabel>
                    <FormDescription>
                      The Microsoft User ID of the user, if applicable
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="Microsoft User ID" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="githubUserId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>GitHub User ID</FormLabel>
                    <FormDescription>
                      The GitHub User ID of the user, if applicable
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="GitHub User ID" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="owner"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Owner</FormLabel>
                    <FormDescription>
                      Check this box if the user should be an owner of the
                      Organization
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Switch
                        className="block"
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>

            <DialogFooter className="justify-end mt-4">
              <Button variant="outline" onClick={handleCancel}>
                Cancel
              </Button>
              <Button
                disabled={
                  !form.formState.isDirty || createUserMutation.isPending
                }
                type="submit"
              >
                {createUserMutation.isPending && (
                  <LoaderCircle className="animate-spin" />
                )}
                {createUserMutation.isPending ? "Creating User" : "Create User"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
