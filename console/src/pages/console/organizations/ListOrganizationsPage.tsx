import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlignLeft,
  Key,
  LoaderCircle,
  Logs,
  Plus,
  Settings,
  Shield,
  Trash,
  Users,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { ValueCopier } from "@/components/core/ValueCopier";
import { PageContent } from "@/components/page";
import { Title } from "@/components/page/Title";
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  createOrganization,
  deleteOrganization,
  getOrganization,
  listOrganizations,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { Organization } from "@/gen/tesseral/backend/v1/models_pb";

export function ListOrganizationsPage() {
  const {
    data: listOrganizationsResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listOrganizations,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const organizations = listOrganizationsResponses?.pages?.flatMap(
    (page) => page.organizations,
  );

  return (
    <PageContent>
      <Title title="Organizations" />

      <div className="flex justify-between items-center gap-x-8">
        <div>
          <h1 className="font-semibold text-xl">Organizations</h1>
          <p className="text-muted-foreground text-sm">
            Manage organizations and their authentication settings.
          </p>
        </div>
        <CreateOrganizationButton />
      </div>
      <div>
        <Card>
          <CardHeader>
            <CardTitle>All Organizations</CardTitle>
            <CardDescription>
              Showing {organizations?.length || 0} Organizations
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <TableSkeleton />
            ) : (
              <>
                {organizations?.length === 0 ? (
                  <div className="text-center text-muted-foreground py-6">
                    No Organizations Found
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Organization</TableHead>
                        <TableHead>Auth Methods</TableHead>
                        <TableHead>MFA</TableHead>
                        <TableHead>Created</TableHead>
                        <TableHead className="text-right">Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {organizations?.map((org) => (
                        <TableRow key={org.id}>
                          <TableCell className="font-medium">
                            <div className="flex flex-col items-start gap-2">
                              <Link to={`/organizations/${org.id}`}>
                                {org.displayName}
                              </Link>
                              <ValueCopier
                                value={org.id}
                                label="Organization ID"
                              />
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center flex-wrap gap-2">
                              {org.logInWithGoogle && (
                                <Badge variant="outline">Google</Badge>
                              )}
                              {org.logInWithMicrosoft && (
                                <Badge variant="outline">Microsoft</Badge>
                              )}
                              {org.logInWithGithub && (
                                <Badge variant="outline">GitHub</Badge>
                              )}
                              {org.logInWithEmail && (
                                <Badge variant="outline">Email</Badge>
                              )}
                              {org.logInWithSaml && (
                                <Badge variant="outline">SAML</Badge>
                              )}
                            </div>
                            {/* {org.authenticationMethods.join(", ")} */}
                          </TableCell>
                          <TableCell>
                            {org.requireMfa ? (
                              <Badge className="bg-green-500">Required</Badge>
                            ) : (
                              <Badge variant="secondary">Not Required</Badge>
                            )}
                          </TableCell>
                          <TableCell>
                            {org.createTime &&
                              DateTime.fromJSDate(
                                timestampDate(org.createTime),
                              ).toRelative()}
                          </TableCell>
                          <TableCell className="text-right">
                            <ManageOrganizationButton organization={org} />
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
                disabled={isFetchingNextPage}
                onClick={() => fetchNextPage()}
              >
                Load More
              </Button>
            </CardFooter>
          )}
        </Card>
      </div>
    </PageContent>
  );
}

function ManageOrganizationButton({
  organization,
}: {
  organization: Organization;
}) {
  const { refetch } = useInfiniteQuery(
    listOrganizations,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organization.id,
  });
  const deleteOrganizationMutation = useMutation(deleteOrganization);

  const [deleteOpen, setDeleteOpen] = useState(false);

  async function handleDelete() {
    await deleteOrganizationMutation.mutateAsync({ id: organization.id });
    setDeleteOpen(false);
    await refetch();
    toast.success("Organization deleted successfully");
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm">
          <Settings />
          Manage
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem>
          <Link to={`/organizations/${organization.id}`}>
            <div className="w-full flex items-center">
              <AlignLeft className="inline mr-2" />
              Details
            </div>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <Link to={`/organizations/${organization.id}/authentication`}>
            <div className="w-full flex items-center">
              <Shield className="inline mr-2" />
              Authentication Settings
            </div>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <Link to={`/organizations/${organization.id}/api-keys`}>
            <div className="w-full flex items-center">
              <Key className="inline mr-2" />
              Managed API Keys
            </div>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <Link to={`/organizations/${organization.id}/users`}>
            <div className="w-full flex items-center">
              <Users className="inline mr-2" />
              Users
            </div>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <Link to={`/organizations/${organization.id}/logs`}>
            <div className="w-full flex items-center">
              <Logs className="inline mr-2" />
              Audit Logs
            </div>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem>
          <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
            <AlertDialogTrigger asChild>
              <Button
                className="group text-destructive hover:text-destructive"
                variant="ghost"
                size="sm"
              >
                <Trash className="text-destructive group-hover:text:destructive" />
                Delete Organization
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                <AlertDialogDescription>
                  <p>
                    You are about to delete the{" "}
                    <span className="font-semibold">
                      {getOrganizationResponse?.organization?.displayName}
                    </span>{" "}
                    Organization.
                  </p>
                  <p className="font-semibold">This action cannot be undone.</p>
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter className="space-x-2 justify-end">
                <Button
                  onClick={() => setDeleteOpen(false)}
                  variant="secondary"
                >
                  Cancel
                </Button>
                <Button variant="destructive" onClick={handleDelete}>
                  Delete Organization
                </Button>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
});

function CreateOrganizationButton() {
  const { refetch } = useInfiniteQuery(
    listOrganizations,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const createOrganizationMutation = useMutation(createOrganization);

  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await createOrganizationMutation.mutateAsync({
      organization: {
        displayName: data.displayName,
      },
    });
    form.reset();
    await refetch();
    setOpen(false);
    toast.success("Organization created successfully");
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus />
          Add Organization
        </Button>
      </DialogTrigger>
      <DialogContent className="space-y-4">
        <DialogHeader>
          <DialogTitle>Create Organization</DialogTitle>
          <DialogDescription>Create a new Organization.</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="space-y-4">
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormDescription>
                      The display name of the Organization. This will be
                      displayed to users during the login process.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input placeholder="ACME Corp" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>
            <DialogFooter className="mt-8 justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setOpen(false);
                  form.reset();
                }}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={
                  !form.formState.isDirty ||
                  createOrganizationMutation.isPending
                }
              >
                {createOrganizationMutation.isPending && (
                  <LoaderCircle className="animate-spin" />
                )}
                {createOrganizationMutation.isPending
                  ? "Creating Organization"
                  : "Create Organization"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
