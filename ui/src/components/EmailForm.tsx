import React, { ChangeEvent, FormEvent, useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import { useMutation } from '@connectrpc/connect-query'

import { setIntermediateSessionToken } from '@/auth'
import { Button } from './ui/button'
import { signInWithEmail } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'

const EmailForm = () => {
  const navigate = useNavigate()
  const signInWithEmailMutation = useMutation(signInWithEmail)

  const [email, setEmail] = useState<string>('')
  const [emailIsValid, setEmailIsValid] = useState<boolean>(false)

  const handleEmail = (e: ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value)
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
      navigate(`/verify-email?challengeId=${challengeId}`)
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/
    setEmailIsValid(emailRegex.test(email))
  }, [email])

  return (
    <form className="flex flex-col justify-center" onSubmit={handleSubmit}>
      <label
        className="text-center uppercase text-foreground font-semibold text-sm mb-6 tracking-wide"
        htmlFor="email"
      >
        Continue with Email
      </label>
      <input
        className="text-sm rounded border border-border focus:border-primary w-[clamp(240px,50%,100%)] mb-2"
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
  )
}

export default EmailForm
