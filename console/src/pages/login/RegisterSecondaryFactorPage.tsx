import React from "react";
import { Link } from "react-router-dom";

import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { Button } from "@/components/ui/button";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Title } from "@/components/Title";

export function RegisterSecondaryFactorPage() {
  return (
    <LoginFlowCard>
      <Title title="Set up secondary authentication factor" />
      <CardHeader>
        <CardTitle>Set up secondary authentication factor</CardTitle>
        <CardDescription>
          To continue logging in, you must set up a secondary authentication
          factor.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <Button className="w-full" variant="outline" asChild>
            <Link to="/register-passkey">Set up a passkey</Link>
          </Button>

          <Button className="w-full" variant="outline" asChild>
            <Link to="/register-authenticator-app">
              Set up an authenticator app
            </Link>
          </Button>
        </div>
      </CardContent>
    </LoginFlowCard>
  );
}
