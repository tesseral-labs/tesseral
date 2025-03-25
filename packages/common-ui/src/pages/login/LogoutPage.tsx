import { useMutation } from "@connectrpc/connect-query";
import { LoaderCircleIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

import { logout } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { clearAccessToken } from "@/lib/access-token";

export function LogoutPage() {
  const { mutateAsync: logoutAsync } = useMutation(logout);
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      await logoutAsync({});
      clearAccessToken();
      toast.success("You have been logged out.");
      navigate("/login");
    })();
  }, [logoutAsync, navigate]);

  return (
    <LoaderCircleIcon className="mx-auto text-muted-foreground h-4 w-4 animate-spin" />
  );
}
