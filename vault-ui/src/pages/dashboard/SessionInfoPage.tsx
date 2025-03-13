import React from "react";

import { Title } from "@/components/Title";
import { useUser } from "@/lib/auth";

const SessionInfoPage = () => {
  const user = useUser();

  return (
    <>
      <Title title="Session Info" />
      <div>
        <h1 className="text-foreground">Hello, {user?.email}</h1>
      </div>
    </>
  );
};

export default SessionInfoPage;
