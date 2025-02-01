import React, { Dispatch, FC, SetStateAction, useEffect } from 'react'

import { Title } from '@/components/Title'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  exchangeIntermediateSessionForSession,
  listOrganizations,
  setOrganization,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import {
  IntermediateSession,
  Organization,
} from '@/gen/openauth/intermediate/v1/intermediate_pb'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Link, useNavigate } from 'react-router-dom'
import { setAccessToken, setRefreshToken } from '@/auth'
import { LoginViews } from '@/lib/views'

interface ChooseOrganizationProps {
  setIntermediateOrganization: Dispatch<
    SetStateAction<Organization | undefined>
  >
  setView: Dispatch<SetStateAction<LoginViews>>
}

const ChooseOrganization: FC<ChooseOrganizationProps> = ({
  setIntermediateOrganization,
  setView,
}) => {
  const navigate = useNavigate()

  const { data: whoamiRes } = useQuery(whoami)
  const { data: listOrganizationsResponse } = useQuery(listOrganizations)
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const setOrganizationMutation = useMutation(setOrganization)

  // This function is effectively the central routing layer for Organization selection.
  // It determines the next view based on the current organization's settings and the user's current state.
  // If the organization requires MFA
  // - if this is an OAuth login
  //   - if the user has not yet set up MFA
  //   - or the the organization requires multiple MFA factors and user has configured both
  //     - return the ChooseAdditionalFactor view
  //   - if the user has set up only one additional factor
  //     - return the Verify{Factor} view
  // - if this is an email login
  //   - return the VerifyPassword view
  // If the organization does not require MFA
  // - if the user has not setup MFA or password
  // - or the user has setup multiple factors (in this context, password included)
  //   - return the ChooseAdditionalFactor view
  // - if the user has setup a single MFA factor (in this context, password included)
  //   - return the Verify{Factor} view
  const deriveNextView = (
    organization: Organization,
    intermdiateSession: IntermediateSession,
  ): LoginViews | undefined => {
    if (organization.requireMfa) {
      if (
        whoamiRes?.intermediateSession?.googleUserId ||
        whoamiRes?.intermediateSession?.microsoftUserId
      ) {
        if (
          organization.userHasAuthenticatorApp &&
          !organization.userHasPasskey
        ) {
          return LoginViews.VerifyAuthenticatorApp
        } else if (
          organization.userHasPasskey &&
          !organization.userHasAuthenticatorApp
        ) {
          return LoginViews.VerifyPasskey
        }

        return LoginViews.ChooseAdditionalFactor
      } else if (organization.logInWithPasswordEnabled) {
        if (organization.userHasPassword) {
          return LoginViews.VerifyPassword
        }

        return LoginViews.RegisterPassword
      }
    } else {
      if (
        organization.logInWithAuthenticatorAppEnabled &&
        organization.userHasAuthenticatorApp &&
        !organization.userHasPasskey &&
        !organization.userHasPassword
      ) {
        return LoginViews.VerifyAuthenticatorApp
      } else if (
        organization.logInWithPasskeyEnabled &&
        organization.userHasPasskey &&
        !organization.userHasAuthenticatorApp &&
        !organization.userHasPassword
      ) {
        return LoginViews.VerifyPasskey
      } else if (
        organization.logInWithPasswordEnabled &&
        organization.userHasPassword &&
        !organization.userHasAuthenticatorApp &&
        !organization.userHasPasskey
      ) {
        return LoginViews.VerifyPassword
      }

      return LoginViews.ChooseAdditionalFactor
    }
  }

  const handleOrganizationClick = async (organization: Organization) => {
    const intermediateOrganization = {
      ...organization,
      requireMfa: true,
      logInWithPasskeyEnabled: true,
    }
    try {
      if (!whoamiRes?.intermediateSession) {
        throw new Error('No intermediate session found')
      }

      await setOrganizationMutation.mutateAsync({
        organizationId: organization.id,
      })

      setIntermediateOrganization(intermediateOrganization)

      // Check if the needs to provide additional factors
      const nextView = deriveNextView(
        intermediateOrganization,
        whoamiRes.intermediateSession,
      )

      console.log(`nextVivew: ${nextView}`)

      if (!!nextView) {
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
      <Title title="Choose an Organization" />

      <Card className="w-[clamp(320px,50%,420px)]">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Choose an Organization
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
          <p className="text-sm text-center w-full">
            Or you can{' '}
            <span
              className="text-primary underline cursor-pointer"
              onClick={() => setView(LoginViews.CreateOrganization)}
            >
              create an organization
            </span>
            .
          </p>
        </CardFooter>
      </Card>
    </>
  )
}

export default ChooseOrganization
