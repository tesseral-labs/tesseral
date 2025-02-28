import * as React from 'react';

import { cn } from '@/lib/utils';

const AuthPreviewInput = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => {
  return (
    <div
      className={cn(
        'flex h-9 w-full cursor-default rounded-sm border border-input bg-background px-3 py-2 text-sm text-muted-foreground',
        className,
      )}
      ref={ref}
      {...props}
    >
      email@example.com
    </div>
  );
});
AuthPreviewInput.displayName = 'AuthPreviewInput';

export default AuthPreviewInput;
