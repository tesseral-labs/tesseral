import { cn } from '@/lib/utils';
import clsx from 'clsx';
import React, { FC, PropsWithChildren } from 'react';
import { Link } from 'react-router-dom';

export const TabBar: FC<PropsWithChildren> = ({
  className,
  children,
}: React.HTMLAttributes<HTMLDivElement>) => {
  return (
    <div
      className={cn('bg-white w-full h-11 fixed t-16 z-10 border-b', className)}
    >
      <div>
        <nav aria-label="Tabs" className="-mb-px flex h-full">
          {children}
        </nav>
      </div>
    </div>
  );
};

interface TabBarLinkProps {
  active?: boolean;
  label: string;
  url: string;
}

export const TabBarLink: FC<TabBarLinkProps> = ({
  active = false,
  label,
  url,
}) => {
  return (
    <Link
      to={url}
      className={clsx(
        active
          ? 'border-indigo-600 text-indigo-600'
          : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
        'whitespace-nowrap border-t-2 py-3 px-4 text-sm font-medium h-11',
      )}
    >
      {label}
    </Link>
  );
};
