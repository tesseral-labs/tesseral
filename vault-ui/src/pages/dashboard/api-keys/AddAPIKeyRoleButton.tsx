import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { CirclePlus } from "lucide-react";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  createAPIKeyRoleAssignment,
  listAPIKeyRoleAssignments,
  listRoles,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  roleId: z.string(),
});

export function AddAPIKeyRoleButton() {
  const [open, setOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
  });
  const { apiKeyId, organizationId } = useParams();
  const { data: listAPIKeyRoleAssignmentsResponse, refetch } = useQuery(
    listAPIKeyRoleAssignments,
    {
      apiKeyId,
    },
  );
  const { data: listProjectRolesResponse } = useQuery(listRoles, {});
  const { data: listOrganizationRolesResponse } = useQuery(
    listRoles,
    {
      organizationId,
    },
    {
      enabled: !!organizationId,
    },
  );
  const createAPIKeyRoleAssignmentMutation = useMutation(
    createAPIKeyRoleAssignment,
  );

  const roles = Array.from(
    new Map(
      [
        ...(listProjectRolesResponse?.roles || []),
        ...(listOrganizationRolesResponse?.roles || []),
      ].map((role) => [role.id, role]),
    ).values(),
  ).filter((role) => {
    return !listAPIKeyRoleAssignmentsResponse?.apiKeyRoleAssignments.some(
      (assignment) => assignment.roleId === role.id,
    );
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    await createAPIKeyRoleAssignmentMutation.mutateAsync({
      apiKeyRoleAssignment: {
        apiKeyId,
        roleId: data.roleId,
      },
    });
    await refetch();
    form.reset();

    toast.success("Role added successfully");
    setOpen(false);
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">
          <CirclePlus className="h-4 w-4" />
          Add Role
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Add Role</AlertDialogTitle>
        </AlertDialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <div className="mb-4">
              <FormField
                control={form.control}
                name="roleId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Role</FormLabel>
                    <FormControl>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                      >
                        <SelectTrigger>
                          <SelectValue
                            className="max-w-full overflow-x-hidden"
                            placeholder="Select a role"
                          />
                        </SelectTrigger>
                        <SelectContent>
                          {roles?.map((role) => (
                            <SelectItem key={role.id} value={role.id}>
                              {role.displayName} -{" "}
                              <span className="bg-muted text-muted-foreground rounded p-1 font-mono text-xs">
                                {role.id}
                              </span>
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>

            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Save</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
