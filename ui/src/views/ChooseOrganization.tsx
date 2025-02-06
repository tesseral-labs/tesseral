import React, { Dispatch, FC, SetStateAction, useEffect, useState } from 'react'

import { Title } from '@/components/Title'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  exchangeIntermediateSessionForSession,
  listOrganizations,
  setOrganization,
  whoami,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { Organization } from '@/gen/openauth/intermediate/v1/intermediate_pb'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { useNavigate } from 'react-router-dom'
import { setAccessToken, setRefreshToken } from '@/auth'
import { LoginLayouts, LoginViews } from '@/lib/views'
import { cn } from '@/lib/utils'
import { useLayout } from '@/lib/settings'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'
import Loader from '@/components/ui/loader'
import {
  isValidPrimaryLoginFactor,
  PrimaryLoginFactor,
} from '@/lib/login-factors'

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
  const layout = useLayout()
  const navigate = useNavigate()

  const [setting, setSetting] = useState<boolean>(false)

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
  ): LoginViews | undefined => {
    const primaryLoginFactor =
      whoamiRes?.intermediateSession?.primaryLoginFactor

    if (
      primaryLoginFactor &&
      !isValidPrimaryLoginFactor(
        primaryLoginFactor as PrimaryLoginFactor,
        organization,
      )
    ) {
      return LoginViews.ChooseOrganizationPrimaryLoginFactor
    }

    if (organization.requireMfa) {
      if (
        whoamiRes?.intermediateSession?.googleUserId ||
        whoamiRes?.intermediateSession?.microsoftUserId
      ) {
        if (organization.logInWithPassword) {
          if (organization.userHasPassword) {
            return LoginViews.VerifyPassword
          }

          return LoginViews.RegisterPassword
        } else if (
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
      }
    } else {
      if (organization.logInWithPassword) {
        if (organization.userHasPassword) {
          return LoginViews.VerifyPassword
        }

        return LoginViews.RegisterPassword
      } else if (
        organization.logInWithAuthenticatorApp &&
        organization.userHasAuthenticatorApp &&
        !organization.userHasPasskey
      ) {
        return LoginViews.VerifyAuthenticatorApp
      } else if (
        organization.logInWithPasskey &&
        organization.userHasPasskey &&
        !organization.userHasAuthenticatorApp
      ) {
        return LoginViews.VerifyPasskey
      }

      return LoginViews.ChooseAdditionalFactor
    }
  }

  const handleOrganizationClick = async (organization: Organization) => {
    setSetting(true)
    const intermediateOrganization = {
      ...organization,
    }
    try {
      if (!whoamiRes?.intermediateSession) {
        throw new Error('No intermediate session found')
      }

      await setOrganizationMutation.mutateAsync({
        organizationId: organization.id,
      })

      setSetting(false)
      setIntermediateOrganization(intermediateOrganization)
    } catch (error) {
      setSetting(false)
      const message = parseErrorMessage(error)
      toast.error('Could not set organization', {
        description: message,
      })
    }

    try {
      // Check if the needs to provide additional factors
      const nextView = deriveNextView(intermediateOrganization)
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
      setSetting(false)

      navigate('/settings')
    } catch (error) {
      setSetting(false)
      const message = parseErrorMessage(error)
      toast.error('Could not set organization', {
        description: message,
      })
    }
  }

  return (
    <>
      <Title title="Choose an Organization" />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          <CardTitle className="text-center">Choose an Organization</CardTitle>
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
                  {setting && <Loader />}
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
