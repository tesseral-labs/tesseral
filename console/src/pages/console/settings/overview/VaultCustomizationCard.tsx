import { Settings2 } from "lucide-react";
import React from "react";
import { Link } from "react-router";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function VaultCustomizationCard() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Settings2 />
          Vault Customization
        </CardTitle>
        <CardDescription>
          Customize the appearance of your Vault, including colors, logos, and
          more.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-grow"></CardContent>
      <CardFooter className="mt-4">
        <Link className="w-full" to="/settings/vault">
          <Button className="w-full" variant="outline">
            Customize Vault
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}
