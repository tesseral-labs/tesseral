import { useEffect } from "react";
import React from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

export function StripeCheckoutSuccessPage() {
  const navigate = useNavigate();
  useEffect(() => {
    // sonner does not honor toasts on mount, so we delay here
    const id = setTimeout(() => {
      toast.success("Your payment has been processed successfully.");
      navigate("/");
    });

    return () => clearTimeout(id);
  }, [navigate]);

  return <></>;
}
