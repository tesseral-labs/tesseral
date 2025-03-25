import * as React from 'react';
import { clsx } from 'clsx';
import {
  CircleCheckIcon,
  CircleXIcon,
  Clock9Icon,
} from 'lucide-react';

export const StatusIndicator = ({
  variant,
  children,
}: {
  variant: 'error' | 'success' | 'pending';
  children?: React.ReactNode;
}) => {
  return (
    <span
      className={clsx(
        'inline-flex items-center gap-x-1',
        variant === 'error' && 'text-red-500',
        variant === 'success' && 'text-green-700',
        variant === 'pending' && 'text-gray-500',
      )}
    >
      {variant === 'error' && <CircleXIcon className="h-4 w-4" />}
      {variant === 'success' && <CircleCheckIcon className="h-4 w-4" />}
      {variant === 'pending' && <Clock9Icon className="h-4 w-4" />}
      {children}
    </span>
  );
};
