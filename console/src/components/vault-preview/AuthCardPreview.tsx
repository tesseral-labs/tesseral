import { useQuery } from "@connectrpc/connect-query";
import React from "react";

import { getProject } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { cn } from "@/lib/utils";

import { AuthPreviewButton } from "./AuthPreviewButton";
import { AuthPreviewInput } from "./AuthPreviewInput";
import { AuthPreviewTextDivider } from "./AuthPreviewTextDivider";

interface AuthCardProps {
  logo?: string;
  noBorder?: boolean;
}

export function AuthCardPreview({ logo, noBorder = false }: AuthCardProps) {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <div>
      {logo && (
        <div className="h-6 w-full flex flex-col justify-center">
          <img className="h-full w-full object-contain" src={logo} />
        </div>
      )}

      <div
        className={cn("rounded-sm p-4 bg-card mt-4", noBorder ? "" : "border")}
      >
        <div className="cursor-default w-full text-center font-semibold text-base mb-4 text-foreground">
          Log in
        </div>
        {getProjectResponse?.project?.logInWithGoogle && (
          <div>
            <AuthPreviewButton className="w-full" variant="outline">
              <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M22.56 12.25C22.56 11.47 22.49 10.72 22.36 10H12V14.26H17.92C17.66 15.63 16.88 16.79 15.71 17.57V20.34H19.28C21.36 18.42 22.56 15.6 22.56 12.25Z"
                  fill="#4285F4"
                />
                <path
                  d="M12 23C14.97 23 17.46 22.02 19.28 20.34L15.71 17.57C14.73 18.23 13.48 18.63 12 18.63C9.13999 18.63 6.70999 16.7 5.83999 14.1H2.17999V16.94C3.98999 20.53 7.69999 23 12 23Z"
                  fill="#34A853"
                />
                <path
                  d="M5.84 14.09C5.62 13.43 5.49 12.73 5.49 12C5.49 11.27 5.62 10.57 5.84 9.91V7.07H2.18C1.43 8.55 1 10.22 1 12C1 13.78 1.43 15.45 2.18 16.93L5.03 14.71L5.84 14.09Z"
                  fill="#FBBC05"
                />
                <path
                  d="M12 5.38C13.62 5.38 15.06 5.94 16.21 7.02L19.36 3.87C17.45 2.09 14.97 1 12 1C7.69999 1 3.98999 3.47 2.17999 7.07L5.83999 9.91C6.70999 7.31 9.13999 5.38 12 5.38Z"
                  fill="#EA4335"
                />
              </svg>
              Log in with Google
            </AuthPreviewButton>
          </div>
        )}
        {getProjectResponse?.project?.logInWithMicrosoft && (
          <div className="mt-2">
            <AuthPreviewButton className="w-full" variant="outline">
              <svg
                width="24"
                height="23"
                viewBox="0 0 24 23"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <g clipPath="url(#clip0_65_191)">
                  <path d="M0 0H24V23H0V0Z" fill="#F3F3F3" />
                  <path d="M1.04347 1H11.4783V11H1.04347V1Z" fill="#F35325" />
                  <path d="M12.5217 1H22.9565V11H12.5217V1Z" fill="#81BC06" />
                  <path d="M1.04347 12H11.4783V22H1.04347V12Z" fill="#05A6F0" />
                  <path d="M12.5217 12H22.9565V22H12.5217V12Z" fill="#FFBA08" />
                </g>
                <defs>
                  <clipPath id="clip0_65_191">
                    <rect width="24" height="23" fill="white" />
                  </clipPath>
                </defs>
              </svg>
              Log in with Microsoft
            </AuthPreviewButton>
          </div>
        )}
        {getProjectResponse?.project?.logInWithGithub && (
          <div className="mt-2">
            <AuthPreviewButton className="w-full" variant="outline">
              <svg
                height="24"
                width="24"
                fill="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  fillRule="evenodd"
                  d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                  clipRule="evenodd"
                />
              </svg>
              Log in with Github
            </AuthPreviewButton>
          </div>
        )}

        <AuthPreviewTextDivider className="mt-6" variant="tight">
          or continue with email
        </AuthPreviewTextDivider>

        <div className="mt-4">
          <div className="text-xs font-semibold mb-1 cursor-default">Email</div>
          <AuthPreviewInput />
          <AuthPreviewButton className="w-full mt-2">Log in</AuthPreviewButton>
        </div>
      </div>

      <div className="mt-4 text-xs text-muted-foreground text-center cursor-default">
        Don't have an account?{" "}
        <span className="text-primary underline">Sign up</span>
      </div>
    </div>
  );
}
