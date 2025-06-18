import { useMutation } from "@connectrpc/connect-query";
import { useCallback } from "react";
import { toast } from "sonner";

import { createStripeCheckoutLink } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";

export function useHandleUpgrade(): () => void {
  const createStripeCheckoutLinkMutation = useMutation(
    createStripeCheckoutLink,
  );

  return useCallback(async () => {
    try {
      const { url } = await createStripeCheckoutLinkMutation.mutateAsync({});
      window.location.href = url;
    } catch {
      toast.error("Somthing went wrong. Please try again later.");
    }
  }, [createStripeCheckoutLinkMutation]);
}
