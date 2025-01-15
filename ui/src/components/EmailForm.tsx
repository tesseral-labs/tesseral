import React, {
  ChangeEvent,
  FormEvent,
  useCallback,
  useEffect,
  useState,
} from 'react'
import { useNavigate } from 'react-router'
import { useMutation } from '@connectrpc/connect-query'
import debounce from 'lodash.debounce'

import { setIntermediateSessionToken } from '@/auth'
import { Button } from './ui/button'
import {
  listSAMLOrganizations,
  signInWithEmail,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { LoginViews } from '@/lib/views'
import { Organization } from '@/gen/openauth/intermediate/v1/intermediate_pb'
import TextDivider from './ui/TextDivider'

const EmailForm = () => {
  const navigate = useNavigate()
  const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/i
  const signInWithEmailMutation = useMutation(signInWithEmail)
  const listSAMLOrganizationsMutation = useMutation(listSAMLOrganizations)

  const [email, setEmail] = useState<string>('')
  const [emailIsValid, setEmailIsValid] = useState<boolean>(false)
  const [samlOrganizations, setSamlOrganizations] = useState<Organization[]>([])

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

    try {
      const { intermediateSessionToken, challengeId } =
        await signInWithEmailMutation.mutateAsync({
          email,
        })

      // set the intermediate sessionToken
      setIntermediateSessionToken(intermediateSessionToken)

      // redirect to challenge page
      navigate(`/login`, {
        state: {
          view: LoginViews.EmailVerification,
          challengeId,
        },
      })
    } catch (error) {
      console.error(error)
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
      <form className="flex flex-col justify-center" onSubmit={handleSubmit}>
        <label
          className="text-center uppercase text-foreground font-semibold text-sm mb-6 tracking-wide"
          htmlFor="email"
        >
          Continue with Email
        </label>
        <input
          className="text-sm bg-input rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
          id="email"
          type="email"
          onChange={handleEmail}
          placeholder="jane.doe@email.com"
          value={email}
        />
        <Button type="submit" disabled={!emailIsValid}>
          Sign In
        </Button>
      </form>

      {samlOrganizations && samlOrganizations.length > 0 && (
        <>
          <TextDivider text="or" />

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
                className="w-full"
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
