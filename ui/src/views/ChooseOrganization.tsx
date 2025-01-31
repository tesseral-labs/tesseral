import React, { Dispatch, FC, SetStateAction, useEffect } from 'react'

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
import { Link, useNavigate } from 'react-router-dom'
import { setAccessToken, setRefreshToken } from '@/auth'
import { LoginViews } from '@/lib/views'

interface ChooseOrganizationProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const ChooseOrganization: FC<ChooseOrganizationProps> = ({ setView }) => {
  const navigate = useNavigate()

  const { data: whoamiRes } = useQuery(whoami)
  const { data: listOrganizationsResponse } = useQuery(listOrganizations)
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )
  const setOrganizationMutation = useMutation(setOrganization)

  const isEmailLogin = (): boolean => {
    return (
      !whoamiRes?.intermediateSession?.googleUserId &&
      !whoamiRes?.intermediateSession?.microsoftUserId &&
      !!whoamiRes?.intermediateSession?.email
    )
  }

  const handleOrganizationClick = async (organization: Organization) => {
    try {
      await setOrganizationMutation.mutateAsync({
        organizationId: organization.id,
      })

      // Check if the user is logging in with an email address and the organization supports passwords
      if (isEmailLogin() && organization.logInWithPasswordEnabled) {
        setView(
          organization.userExists
            ? LoginViews.VerifyPassword
            : LoginViews.RegisterPassword,
        )
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
