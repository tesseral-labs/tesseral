import React, { FC, PropsWithChildren } from 'react';

const CenterLayout: FC<PropsWithChildren> = ({ children }) => {
  return (
    <div className="bg-background w-full flex flex-row items-center justify-center p-8 py-16">
      {children}
    </div>
  );
};

export default CenterLayout;
