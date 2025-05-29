import React, { FC, useState } from 'react';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createUser,
  getOrganization,
  listUsers,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  ConsoleCard,
  ConsoleCardDetails,
  ConsoleCardContent,
  ConsoleCardDescription,
  ConsoleCardHeader,
  ConsoleCardTitle,
  ConsoleCardTableContent,
} from '@/components/ui/console-card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Link } from 'react-router-dom';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { useNavigate, useParams } from 'react-router';
import { Badge } from '@/components/ui/badge';
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import {
  Form,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { AlertDialogCancel } from '@radix-ui/react-alert-dialog';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import { CirclePlus } from 'lucide-react';
import { toast } from 'sonner';
import { User } from '@/gen/tesseral/backend/v1/models_pb';

export const OrganizationUsersTab = () => {
  const { organizationId } = useParams();
  const { data: listUsersResponse } = useQuery(listUsers, {
    organizationId,
  });

  return (
    <ConsoleCard>
      <ConsoleCardHeader>
        <ConsoleCardDetails>
          <ConsoleCardTitle>Users</ConsoleCardTitle>
          <ConsoleCardDescription>
            A user is what people using your product log into.
          </ConsoleCardDescription>
        </ConsoleCardDetails>

        <CreateUserButton />
      </ConsoleCardHeader>
      <ConsoleCardTableContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Email</TableHead>
              <TableHead>ID</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead>Updated At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {listUsersResponse?.users?.map((user) => (
              <TableRow key={user.id}>
                <TableCell>
                  <Link
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                    to={`/organizations/${organizationId}/users/${user.id}`}
                  >
                    {user.email}
                  </Link>

                  {user.owner && (
                    <Badge variant="outline" className="ml-2">
                      Owner
                    </Badge>
                  )}
                </TableCell>
                <TableCell className="font-mono">{user.id}</TableCell>
                <TableCell>
                  {user.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(user.createTime),
                    ).toRelative()}
                </TableCell>
                <TableCell>
                  {user.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(user.updateTime),
                    ).toRelative()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </ConsoleCardTableContent>
    </ConsoleCard>
  );
};

const schema = z.object({
  email: z.string().email(),
  googleUserId: z.string().optional(),
  microsoftUserId: z.string().optional(),
  owner: z.boolean(),
});

const CreateUserButton: FC = () => {
  const navigate = useNavigate();

  const { organizationId } = useParams();
  const [open, setOpen] = useState(false);

  const { data: organizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });

  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      email: '',
      googleUserId: '',
      microsoftUserId: '',
      owner: false,
    },
  });

  const createUserMutation = useMutation(createUser);

  const handleSubmit = async (data: z.infer<typeof schema>) => {
    try {
      const newUser: Partial<User> = {
        organizationId: organizationId as string,
        email: data.email,
        owner: data.owner,
      };

      if (data.googleUserId) {
        newUser.googleUserId = data.googleUserId;
      }

      if (data.microsoftUserId) {
        newUser.microsoftUserId = data.microsoftUserId;
      }

      const createUserResponse = await createUserMutation.mutateAsync({
        user: newUser as User,
      });

      setOpen(false);

      navigate(
        `/organizations/${organizationId}/users/${createUserResponse.user?.id}`,
      );

      toast.success(`User created successfully!`);
    } catch (error) {
      console.error('Error creating user:', error);
    }
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button variant="outline">
          <CirclePlus />
          Create User
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create User</AlertDialogTitle>
          <AlertDialogDescription>
            Create a new User in the{' '}
            <span className="text-semibold">
              {organizationResponse?.organization?.displayName}
            </span>{' '}
            Organization.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form
            className="space-y-4"
            onSubmit={form.handleSubmit(handleSubmit)}
          >
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <Input
                    type="email"
                    placeholder="jane.doe@example.com"
                    {...field}
                  />
                  <FormDescription>
                    The email address of the User being created.
                  </FormDescription>
                </FormItem>
              )}
            />

            {organizationResponse?.organization?.logInWithGoogle && (
              <FormField
                control={form.control}
                name="googleUserId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      Google User ID{' '}
                      <span className="font-normal text-sm">(optional)</span>
                    </FormLabel>
                    <Input placeholder="Google User ID" {...field} />
                    <FormDescription>
                      The Google User ID belonging to the User. This is
                      optional, and will be set on the User automatically on a
                      successful login attempt.
                    </FormDescription>
                  </FormItem>
                )}
              />
            )}

            {organizationResponse?.organization?.logInWithMicrosoft && (
              <FormField
                control={form.control}
                name="microsoftUserId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      Microsoft User ID{' '}
                      <span className="font-normal text-sm">(optional)</span>
                    </FormLabel>
                    <Input placeholder="Microsoft User ID" {...field} />
                    <FormDescription>
                      The Microsoft User ID belonging to the User. This is
                      optional, and will be set on the User automatically on a
                      successful login attempt.
                    </FormDescription>
                  </FormItem>
                )}
              />
            )}

            <FormField
              control={form.control}
              name="owner"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Owner</FormLabel>
                  <Switch
                    className="block"
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                  <FormDescription>
                    Whether the User should be an owner of the organization.
                    This will give them full access to the organization and all
                    its resources.
                  </FormDescription>
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Create</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
};
