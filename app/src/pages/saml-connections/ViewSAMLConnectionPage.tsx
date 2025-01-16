import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Link } from 'react-router-dom'
import React, { useState } from 'react'
import { useNavigate, useParams } from 'react-router'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import {
  deleteSAMLConnection,
  getOrganization,
  getSAMLConnection,
} from '@/gen/openauth/backend/v1/backend-BackendService_connectquery'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { DateTime } from 'luxon'
import { timestampDate } from '@bufbuild/protobuf/wkt'
import { toast } from 'sonner'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'

export function ViewSAMLConnectionPage() {
  const { organizationId, samlConnectionId } = useParams()
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  })
  const { data: getSAMLConnectionResponse } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  })
  return (
    // TODO remove padding when app shell in place
    <div className="pt-8">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/">Home</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to="/organizations">Organizations</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to={`/organizations/${organizationId}`}>
                {getOrganizationResponse?.organization?.displayName}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link to={`/organizations/${organizationId}/saml-connections`}>
                SAML Connections
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{samlConnectionId}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <h1 className="mt-4 font-semibold text-2xl">SAML Connection</h1>
      <span className="mt-1 inline-block border rounded bg-gray-100 py-1 px-2 font-mono text-xs text-muted-foreground">
        {samlConnectionId}
      </span>
      <div className="mt-4">
        A SAML connection is a link between Tesseral and your customer's
        enterprise Identity Provider. Lorem ipsum dolor.
      </div>

      <Card className="my-8">
        <CardHeader className="flex-row justify-between items-center">
          <div className="flex flex-col space-y-1 5">
            <CardTitle>Configuration</CardTitle>
            <CardDescription>Lorem ipsum dolor.</CardDescription>
          </div>
          <Button variant="outline" asChild>
            <Link
              to={`/organizations/${organizationId}/saml-connections/${samlConnectionId}/edit`}
            >
              Edit
            </Link>
          </Button>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-x-2 text-sm">
            <div className="border-r border-gray-200 pr-8 flex flex-col gap-4">
              <div>
                <div className="font-semibold">
                  Assertion Consumer Service (ACS) URL
                </div>
                <div className="truncate">
                  {getSAMLConnectionResponse?.samlConnection?.spAcsUrl}
                </div>
              </div>

              <div>
                <div className="font-semibold">SP Entity ID</div>
                <div className="truncate">
                  {getSAMLConnectionResponse?.samlConnection?.spEntityId}
                </div>
              </div>
            </div>
            <div className="border-r border-gray-200 pr-8 pl-8 flex flex-col gap-4">
              <div>
                <div className="font-semibold">IDP Entity ID</div>
                <div className="truncate">
                  {getSAMLConnectionResponse?.samlConnection?.idpEntityId ||
                    '-'}
                </div>
              </div>
              <div>
                <div className="font-semibold">IDP Redirect URL</div>
                <div className="truncate">
                  {getSAMLConnectionResponse?.samlConnection?.idpRedirectUrl ||
                    '-'}
                </div>
              </div>
              <div>
                <div className="font-semibold">IDP Certificate</div>
                <div>
                  {getSAMLConnectionResponse?.samlConnection
                    ?.idpX509Certificate ? (
                    <a
                      className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                      download={`Certificate ${samlConnectionId}.crt`}
                      href={`data:text/plain;base64,${btoa(getSAMLConnectionResponse.samlConnection.idpX509Certificate)}`}
                    >
                      Download (.crt)
                    </a>
                  ) : (
                    '-'
                  )}
                </div>
              </div>
            </div>
            <div className="border-gray-200 pl-8 flex flex-col gap-4">
              <div>
                <div className="font-semibold">Primary</div>
                <div className="truncate">
                  {getSAMLConnectionResponse?.samlConnection?.primary
                    ? 'Yes'
                    : 'No'}
                </div>
              </div>

              <div>
                <div className="font-semibold">Created</div>
                <div className="truncate">
                  {getSAMLConnectionResponse?.samlConnection?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getSAMLConnectionResponse?.samlConnection?.createTime,
                      ),
                    ).toRelative()}
                </div>
              </div>

              <div>
                <div className="font-semibold">Updated</div>
                <div className="truncate">
                  {getSAMLConnectionResponse?.samlConnection?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getSAMLConnectionResponse?.samlConnection?.updateTime,
                      ),
                    ).toRelative()}
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <DangerZoneCard />
    </div>
  )
}

function DangerZoneCard() {
  const { organizationId, samlConnectionId } = useParams()
  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false)

  const handleDelete = () => {
    setConfirmDeleteOpen(true)
  }

  const deleteSAMLConnectionMutation = useMutation(deleteSAMLConnection)
  const navigate = useNavigate()
  const handleConfirmDelete = async () => {
    await deleteSAMLConnectionMutation.mutateAsync({
      id: samlConnectionId,
    })

    toast.success('SAML connection deleted')
    navigate(`/organizations/${organizationId}/saml-connections`)
  }

  return (
    <>
      <AlertDialog open={confirmDeleteOpen} onOpenChange={setConfirmDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete SAML Connection?</AlertDialogTitle>
            <AlertDialogDescription>
              Deleting a SAML connection cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button variant="destructive" onClick={handleConfirmDelete}>
              Permanently Delete SAML Connection
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Card className="border-destructive">
        <CardHeader>
          <CardTitle>Danger Zone</CardTitle>
        </CardHeader>

        <CardContent>
          <div className="flex justify-between items-center">
            <div>
              <div className="text-sm font-semibold">
                Delete SAML Connection
              </div>
              <p className="text-sm">
                Delete this SAML Connection. This cannot be undone.
              </p>
            </div>

            <Button variant="destructive" onClick={handleDelete}>
              Delete SAML Connection
            </Button>
          </div>
        </CardContent>
      </Card>
    </>
  )
}
