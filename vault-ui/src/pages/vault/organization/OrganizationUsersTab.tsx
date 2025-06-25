import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useMutation } from "@connectrpc/connect-query";
import { AlignLeft, Settings, Trash, UserLock } from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { Link } from "react-router";
import { toast } from "sonner";

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
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  deleteUser,
  listUsers,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { User } from "@/gen/tesseral/frontend/v1/models_pb";

export function OrganizationUsersTab() {
  const {
    data: listUsersResponses,
    hasNextPage,
    fetchNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery(
    listUsers,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
    },
  );

  const users = listUsersResponses?.pages.flatMap((page) => page.users) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Users</CardTitle>
        <CardDescription>Manage your users.</CardDescription>
      </CardHeader>
      <CardContent>
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
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>
                  <Link to={`/organization/users/${user.id}`}>
                    {user.displayName && (
                      <div className="font-medium">{user.displayName}</div>
                    )}
                    <div className="text-muted-foreground">{user.email}</div>
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
                <TableCell className="space-y-2 space-x-2">
                  {user.email && <Badge variant="outline">Email</Badge>}
                  {user.googleUserId && <Badge variant="outline">Google</Badge>}
                  {user.microsoftUserId && (
                    <Badge variant="outline">Microsoft</Badge>
                  )}
                  {user.githubUserId && <Badge variant="outline">GitHub</Badge>}
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
      </CardContent>
      {hasNextPage && (
        <CardFooter>
          <Button
            variant="outline"
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

function ManageUserButton({ user }: { user: User }) {
  const deleteUserMutation = useMutation(deleteUser);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    try {
      await deleteUserMutation.mutateAsync({
        id: user.id,
      });
      toast.success("User deleted successfully.");
      setDeleteOpen(false);
    } catch {
      toast.error("Failed to delete user. Please try again.");
      setDeleteOpen(false);
      return;
    }
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
          <DropdownMenuItem asChild>
            <Link to={`/organization/users/${user.id}`}>
              <AlignLeft />
              Details
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <Link to={`/organization/users/${user.id}/roles`}>
              <UserLock />
              Roles
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="group"
            onClick={() => setDeleteOpen(true)}
          >
            <Trash className="text-destructive group-hover:text-destructive" />
            <span className="text-destructive group-hover:text-destructive">
              Delete User
            </span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <Trash />
              Are you sure?
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete{" "}
              <span className="font-semibold">{user.email}</span>. This action
              cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
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
