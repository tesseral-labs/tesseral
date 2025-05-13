import { ArrowRightIcon, BookOpen, Building2, Settings2 } from 'lucide-react';

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import React from 'react';
import { Link } from 'react-router-dom';
import {
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createStripeCheckoutLink,
  getProjectEntitlements,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { Button } from '@/components/ui/button';

export function HomePage() {
  const { data: getProjectEntitlementsResponse } = useQuery(
    getProjectEntitlements,
  );
  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  const handleUpgrade = async () => {
    const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
    window.location.href = url;
  };

  return (
    <>
      {!getProjectEntitlementsResponse?.entitledBackendApiKeys && (
        <div className="p-2 bg-indigo-400 z-10 fixed t-16 w-full">
          <div className="container m-auto flex items-center justify-between">
            <div className="text-sm text-white/90">
              <span className="text-white font-semibold">
                Get more out of Tesseral.
              </span>{' '}
              Upgrade to the Growth tier to unlock more features and support.
            </div>

            <div className="ml-auto">
              <Button
                className="bg-white/90 text-zinc-950 hover:bg-white hover:text-zinc-950 transition-colors"
                onClick={handleUpgrade}
              >
                Upgrade to Growth Tier
              </Button>
            </div>
          </div>
        </div>
      )}
      <PageHeader>
        <PageTitle>Welcome to Tesseral</PageTitle>
        <PageDescription className="mt-2">
          Tesseral is currently in beta.
        </PageDescription>
      </PageHeader>

      <PageContent>
        <Card className="overflow-hidden mt-8">
          <div className="grid grid-cols-1 md:grid-cols-2">
            <div>
              <CardHeader className="pb-4">
                <CardTitle className="text-xl font-medium">
                  A Note From Our Founders 👋
                </CardTitle>
              </CardHeader>
              <CardContent className="text-sm">
                <p className="text-muted-foreground">
                  Welcome to Tesseral! We consider it a privilege to support
                  auth for your app.
                </p>
                <p className="mt-4 text-muted-foreground">
                  We've designed Tesseral to be powerful{' '}
                  <span className="italic">and</span> easy to use. If you
                  encounter issues with your implementation, please let us know
                  immediately. We want to get this right.
                </p>
                <p className="mt-4 text-muted-foreground">
                  This project is still in its early days. We're always grateful
                  to receive feedback from the community as we continue to
                  advance the project.
                </p>
                <p className="mt-4 text-muted-foreground">
                  &mdash;{' '}
                  <Link
                    to="https://www.linkedin.com/in/ucarion/"
                    className="underline underline-offset-2 decoration-muted-foreground/40 hover:text-primary hover:decoration-primary transition-colors"
                  >
                    Ulysse
                  </Link>{' '}
                  and{' '}
                  <Link
                    to="https://www.linkedin.com/in/nedoleary/"
                    className="underline underline-offset-2 decoration-muted-foreground/40 hover:text-primary hover:decoration-primary transition-colors"
                  >
                    Ned
                  </Link>
                </p>
              </CardContent>
            </div>

            <div className="bg-indigo-50 flex items-center justify-center p-6">
              <img className="max-h-16" src="/images/tesseral-beta.svg" />
            </div>
          </div>
        </Card>

        <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <Link to="/project-settings" className="group">
            <Card className="h-full transition-all hover:shadow-md">
              <CardHeader className="space-y-0 pb-2">
                <CardTitle className="flex items-center text-xl font-medium">
                  <Settings2 className="mr-2 inline h-5 w-5 text-primary" />
                  Project Settings
                </CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-sm text-muted-foreground">
                  Manage your Tesseral implementation.
                </CardDescription>
              </CardContent>
              <CardFooter>
                <span className="text-sm inline-flex items-center gap-x-2">
                  Go to settings <ArrowRightIcon className="h-4 w-4" />
                </span>
              </CardFooter>
            </Card>
          </Link>

          <Link to="/organizations" className="group">
            <Card className="h-full transition-all hover:shadow-md">
              <CardHeader className="pb-2">
                <CardTitle className="text-xl font-medium flex items-center">
                  <Building2 className="mr-2 inline h-5 w-5 text-primary" />
                  Organizations
                </CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-sm text-muted-foreground">
                  Manage your organizations and their users.
                </CardDescription>
              </CardContent>
              <CardFooter>
                <span className="text-sm inline-flex items-center gap-x-2">
                  Go to organizations <ArrowRightIcon className="h-4 w-4" />
                </span>
              </CardFooter>
            </Card>
          </Link>

          <a href="https://tesseral.com/docs/quickstart" className="group">
            <Card className="h-full transition-all hover:shadow-md">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-xl font-medium">
                  <BookOpen className="mr-2 inline h-5 w-5 text-primary" />
                  Documentation
                </CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-sm text-muted-foreground">
                  Learn how to use Tesseral effectively.
                </CardDescription>
              </CardContent>
              <CardFooter>
                <span className="text-sm inline-flex items-center gap-x-2">
                  Read the docs <ArrowRightIcon className="h-4 w-4" />
                </span>
              </CardFooter>
            </Card>
          </a>
        </div>
      </PageContent>
    </>
  );
}
