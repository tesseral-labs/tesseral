import React from 'react';
import { cn } from '@/lib/utils';
import { Outlet } from 'react-router';
import ConsoleSidebar from './ConsoleSidebar';
import { SidebarInset, SidebarProvider } from './ui/sidebar';

export const PageShell = () => {
  return (
    <SidebarProvider>
      <ConsoleSidebar />
      <SidebarInset>
        <main className="bg-body w-full">
          <div className="bg-indigo-600 pb-64 w-full" />
          <div className="-mt-64 mx-auto max-w-7xl sm:px-6 lg:px-8 pt-8">
            <Outlet />
          </div>
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
};

export const PageTitle = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLHeadingElement>) => (
  <h1
    className={cn('mt-4 font-semibold text-3xl text-white', className)}
    {...props}
  />
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
  <div className={cn('mt-4 text-white', className)} {...props} />
);
PageDescription.displayName = 'PageDescription';
