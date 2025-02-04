import React, {
  ChangeEvent,
  Dispatch,
  FC,
  FormEvent,
  SetStateAction,
  useCallback,
  useEffect,
  useState,
} from 'react'
import { useMutation } from '@connectrpc/connect-query'
import debounce from 'lodash.debounce'

import { setIntermediateSessionToken } from '@/auth'
import { Button } from './ui/button'
import {
  createIntermediateSession,
  issueEmailVerificationChallenge,
  listSAMLOrganizations,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { LoginViews } from '@/lib/views'
import { Organization } from '@/gen/openauth/intermediate/v1/intermediate_pb'
import TextDivider from './ui/test-divider'
import { Input } from './ui/input'
import { Label } from './ui/label'
import Loader from './ui/loader'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'

interface EmailFormProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const EmailForm: FC<EmailFormProps> = ({ setView }) => {
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/i
  const createIntermediateSessionMutation = useMutation(
    createIntermediateSession,
  )
  const issueEmailVerificationChallengeMutation = useMutation(
    issueEmailVerificationChallenge,
  )
  const listSAMLOrganizationsMutation = useMutation(listSAMLOrganizations)

  const [email, setEmail] = useState<string>('')
  const [emailIsValid, setEmailIsValid] = useState<boolean>(false)
  const [samlOrganizations, setSamlOrganizations] = useState<Organization[]>([])
  const [submitting, setSubmitting] = useState<boolean>(false)

  const fetchSamlOrganizations = useCallback(
    debounce(async () => {
      const { organizations } = await listSAMLOrganizationsMutation.mutateAsync(
        {
          email,
        },
      )

      setSamlOrganizations(organizations)
    }, 300),
    [email],
  )

  const handleEmail = (e: ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value)
  }

  const handleSAMLLogin = async (samlConnectId: string) => {
    location.href = `/api/saml/v1/${samlConnectId}/init`
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setSubmitting(true)

    try {
      // this sets a cookie that subsequent requests use
      const { intermediateSessionSecretToken } =
        await createIntermediateSessionMutation.mutateAsync({})

      // set the intermediate sessionToken
      setIntermediateSessionToken(intermediateSessionSecretToken)

      await issueEmailVerificationChallengeMutation.mutateAsync({
        email,
      })

      setSubmitting(false)

      // redirect to challenge page
      setView(LoginViews.VerifyEmail)
    } catch (error) {
      setSubmitting(false)
      const message = parseErrorMessage(error)
      toast.error('Could not initiate login', {
        description: message,
      })
    }
  }

  useEffect(() => {
    const valid = emailRegex.test(email)
    setEmailIsValid(valid)

    if (valid) {
      fetchSamlOrganizations()
    }
  }, [email])

  return (
    <>
      <form
        className="flex flex-col justify-center w-full"
        onSubmit={handleSubmit}
      >
        <div className="grid gap-2">
          <Label htmlFor="email">Email</Label>
          <Input
            className="w-full mb-2"
            id="email"
            type="email"
            onChange={handleEmail}
            placeholder="jane.doe@email.com"
            value={email}
          />
        </div>

        <Button type="submit" disabled={!emailIsValid || submitting}>
          {submitting && <Loader />}
          Sign In
        </Button>
      </form>

      {samlOrganizations && samlOrganizations.length > 0 && (
        <>
          <TextDivider>or continue with SAML</TextDivider>

          {samlOrganizations.map((organization) => (
            <div key={organization.id} className="flex flex-col items-center">
              <label
                className="text-center uppercase text-foreground font-semibold text-sm mb-3 tracking-wide w-full"
                htmlFor="email"
              >
                Continue with SAML
              </label>
              <a
                href={`/api/saml/v1/${organization.primarySamlConnectionId}/init`}
                className="w-[clamp(240px,50%,100%)]"
              >
                <Button variant="outline">{organization.displayName}</Button>
              </a>
            </div>
          ))}
        </>
      )}
    </>
  )
}

export default EmailForm
