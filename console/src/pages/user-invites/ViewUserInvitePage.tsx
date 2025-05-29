import { useNavigate, useParams } from 'react-router';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  deleteUserInvite,
  getOrganization,
  getUserInvite,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import React from 'react';
import {
  PageCodeSubtitle,
  PageContent,
  PageDescription,
  PageHeader,
  PageTitle,
} from '@/components/page';
import {
  DetailsGrid,
  DetailsGridColumn,
  DetailsGridEntry,
  DetailsGridKey,
  DetailsGridValue,
} from '@/components/details-grid';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';

export const ViewUserInvitePage = () => {
  const { organizationId, userInviteId } = useParams();
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: getUserInviteResponse } = useQuery(getUserInvite, {
    id: userInviteId,
  });

  return (
    <>
      <PageHeader>
        <PageTitle>
          User Invite for {getUserInviteResponse?.userInvite?.email}
        </PageTitle>
        <PageCodeSubtitle>{userInviteId}</PageCodeSubtitle>
        <PageDescription>
          A User Invite lets outside collaborators join an organization.
        </PageDescription>
      </PageHeader>

      <PageContent>
        <Card className="my-8">
          <CardHeader>
            <CardTitle>General settings</CardTitle>
            <CardDescription>
              Basic settings for this user invite.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <DetailsGrid>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Email</DetailsGridKey>
                  <DetailsGridValue>
                    {getUserInviteResponse?.userInvite?.email}
                  </DetailsGridValue>
                </DetailsGridEntry>
                <DetailsGridEntry>
                  <DetailsGridKey>Owner</DetailsGridKey>
                  <DetailsGridValue>
                    {getUserInviteResponse?.userInvite?.owner ? 'Yes' : 'No'}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Created</DetailsGridKey>
                  <DetailsGridValue>
                    {getUserInviteResponse?.userInvite?.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(
                          getUserInviteResponse?.userInvite?.createTime,
                        ),
                      ).toRelative()}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
              <DetailsGridColumn>
                <DetailsGridEntry>
                  <DetailsGridKey>Updated</DetailsGridKey>
                  <DetailsGridValue>
                    {getUserInviteResponse?.userInvite?.updateTime &&
                      DateTime.fromJSDate(
                        timestampDate(
                          getUserInviteResponse?.userInvite?.updateTime,
                        ),
                      ).toRelative()}
                  </DetailsGridValue>
                </DetailsGridEntry>
              </DetailsGridColumn>
            </DetailsGrid>
          </CardContent>
        </Card>

        <DangerZoneCard />
      </PageContent>
    </>
  );
};

const DangerZoneCard = () => {
  const { organizationId, userInviteId } = useParams();

  const deleteUserInviteMutation = useMutation(deleteUserInvite);
  const navigate = useNavigate();
  const handleDelete = async () => {
    await deleteUserInviteMutation.mutateAsync({
      id: userInviteId,
    });

    toast.success('User invite deleted');
    navigate(`/organizations/${organizationId}/user-invites`);
  };

  return (
    <>
      <Card className="border-destructive">
        <CardHeader>
          <CardTitle>Danger Zone</CardTitle>
        </CardHeader>

        <CardContent>
          <div className="flex justify-between items-center">
            <div>
              <div className="text-sm font-semibold">Delete User Invite</div>
              <p className="text-sm">
                Delete this user invite. You can recreate a new user invite with
                the same email at any time.
              </p>
            </div>

            <Button variant="destructive" onClick={handleDelete}>
              Delete User Invite
            </Button>
          </div>
        </CardContent>
      </Card>
    </>
  );
};
