import * as React from 'react';

import { cn } from '@/lib/utils';

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
ConsoleCardDetails.displayName = 'ConsoleCardDetails';

const ConsoleCardTableContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      'px-6 overflow-x-auto [&_tr>th]:py-4 [&_tr>td]:py-4',
      className,
    )}
    {...props}
  />
));
ConsoleCardTableContent.displayName = 'ConsoleCardTableContent';

export { ConsoleCardDetails, ConsoleCardTableContent };
