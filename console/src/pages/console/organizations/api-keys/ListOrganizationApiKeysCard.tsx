import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { format } from "date-fns";
import {
  CalendarIcon,
  ExternalLink,
  GlobeLock,
  LoaderCircle,
  Logs,
  Plus,
  Settings,
  ShieldBan,
  Trash,
  TriangleAlert,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { SecretCopier } from "@/components/core/SecretCopier";
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
import { Calendar } from "@/components/ui/calendar";
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
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Select,
  SelectContent,
  SelectItem,
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
  createAPIKey,
  deleteAPIKey,
  getOrganization,
  getProject,
  listAPIKeys,
  revokeAPIKey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { APIKey } from "@/gen/tesseral/backend/v1/models_pb";
import { cn } from "@/lib/utils";

export function ListOrganizationApiKeysCard() {
  const { organizationId } = useParams();

  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getProjectResponse } = useQuery(getProject);

  const {
    data: listApiKeysResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listAPIKeys,
    {
      organizationId: organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const apiKeys =
    listApiKeysResponses?.pages?.flatMap((page) => page.apiKeys) || [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Managed API Keys</CardTitle>
        <CardDescription>
          Managed API keys for{" "}
          <span className="font-semibold">
            {getOrganizationResponse?.organization?.displayName}
          </span>{" "}
          to authenticate to your service.
        </CardDescription>
        <CardAction>
          <CreateApiKeyButton />
        </CardAction>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <TableSkeleton columns={4} />
        ) : (
          <>
            {apiKeys.length === 0 ? (
              <div className="text-center text-muted-foreground text-sm py-6">
                No Managed API Keys found. Create a new API Key to get started.
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>API Key</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Created</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {apiKeys?.map((apiKey) => (
                    <TableRow key={apiKey.id}>
                      <TableCell>
                        <Link
                          to={`/organizations/${organizationId}/api-keys/${apiKey.id}`}
                        >
                          <div className="space-y-2">
                            <h3 className="font-semibold">
                              {apiKey.displayName}
                            </h3>
                            <div className="inline bg-muted text-muted-foreground font-mono text-xs py-1 px-2 rounded-sm">
                              {
                                getProjectResponse?.project
                                  ?.apiKeySecretTokenPrefix
                              }
                              ...
                              {apiKey.secretTokenSuffix}
                            </div>
                          </div>
                        </Link>
                      </TableCell>
                      <TableCell>
                        {apiKey.revoked ? (
                          <Badge>Active</Badge>
                        ) : (
                          <Badge variant="secondary">Revoked</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        {apiKey.createTime &&
                          DateTime.fromJSDate(
                            timestampDate(apiKey.createTime),
                          ).toRelative()}
                      </TableCell>
                      <TableCell className="text-right">
                        <ManageApiKeyButton apiKey={apiKey} />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </>
        )}
      </CardContent>
      <CardFooter className="flex justify-center">
        {hasNextPage && (
          <Button
            variant="outline"
            size="sm"
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
          >
            {isFetchingNextPage ? "Loading more..." : "Load More"}
          </Button>
        )}
      </CardFooter>
    </Card>
  );
}

const schema = z.object({
  displayName: z.string().min(1, "Display name is required"),
  expireTime: z.string().optional(),
});

function CreateApiKeyButton() {
  const [createOpen, setCreateOpen] = useState(false);
  const [secretOpen, setSecretOpen] = useState(false);
  const [apiKey, setApiKey] = useState<APIKey>();
  const navigate = useNavigate();

  const [customDate, setCustomDate] = useState<Date>();

  const { organizationId } = useParams();
  const { data: getProjectResponse } = useQuery(getProject);
  const { refetch } = useInfiniteQuery(
    listAPIKeys,
    {
      organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const createApiKeyMutation = useMutation(createAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
      expireTime: "1 day",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const createParams: Record<string, any> = {
      organizationId: organizationId!,
      displayName: data.displayName,
    };

    switch (data.expireTime) {
      case "1 day":
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 24 * 60 * 60 * 1000),
        );
        break;
      case "7 days":
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
        );
        break;
      case "30 days":
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
        );
        break;
      case "custom":
        if (customDate) {
          createParams.expireTime = timestampFromDate(customDate);
        }
        break;
      case "noexpire":
        break;
    }

    const { apiKey } = await createApiKeyMutation.mutateAsync({
      apiKey: createParams,
    });

    if (apiKey) {
      setApiKey(apiKey);
      setCreateOpen(false);
      setSecretOpen(true);

      toast.success("API Key created successfully");

      await refetch();
      navigate(`/organizations/${organizationId}/api-keys/${apiKey.id}`);
    }
  }

  return (
    <>
      <Dialog
        open={!!apiKey?.secretToken && secretOpen}
        onOpenChange={setSecretOpen}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Managed API Key Created</DialogTitle>
            <DialogDescription>
              Managed API Key was created successfully.
            </DialogDescription>
          </DialogHeader>

          <div className="text-sm font-medium leading-none">
            API Key Secret Token
          </div>

          {apiKey?.secretToken && (
            <SecretCopier
              placeholder={`${getProjectResponse?.project?.apiKeySecretTokenPrefix}•••••••••••••••••••••••••••••••••••••••••••••••••••••••`}
              secret={apiKey.secretToken}
            />
          )}

          <div className="text-sm text-muted-foreground">
            Store this secret in your secrets manager. You will not be able to
            see this secret token again later.
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setSecretOpen(false)}>
              Close
            </Button>
            {!!apiKey?.id && (
              <Link
                to={`/organizations/${organizationId}/api-keys/${apiKey.id}`}
              >
                <Button>View API Key</Button>
              </Link>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogTrigger asChild>
          <Button>
            <Plus className="h-4 w-4" />
            Create Managed API Key
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="font-semibold">Create API Key</DialogTitle>
            <DialogDescription>
              Create a new Managed API Key for this organization. This key can
              be used to authenticate to your service.
            </DialogDescription>
          </DialogHeader>

          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(handleSubmit)}
              className="space-y-4"
            >
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormDescription>
                      A human-friendly name for the API Key.
                    </FormDescription>
                    <FormControl>
                      <Input placeholder="Display name" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="expireTime"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Expire time</FormLabel>
                    <FormDescription>
                      The expiration time for the API Key. After this time, the
                      API Key will no longer be valid.
                    </FormDescription>
                    <FormControl>
                      <div className="flex flex-row gap-2">
                        <Select
                          {...field}
                          onValueChange={(value) => {
                            field.onChange(value);
                          }}
                        >
                          <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="Pick a custom date" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="1 day">1 day</SelectItem>
                            <SelectItem value="7 days">7 days</SelectItem>
                            <SelectItem value="30 days">30 days</SelectItem>
                            <SelectItem value="custom">Custom</SelectItem>
                            <SelectItem value="noexpire">
                              No expiration
                            </SelectItem>
                          </SelectContent>
                        </Select>

                        {field.value === "custom" && (
                          <Popover>
                            <PopoverTrigger asChild>
                              <Button
                                variant={"outline"}
                                className={cn(
                                  "w-[270px] justify-start text-left font-normal",
                                  !customDate && "text-muted-foreground",
                                )}
                              >
                                <CalendarIcon className="mr-2 h-4 w-4" />
                                {customDate ? (
                                  format(customDate, "PPP")
                                ) : (
                                  <span>Pick a date</span>
                                )}
                              </Button>
                            </PopoverTrigger>
                            <PopoverContent className="w-auto p-0">
                              <Calendar
                                mode="single"
                                selected={customDate}
                                onSelect={setCustomDate}
                              />
                            </PopoverContent>
                          </Popover>
                        )}
                      </div>
                    </FormControl>
                  </FormItem>
                )}
              />

              <DialogFooter>
                <Button variant="outline" onClick={() => setCreateOpen(false)}>
                  Cancel
                </Button>
                <Button
                  disabled={
                    !form.formState.isDirty || createApiKeyMutation.isPending
                  }
                  type="submit"
                >
                  {createApiKeyMutation.isPending && (
                    <LoaderCircle className="animate-spin" />
                  )}
                  {createApiKeyMutation.isPending
                    ? "Creating API Key"
                    : "Create API Key"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    </>
  );
}

function ManageApiKeyButton({ apiKey }: { apiKey: APIKey }) {
  const { organizationId } = useParams();
  const { refetch } = useInfiniteQuery(
    listAPIKeys,
    {
      organizationId: organizationId,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const deleteApiKeyMutation = useMutation(deleteAPIKey);
  const revokeApiKeyMutation = useMutation(revokeAPIKey);

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [revokeOpen, setRevokeOpen] = useState(false);

  async function handleDelete() {
    await deleteApiKeyMutation.mutateAsync({
      id: apiKey.id,
    });
    await refetch();
    setDeleteOpen(false);
    toast.success("API Key deleted successfully");
  }

  async function handleRevoke() {
    await revokeApiKeyMutation.mutateAsync({
      id: apiKey.id,
    });
    await refetch();
    setRevokeOpen(false);
    toast.success("API Key revoked successfully");
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
              className="flex items-center gap-2"
              to={`/organizations/${organizationId}/api-keys/${apiKey.id}`}
            >
              <ExternalLink />
              <span>Details</span>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link
              className="flex items-center gap-2"
              to={`/organizations/${organizationId}/api-keys/${apiKey.id}/roles`}
            >
              <GlobeLock />
              <span>Roles</span>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link
              className="flex items-center gap-2"
              to={`/organizations/${organizationId}/api-keys/${apiKey.id}/logs`}
            >
              <Logs />
              <span>Audit Logs</span>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className="group hover:bg-destructive/10"
            onClick={() =>
              apiKey.revoked ? setDeleteOpen(true) : setRevokeOpen(true)
            }
          >
            {apiKey.revoked ? (
              <div className="w-full flex items-center gap-2 text-sm text-destructive group-hover:text-destructive">
                <Trash className="text-destructive" />
                Delete API Key
              </div>
            ) : (
              <div className="w-full flex items-center gap-2 text-sm text-destructive group-hover:text-destructive">
                <ShieldBan className="text-destructive" />
                Revoke API Key
              </div>
            )}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Revoke Confirmation Dialog */}
      <AlertDialog open={revokeOpen} onOpenChange={setRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This can not be undone. This will revoke the{" "}
              <span>{apiKey.displayName || apiKey.id}</span> API Key.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setRevokeOpen(false)}>
              Cancel
            </Button>
            <Button
              className="bg-destructive text-white hover:bg-destructive/90"
              onClick={handleRevoke}
            >
              Revoke
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This can not be undone. This will permanently delete the{" "}
              <span>{apiKey.displayName || apiKey.id}</span> API Key.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setDeleteOpen(false)}>
              Cancel
            </Button>
            <Button
              className="bg-destructive text-white hover:bg-destructive/90"
              onClick={handleDelete}
            >
              Delete
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
