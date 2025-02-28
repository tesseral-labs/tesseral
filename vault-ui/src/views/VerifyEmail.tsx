import React, {
  Dispatch,
  FC,
  SetStateAction,
  useEffect,
  useState,
} from 'react';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { toast } from 'sonner';

import { parseErrorMessage } from '@/lib/errors';
import { useLayout } from '@/lib/settings';
import { cn } from '@/lib/utils';
import { LoginLayouts, LoginViews } from '@/lib/views';

import { Title } from '@/components/Title';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  issueEmailVerificationChallenge,
  verifyEmailChallenge,
  whoami,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from '@/components/ui/input-otp';
import {
  ArrowRightIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  MailIcon,
} from 'lucide-react';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { useNavigate } from 'react-router';

interface VerifyEmailProps {
  setView: Dispatch<SetStateAction<LoginViews>>;
}

const VerifyEmail: FC<VerifyEmailProps> = () => {
  const { data: whoamiRes } = useQuery(whoami);

  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  );
  const [hasResent, setHasResent] = useState(false);
  const handleResend = async () => {
    await issueEmailVerificationChallengeMutation.mutateAsync({
      email: whoamiRes?.intermediateSession?.email,
    });

    toast.success('New verification link sent');
    setHasResent(true);
  };

  useEffect(() => {
    // allow another send after 10 seconds
    setTimeout(() => {
      setHasResent(false);
    }, 10000);
  }, [hasResent]);

  const [isAdvancedOpen, setIsAdvancedOpen] = useState(false);
  const [verificationCode, setVerificationCode] = useState('');
  const verifyEmailChallengeMutation = useMutation(verifyEmailChallenge);
  const navigate = useNavigate();
  const handleManualVerification = async () => {
    let code = verificationCode
    if (!code.startsWith("email_verification_challenge_code_")) {
      // This code is definitely incorrect. The user may have copy-pasted the
      // URL instead. Try to fall back to that.
      try {
        code = new URL(code).searchParams.get("code") ?? '';
      } catch {
        // ignore
      }
    }

    await verifyEmailChallengeMutation.mutateAsync({
      code,
    });

    toast.success('Email verified');

    navigate(`/login?view=${LoginViews.ChooseOrganization}`, {
      replace: true,
    });
  };

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="space-y-1">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
          <MailIcon className="h-6 w-6 text-primary" />
        </div>
        <CardTitle className="text-center text-2xl font-bold">
          Check your email
        </CardTitle>
        <CardDescription className="text-center">
          We've sent a verification link to{' '}
          <span className="font-medium">
            {whoamiRes?.intermediateSession?.email}
          </span>
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="text-center text-sm text-muted-foreground">
          <p>Didn't receive an email? Check your spam folder or</p>
        </div>

        <Button
          variant="outline"
          className="w-full"
          onClick={handleResend}
          disabled={hasResent}
        >
          {hasResent
            ? 'Email verification resent!'
            : 'Resend verification link'}
        </Button>

        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t" />
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-card px-2 text-muted-foreground">Or</span>
          </div>
        </div>

        <Button
          variant="ghost"
          className="w-full justify-between"
          onClick={() => setIsAdvancedOpen(!isAdvancedOpen)}
        >
          <span>Enter verification code manually</span>
          {isAdvancedOpen ? (
            <ChevronUpIcon className="h-4 w-4" />
          ) : (
            <ChevronDownIcon className="h-4 w-4" />
          )}
        </Button>

        {isAdvancedOpen && (
          <div className="space-y-3">
            <div className="space-y-2">
              <Label htmlFor="verification-code">Verification Code</Label>
              <Input
                id="verification-code"
                placeholder="email_verification_challenge_code_..."
                value={verificationCode}
                onChange={(e) => setVerificationCode(e.target.value)}
                className="font-mono"
              />
              <p className="text-xs text-muted-foreground">
                Paste the full verification code from the email.
              </p>
            </div>
            <Button onClick={handleManualVerification} className="w-full">
              Verify <ArrowRightIcon className="ml-2 h-4 w-4" />
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default VerifyEmail;
