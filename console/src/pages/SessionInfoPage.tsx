import React from 'react';
import { Title } from '@/components/Title';
import { useUser } from '@/lib/use-access-token';

const SessionInfoPage = () => {
  const user = useUser();

  return (
    <>
      <Title title="Session Info" />
      <div>
        <h1>
          Hello, {user?.id} {user?.email}
        </h1>
      </div>
    </>
  );
};

export default SessionInfoPage;
