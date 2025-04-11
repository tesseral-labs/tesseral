import React, { FC } from 'react';
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from '@/components/ui/navigation-menu';
import { useQuery } from '@connectrpc/connect-query';
import {
  getProject,
  listSwitchableOrganizations,
} from '@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery';
import { Link } from 'react-router-dom';

const ConsoleNavigation: FC = () => {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: listSwitchableOrganizationsResponse } = useQuery(
    listSwitchableOrganizations,
    {},
  );

  return (
    <header className="bg-white w-full p-2 pb-4">
      <nav className="flex flex-row items-center space-x-2">
        <Link to="/">
          <img className="max-h-8" src="/images/tesseral-icon-black.svg" />
        </Link>
        <div>
          <NavigationMenu>
            <NavigationMenuList>
              <NavigationMenuItem>
                <NavigationMenuTrigger>
                  {getProjectResponse?.project?.displayName}
                </NavigationMenuTrigger>
                <NavigationMenuContent className="w-[240px]">
                  {listSwitchableOrganizationsResponse?.switchableOrganizations?.map(
                    (org) => (
                      <Link id={org.id} to={`/switch-organizations/${org.id}`}>
                        <NavigationMenuLink className="font-medium">
                          {org.displayName}
                        </NavigationMenuLink>
                      </Link>
                    ),
                  )}
                </NavigationMenuContent>
              </NavigationMenuItem>
            </NavigationMenuList>
          </NavigationMenu>
        </div>
      </nav>
    </header>
  );
};

export default ConsoleNavigation;
