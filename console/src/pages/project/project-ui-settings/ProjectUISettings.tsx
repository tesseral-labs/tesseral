import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import AuthCardPreview from '@/components/vault-preview/auth-card-preview';
import React, {
  ChangeEvent,
  createRef,
  FC,
  FormEvent,
  SyntheticEvent,
  useEffect,
  useState,
} from 'react';
import SideBySideLayout from '@/components/vault-preview/layouts/side-by-side';
import CenterLayout from '@/components/vault-preview/layouts/center';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import {
  PageCodeSubtitle,
  PageDescription,
  PageTitle,
} from '@/components/page';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  getProject,
  getProjectUISettings,
  updateProjectUISettings,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import { Link } from 'react-router-dom';
import { Label } from '@/components/ui/label';
import {
  Check,
  Moon,
  SquareSplitHorizontal,
  SquareSquare,
  Sun,
} from 'lucide-react';
import { cn, hexToHSL, isColorDark } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { ColorPicker } from '@/components/ui/color-picker';
import { Switch } from '@/components/ui/switch';
import { parseErrorMessage } from '@/lib/errors';
import { toast } from 'sonner';

const settingsPage: FC = () => {
  const darkModeLogoPickerRef = createRef<HTMLInputElement>();
  const logoPickerRef = createRef<HTMLInputElement>();
  const previewRef = createRef<HTMLDivElement>();
  const { data: getProjectResponse } = useQuery(getProject, {});
  const { data: getProjectUISettingsResponse } = useQuery(getProjectUISettings);

  const updateProjectUISettingsMutation = useMutation(updateProjectUISettings);

  const [darkMode, setDarkMode] = useState<boolean>(false);
  const [darkModeLogo, setDarkModeLogo] = useState<string>();
  const [darkModeLogoFile, setDarkModeLogoFile] = useState<File | null>(null);
  const [darkModePrimaryColor, setDarkModePrimaryColor] =
    useState<string>('#ffffff');
  const [detectDarkModeEnabled, setDetectDarkModeEnabled] =
    useState<boolean>(false);
  const [layout, setLayout] = useState<string>('centered');
  const [logo, setLogo] = useState<string>();
  const [logoFile, setLogoFile] = useState<File | null>(null);
  const [primaryColor, setPrimaryColor] = useState('#0f172a');

  const applyTheme = () => {
    const root = previewRef.current as HTMLElement;

    const primary = primaryColor;
    const darkPrimary = darkModePrimaryColor;

    if (!darkMode && primary) {
      const foreground = isColorDark(primary) ? '0 0% 100%' : '0 0% 0%';

      root.style.setProperty('--primary', hexToHSL(primary));
      root.style.setProperty('--primary-foreground', foreground);
    }

    if (darkPrimary && darkMode) {
      const darkForeground = isColorDark(darkPrimary) ? '0 0% 100%' : '0 0% 0%';

      root.style.setProperty('--primary', hexToHSL(darkPrimary));
      root.style.setProperty('--primary-foreground', darkForeground);
    }
  };

  const handleDarkModeLogoChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) {
      return;
    }

    setDarkModeLogoFile(file);

    const reader = new FileReader();
    reader.onload = (e) => {
      const base64String = e.target?.result;

      if (typeof base64String !== 'string') {
        return;
      }

      setDarkModeLogo(base64String as string);
    };

    reader.readAsDataURL(file);
  };

  const handleLogoChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) {
      return;
    }

    setLogoFile(file);

    const reader = new FileReader();
    reader.onload = (e) => {
      const base64String = e.target?.result;

      if (typeof base64String !== 'string') {
        return;
      }

      setLogo(base64String as string);
    };

    reader.readAsDataURL(file);
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();

    let logoUploadUrl, darkModeLogoUploadUrl;
    try {
      const { logoPresignedUploadUrl, darkModeLogoPresignedUploadUrl } =
        await updateProjectUISettingsMutation.mutateAsync({
          logInLayout: layout,
          detectDarkModeEnabled,
          primaryColor,
          darkModePrimaryColor,
          logoContentType: logoFile?.type,
          darkModeLogoContentType: darkModeLogoFile?.type,
        });

      logoUploadUrl = logoPresignedUploadUrl;
      darkModeLogoUploadUrl = darkModeLogoPresignedUploadUrl;

      // special-case local development, where the s3 that api can dial isn't
      // the same s3 that the host can dial
      if (logoPresignedUploadUrl.startsWith('http://s3:9090/')) {
        logoUploadUrl = logoPresignedUploadUrl.replace(
          'http://s3:9090/',
          'https://tesseralusercontent.example.com/',
        );
      }
      if (darkModeLogoPresignedUploadUrl.startsWith('http://s3:9090/')) {
        darkModeLogoUploadUrl = darkModeLogoPresignedUploadUrl.replace(
          'http://s3:9090/',
          'https://tesseralusercontent.example.com/',
        );
      }
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error('Failed to update vault UI settings', {
        description: message,
      });
    }

    try {
      if (logoUploadUrl && logo) {
        const response = await fetch(logoUploadUrl, {
          body: logoFile,
          headers: {
            'Content-Type': logoFile?.type || 'image/png',
            'x-amz-meta-trigger': 'true',
          },
          method: 'PUT',
        });

        if (!response.ok) {
          throw new Error('Failed to upload logo');
        }
      }
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error('Failed to update logo', {
        description: message,
      });
    }

    try {
      if (darkModeLogoUploadUrl && darkModeLogo) {
        const response = await fetch(darkModeLogoUploadUrl, {
          body: darkModeLogoFile,
          headers: {
            'Content-Type': darkModeLogoFile?.type || 'image/png',
            'x-amz-meta-trigger': 'true',
          },
          method: 'PUT',
        });

        if (!response.ok) {
          throw new Error('Failed to upload dark mode logo');
        }
      }

      toast.success('Vault UI settings updated successfully');
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error('Failed to update dark mode logo', {
        description: message,
      });
    }
  };

  useEffect(() => {
    if (getProjectUISettingsResponse) {
      setLayout(
        getProjectUISettingsResponse.projectUiSettings?.logInLayout ||
          'centered',
      );
      setDetectDarkModeEnabled(
        getProjectUISettingsResponse.projectUiSettings?.detectDarkModeEnabled ||
          false,
      );
      setPrimaryColor(
        getProjectUISettingsResponse.projectUiSettings?.primaryColor ||
          '#0f172a',
      );
      setDarkModePrimaryColor(
        getProjectUISettingsResponse.projectUiSettings?.darkModePrimaryColor ||
          '#ffffff',
      );

      // Light mode logo
      if (getProjectUISettingsResponse.projectUiSettings?.logoUrl) {
        setLogo(getProjectUISettingsResponse.projectUiSettings.logoUrl);
      } else {
        setLogo('/images/tesseral-logo-black.svg');
      }

      // Dark mode logo
      if (getProjectUISettingsResponse.projectUiSettings?.darkModeLogoUrl) {
        setDarkModeLogo(
          getProjectUISettingsResponse.projectUiSettings.darkModeLogoUrl,
        );
      } else {
        setDarkModeLogo('/images/tesseral-logo-white.svg');
      }
    }
  }, [getProjectUISettingsResponse]);

  useEffect(() => {
    applyTheme();
  }, [darkMode, primaryColor, darkModePrimaryColor]);

  return (
    <div>
      <div className="mt-8">
        <form onSubmit={handleSubmit}>
          <Card>
            <CardHeader>
              <div className="grid grid-cols-2 gap-8">
                <div>
                  <CardTitle>Vault UI Settings</CardTitle>
                  <CardDescription>
                    This controls the layout, logo, and colors for your vault
                    login pages.
                  </CardDescription>
                </div>
                <div className="text-right">
                  <Button
                    disabled={
                      detectDarkModeEnabled ===
                        getProjectUISettingsResponse?.projectUiSettings
                          ?.detectDarkModeEnabled &&
                      layout ===
                        getProjectUISettingsResponse?.projectUiSettings
                          ?.logInLayout &&
                      (logo ===
                        getProjectUISettingsResponse?.projectUiSettings
                          ?.logoUrl ||
                        logo === '/images/tesseral-logo-black.svg') &&
                      (primaryColor ===
                        getProjectUISettingsResponse?.projectUiSettings
                          ?.primaryColor ||
                        primaryColor === '#0f172a') &&
                      (darkModePrimaryColor ===
                        getProjectUISettingsResponse?.projectUiSettings
                          ?.darkModePrimaryColor ||
                        darkModePrimaryColor === '#ffffff') &&
                      (darkModeLogo ===
                        getProjectUISettingsResponse?.projectUiSettings
                          ?.darkModeLogoUrl ||
                        darkModeLogo === '/images/tesseral-logo-white.svg')
                    }
                    type="submit"
                  >
                    Save changes
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-3 gap-8 max-h-[650px]">
                <div className="relative">
                  <div className="overflow-y-scroll pr-8 pb-24 h-[630px]">
                    <div>
                      <Label>Layout</Label>
                      <p className="text-sm text-muted-foreground">
                        The page layout to use for vault login pages.
                      </p>

                      <div className="grid grid-cols-2 gap-4 mt-4">
                        <div
                          className={cn(
                            'p-4 border rounded-sm relative',
                            layout === 'centered'
                              ? 'border-primary border-2 cursor-default'
                              : 'cursor-pointer',
                          )}
                          onClick={() => setLayout('centered')}
                        >
                          <div
                            className={cn(
                              'font-semibold text-sm',
                              layout === 'centered'
                                ? 'text-primary'
                                : 'text-muted-foreground',
                            )}
                          >
                            <SquareSquare
                              className="inline-block mr-2"
                              size={16}
                            />
                            Center card
                            {layout === 'centered' && (
                              <div className="h-5 w-5 text-white bg-primary rounded-full flex justify-center items-center absolute top-2 right-2">
                                <Check size={12} />
                              </div>
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-2">
                            A center-aligned card layout.
                          </p>
                        </div>

                        <div
                          className={cn(
                            'p-4 border rounded-sm relative',
                            layout === 'side_by_side'
                              ? 'border-primary border-2 cursor-default'
                              : 'cursor-pointer',
                          )}
                          onClick={() => setLayout('side_by_side')}
                        >
                          <div
                            className={cn(
                              'font-semibold text-sm',
                              layout === 'side_by_side'
                                ? 'text-primary'
                                : 'text-muted-foreground',
                            )}
                          >
                            <SquareSplitHorizontal
                              className="inline-block mr-2"
                              size={16}
                            />
                            Side by side
                            {layout === 'side_by_side' && (
                              <div className="h-5 w-5 text-white bg-primary rounded-full flex justify-center items-center absolute top-2 right-2">
                                <Check size={12} />
                              </div>
                            )}
                          </div>
                          <p className="text-xs text-muted-foreground mt-2">
                            A horizontally split layout.
                          </p>
                        </div>
                      </div>
                    </div>

                    <div className="mt-8">
                      <Label>Auto-detect dark mode</Label>
                      <p className="text-sm text-muted-foreground">
                        Automatically switch to dark mode based on user
                        preferences.
                      </p>

                      <Switch
                        checked={detectDarkModeEnabled}
                        className="mt-2"
                        onCheckedChange={(checked) => {
                          if (!checked && darkMode) {
                            setDarkMode(false);
                          }
                          setDetectDarkModeEnabled(checked);
                        }}
                      />
                    </div>

                    <div className="mt-8">
                      <div>
                        <Label>Logo</Label>
                        <p className="text-sm text-muted-foreground">
                          The logo to display on the vault login page.
                        </p>

                        <div
                          className="group cursor-pointer mt-2 p-4 rounded-sm border inline-block relative"
                          onClick={() => logoPickerRef.current?.click()}
                        >
                          {logo && (
                            <img
                              className="max-h-8"
                              src={logo}
                              onError={(
                                e: SyntheticEvent<HTMLImageElement, Event>,
                              ) => {
                                const target = e.target as HTMLImageElement;
                                target.onerror = null;
                                target.src = '/images/tesseral-logo-black.svg';
                              }}
                            />
                          )}

                          <div className="logo-overlay absolute top-0 left-0 w-full h-full hidden group-hover:block">
                            <div className="bg-black bg-opacity-75 h-full w-full rounded-sm" />
                            <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 text-white text-sm font-semibold">
                              Update logo
                            </div>
                          </div>
                        </div>
                        <input
                          className="hidden"
                          onChange={handleLogoChange}
                          ref={logoPickerRef}
                          type="file"
                        />
                      </div>
                      <div className="mt-8">
                        <Label>Primary color</Label>
                        <p className="text-sm text-muted-foreground">
                          The accent color used in the vault UI.
                        </p>

                        <ColorPicker
                          className="mt-2"
                          onChange={setPrimaryColor}
                          value={primaryColor}
                        />
                      </div>

                      {detectDarkModeEnabled && (
                        <>
                          <div className="mt-8">
                            <Label>Dark mode logo</Label>
                            <p className="text-sm text-muted-foreground">
                              The logo to display on the vault login page in
                              dark mode.
                            </p>

                            <div
                              className="group mt-2 p-4 rounded-sm border inline-block bg-primary cursor-pointer relative"
                              onClick={() =>
                                darkModeLogoPickerRef.current?.click()
                              }
                            >
                              {darkModeLogo && (
                                <img
                                  className="max-h-8"
                                  src={darkModeLogo}
                                  onError={(
                                    e: SyntheticEvent<HTMLImageElement, Event>,
                                  ) => {
                                    const target = e.target as HTMLImageElement;
                                    target.onerror = null;
                                    target.src =
                                      '/images/tesseral-logo-white.svg';
                                  }}
                                />
                              )}

                              <div className="logo-overlay absolute top-0 left-0 w-full h-full hidden group-hover:block">
                                <div className="bg-black bg-opacity-75 h-full w-full rounded-sm" />
                                <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 text-white text-sm font-semibold">
                                  Update logo
                                </div>
                              </div>
                            </div>

                            <input
                              className="hidden"
                              onChange={handleDarkModeLogoChange}
                              ref={darkModeLogoPickerRef}
                              type="file"
                            />
                          </div>
                          <div className="mt-8">
                            <Label>Dark mode primary color</Label>
                            <p className="text-sm text-muted-foreground">
                              The accent color used in the vault UI in dark
                              mode.
                            </p>

                            <ColorPicker
                              className="mt-2"
                              onChange={setDarkModePrimaryColor}
                              value={darkModePrimaryColor}
                            />
                          </div>
                        </>
                      )}
                    </div>
                  </div>
                  <div className="absolute bottom-0 h-24 bg-gradient-to-b from-white/0 to-white p-8 w-full"></div>
                </div>

                <div className="col-span-2">
                  <Card>
                    <CardHeader>
                      <CardTitle className="grid grid-cols-2">
                        <div className="text-lg">Preview</div>
                        {detectDarkModeEnabled && (
                          <div className="relative text-right">
                            <Switch
                              checked={darkMode}
                              onCheckedChange={setDarkMode}
                            />
                            {darkMode ? (
                              <div className="absolute top-1 right-6 text-white">
                                <Moon size={16} />
                              </div>
                            ) : (
                              <div className="absolute top-1 right-1 text-muted-foreground">
                                <Sun size={16} />
                              </div>
                            )}
                          </div>
                        )}
                      </CardTitle>
                      <CardDescription>
                        A preview of how your vault login page will look.
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div
                        className={cn('rounded border', darkMode ? 'dark' : '')}
                        ref={previewRef}
                      >
                        {layout === 'side_by_side' ? (
                          <SideBySideLayout>
                            <AuthCardPreview
                              darkMode={darkMode}
                              logo={darkMode ? darkModeLogo : logo}
                              noBorder
                            />
                          </SideBySideLayout>
                        ) : (
                          <CenterLayout>
                            <AuthCardPreview
                              darkMode={darkMode}
                              logo={darkMode ? darkModeLogo : logo}
                            />
                          </CenterLayout>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </CardContent>
          </Card>
        </form>
      </div>
    </div>
  );
};

export default settingsPage;
