import React, { useEffect } from 'react'

import { Title } from '@/components/Title'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  exchangeIntermediateSessionForSession,
  listOrganizations,
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
import { Button } from '@/components/ui/button'
import { setAccessToken, setRefreshToken } from '@/auth'
import { LoginViews } from '@/lib/views'

const Organizations = () => {
  const navigate = useNavigate()

  const [organizations, setOrganizations] = React.useState<Organization[]>([])
  const [pageToken, setPageToken] = React.useState('')

  const { data: whoamiRes } = useQuery(whoami)
  const listOrganizationsMutation = useMutation(listOrganizations)
  const exchangeIntermediateSessionForSessionMutation = useMutation(
    exchangeIntermediateSessionForSession,
  )

  const fetchOrganizations = async () => {
    const { nextPageToken, organizations: organizationsRes } =
      await listOrganizationsMutation.mutateAsync({
        pageToken,
      })

    setOrganizations([...organizations, ...organizationsRes])

    if (nextPageToken) {
      setPageToken(nextPageToken)
    } else {
      setPageToken('')
    }
  }

  const handleOrganizationClick = async (organization: Organization) => {
    try {
      // Check if the user is logging in with an email address and the organization supports passwords
      if (
        !whoamiRes?.googleUserId &&
        !whoamiRes?.microsoftUserId &&
        whoamiRes?.email &&
        organization.logInWithPasswordEnabled
      ) {
        navigate(`/login`, {
          state: {
            view: LoginViews.PasswordVerification,
            organizationId: organization.id,
          },
        })
        return
      }

      const { accessToken, refreshToken } =
        await exchangeIntermediateSessionForSessionMutation.mutateAsync({
          organizationId: organization.id,
        })

      setAccessToken(accessToken)
      setRefreshToken(refreshToken)

      navigate('/session-info')
    } catch (e) {
      // TODO: Display error message to user
      console.error('Error exchanging session for tokens', e)
    }
  }

  useEffect(() => {
    fetchOrganizations()
  }, [])

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
            {organizations.map((organization, idx) => (
              <li
                className={`py-2 px-4 border-b ${idx === organizations.length ? 'rounded-b-md' : ''} hover:bg-gray-50 hover:text-dark cursor-pointer font-semibold`}
                key={organization.id}
                onClick={() => handleOrganizationClick(organization)}
              >
                {organization.displayName}
              </li>
            ))}
          </ul>

          {pageToken && pageToken.length > 0 && (
            <Button
              className="mt-4"
              onClick={fetchOrganizations}
              variant="outline"
            >
              Load More
            </Button>
          )}
        </CardContent>
        <CardFooter>
          <p className="text-sm text-center w-full">
            Or you can{' '}
            <Link
              className="text-primary underline"
              to="/login"
              state={{ view: LoginViews.CreateOrganization }}
            >
              create an organization
            </Link>
            .
          </p>
        </CardFooter>
      </Card>
    </>
  )
}

export default Organizations
