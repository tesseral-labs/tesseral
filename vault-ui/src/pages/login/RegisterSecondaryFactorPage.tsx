import React from "react";

import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { useDarkMode } from "@/lib/dark-mode";

export function RegisterSecondaryFactorPage() {
  const darkMode = useDarkMode()
  return (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Set up secondary authentication factor</CardTitle>
        <CardDescription>
          To continue logging in, you must set up a secondary authentication
          factor.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <Button
            className="w-full"
            variant={darkMode ? "default" : "outline"}
            asChild
          >
            <Link to="/register-passkey">Set up a passkey</Link>
          </Button>

          <Button
            className="w-full"
            variant={darkMode ? "default" : "outline"}
            asChild
          >
            <Link to="/register-authenticator-app">Set up an authenticator app</Link>
          </Button>
        </div>
      </CardContent>
    </LoginFlowCard>
  );
}
