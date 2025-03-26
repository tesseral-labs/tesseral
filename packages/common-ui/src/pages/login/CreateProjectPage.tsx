import React, { FC, useState } from "react";
import { Title } from "../../components/Title";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { useNavigate } from "react-router";
import { useMutation } from "@connectrpc/connect-query";
import { refresh } from "../../gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { Input } from "../../components/ui/input";
import {
  createProject,
  exchangeIntermediateSessionForSession,
  setOrganization,
} from "../../gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { parseErrorMessage } from "../../lib/errors";
import { toast } from "sonner";
import Loader from "../../components/ui/loader";

interface CreateProjectPageProps {
  setAccessToken: (accessToken: string) => void;
  setRefreshToken: (refreshToken: string) => void;
}

export const CreateProjectPage: FC<CreateProjectPageProps> = ({
  setAccessToken,
  setRefreshToken,
}) => {
  const navigate = useNavigate();

  const [displayName, setDisplayName] = useState<string>("");
  const [redirectUri, setRedirectUri] = useState<string>("");
  const [submitting, setSubmitting] = useState<boolean>(false);

  const createProjectMutation = useMutation(createProject);
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession
  );
  const refreshMutation = useMutation(refresh);
  const setOrganizationMutation = useMutation(setOrganization);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      const projectRes = await createProjectMutation.mutateAsync({
        displayName,
        redirectUri,
      });

      await setOrganizationMutation.mutateAsync({
        organizationId: projectRes?.project?.organizationId,
      });

      const { refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({});

      const { accessToken } = await refreshMutation.mutateAsync({});

      setRefreshToken(refreshToken);
      setAccessToken(accessToken);

      setSubmitting(false);
      navigate("/");
    } catch (error) {
      setSubmitting(false);
      const message = parseErrorMessage(error);
      toast.error(message);
    }
  };

  return (
    <>
      <Title title="Create a new Project" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Create a new Project</CardTitle>
        </CardHeader>
        <CardContent>
          <form className="w-full" onSubmit={handleSubmit}>
            <Input
              id="displayName"
              placeholder="Acme, Inc."
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
            />
            <Input
              id="redirectUri"
              placeholder="https://app.company.com/"
              value={redirectUri}
              onChange={(e) => setRedirectUri(e.target.value)}
            />
            <Button
              className="mt-2 w-full"
              disabled={displayName.length < 1 || submitting}
              type="submit"
            >
              {submitting && <Loader />}
              Create Project
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  );
};
