import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { RotateCcwKey, TriangleAlert } from "lucide-react";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { TabContent } from "@/components/page/Tabs";
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
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
import {
  setPassword,
  updateMe,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  displayName: z.string().optional(),
});

export function UserDetailsTab() {
  const { data: whoamiResponse, refetch } = useQuery(whoami);
  const updateMeMutation = useMutation(updateMe);

  const user = whoamiResponse?.user;

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
    },
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    try {
      await updateMeMutation.mutateAsync({
        user: {
          displayName: data.displayName || undefined,
        },
      });
      await refetch();
      form.reset(data);
      toast.success("Account details updated successfully.");
    } catch {
      toast.error("Failed to update account details. Please try again.");
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
    <TabContent>
      <Form {...form}>
        <form className="w-full" onSubmit={form.handleSubmit(handleSubmit)}>
          <Card>
            <CardHeader>
              <CardTitle>Account Details</CardTitle>
              <CardDescription>Update your account details.</CardDescription>
              <CardAction>
                <Button
                  size="sm"
                  type="submit"
                  disabled={
                    !form.formState.isDirty || updateMeMutation.isPending
                  }
                >
                  Save changes
                </Button>
              </CardAction>
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

              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormDescription>
                  Your account email address. If you need to change this, please
                  reach out to a system administrator.
                </FormDescription>
                <FormControl>
                  <Input
                    className="max-w-lg"
                    disabled
                    readOnly
                    value={user?.email || ""}
                    type="email"
                  />
                </FormControl>
              </FormItem>

              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormDescription>
                      Your full name or preferred display name.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        {...field}
                        className="max-w-lg"
                        placeholder="Jane Doe"
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>
        </form>
      </Form>

      <DangerZoneCard />
    </TabContent>
  );
}

const resetSchema = z.object({
  password: z.string().min(8, "Password must be at least 8 characters long"),
});

function DangerZoneCard() {
  const [resetOpen, setResetOpen] = useState(false);

  const setPasswordMutation = useMutation(setPassword);

  const form = useForm<z.infer<typeof resetSchema>>({
    resolver: zodResolver(resetSchema),
    defaultValues: {
      password: "",
    },
  });

  async function handleSubmit(data: z.infer<typeof resetSchema>) {
    try {
      await setPasswordMutation.mutateAsync({
        password: data.password,
      });
      form.reset();
      toast.success("Password changed successfully.");
      setResetOpen(false);
    } catch {
      toast.error("Failed to change password. Please try again.");
      setResetOpen(false);
    }
  }

  return (
    <>
      <Card className="bg-red-50/30 border-red-200">
        <CardHeader>
          <CardTitle className="text-destructive flex items-center gap-2">
            <TriangleAlert />
            Danger Zone
          </CardTitle>
          <CardDescription>
            Actions in this section are irreversible and can lead to data loss.
            Please proceed with caution.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between gap-8 w-full lg:w-auto flex-wrap lg:flex-nowrap">
            <div className="space-y-1">
              <div className="text-sm font-semibold flex items-center gap-2">
                <RotateCcwKey className="w-6 h-6" />
                <span>Change password</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Update the password you use to log in. This cannot be undone.
              </div>
            </div>
            <Button
              className="border-destructive text-destructive hover:bg-destructive hover:text-white"
              variant="outline"
              size="sm"
              onClick={() => setResetOpen(true)}
            >
              Change password
            </Button>
          </div>
        </CardContent>
      </Card>

      <Dialog open={resetOpen} onOpenChange={setResetOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Change password</DialogTitle>
            <DialogDescription>
              Update the password you use to log in. This cannot be undone.
            </DialogDescription>
          </DialogHeader>

          <Form {...form}>
            <form
              className="space-y-6"
              onSubmit={form.handleSubmit(handleSubmit)}
            >
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>New Password</FormLabel>
                    <FormDescription>
                      Enter your new password. It must be at least 8 characters
                      long.
                    </FormDescription>
                    <FormMessage />
                    <FormControl>
                      <Input
                        {...field}
                        type="password"
                        placeholder="••••••••"
                        autoComplete="new-password"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </form>
          </Form>

          <DialogFooter>
            <Button variant="outline" onClick={() => setResetOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={!form.formState.isDirty}>
              Change password
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
