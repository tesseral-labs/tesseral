import React, { FC, useState } from 'react';
import {
  NavigationMenu,
  NavigationMenuList,
} from '@/components/ui/navigation-menu';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  listSwitchableOrganizations,
  logout,
  whoami,
} from '@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery';
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom';
import {
  BookOpen,
  Bug,
  Building2,
  ChevronDown,
  LifeBuoy,
  LogOut,
  Settings2,
  User,
} from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
} from './ui/dropdown-menu';
import { DropdownMenuTrigger } from '@radix-ui/react-dropdown-menu';
import { Avatar, AvatarFallback } from './ui/avatar';
import { API_URL } from '@/config';
import { toast } from 'sonner';
import {
  getOrganization,
  getProject,
  listOrganizations,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { cn, titleCaseSlug } from '@/lib/utils';

const ConsoleNavigation: FC = () => {
  const navigate = useNavigate();
  const { pathname } = useLocation();

  const { data: whoamiResponse } = useQuery(whoami);
  const { mutateAsync: logoutAsync } = useMutation(logout);

  const handleLogout = async () => {
    await logoutAsync({});
    toast.success('You have been logged out.');
    navigate('/login');
  };

  return (
    <header className="bg-white h-16 w-full px-2 py-4 border-b fixed z-10">
      <nav className="flex flex-row items-center justify-between space-x-2 w-full px-2">
        <div className="flex flex-row items-center space-x-2">
          <Link to="/">
            <img className="max-h-8" src="/images/tesseral-icon-black.svg" />
          </Link>
          <div className="mr-auto">
            <NavigationMenu>
              <NavigationMenuList>
                <NavigationProjects />
                {pathname.split('/').map((slug, index) => (
                  <>
                    {slug !== '' && (
                      <>
                        <div className="font-thin text-foreground-muted mx-2">
                          /
                        </div>
                        {index === 1 ? (
                          <NavigationProjectPages
                            slug={slug as 'project-settings' | 'organizations'}
                          />
                        ) : index === 3 &&
                          [
                            'users',
                            'user-invites',
                            'saml-connections',
                            'scim-api-keys',
                          ].includes(slug) ? (
                          <NavigationOrganizationPages slug={slug as any} />
                        ) : (
                          <div className="px-2 text-sm font-medium">
                            {index === pathname.split('/').length - 1 ||
                            [
                              'publishable-keys',
                              'backend-api-keys',
                              'rbac-policy',
                            ].includes(slug) ? (
                              <span className="text-muted-foreground">
                                {titleCaseSlug(
                                  slug,
                                  index === pathname.split('/').length - 1,
                                )}
                              </span>
                            ) : (
                              <Link
                                className="text-sm font-medium"
                                to={`${pathname
                                  .split('/')
                                  .slice(0, index + 1)
                                  .join('/')}`}
                              >
                                {titleCaseSlug(
                                  slug,
                                  index === pathname.split('/').length - 1,
                                )}
                              </Link>
                            )}
                          </div>
                        )}
                      </>
                    )}
                  </>
                ))}
              </NavigationMenuList>
            </NavigationMenu>
          </div>
        </div>
        <div className="ml-auto flex-col space-x-2 text-sm items-end justify-center">
          <div className="hidden xl:inline-flex space-x-2 mr-2">
            <Link
              className="text-muted-foreground hover:text-foreground"
              target="_blank"
              to="https://tesseral.com/docs"
            >
              <BookOpen className="inline max-h-4" />
              Docs
            </Link>
            <Link
              className="text-muted-foreground hover:text-foreground"
              target="_blank"
              to="https://github.com/tesseral-labs/tesseral/issues/new"
            >
              <Bug className="inline max-h-4" />
              Report
            </Link>
            <Link
              className="text-muted-foreground hover:text-foreground"
              target="_blank"
              to="mailto:support@tesseral.com"
            >
              <LifeBuoy className="inline max-h-4" />
              Support
            </Link>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger className="inline-flex items-center">
              <Avatar className="h-8 w-8 rounded-full">
                <AvatarFallback className="rounded-full bg-indigo-600 text-white font-semibold">
                  {whoamiResponse?.user?.email?.substring(0, 1)?.toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <ChevronDown className="max-h-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel className="p-0 font-normal">
                <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                  <Avatar className="h-8 w-8 rounded-full">
                    <AvatarFallback className="rounded-full bg-indigo-600 text-white font-semibold">
                      {whoamiResponse?.user?.email
                        ?.substring(0, 1)
                        ?.toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">
                      {whoamiResponse?.user?.email}
                    </span>
                    <span className="truncate text-xs">
                      {whoamiResponse?.user?.email}
                    </span>
                  </div>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuGroup>
                <DropdownMenuLabel>Settings</DropdownMenuLabel>
                <DropdownMenuItem>
                  <Link to={`${API_URL}/user-settings`}>
                    <User className="inline max-h-4" />
                    User Settings
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Link to={`${API_URL}/organization-settings`}>
                    <Building2 className="inline max-h-4" />
                    Organization Settings
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <div className="block xl:hidden">
                <DropdownMenuGroup>
                  <DropdownMenuLabel>Resources</DropdownMenuLabel>
                  <DropdownMenuItem>
                    <BookOpen className="inline max-h-4" />
                    <Link
                      className="text-sm font-medium"
                      target="_blank"
                      to="https://tesseral.com/docs"
                    >
                      Docs
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <Bug className="inline max-h-4" />
                    <Link
                      className="text-sm font-medium"
                      target="_blank"
                      to="https://github.com/tesseral-labs/tesseral/issues/new"
                    >
                      Report
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <LifeBuoy className="inline max-h-4" />
                    <Link
                      className="text-sm font-medium"
                      target="_blank"
                      to="mailto:support@tesseral.com"
                    >
                      Support
                    </Link>
                  </DropdownMenuItem>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
              </div>
              <DropdownMenuGroup>
                <DropdownMenuItem onClick={handleLogout}>
                  <LogOut className="inline max-h-4" />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuGroup>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </nav>
    </header>
  );
};

const NavigationProjects = () => {
  const [open, setOpen] = useState(false);

  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: listSwitchableOrganizationsResponse } = useQuery(
    listSwitchableOrganizations,
    {},
  );

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger className="text-sm font-medium ring-0 active:ring-0 focus:ring-0">
        {getProjectResponse?.project?.displayName}
        <ChevronDown
          className={cn(
            'inline max-h-3 transition-transform',
            open ? 'rotate-180' : 'rotate-none',
          )}
        />
      </DropdownMenuTrigger>
      <DropdownMenuContent className="block w-[300px]">
        {listSwitchableOrganizationsResponse?.switchableOrganizations?.map(
          (org) => (
            <DropdownMenuItem
              asChild
              className="w-full font-medium text-sm p-2"
              key={org.id}
            >
              <Link
                className="h-full w-full"
                id={org.id}
                to={`/switch-organizations/${org.id}`}
              >
                <Avatar className="h-8 w-8 rounded-full">
                  <AvatarFallback className="rounded-full bg-muted-foreground/15 text-muted-foreground font-semibold">
                    {org.displayName?.substring(0, 1)?.toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                {org.displayName}
              </Link>
            </DropdownMenuItem>
          ),
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

const NavigationProjectPages = ({
  slug,
}: {
  slug: 'project-settings' | 'organizations';
}) => {
  const [open, setOpen] = useState(false);
  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger className="text-sm font-medium ring-0 active:ring-0 focus:ring-0 px-2">
        {slug === 'project-settings' ? 'Project Settings' : 'Organizations'}
        <ChevronDown
          className={cn(
            'inline max-h-3 transition-transform',
            open ? 'rotate-180' : 'rotate-none',
          )}
        />
      </DropdownMenuTrigger>
      <DropdownMenuContent className="block w-[300px]">
        <DropdownMenuItem
          asChild
          className="block w-full font-medium text-sm  p-2"
        >
          <Link className="h-full w-full" to="/project-settings">
            <Settings2 className="inline max-h-4 mr-2" />
            Project Settings
          </Link>
        </DropdownMenuItem>

        <DropdownMenuItem
          asChild
          className="block w-full font-medium text-sm  p-2"
        >
          <Link className="h-full w-full" to="/Organizations">
            <Building2 className="inline max-h-4 mr-2" />
            Organizations
          </Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

const NavigationOrganizations = () => {
  const { organizationId } = useParams();
  const [open, setOpen] = useState(false);

  const { data: organizationResponse } = useQuery(getOrganization, {
    id: organizationId,
  });
  const { data: listOrganizationsResponse } = useQuery(listOrganizations, {
    pageToken: '',
  });

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger className="text-sm font-medium ring-0 active:ring-0 focus:ring-0 px-2">
        {organizationResponse?.organization?.displayName}
        <ChevronDown
          className={cn(
            'inline max-h-3 transition-transform',
            open ? 'rotate-180' : 'rotate-none',
          )}
        />
      </DropdownMenuTrigger>
      <DropdownMenuContent className="block w-[300px]">
        {listOrganizationsResponse?.organizations?.map((org) => (
          <DropdownMenuItem
            asChild
            className="block w-full font-medium text-sm p-2"
            key={org.id}
          >
            <Link
              className="h-full w-full"
              id={org.id}
              to={`/organizations/${org.id}`}
            >
              {org.displayName}
            </Link>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

const NavigationOrganizationPages = ({
  slug,
}: {
  slug:
    | 'details'
    | 'users'
    | 'user-invites'
    | 'saml-connections'
    | 'scim-api-keys';
}) => {
  const { organizationId } = useParams();
  const [open, setOpen] = useState(false);

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger className="text-sm font-medium ring-0 active:ring-0 focus:ring-0 px-2">
        {titleCaseSlug(slug)}
        <ChevronDown
          className={cn(
            'inline max-h-3 transition-transform',
            open ? 'rotate-180' : 'rotate-none',
          )}
        />
      </DropdownMenuTrigger>
      <DropdownMenuContent className="block w-[300px]">
        <DropdownMenuItem
          asChild
          className="block w-full font-medium text-sm p-2"
        >
          <Link
            className="h-full w-full"
            id={organizationId}
            to={`/organizations/${organizationId}`}
          >
            Details
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem
          asChild
          className="block w-full font-medium text-sm p-2"
        >
          <Link
            className="h-full w-full"
            id={organizationId}
            to={`/organizations/${organizationId}/users`}
          >
            Users
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem
          asChild
          className="block w-full font-medium text-sm p-2"
        >
          <Link
            className="h-full w-full"
            id={organizationId}
            to={`/organizations/${organizationId}/user-invites`}
          >
            User Invites
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem
          asChild
          className="block w-full font-medium text-sm p-2"
        >
          <Link
            className="h-full w-full"
            id={organizationId}
            to={`/organizations/${organizationId}/saml-connections`}
          >
            SAML Connections
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem
          asChild
          className="block w-full font-medium text-sm p-2"
        >
          <Link
            className="h-full w-full"
            id={organizationId}
            to={`/organizations/${organizationId}/scim-api-keys`}
          >
            SCIM API Keys
          </Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default ConsoleNavigation;
