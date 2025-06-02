import React from 'react';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  ConsoleCard,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardDetails,
  ConsoleCardHeader,
  ConsoleCardTitle,
} from '@/components/ui/console-card';
import { getRBACPolicy } from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { useQuery } from '@connectrpc/connect-query';
import { Checkbox } from '@/components/ui/checkbox';

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
        <ConsoleCard>
          <ConsoleCardHeader>
            <ConsoleCardDetails>
              <ConsoleCardTitle>Role Details</ConsoleCardTitle>
              <ConsoleCardDescription>
                Basic information about the Role.
              </ConsoleCardDescription>
            </ConsoleCardDetails>
          </ConsoleCardHeader>
          <ConsoleCardContent className="space-y-4">
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Display Name</FormLabel>
                  <FormDescription>
                    The display name of the Role. This will be displayed to
                    users.
                  </FormDescription>
                  <FormControl>
                    <Input placeholder="Engineering" {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormDescription>
                    Description of the Role. This will be displayed to users.
                  </FormDescription>
                  <FormControl>
                    <Input
                      placeholder="Grants read/write access to databases and logs."
                      {...field}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
          </ConsoleCardContent>
        </ConsoleCard>

        <ConsoleCard>
          <ConsoleCardHeader>
            <ConsoleCardDetails>
              <ConsoleCardTitle>Actions</ConsoleCardTitle>
              <ConsoleCardDescription>
                Actions that this Role can perform.
              </ConsoleCardDescription>
            </ConsoleCardDetails>
          </ConsoleCardHeader>

          <ConsoleCardContent className="space-y-4">
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
          </ConsoleCardContent>
        </ConsoleCard>

        <Button type="submit">Save</Button>
      </form>
    </Form>
  );
}
