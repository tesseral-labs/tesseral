import React from "react";
import { Link } from "react-router";

import { Title } from "@/components/core/Title";
import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { Button } from "@/components/ui/button";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function VerifySecondaryFactorPage() {
  return (
    <LoginFlowCard>
      <Title title="Verify secondary authentication factor" />
      <CardHeader>
        <CardTitle>Verify secondary authentication factor</CardTitle>
        <CardDescription>
          To continue logging in, you must verify your identity using a
          secondary authentication factor.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <Button className="w-full" variant={"outline"} asChild>
            <Link to="/verify-passkey">Verify using a passkey</Link>
          </Button>

          <Button className="w-full" variant={"outline"} asChild>
            <Link to="/verify-authenticator-app">
              Verify using an authenticator app
            </Link>
          </Button>
        </div>
      </CardContent>
    </LoginFlowCard>
  );
}
