import React, { Dispatch, FC, SetStateAction } from 'react'
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

      navigate('/settings')
    } catch (e) {
      // TODO: Display error message to user
      console.error('Error exchanging session for tokens', e)
    }
  }

  return (
    <>
      <Title title="Choose an Project" />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Choose a Project
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center w-full">
          <ul className="w-full p-0 border border-b-0 rounded-md">
            {listOrganizationsResponse?.organizations?.map(
              (organization, idx) => (
                <li
                  className={`py-2 px-4 border-b ${idx === listOrganizationsResponse?.organizations?.length ? 'rounded-b-md' : ''} hover:bg-gray-50 hover:text-dark cursor-pointer font-semibold`}
                  key={organization.id}
                  onClick={() => handleOrganizationClick(organization)}
                >
                  {organization.displayName}
                </li>
              ),
            )}
          </ul>
        </CardContent>
        <CardFooter>
          <p className="text-sm text-center w-full cursor-pointer">
            Or you can{' '}
            <span
              className="text-primary underline"
              onClick={() => setView(LoginView.CreateProject)}
            >
              create a project
            </span>
            .
          </p>
        </CardFooter>
      </Card>
    </>
  )
}

export default ChooseProjectView
