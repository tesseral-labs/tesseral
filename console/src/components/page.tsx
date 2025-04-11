import React from 'react';
import { cn } from '@/lib/utils';
import { Outlet, useNavigate } from 'react-router';
import ConsoleSidebar from './ConsoleSidebar';
import { SidebarInset, SidebarProvider } from './ui/sidebar';
import { AccessTokenProvider, useAccessToken } from '@/lib/AccessTokenProvider';
import ConsoleNavigation from './ConsoleNavigation';

export const PageShell = () => {
  return (
    <AccessTokenProvider>
      <PageShellInner />
    </AccessTokenProvider>
  );
};

function PageShellInner() {
  const accessToken = useAccessToken();
  if (!accessToken) {
    window.location.href = '/login';
    return null;
  }

  return (
    <>
      {/* <SidebarProvider>
        <ConsoleSidebar />
        <SidebarInset> */}
      <main className="bg-gray-100 w-full min-h-screen">
        <ConsoleNavigation />
        <Outlet />
      </main>
      {/* </SidebarInset>
      </SidebarProvider> */}
    </>
  );
}

export const PageTitle = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLHeadingElement>) => (
  <h1 className={cn('mt-4 font-semibold text-3xl ', className)} {...props} />
);
PageTitle.displayName = 'PageTitle';

export const PageCodeSubtitle = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
  <div
    className={cn(
      'mt-2 inline-block rounded py-1 px-2 font-mono text-xs bg-indigo-700 text-gray-100',
      className,
    )}
    {...props}
  />
);
PageCodeSubtitle.displayName = 'PageCodeSubtitle';

export const PageDescription = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
  <div className={cn('mt-4', className)} {...props} />
);
PageDescription.displayName = 'PageDescription';
