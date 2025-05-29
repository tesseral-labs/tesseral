import * as React from 'react';

import { cn } from '@/lib/utils';

const ConsoleCard = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('rounded-lg border bg-card text-card-foreground', className)}
    {...props}
  />
));
ConsoleCard.displayName = 'ConsoleCard';

const ConsoleCardHeader = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      'flex flex-row items-center justify-between space-x-4 p-6',
      className,
    )}
    {...props}
  />
));
ConsoleCardHeader.displayName = 'ConsoleCardHeader';

const ConsoleCardDetails = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('flex flex-col space-y-1.5', className)}
    {...props}
  />
));
ConsoleCardDetails.displayName = 'ConsoleCardHeaderDetails';

const ConsoleCardTitle = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('font-semibold leading-none tracking-tight', className)}
    {...props}
  />
));
ConsoleCardTitle.displayName = 'ConsoleCardTitle';

const ConsoleCardDescription = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('text-sm text-muted-foreground', className)}
    {...props}
  />
));
ConsoleCardDescription.displayName = 'ConsoleCardDescription';

const ConsoleCardContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div ref={ref} className={cn('p-6 pt-0', className)} {...props} />
));
ConsoleCardContent.displayName = 'ConsoleCardContent';

const ConsoleCardTableContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      'p-0 pb-6 overflow-x-auto [&_tr>th:first-child]:pl-6 [&_tr>th:last-child]:pr-6 [&_tr>td:first-child]:pl-6 [&_tr>td:last-child]:pr-6',
      className,
    )}
    {...props}
  />
));
ConsoleCardTableContent.displayName = 'ConsoleCardTableContent';

const ConsoleCardFooter = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('flex items-center p-6 pt-0', className)}
    {...props}
  />
));
ConsoleCardFooter.displayName = 'ConsoleCardFooter';

export {
  ConsoleCard,
  ConsoleCardDetails,
  ConsoleCardHeader,
  ConsoleCardFooter,
  ConsoleCardTitle,
  ConsoleCardDescription,
  ConsoleCardContent,
  ConsoleCardTableContent,
};
