import { Link } from 'react-router-dom';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import React, { ReactNode, useEffect, useState } from 'react';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getProject,
  getRBACPolicy,
  updateRBACPolicy,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  ConsoleCard,
  ConsoleCardDescription,
  ConsoleCardHeader,
  ConsoleCardTitle,
  ConsoleCardTableContent,
  ConsoleCardDetails,
} from '@/components/ui/console-card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { toast } from 'sonner';
import { useNavigate } from 'react-router';

export function EditRBACPolicyPage() {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getRBACPolicyResponse } = useQuery(getRBACPolicy, {});

  const [actions, setActions] = useState<z.infer<typeof schema>[]>([]);
  useEffect(() => {
    if (getRBACPolicyResponse?.rbacPolicy) {
      setActions(getRBACPolicyResponse.rbacPolicy.actions);
    }
  }, [getRBACPolicyResponse]);

  function addAction(action: z.infer<typeof schema>) {
    setActions([...actions, action]);
  }

  function updateAction(index: number, action: z.infer<typeof schema>) {
    const updatedActions = [...actions];
    updatedActions[index] = action;
    setActions(updatedActions);
  }

  function deleteAction(index: number) {
    const updatedActions = [...actions];
    updatedActions.splice(index, 1);
    setActions(updatedActions);
  }

  const { mutateAsync: updateRBACPolicyAsync } = useMutation(updateRBACPolicy);
  const navigate = useNavigate();
  async function handleSave() {
    await updateRBACPolicyAsync({
      rbacPolicy: { actions },
    });
    toast.success('RBAC Policy updated');
    navigate('/project-settings/rbac-settings');
  }

  return (
    <>
      <PageHeader>
        <PageTitle>Edit RBAC Policy</PageTitle>
        <PageCodeSubtitle>{getProjectResponse?.project?.id}</PageCodeSubtitle>
        <PageDescription>
          Edit the Role-Based Access Control policy for your Project.
        </PageDescription>
      </PageHeader>
      <PageContent>
        <ConsoleCard className="mt-8">
          <ConsoleCardHeader>
            <ConsoleCardDetails>
              <ConsoleCardTitle>
                Role-Based Access Control Policy
              </ConsoleCardTitle>
              <ConsoleCardDescription>
                A Role-Based Access Control Policy is the set of fine-grained
                Actions in a Project.
              </ConsoleCardDescription>
            </ConsoleCardDetails>

            <div className="shrink-0 space-x-4">
              <Link to="/project-settings/rbac-settings">
                <Button variant="outline">Cancel</Button>
              </Link>
              <Button onClick={handleSave}>Save</Button>
            </div>
          </ConsoleCardHeader>

          <ConsoleCardTableContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Action Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {actions.map((action, index) => (
                  <TableRow key={action.name}>
                    <TableCell className="font-medium font-mono">
                      {action.name}
                    </TableCell>
                    <TableCell>{action.description}</TableCell>
                    <TableCell className="text-right">
                      <Button
                        onClick={() => deleteAction(index)}
                        variant="link"
                      >
                        Delete
                      </Button>

                      <EditActionButton
                        action={action}
                        onSubmit={(action) => updateAction(index, action)}
                      >
                        <Button variant="link">Edit</Button>
                      </EditActionButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            <EditActionButton
              action={{ name: '', description: '' }}
              onSubmit={addAction}
            >
              <Button className="mt-4 mb-6" variant="outline">
                Add Action
              </Button>
            </EditActionButton>
          </ConsoleCardTableContent>
        </ConsoleCard>
      </PageContent>
    </>
  );
}

const schema = z.object({
  name: z.string().regex(/^[a-z0-9_]+\.[a-z0-9_]+\.[a-z0-9_]+$/i, {
    message:
      "Action name must contain only lowercase letters, numbers, and underscores, and must be of the form 'x.y.z'.",
  }),
  description: z.string(),
});

function EditActionButton({
  children,
  action,
  onSubmit,
}: {
  children: ReactNode;
  action: z.infer<typeof schema>;
  onSubmit: (action: z.infer<typeof schema>) => void;
}) {
  const [open, setOpen] = useState(false);
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: action.name,
      description: action.description,
    },
  });

  function handleSubmit(values: z.infer<typeof schema>) {
    onSubmit(values);
    setOpen(false);
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>{children}</AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Edit Action</AlertDialogTitle>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-8"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Action Name</FormLabel>
                  <FormDescription>
                    The name of the Action. Must be of the form "x.y.z".
                  </FormDescription>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="acme.workspaces.create"
                      {...field}
                    />
                  </FormControl>

                  <FormMessage />
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
                    A human-readable description of what the Action lets Users
                    perform in your product.
                  </FormDescription>
                  <FormControl>
                    <Input
                      type="text"
                      placeholder="Create new workspaces"
                      {...field}
                    />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

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
