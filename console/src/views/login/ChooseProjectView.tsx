import React, { Dispatch, FC, SetStateAction, useState } from 'react'
import { LoginView } from '@/lib/views'
import { useNavigate } from 'react-router'
import {
  exchangeIntermediateSessionForSession,
  listOrganizations,
  setOrganization,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import { Organization } from '@/gen/openauth/intermediate/v1/intermediate_pb'
import { setAccessToken, setRefreshToken } from '@/auth'
import { Title } from '@/components/Title'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import TextDivider from '@/components/ui/text-divider'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'
import Loader from '@/components/ui/loader'

interface ChooseProjectViewProps {
  setIntermediateOrganization: Dispatch<
    SetStateAction<Organization | undefined>
  >
  setView: Dispatch<SetStateAction<LoginView>>
}

const ChooseProjectView: FC<ChooseProjectViewProps> = ({
  setIntermediateOrganization,
  setView,
}) => {
  const navigate = useNavigate()

  const [submitting, setSubmitting] = useState<boolean>(false)

  const { data: listOrganizationsResponse } = useQuery(listOrganizations)
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const setOrganizationMutation = useMutation(setOrganization)

  // This function is effectively the central routing layer for Organization selection.
  // It determines the next view based on the current organization's settings and the user's current state.
  // If the organization requires MFA
  // - if the primary login factor is not valid, the user is redirected to the ChooseOrganizationPrimaryLoginFactor view
  // - if the organization has passwords enabled
  //   - if the user has a password, they are redirected to the VerifyPassword view
  //   - if the user does not have a password, they are redirected to the RegisterPassword view
  // - if the organization's only secondary factor is authenticator apps
  //   - if the user has an authenticator app registered, they are redirected to the VerifyAuthenticatorApp view
  //   - if the user does not have an authenticator app registered, they are redirected to the RegisterAuthenticatorApp view
  // - if the organization's only secondary factor is passkeys
  //   - if the user has a passkey registered, they are redirected to the VerifyPasskey view
  //   - if the user does not have a passkey registered, they are redirected to the RegisterPasskey view
  // - if the organization has multiple secondary factors
  //   - the user is redirected to the ChooseAdditionalFactor view
  const deriveNextView = (
    organization: Organization,
  ): LoginView | undefined => {
    if (organization.logInWithPassword) {
      if (organization.userHasPassword) {
        return LoginView.VerifyPassword
      }

      return LoginView.RegisterPassword
    } else if (
      organization.userHasAuthenticatorApp &&
      !organization.userHasPasskey
    ) {
      return LoginView.VerifyAuthenticatorApp
    } else if (
      organization.userHasPasskey &&
      !organization.userHasAuthenticatorApp
    ) {
      return LoginView.VerifyPasskey
    } else if (organization.requireMfa) {
      // this is the case where the organization has multiple secondary factors and requires mfa
      return LoginView.ChooseAdditionalFactor
    }
  }

  const handleOrganizationClick = async (organization: Organization) => {
    setSubmitting(true)
    try {
      await setOrganizationMutation.mutateAsync({
        organizationId: organization.id,
      })

      setIntermediateOrganization(organization)

      const nextView = deriveNextView(organization)

      console.log(`nextView:`, nextView)

      if (nextView) {
        setView(nextView)
        return
      }

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({
          organizationId: organization.id,
        })

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      setSubmitting(false)
      navigate('/settings')
    } catch (error) {
      setSubmitting(false)
      const message = parseErrorMessage(error)
      toast.error('Could not select Project', {
        description: message,
      })
    }
  }

  return (
    <>
      <Title title="Choose an Project" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Choose a Project</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <ul className="w-full p-0">
            {listOrganizationsResponse?.organizations?.map((organization) => (
              <li key={organization.id}>
                <Button
                  className="w-full"
                  onClick={() => handleOrganizationClick(organization)}
                  variant="outline"
                >
                  {submitting && <Loader />}
                  {organization.displayName}
                </Button>
              </li>
            ))}
          </ul>
        </CardContent>
        <CardFooter>
          <div className="w-full">
            <TextDivider className="w-full">Or you can</TextDivider>

            <Button
              className="w-full"
              onClick={() => setView(LoginView.CreateProject)}
            >
              Create a Project
            </Button>
          </div>
        </CardFooter>
      </Card>
    </>
  )
}

export default ChooseProjectView
