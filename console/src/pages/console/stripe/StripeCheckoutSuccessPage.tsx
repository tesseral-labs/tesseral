import { useEffect } from 'react';
import { toast } from 'sonner';
import { useNavigate } from 'react-router';
import React from 'react';

export function StripeCheckoutSuccessPage() {
  const navigate = useNavigate();
  useEffect(() => {
    // sonner does not honor toasts on mount, so we delay here
    const id = setTimeout(() => {
      toast.success("Your payment has been processed successfully.");
      navigate("/");
    })

    return () => clearTimeout(id);
  }, [navigate]);

  return <></>
}
