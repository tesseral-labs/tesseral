import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
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
import {
  getOrganization,
  getUser,
  updateUser,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  displayName: z.string().optional(),
});

export function OrganizationUserDetailsTab() {
  const { userId } = useParams();

  const { data: organizationResponse } = useQuery(getOrganization);
  const { data: getUserResponse } = useQuery(getUser, {
    id: userId,
  });
  const { data: whoamiResponse } = useQuery(whoami);
  const updateUserMutation = useMutation(updateUser);

  const organization = organizationResponse?.organization;
  const sessionUser = whoamiResponse?.user;
  const user = getUserResponse?.user;

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: user?.displayName || "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await updateUserMutation.mutateAsync({
        id: userId,
        user: {
          displayName: data.displayName,
        },
      });
      toast.success("User details updated successfully");
      form.reset(data);
    } catch {
      toast.error("Failed to update user details");
    }
  }

  useEffect(() => {
    if (user) {
      form.reset({
        displayName: user.displayName || "",
      });
    }
  }, [user, form]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)}>
        <Card>
          <CardHeader>
            <CardTitle>User Details</CardTitle>
            {sessionUser?.owner && (
              <CardAction>
                <Button
                  size="sm"
                  type="submit"
                  disabled={
                    !form.formState.isDirty || updateUserMutation.isPending
                  }
                >
                  Save changes
                </Button>
              </CardAction>
            )}
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              {user?.profilePictureUrl ? (
                <img
                  src={user.profilePictureUrl}
                  alt="Profile picture"
                  className="w-16 h-16 rounded-full"
                />
              ) : (
                <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center">
                  <span className="text-muted-foreground">
                    {(user?.displayName || user?.email || "user")
                      ?.charAt(0)
                      .toUpperCase()}
                  </span>
                </div>
              )}
            </div>

            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormDescription>
                    The given or preferred name of the user.
                  </FormDescription>
                  <FormMessage />
                  <FormControl>
                    <Input
                      {...field}
                      className="max-w-lg"
                      disabled={!sessionUser?.owner}
                      placeholder="Jane Doe"
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormDescription>
                This user's account email address.
              </FormDescription>
              <FormControl>
                <Input
                  className="max-w-lg"
                  disabled
                  value={user?.email || ""}
                  type="email"
                />
              </FormControl>
            </FormItem>

            {organization?.logInWithGoogle && (
              <FormItem>
                <FormLabel>Google User ID</FormLabel>
                <FormDescription>
                  The Google User ID associated with this user, if applicable.
                </FormDescription>
                <FormControl>
                  <Input
                    className="max-w-lg"
                    disabled
                    value={user?.googleUserId || ""}
                  />
                </FormControl>
              </FormItem>
            )}
            {organization?.logInWithMicrosoft && (
              <FormItem>
                <FormLabel>Microsoft User ID</FormLabel>
                <FormDescription>
                  The Microsoft User ID associated with this user, if
                  applicable.
                </FormDescription>
                <FormControl>
                  <Input
                    className="max-w-lg"
                    disabled
                    value={user?.microsoftUserId || ""}
                  />
                </FormControl>
              </FormItem>
            )}
            {organization?.logInWithGithub && (
              <FormItem>
                <FormLabel>GitHub User ID</FormLabel>
                <FormDescription>
                  The GitHub User ID associated with this user, if applicable.
                </FormDescription>
                <FormControl>
                  <Input
                    className="max-w-lg"
                    disabled
                    value={user?.githubUserId || ""}
                  />
                </FormControl>
              </FormItem>
            )}
          </CardContent>
        </Card>
      </form>
    </Form>
  );
}
