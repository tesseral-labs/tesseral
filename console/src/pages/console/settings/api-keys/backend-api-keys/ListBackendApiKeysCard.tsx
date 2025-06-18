import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlignLeft,
  ArrowRight,
  Crown,
  Info,
  LoaderCircle,
  Plus,
  Settings,
  ShieldBan,
  Trash,
  TriangleAlert,
} from "lucide-react";
import { DateTime } from "luxon";
import React, { MouseEvent, useState } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { SecretCopier } from "@/components/core/SecretCopier";
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
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
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
  createBackendAPIKey,
  createStripeCheckoutLink,
  deleteBackendAPIKey,
  getProjectEntitlements,
  listBackendAPIKeys,
  revokeBackendAPIKey,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { BackendAPIKey } from "@/gen/tesseral/backend/v1/models_pb";

export function ListBackendApiKeysCard() {
  const {
    data: getProjectEntitlementsResponse,
    isLoading: isLoadingEntitlements,
  } = useQuery(getProjectEntitlements);
  const {
    data: listBackendApiKeysResponses,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isLoading,
  } = useInfiniteQuery(
    listBackendAPIKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  const backendApiKeys =
    listBackendApiKeysResponses?.pages.flatMap((page) => page.backendApiKeys) ||
    [];

  async function handleUpgrade() {
    const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
    window.location.href = url;
  }

  return (
    <>
      {!isLoadingEntitlements &&
        !getProjectEntitlementsResponse?.entitledBackendApiKeys && (
          <div className="bg-gradient-to-br from-violet-500 via-purple-500 to-blue-500 border-0 text-white relative overflow-hidden shadow-xl p-8 rounded-lg">
            <div className="absolute inset-0 bg-gradient-to-br from-white/10 to-transparent" />

            <div className="flex flex-wrap w-full gap-8">
              <div className="w-full space-y-4 md:flex-grow">
                <div className="flex items-center space-x-3">
                  <div className="p-2 rounded-full bg-white/20 backdrop-blur-sm">
                    <Crown className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-white">
                      Upgrade to Growth
                    </h3>
                    <p className="text-xs text-white/80">
                      Unlock advanced features
                    </p>
                  </div>
                </div>
                <p className="font-semibold text-sm">
                  Backend API Keys are available on the Growth Tier.
                </p>
                <p className="text-sm text-white/80">
                  When you upgrade, you'll also unlock custom domains, Managed
                  API Keys allow your customers to authenticate to your service
                  without a session, and dedicated email support.
                </p>
              </div>

              <div className="mt-8 md:mt-auto w-full">
                <Button
                  className="bg-white text-purple-600 hover:bg-white/90 font-medium cursor-pointer"
                  onClick={handleUpgrade}
                  size="lg"
                >
                  Upgrade Now
                  <ArrowRight className="h-4 w-4 ml-2" />
                </Button>
              </div>
            </div>
          </div>
        )}
      {getProjectEntitlementsResponse?.entitledBackendApiKeys && (
        <Card>
          <CardHeader>
            <CardTitle>Backend API Keys</CardTitle>
            <CardDescription>
              Backend API keys are how your backend can automate operations in
              Tesseral using the Tesseral Backend API.
            </CardDescription>
            <CardAction>
              <CreateBackendApiKeyButton />
            </CardAction>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <TableSkeleton />
            ) : (
              <>
                {!backendApiKeys.length ? (
                  <div className="text-center text-muted-foreground text-sm py-6">
                    No API keys found. Create a new key to get started.
                  </div>
                ) : (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Name</TableHead>
                        <TableHead>
                          <div className="flex items-center">
                            <span>ID</span>
                            <HoverCard>
                              <HoverCardTrigger asChild>
                                <div className="inline-block ml-2 text-muted-foreground">
                                  <Info className="h-4 w-4" />
                                </div>
                              </HoverCardTrigger>
                              <HoverCardContent className="bg-primary text-white space-y-4 text-sm">
                                <p className="font-semibold">
                                  Not the secret token
                                </p>
                                <p className="text-xs">
                                  The secret token used to authenticate to the
                                  Backend API is only available immediately
                                  after creation, and is never shown again.
                                </p>
                              </HoverCardContent>
                            </HoverCard>
                          </div>
                        </TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Created</TableHead>
                        <TableHead className="text-right">Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {backendApiKeys.map((key) => (
                        <TableRow key={key.id}>
                          <TableCell>
                            <Link
                              className="font-medium"
                              to={`/settings/api-keys/backend-api-keys/${key.id}`}
                            >
                              {key.displayName || "—"}
                            </Link>
                          </TableCell>
                          <TableCell>
                            <ValueCopier
                              value={key.id}
                              label="Backend API Key ID"
                            />
                          </TableCell>
                          <TableCell>
                            {key.revoked ? (
                              <Badge variant="secondary">Revoked</Badge>
                            ) : (
                              <Badge>Active</Badge>
                            )}
                          </TableCell>
                          <TableCell>
                            {key.createTime &&
                              DateTime.fromJSDate(
                                timestampDate(key.createTime),
                              ).toRelative()}
                          </TableCell>
                          <TableCell className="text-right">
                            <ManageBackendApiKeyButton backendApiKey={key} />
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
                disabled={isFetching}
                variant="outline"
                onClick={() => fetchNextPage()}
                size="sm"
              >
                Load More
              </Button>
            </CardFooter>
          )}
        </Card>
      )}
    </>
  );
}

function ManageBackendApiKeyButton({
  backendApiKey,
}: {
  backendApiKey: BackendAPIKey;
}) {
  const { refetch } = useInfiniteQuery(
    listBackendAPIKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const deleteBackendApiKeyMutation = useMutation(deleteBackendAPIKey);
  const revokeBackendApiKeyMutation = useMutation(revokeBackendAPIKey);

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [revokeOpen, setRevokeOpen] = useState(false);

  async function handleDelete() {
    await deleteBackendApiKeyMutation.mutateAsync({
      id: backendApiKey.id,
    });
    await refetch();
    setDeleteOpen(false);
    toast.success("Backend API Key deleted successfully");
  }

  async function handleRevoke() {
    await revokeBackendApiKeyMutation.mutateAsync({
      id: backendApiKey.id,
    });
    await refetch();
    setRevokeOpen(false);
    toast.success("Backend API Key revoked successfully");
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
              to={`/settings/api-keys/backend-api-keys/${backendApiKey.id}`}
            >
              <div className="w-full flex items-center gap-2">
                <AlignLeft />
                Details
              </div>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          {backendApiKey.revoked ? (
            <DropdownMenuItem
              className="group"
              onClick={() => setDeleteOpen(true)}
            >
              <Trash className="text-destructive group-hover:text-destructive" />
              <span className="text-destructive group-hover:text-destructive">
                Delete API Key
              </span>
            </DropdownMenuItem>
          ) : (
            <DropdownMenuItem
              className="group"
              onClick={() => setRevokeOpen(true)}
            >
              <ShieldBan className="text-destructive group-hover:text-destructive" />
              <span className="text-destructive group-hover:text-destructive">
                Revoke API Key
              </span>
            </DropdownMenuItem>
          )}
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Revoke Confirmation Dialog */}
      <AlertDialog open={revokeOpen} onOpenChange={setRevokeOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <TriangleAlert />
              <span>Are your sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will prevent this Backend API Key from being used to
              authenticate to the Backend API. This cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button variant="outline" onClick={() => setRevokeOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleRevoke}>
              Revoke
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

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
});

function CreateBackendApiKeyButton() {
  const { refetch } = useInfiniteQuery(
    listBackendAPIKeys,
    {
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );
  const createBackendApiKeyMutation = useMutation(createBackendAPIKey);

  const [backendApiKey, setBackendApiKey] = useState<BackendAPIKey>();
  const [open, setOpen] = useState(false);
  const [secretOpen, setSecretOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  function handleCancel(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setOpen(false);
    return false;
  }

  function handleDone(e: MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    e.stopPropagation();
    setSecretOpen(false);
    setOpen(false);
    return false;
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    const { backendApiKey } = await createBackendApiKeyMutation.mutateAsync({
      backendApiKey: {
        displayName: data.displayName,
      },
    });
    if (backendApiKey) {
      form.reset();
      await refetch();
      setBackendApiKey(backendApiKey);
      setOpen(false);
      setSecretOpen(true);
      toast.success("Backend API Key created successfully");
    } else {
      toast.error("Failed to create Backend API Key. Please try again.");
    }
  }

  return (
    <>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button>
            <Plus />
            Create Backend API Key
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create Backend API Key</DialogTitle>
            <DialogDescription>
              Backend API keys are how your backend can automate operations in
              Tesseral using the Tesseral Backend API.
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
                        A human-readable name for the backend API key.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input placeholder="My Backend API Key" {...field} />
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
                    createBackendApiKeyMutation.isPending
                  }
                  type="submit"
                >
                  {createBackendApiKeyMutation.isPending && (
                    <LoaderCircle className="animate-spin" />
                  )}
                  {createBackendApiKeyMutation.isPending
                    ? "Creating Backend API Key"
                    : "Create Backend API Key"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
      {backendApiKey && (
        <Dialog open={secretOpen} onOpenChange={setSecretOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Backend API Key Created</DialogTitle>
              <DialogDescription>
                Backend API Key was created successfully.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div className="space-y-2">
                <span className="text-sm font-semibold">
                  Backend API Key Secret Token
                </span>
                <SecretCopier
                  placeholder="tesseral_secret_key_•••••••••••••••••••••••••"
                  secret={backendApiKey.secretToken || ""}
                />
              </div>
              <div className="text-muted-foreground text-sm">
                Store this secret as TESSERAL_API_KEY in your secrets manager.
                You will not be able to see this secret token again later.
              </div>
            </div>
            <DialogFooter className="mt-8">
              <Button variant="outline" onClick={handleDone}>
                Close
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}
