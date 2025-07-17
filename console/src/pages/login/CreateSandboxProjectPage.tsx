import { useMutation } from "@connectrpc/connect-query";
import React, { useEffect } from "react";
import { useNavigate } from "react-router";

import { Loader } from "@/components/core/Loader";
import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { Title } from "@/components/page/Title";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { onboardingCreateProjects } from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";

export function CreateSandboxProjectPage() {
  const { mutateAsync: onboardingCreateProjectsAsync } = useMutation(
    onboardingCreateProjects,
  );
  const navigate = useNavigate();

  useEffect(() => {
    async function createProject() {
      try {
        await onboardingCreateProjectsAsync({
          displayName: "Sandbox",
          appUrl: "http://localhost:3000",
        });

        navigate("/");
      } catch {
        // Fall back to manual creation if auto-creation fails
        navigate("/create-organization");
      }
    }

    createProject();
  }, [onboardingCreateProjectsAsync, navigate]);

  return (
    <LoginFlowCard>
      <Title title="Setting up your Project" />
      <CardHeader>
        <CardTitle>Setting up your Project</CardTitle>
        <CardDescription>
          We're creating your Project. You'll be automatically redirected to
          your new Project once it's ready.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Loader />
      </CardContent>
    </LoginFlowCard>
  );
}
