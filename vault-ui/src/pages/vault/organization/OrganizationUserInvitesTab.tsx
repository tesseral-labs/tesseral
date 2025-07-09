import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { DialogTitle } from "@radix-ui/react-dialog";
import { Plus, TriangleAlert } from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
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
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import {
  createUserInvite,
  deleteUserInvite,
  listUserInvites,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { UserInvite } from "@/gen/tesseral/frontend/v1/models_pb";

export function OrganizationUserInvitesTab() {
  const {
    data: listUserInvitesResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listUserInvites,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );

  const userInvites = listUserInvitesResponses?.pages.flatMap(
    (page) => page.userInvites || [],
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>User Invites</CardTitle>
        <CardDescription>
          Pending invites for users to join your organization.
        </CardDescription>
        <CardAction>
          <InviteUserButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton columns={3} />
        ) : (
          <>
            {userInvites?.length === 0 ? (
              <div className="text-center text-muted-foreground text-sm pt-8">
                No pending user invites.
              </div>
            ) : (
              <Table>
                <TableBody>
                  {userInvites?.map((userInvite) => (
                    <TableRow key={userInvite.id}>
                      <TableCell>
                        <div className="space-x-2">
                          <span className="font-medium">
                            {userInvite.email}
                          </span>
                          {userInvite.owner && (
                            <Badge variant="outline">Owner</Badge>
                          )}
                        </div>
                        <div className="text-muted-foreground text-sm">
                          Invited{" "}
                          {userInvite.createTime &&
                            DateTime.fromJSDate(
                              timestampDate(userInvite.createTime),
                            ).toRelative()}
                        </div>
                      </TableCell>
                      <TableCell className="text-right">
                        <DeleteUserInviteButton userInvite={userInvite} />
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
            Load more
          </Button>
        </CardFooter>
      )}
    </Card>
  );
}

const schema = z.object({
  email: z.string().email({ message: "Invalid email address" }),
  owner: z.boolean().optional(),
  sendEmail: z.boolean().optional(),
});

function DeleteUserInviteButton({ userInvite }: { userInvite: UserInvite }) {
  const { refetch } = useInfiniteQuery(
    listUserInvites,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const deleteUserInviteMutation = useMutation(deleteUserInvite);

  const [open, setOpen] = useState(false);

  async function handleDelete() {
    try {
      await deleteUserInviteMutation.mutateAsync({
        id: userInvite.id,
      });
      await refetch();
      setOpen(false);
      toast.success("User invite deleted successfully");
    } catch {
      toast.error("Failed to delete user invite. Please try again.");
    }
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button
          className="text-destructive border-destructive/10 hover:bg-destructive hover:text-white"
          variant="outline"
          size="sm"
        >
          Delete Invite
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle className="flex items-center gap-2">
            <TriangleAlert />
            Are you sure?
          </AlertDialogTitle>
          <AlertDialogDescription>
            This will permanently delete the invite for{" "}
            <span className="font-semibold">{userInvite.email}</span>.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteUserInviteMutation.isPending}
          >
            Delete Invite
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function InviteUserButton() {
  const { refetch } = useInfiniteQuery(
    listUserInvites,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );
  const createUserInviteMutation = useMutation(createUserInvite);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
      owner: false,
      sendEmail: true,
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await createUserInviteMutation.mutateAsync({
        userInvite: {
          email: data.email,
          owner: data.owner || false,
        },
        sendEmail: data.sendEmail || true,
      });
      await refetch();
      form.reset();
      setOpen(false);
      toast.success("User invite created successfully");
    } catch {
      toast.error("Failed to create user invite. Please try again.");
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus />
          Invite User
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Invite User</DialogTitle>
          <DialogDescription>
            When you invite a user, they will be able to join your organization
            after verifying their email address. They do not need a special link
            to join.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormDescription>
                      The email address of the user you want to invite
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        {...field}
                        type="email"
                        placeholder="email@example.com"
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="owner"
                render={({ field }) => (
                  <FormItem className="flex justify-between items-center space-x-2">
                    <div className="space-y-2">
                      <FormLabel>Make this user an owner</FormLabel>
                      <FormDescription>
                        Whether the new user will join as an owner.
                      </FormDescription>
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
              <FormField
                control={form.control}
                name="sendEmail"
                render={({ field }) => (
                  <FormItem className="flex justify-between items-center space-x-2">
                    <div className="space-y-2">
                      <FormLabel>Send invitation email</FormLabel>
                      <FormDescription>
                        Send the new user an email notifying them of their
                        invite.
                      </FormDescription>
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

            <DialogFooter className="mt-6">
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={
                  !form.formState.isDirty || createUserInviteMutation.isPending
                }
              >
                Invite User
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
