import { useMutation, useQuery } from "@connectrpc/connect-query";
import {
  BookOpen,
  Bug,
  Building2,
  ChevronDown,
  Home,
  Key,
  LifeBuoy,
  Lock,
  LogOut,
  Settings,
  Settings2,
  Shield,
  User,
  Webhook,
} from "lucide-react";
import React from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { toast } from "sonner";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import { API_URL } from "@/config";
import {
  getProject,
  getProjectWebhookManagementURL,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import {
  listSwitchableOrganizations,
  logout,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { cn } from "@/lib/utils";

import { Separator } from "../ui/separator";
import { BreadcrumbBar } from "./BreadcrumbBar";

export function Navigation() {
  const { pathname } = useLocation();

  return (
    <header className="w-full sticky top-0 z-10">
      <nav className="p-4 w-full z-50 bg-white/90 backdrop-blur supports-[backdrop-filter]:bg-white/80 flex flex-row items-center justify-between border-b lg:border-0">
        <div className="flex items-center">
          <NavigationMenu className="relative ">
            <NavigationMenuList className="relative mr-auto">
              <NavigationMenuItem>
                <Link to="/">
                  <img
                    className="max-h-8"
                    src="/images/tesseral-icon-black.svg"
                  />
                </Link>
              </NavigationMenuItem>
              <NavigationProjects />
              <NavigationMenuItem>
                <NavigationMenuLink
                  active={pathname === "/"}
                  asChild
                  className={navigationMenuTriggerStyle()}
                >
                  <Link to="/">
                    <div className="flex items-center">
                      <Home className="inline h-4 w-4 mr-2" />
                      Home
                    </div>
                  </Link>
                </NavigationMenuLink>
              </NavigationMenuItem>
              <NavigationMenuItem>
                <NavigationMenuLink
                  active={pathname.startsWith("/organizations")}
                  className={navigationMenuTriggerStyle()}
                  asChild
                >
                  <Link to="/organizations">
                    <div className="flex items-center">
                      <Building2 className="h-4 w-4 mr-2" />
                      Organizations
                    </div>
                  </Link>
                </NavigationMenuLink>
              </NavigationMenuItem>
            </NavigationMenuList>
          </NavigationMenu>
          <NavigationMenu className="relative">
            <NavigationMenuList>
              <NavigationSettings />
            </NavigationMenuList>
          </NavigationMenu>
        </div>
        <NavigationMenu>
          <NavigationMenuList className="ml-auto">
            <NavigationUser />
          </NavigationMenuList>
        </NavigationMenu>
      </nav>
      <BreadcrumbBar />
    </header>
  );
}

function NavigationProjects() {
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: listSwitchableOrganizationsResponse } = useQuery(
    listSwitchableOrganizations,
    {},
  );

  return (
    <NavigationMenuItem>
      <NavigationMenuTrigger className="text-sm font-medium ring-0 active:ring-0 focus:ring-0">
        {getProjectResponse?.project?.displayName}
      </NavigationMenuTrigger>
      <NavigationMenuContent>
        <div className="font-semibold mb-2 text-xs px-2">Projects</div>
        <Separator className="mb-2" />
        <div className="w-[300px] space-y-2">
          {listSwitchableOrganizationsResponse?.switchableOrganizations?.map(
            (org) => (
              <NavigationMenuLink
                key={org.id}
                asChild
                className={cn(navigationMenuTriggerStyle(), "w-full")}
              >
                <Link
                  className="h-full w-full"
                  id={org.id}
                  to={`/switch-organizations/${org.id}`}
                >
                  <div className="flex items-center justify-start w-full font-medium text-xs">
                    <Avatar className="mr-4 h-6 w-6 rounded-full">
                      <AvatarFallback className="rounded-full bg-muted-foreground/15 text-muted-foreground text-sm font-semibold">
                        {org.displayName?.substring(0, 1)?.toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    {org.displayName}
                  </div>
                </Link>
              </NavigationMenuLink>
            ),
          )}
        </div>
      </NavigationMenuContent>
    </NavigationMenuItem>
  );
}

function NavigationSettings() {
  const { data: getProjectWebhookManagementUrlResponse } = useQuery(
    getProjectWebhookManagementURL,
  );

  return (
    <NavigationMenuItem>
      <NavigationMenuTrigger className={navigationMenuTriggerStyle()}>
        <Settings className="mr-1 h-4 w-4" />
        Settings
      </NavigationMenuTrigger>
      <NavigationMenuContent>
        <div className="font-semibold mb-2 text-xs px-2">Project Settings</div>
        <Separator className="mb-2" />
        <ul className="grid gap-2 w-[300px] grid-cols-1">
          <ListItem to="/settings/authentication" title="Authentication">
            <div className="text-sm leading-none font-medium">
              <Shield className="inline mr-2" />
              Authentication
            </div>
            <p className="text-muted-foreground line-clamp-2 text-xs leading-snug">
              Configure SAML, SCIM, OAuth, and MFA
            </p>
          </ListItem>
          <ListItem to="/settings/api-keys">
            <div className="text-sm leading-none font-medium">
              <Key className="inline mr-2" />
              API Keys
            </div>
            <p className="text-muted-foreground line-clamp-2 text-xs leading-snug">
              Manage API Keys and Publishable Keys
            </p>
          </ListItem>
          <ListItem to="/settings/vault">
            <div className="text-sm leading-none font-medium">
              <Settings2 className="inline mr-2" />
              Vault Customization
            </div>
            <p className="text-muted-foreground line-clamp-2 text-xs leading-snug">
              Customize the appearance and configuration of your Vault pages
            </p>
          </ListItem>
          <ListItem to="/settings/access">
            <div className="text-sm leading-none font-medium">
              <Lock className="inline mr-2" />
              Access Control
            </div>
            <p className="text-muted-foreground line-clamp-2 text-xs leading-snug">
              Configure Role-based Access Control (RBAC) for your project.
            </p>
          </ListItem>
          <ListItem to={getProjectWebhookManagementUrlResponse?.url || ""}>
            <div className="text-sm leading-none font-medium">
              <Webhook className="inline mr-2" />
              Webhooks
            </div>
            <p className="text-muted-foreground line-clamp-2 text-xs leading-snug">
              Configure Webhook endpoints for sync events
            </p>
          </ListItem>
        </ul>
      </NavigationMenuContent>
    </NavigationMenuItem>
  );
}

function NavigationUser() {
  const navigate = useNavigate();
  const { mutateAsync: logoutAsync } = useMutation(logout);
  const { data: whoamiResponse } = useQuery(whoami);

  async function handleLogout() {
    await logoutAsync({});
    toast.success("You have been logged out.");
    navigate("/login");
  }

  return (
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
                {whoamiResponse?.user?.email?.substring(0, 1)?.toUpperCase()}
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
              <User className="inline max-h-4 mr-2" />
              User Settings
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <Link to={`${API_URL}/organization-settings`}>
              <Building2 className="inline max-h-4 mr-2" />
              Organization Settings
            </Link>
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <div>
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
  );
}

function ListItem({
  children,
  to,
  ...props
}: React.ComponentPropsWithoutRef<"li"> & { to: string }) {
  return (
    <li {...props}>
      <NavigationMenuLink asChild>
        <Link to={to}>{children}</Link>
      </NavigationMenuLink>
    </li>
  );
}
