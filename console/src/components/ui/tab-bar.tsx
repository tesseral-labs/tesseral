import clsx from 'clsx';
import React, { FC, PropsWithChildren } from 'react';
import { Link } from 'react-router-dom';

export const TabBar: FC<PropsWithChildren> = ({ children }) => {
  return (
    <div className="bg-white w-full">
      <div className="border-b border-gray-200">
        <nav aria-label="Tabs" className="-mb-px flex">
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
          ? 'border-indigo-500 text-indigo-600'
          : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
        'whitespace-nowrap border-b-2 py-2 px-4 text-sm font-medium',
      )}
    >
      {label}
    </Link>
  );
};
