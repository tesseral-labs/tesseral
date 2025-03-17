import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router-dom";
import { z } from "zod";


import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form, FormControl, FormDescription,
  FormItem, FormDescription, FormMessage,
  FormField,
  FormLabel,
  FormItem,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  createUserInvite,
  deleteUserInvite,
  listUserInvites,
  listUsers,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { User, UserInvite } from "@/gen/tesseral/frontend/v1/models_pb";
import {
  AlertDialog,
  AlertDialog,
  AlertDialogHeaderContent,
  AlertDialog,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTrigger,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogFooter,
  AlertDialogCancel,
} from "@/components/ui/alert-dialog";
import { toast } from "sonner";
import { Switch } from "@/components/ui/switch";


export function OrganizationUsersTab() {
  const {
    data: listUsersResponses,
    fetchNextPage,
    hasNextPage,
  } = useInfiniteQuery(
    listUsers,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  return (
    <div className="space-y-4">
      <InvitesCard />

      <Card>
        <CardHeader>
          <CardTitle>Users</CardTitle>
          <CardDescription>List of users in your organization.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {listUsersResponses?.pages
              .flatMap((page) => page.users)
              .map((user) => <UserRow key={user.id} user={user} />)}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function InvitesCard() {
  const {
    data: listUserInvitesResponses,
    fetchNextPage,
    hasNextPage,
  } = useInfiniteQuery(
    listUserInvites,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const userInvites = listUserInvitesResponses?.pages?.flatMap(
    (page) => page.userInvites,
  );

  return (
    <Card>
      <div className="flex items-center justify-between">
        <CardHeader>
          <CardTitle>User Invites</CardTitle>
          <CardDescription>
            List of pending user invites to join your organization.
          </CardDescription>
        </CardHeader>

        <div className="pr-6">
          <CreateUserInviteButton />
        </div>
      </div>

      <CardContent>
        {userInvites && userInvites.length > 0 ? (
          <div className="space-y-4">
            {userInvites.map((userInvite) => (
              <UserInviteRow key={userInvite.id} userInvite={userInvite} />
            ))}
          </div>
        ) : (
          <div className="flex justify-center">
            <div className="flex flex-col items-center gap-y-2">
              <div className="text-sm">No pending user invites.</div>
              <CreateUserInviteButton />
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

const schema = z.object({
  email: z.string().email(),
  owner: z.boolean(),
});

function CreateUserInviteButton() {
  const { refetch: refetchListUserInvites } = useInfiniteQuery(
    listUserInvites,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const [open, setOpen] = useState(false);
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
      owner: false,
    },
  });

  const { mutateAsync: createUserInviteAsync } = useMutation(createUserInvite);

  async function handleSubmit(values: z.infer<typeof schema>) {
    await createUserInviteAsync({
      userInvite: {
        email: values.email,
        owner: values.owner,
      },
    });
    await refetchListUserInvites();
    setOpen(false);
    toast.success(`${values.email} can now join your organization.`);
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Invite user</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Invite User</AlertDialogTitle>
          <AlertDialogDescription>
            When you invite a user, they will be able to join your organization
            after verifying their email address. They do not need a special link
            to join.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}
                className="space-y-8">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input placeholder="alice@example.com" {...field} />
                  </FormControl>
                  <FormDescription>
                    The email address of the user you want to invite.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="owner"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Invite as owner</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether the new user will join as an owner.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Invite User</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}

function UserInviteRow({ userInvite }: { userInvite: UserInvite }) {
  const { refetch: refetchListUserInvites } = useInfiniteQuery(
    listUserInvites,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const { mutateAsync: deleteUserInviteAsync } = useMutation(deleteUserInvite);

  async function handleRevoke() {
    await deleteUserInviteAsync({
      id: userInvite.id,
    });
    await refetchListUserInvites();
    toast.success("User invite revoked");
  }

  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-x-4">
        <Avatar>
          <AvatarFallback>
            {userInvite.email.substring(0, 1).toUpperCase()}
          </AvatarFallback>
        </Avatar>

        <div>
          <div className="text-sm font-medium flex items-center gap-x-2">
            {userInvite.email}
            {userInvite.owner && <Badge variant="outline">Owner</Badge>}
          </div>
          <div className="text-sm text-muted-foreground">
            Invited{" "}
            {DateTime.fromJSDate(
              timestampDate(userInvite.createTime!),
            ).toRelative()}
          </div>
        </div>
      </div>

      <Button onClick={handleRevoke} variant="outline">Revoke</Button>
    </div>
  );
}

function UserRow({ user }: { user: User }) {
  const { data: whoamiResponse } = useQuery(whoami);
  const isYou = whoamiResponse?.user?.id === user.id;

  return (
    <div className="flex items-center gap-x-4">
      <Avatar>
        <AvatarFallback>
          {user.email.substring(0, 1).toUpperCase()}
        </AvatarFallback>
      </Avatar>

      <div>
        <div className="text-sm font-medium flex items-center gap-x-2">
          <Link
            to={`/organization-settings/users/${user.id}`}
            className="underline underline-offset-2 decoration-muted-foreground/50"
          >
            {user.email}
          </Link>
          {isYou && <Badge variant="outline">You</Badge>}
        </div>
        <div className="text-sm">{user.email}</div>
      </div>
    </div>
  );
}
