import React from "react";
import { Helmet } from "react-helmet";

export const Title = ({ title }: { title?: string }) => {
  return (
    <>
      <Helmet>
        {/* TODO: Make this conditionally load an organization's configured Display Name */}
        {title ? <title>{title} | Tesseral</title> : <title>Tesseral</title>}
      </Helmet>
    </>
  );
};
