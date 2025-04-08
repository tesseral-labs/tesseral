import {
  ArrowRightIcon,
  BookOpen,
  ChevronRightIcon,
  Settings,
} from 'lucide-react';

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
import { PageDescription, PageTitle } from '@/components/page';

export function HomePage() {
  return (
    <div className="">
      <PageTitle>Welcome to Tesseral</PageTitle>
      <PageDescription className="mt-2">
        Tesseral is currently in beta.
      </PageDescription>

      <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <a href="https://tesseral.com/docs/quickstart" className="group">
          <Card className="h-full transition-all hover:shadow-md">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-xl font-medium">
                Documentation
              </CardTitle>
              <BookOpen className="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors" />
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

        <Link to="/project-settings" className="group">
          <Card className="h-full transition-all hover:shadow-md">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-xl font-medium">Settings</CardTitle>
              <Settings className="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors" />
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
      </div>

      <Card className="overflow-hidden">
        <div className="grid grid-cols-1 md:grid-cols-2">
          <div>
            <CardHeader className="pb-4">
              <CardTitle className="text-xl font-medium">
                A Note From Our Founders 👋
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground">
                Welcome to Tesseral! We consider it a privilege to support auth
                for your app.
              </p>
              <p className="mt-4 text-muted-foreground">
                We've designed Tesseral to be powerful{' '}
                <span className="italic">and</span> easy to use. If you
                encounter issues with your implementation, please let us know
                immediately. We want to get this right.
              </p>
              <p className="mt-4 text-muted-foreground">
                This project is still in its early days. We're always grateful
                to receive feedback from the community as we continue to advance
                the project.
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

          <div className="bg-muted flex items-center justify-center p-6">
            <img className="max-h-16" src="/images/tesseral-beta.svg" />
          </div>
        </div>
      </Card>
    </div>
  );
}
