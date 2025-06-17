import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertTriangle,
  LoaderCircle,
  Trash,
  TriangleAlert,
  VenetianMask,
} from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

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
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import {
  createUserImpersonationToken,
  deleteUser,
  getProject,
  getUser,
  updateUser,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

const schmema = z.object({
  displayName: z.string().optional(),
  email: z.string().email("Invalid email address"),
  googleUserId: z.string().optional(),
  githubUserId: z.string().optional(),
  microsoftUserId: z.string().optional(),
  owner: z.boolean(),
  profilePictureUrl: z.string().optional(),
});

export function UserDetailsTab() {
  const { userId } = useParams();
  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getUserResponse, refetch } = useQuery(getUser, {
    id: userId,
  });
  const updateUserMutation = useMutation(updateUser);

  const form = useForm<z.infer<typeof schmema>>({
    resolver: zodResolver(schmema),
    defaultValues: {
      displayName: getUserResponse?.user?.displayName || "",
      email: getUserResponse?.user?.email || "",
      googleUserId: getUserResponse?.user?.googleUserId || "",
      githubUserId: getUserResponse?.user?.githubUserId || "",
      microsoftUserId: getUserResponse?.user?.microsoftUserId || "",
      owner: getUserResponse?.user?.owner || false,
      profilePictureUrl: getUserResponse?.user?.profilePictureUrl || "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schmema>) {
    await updateUserMutation.mutateAsync({
      id: userId,
      user: {
        displayName: data.displayName,
        email: data.email,
        googleUserId: data.googleUserId,
        githubUserId: data.githubUserId,
        microsoftUserId: data.microsoftUserId,
        owner: data.owner,
        profilePictureUrl: data.profilePictureUrl,
      },
    });
    refetch();
    form.reset(data);
    toast.success("User details updated successfully");
  }

  useEffect(() => {
    if (getUserResponse?.user) {
      form.reset({
        displayName: getUserResponse.user.displayName || "",
        email: getUserResponse.user.email || "",
        googleUserId: getUserResponse.user.googleUserId || "",
        githubUserId: getUserResponse.user.githubUserId || "",
        microsoftUserId: getUserResponse.user.microsoftUserId || "",
        owner: getUserResponse.user.owner || false,
        profilePictureUrl: getUserResponse.user.profilePictureUrl || "",
      });
    }
  }, [getUserResponse, form]);

  return (
    <div className="space-y-8">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)}>
          <Card>
            <CardHeader>
              <CardTitle>User details</CardTitle>
              <CardDescription>
                General settings for{" "}
                <span className="font-semibold">
                  {getUserResponse?.user?.displayName ||
                    getUserResponse?.user?.email}
                </span>
                .
              </CardDescription>
              <CardAction>
                <Button
                  type="submit"
                  disabled={
                    !form.formState.isDirty || updateUserMutation.isPending
                  }
                >
                  {updateUserMutation.isPending && (
                    <LoaderCircle className="animate-spin" />
                  )}
                  {updateUserMutation.isPending
                    ? "Saving changes"
                    : "Save changes"}
                </Button>
              </CardAction>
            </CardHeader>
            <CardContent className="space-y-8">
              <FormField
                control={form.control}
                name="owner"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between gap-4">
                    <div className="space-y-2">
                      <FormLabel>Owner</FormLabel>
                      <FormDescription>
                        Whether the User is an Owner of their Organization.
                        Optional.
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
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <FormField
                  control={form.control}
                  name="displayName"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Display name</FormLabel>
                      <FormDescription>
                        The User's display name. This is typically their full
                        personal name. Optional.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input
                          className="max-w-xl"
                          placeholder="Jane Doe"
                          {...field}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormDescription>
                        The User's email address. This is used for login and
                        notifications.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input
                          type="email"
                          className="max-w-xl"
                          placeholder="Email"
                          {...field}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="profilePictureUrl"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Profile Picture URL</FormLabel>
                      <FormDescription>
                        The URL of the User's profile picture. Optional.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <Input
                          className="max-w-xl"
                          placeholder="https://..."
                          {...field}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
                {getProjectResponse?.project?.logInWithGoogle && (
                  <FormField
                    control={form.control}
                    name="googleUserId"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Google User ID</FormLabel>
                        <FormDescription>
                          The User's Google-assigned ID. Optional.
                        </FormDescription>
                        <FormMessage />
                        <FormControl>
                          <Input
                            className="max-w-xl"
                            placeholder="Google User ID"
                            {...field}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                )}
                {getProjectResponse?.project?.logInWithMicrosoft && (
                  <FormField
                    control={form.control}
                    name="microsoftUserId"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Microsoft User ID</FormLabel>
                        <FormDescription>
                          The User's Microsoft-assigned ID. Optional.
                        </FormDescription>
                        <FormMessage />
                        <FormControl>
                          <Input
                            className="max-w-xl"
                            placeholder="Microsoft User ID"
                            {...field}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                )}
                {getProjectResponse?.project?.logInWithGithub && (
                  <FormField
                    control={form.control}
                    name="githubUserId"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>GitHub User ID</FormLabel>
                        <FormDescription>
                          The User's GitHub-assigned ID. Optional.
                        </FormDescription>
                        <FormMessage />
                        <FormControl>
                          <Input
                            className="max-w-xl"
                            placeholder="GitHub User ID"
                            {...field}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />
                )}
              </div>
            </CardContent>
          </Card>
        </form>
      </Form>

      <DangerZoneCard />
    </div>
  );
}

function DangerZoneCard() {
  const { organizationId, userId } = useParams();
  const navigate = useNavigate();

  const { data: getProjectResponse } = useQuery(getProject);
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const createUserImpersonationTokenMutation = useMutation(
    createUserImpersonationToken,
  );
  const deleteUserMutation = useMutation(deleteUser);

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [impersonateOpen, setImpersonateOpen] = useState(false);

  async function handleDelete() {
    await deleteUserMutation.mutateAsync({
      id: userId,
    });
    setDeleteOpen(false);
    toast.success("User deleted successfully");
    navigate(`/organizations/${organizationId}/users`);
  }

  async function handleImpersonate() {
    const { userImpersonationToken } =
      await createUserImpersonationTokenMutation.mutateAsync({
        userImpersonationToken: {
          impersonatedId: userId,
        },
      });

    window.location.href = `https://${getProjectResponse?.project?.vaultDomain}/impersonate?secret-user-impersonation-token=${userImpersonationToken?.secretToken}`;
  }

  return (
    <>
      <Card className="bg-red-50/50 border-red-200">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-destructive">
            <TriangleAlert className="w-4 h-4" />
            <span>Danger Zone</span>
          </CardTitle>
          <CardDescription>
            This section contains actions that can have significant
            consequences. Proceed with caution.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex items-center justify-between gap-8">
            <div className="space-y-1">
              <div className="text-sm font-semibold flex items-center gap-2">
                <VenetianMask className="w-4 h-4" />
                <span>Impersonate User</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Create an impersonated session as this user. This will allow you
                to act as this user within your Vault.
              </div>
            </div>
            <Button
              className="border-destructive text-destructive hover:bg-destructive hover:text-white"
              variant="outline"
              size="sm"
              onClick={() => setImpersonateOpen(true)}
            >
              Impersonate User
            </Button>
          </div>

          <Separator />

          <div className="flex items-center justify-between gap-8">
            <div className="space-y-1">
              <div className="text-sm font-semibold flex items-center gap-2">
                <Trash className="w-4 h-4" />
                <span>Delete User</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Completely delete the user and all information associated with
                them. This cannot be undone.
              </div>
            </div>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setDeleteOpen(true)}
            >
              Delete User
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Impersonation Confirmation Dialog */}
      <AlertDialog open={impersonateOpen} onOpenChange={setImpersonateOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <AlertTriangle />
              <span>Are you sure?</span>
            </AlertDialogTitle>
            <AlertDialogDescription>
              This will create an impersonated session for{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.displayName ||
                  getUserResponse?.user?.email}
              </span>{" "}
              in your Vault.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button
              variant="outline"
              onClick={() => setImpersonateOpen(false)}
              className="mr-2"
            >
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleImpersonate}>
              Impersonate User
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="text-destructive">
              Delete User
            </AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete{" "}
              <span className="font-semibold">
                {getUserResponse?.user?.displayName ||
                  getUserResponse?.user?.email}
              </span>
              ? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteOpen(false)}
              className="mr-2"
            >
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
