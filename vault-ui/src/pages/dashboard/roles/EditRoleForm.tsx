import { useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { getRBACPolicy } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  displayName: z.string(),
  description: z.string(),
  actions: z.array(z.string()),
});

export function EditRoleForm({
  role,
  onSubmit,
}: {
  role: z.infer<typeof schema>;
  onSubmit: (role: z.infer<typeof schema>) => void;
}) {
  const { data: getRBACPolicyResponse } = useQuery(getRBACPolicy);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: role,
  });

  return (
    <Form {...form}>
      <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
        <Card>
          <CardHeader>
            <CardTitle>Role Details</CardTitle>
            <CardDescription>Basic information about the Role.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Engineering" {...field} />
                  </FormControl>
                  <FormDescription>
                    The display name of the Role. This will be displayed to
                    users.
                  </FormDescription>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Grants read/write access to databases and logs."
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Description of the Role. This will be displayed to users.
                  </FormDescription>
                </FormItem>
              )}
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Actions</CardTitle>
            <CardDescription>
              Actions that this Role can perform.
            </CardDescription>
          </CardHeader>

          <CardContent className="space-y-4">
            {getRBACPolicyResponse?.rbacPolicy?.actions?.map((action) => (
              <FormField
                key={action.name}
                control={form.control}
                name="actions"
                render={({ field }) => {
                  return (
                    <FormItem
                      key={action.name}
                      className="flex flex-row items-start space-x-3 space-y-0"
                    >
                      <FormControl>
                        <Checkbox
                          checked={field.value?.includes(action.name)}
                          onCheckedChange={(checked) => {
                            return checked
                              ? field.onChange([...field.value, action.name])
                              : field.onChange(
                                  field.value?.filter(
                                    (value) => value !== action.name,
                                  ),
                                );
                          }}
                        />
                      </FormControl>
                      <div className="grid gap-1.5 leading-none">
                        <FormLabel className="font-normal">
                          {action.name}
                        </FormLabel>
                        <FormDescription>{action.description}</FormDescription>
                      </div>
                    </FormItem>
                  );
                }}
              />
            ))}
          </CardContent>
        </Card>

        <Button type="submit">Save</Button>
      </form>
    </Form>
  );
}
