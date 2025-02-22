import React, { Dispatch, FC, useState } from 'react';
import { LoginView } from '@/lib/views';
import { Title } from '@/components/Title';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useIntermediateSession } from '@/lib/auth';
import { useMutation } from '@connectrpc/connect-query';
import { verifyEmailChallenge } from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from '@/components/ui/input-otp';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';

interface VerifyEmailViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>;
}

const VerifyEmailView: FC<VerifyEmailViewProps> = ({ setView }) => {
  const intermediateSession = useIntermediateSession();

  const [challengeCode, setChallengeCode] = useState<string>('');

  const verifyEmailChallengeMutation = useMutation(verifyEmailChallenge);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await verifyEmailChallengeMutation.mutateAsync({
        code: challengeCode,
      });

      setView(LoginView.ChooseProject);
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error('Could not verify email address', {
        description: message,
      });
    }
  };

  return (
    <>
      <Title title="Verify Email Address" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Verify Email Address</CardTitle>
          <CardDescription>
            Please enter the verification code sent to{' '}
            <b>{intermediateSession?.email}</b> below.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="flex flex-col items-center" onSubmit={handleSubmit}>
            <InputOTP maxLength={6} onChange={setChallengeCode}>
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>

            <Button
              className="mt-4"
              disabled={challengeCode.length < 6}
              type="submit"
            >
              Continue
            </Button>
          </form>
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </>
  );
};

export default VerifyEmailView;
