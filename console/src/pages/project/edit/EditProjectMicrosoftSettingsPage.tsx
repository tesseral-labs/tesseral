import {
  PageCodeSubtitle,
  PageDescription,
  PageTitle,
} from '@/components/page';
import { Title } from '@/components/Title';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import {
  getProject,
  updateProject,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { parseErrorMessage } from '@/lib/errors';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import React, { FC, FormEvent, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { toast } from 'sonner';

const EditProjectMicrosoftSettingsPage: FC = () => {
  const { data: getProjectResponse, refetch: refetchProject } = useQuery(
    getProject,
    {},
  );
  const updateProjectMutation = useMutation(updateProject);

  const [logInWithMicrosoft, setLogInWithMicrosoft] = useState(
    getProjectResponse?.project?.logInWithMicrosoft,
  );
  const [microsoftOauthClientId, setMicrosoftOauthClientId] = useState('');
  const [microsoftOauthClientSecret, setMicrosoftOauthClientSecret] =
    useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();

    try {
      await updateProjectMutation.mutateAsync({
        project: {
          logInWithMicrosoft: true,
          microsoftOauthClientId,
          microsoftOauthClientSecret,
        },
      });
      const { data: refetchedProjectResponse } = await refetchProject();
      setLogInWithMicrosoft(
        refetchedProjectResponse?.project?.logInWithMicrosoft,
      );
      toast.success('Microsoft settings saved');
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error('Failed to update Microsoft settings', {
        description: message,
      });
    }
  };

  useEffect(() => {
    setLogInWithMicrosoft(getProjectResponse?.project?.logInWithMicrosoft);
  }, [getProjectResponse]);

  return (
    <div>
      <Title title="Edit Project Microsoft Settings" />

      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/project-settings">Project settings</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Log in with Microsoft settings</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <PageTitle>Log in with Microsoft settings</PageTitle>
      <PageCodeSubtitle>{getProjectResponse?.project?.id}</PageCodeSubtitle>
      <PageDescription>
        Edit the Microsoft log in settings for your Project.
      </PageDescription>

      <Card className="mt-8">
        <CardHeader>
          <CardTitle>Log in with Microsoft settings</CardTitle>
          <CardDescription>Log in with Microsoft settings</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div>
              <Label>Log in with Microsoft</Label>
              <p className="text-sm text-muted-foreground">
                Enable or disable log in with Microsoft within your Project.
              </p>
              <Switch
                checked={logInWithMicrosoft}
                className="mt-2"
                onCheckedChange={setLogInWithMicrosoft}
              />
            </div>
            <div className="mt-4 pt-4 border-t">
              <Label>Microsoft OAuth Client ID</Label>
              <p className="text-sm text-muted-foreground">
                The OAuth Client ID for your Microsoft application.
              </p>
              <Input
                className="max-w-xl mt-2"
                onChange={(e) => setMicrosoftOauthClientId(e.target.value)}
                placeholder={
                  getProjectResponse?.project?.microsoftOauthClientId
                }
                value={microsoftOauthClientId}
              />
            </div>
            <div className="mt-4 pt-4 border-t">
              <Label>Microsoft OAuth Client Secret</Label>
              <p className="text-sm text-muted-foreground">
                The OAuth Client Secret for your Microsoft application.
              </p>
              <Input
                className="max-w-xl mt-2"
                onChange={(e) => setMicrosoftOauthClientSecret(e.target.value)}
                placeholder={
                  getProjectResponse?.project?.microsoftOauthClientId
                    ? '<encrypted>'
                    : ''
                }
                value={microsoftOauthClientSecret}
              />
            </div>
            <div className="text-right mt-8">
              <Link to="/project-settings">
                <Button variant="outline" className="mr-4">
                  Cancel
                </Button>
              </Link>
              <Button>Save</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default EditProjectMicrosoftSettingsPage;
