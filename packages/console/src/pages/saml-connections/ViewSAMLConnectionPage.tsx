import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { Link } from 'react-router-dom';
import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  deleteSAMLConnection,
  getOrganization,
  getSAMLConnection,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { toast } from 'sonner';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import {
  PageCodeSubtitle,
  PageDescription,
  PageTitle,
} from '@/components/page';

export const ViewSAMLConnectionPage = () => {
  const { organizationId, samlConnectionId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getSAMLConnectionResponse } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  });
  return (
    <div>
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

      <PageTitle>SAML Connection</PageTitle>
      <PageCodeSubtitle>{samlConnectionId}</PageCodeSubtitle>
      <PageDescription>
        A SAML connection is a link between Tesseral and your customer's
        enterprise Identity Provider. Lorem ipsum dolor.
      </PageDescription>

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
          <DetailsGrid>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>
                  Assertion Consumer Service (ACS) URL
                </DetailsGridKey>
                <DetailsGridValue>
                  {getSAMLConnectionResponse?.samlConnection?.spAcsUrl}
                </DetailsGridValue>
              </DetailsGridEntry>

              <DetailsGridEntry>
                <DetailsGridKey>SP Entity ID</DetailsGridKey>
                <DetailsGridValue>
                  {getSAMLConnectionResponse?.samlConnection?.spEntityId}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>IDP Entity ID</DetailsGridKey>
                <DetailsGridValue>
                  {getSAMLConnectionResponse?.samlConnection?.idpEntityId ||
                    '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>IDP Redirect URL</DetailsGridKey>
                <DetailsGridValue>
                  {getSAMLConnectionResponse?.samlConnection?.idpRedirectUrl ||
                    '-'}
                </DetailsGridValue>
              </DetailsGridEntry>
              <DetailsGridEntry>
                <DetailsGridKey>IDP Certificate</DetailsGridKey>
                <DetailsGridValue>
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
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
            <DetailsGridColumn>
              <DetailsGridEntry>
                <DetailsGridKey>Primary</DetailsGridKey>
                <DetailsGridValue>
                  {getSAMLConnectionResponse?.samlConnection?.primary
                    ? 'Yes'
                    : 'No'}
                </DetailsGridValue>
              </DetailsGridEntry>

              <DetailsGridEntry>
                <DetailsGridKey>Created</DetailsGridKey>
                <DetailsGridValue>
                  {getSAMLConnectionResponse?.samlConnection?.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getSAMLConnectionResponse?.samlConnection?.createTime,
                      ),
                    ).toRelative()}
                </DetailsGridValue>
              </DetailsGridEntry>

              <DetailsGridEntry>
                <DetailsGridKey>Updated</DetailsGridKey>
                <DetailsGridValue>
                  {getSAMLConnectionResponse?.samlConnection?.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(
                        getSAMLConnectionResponse?.samlConnection?.updateTime,
                      ),
                    ).toRelative()}
                </DetailsGridValue>
              </DetailsGridEntry>
            </DetailsGridColumn>
          </DetailsGrid>
        </CardContent>
      </Card>

      <DangerZoneCard />
    </div>
  );
};

const DangerZoneCard = () => {
  const { organizationId, samlConnectionId } = useParams();
  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false);

  const handleDelete = () => {
    setConfirmDeleteOpen(true);
  };

  const deleteSAMLConnectionMutation = useMutation(deleteSAMLConnection);
  const navigate = useNavigate();
  const handleConfirmDelete = async () => {
    await deleteSAMLConnectionMutation.mutateAsync({
      id: samlConnectionId,
    });

    toast.success('SAML connection deleted');
    navigate(`/organizations/${organizationId}/saml-connections`);
  };

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
                Delete this SAML connection. This cannot be undone.
              </p>
            </div>

            <Button variant="destructive" onClick={handleDelete}>
              Delete SAML Connection
            </Button>
          </div>
        </CardContent>
      </Card>
    </>
  );
};
