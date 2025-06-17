import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Check,
  Moon,
  SquareSplitHorizontal,
  SquareSquare,
  Sun,
} from "lucide-react";
import React, { ChangeEvent, createRef, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ColorPicker } from "@/components/ui/color-picker";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { AuthCardPreview } from "@/components/vault-preview/AuthCardPreview";
import {
  getProject,
  getProjectUISettings,
  updateProjectUISettings,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { parseErrorMessage } from "@/lib/errors";
import { cn, isColorDark } from "@/lib/utils";

const hexRegexp = /^#(?:[0-9a-fA-F]{3}){1,2}$/;

const schema = z.object({
  darkModePrimaryColor: z
    .string()
    .min(1, "Dark mode primary color is required"),
  detectDarkModeEnabled: z.boolean(),
  primaryColor: z.string().min(1, "Primary color is required"),
  logInLayout: z.enum(["centered", "side_by_side"]),
  logo: z.string().optional(),
  darkModeLogo: z.string().optional(),
});

export function VaultBrandingSettingsTab() {
  const darkModeLogoPickerRef = createRef<HTMLInputElement>();
  const logoPickerRef = createRef<HTMLInputElement>();
  const previewRef = createRef<HTMLDivElement>();

  const { data: getProjectUISettingsResponse, refetch } =
    useQuery(getProjectUISettings);
  const updateProjectUISettingsMutation = useMutation(updateProjectUISettings);

  const [darkMode, setDarkMode] = useState<boolean>(false);
  const [darkModeLogo, setDarkModeLogo] = useState<string>();
  const [darkModeLogoFile, setDarkModeLogoFile] = useState<File | null>(null);
  const [darkModePrimaryColor, setDarkModePrimaryColor] =
    useState<string>("#ffffff");
  const [detectDarkModeEnabled, setDetectDarkModeEnabled] =
    useState<boolean>(false);
  const [logo, setLogo] = useState<string>();
  const [logoFile, setLogoFile] = useState<File | null>(null);
  const [primaryColor, setPrimaryColor] = useState("#0f172a");

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      detectDarkModeEnabled: false,
    },
  });

  // Handle form state changes and update state variables accordingly
  form.subscribe({
    formState: { values: true },
    callback: function ({ values }) {
      if (values.detectDarkModeEnabled !== detectDarkModeEnabled) {
        setDetectDarkModeEnabled(values.detectDarkModeEnabled);
      }

      if (
        values.primaryColor &&
        values.primaryColor !== primaryColor &&
        values.primaryColor.length >= 7
      ) {
        if (hexRegexp.test(values.primaryColor)) {
          setPrimaryColor(values.primaryColor);
        } else if (
          hexRegexp.test(`#${values.primaryColor}`) &&
          values.primaryColor.length >= 6
        ) {
          setPrimaryColor(`#${values.primaryColor}`);
          form.setValue("primaryColor", `#${values.primaryColor}`);
        }
      }

      if (
        values.darkModePrimaryColor &&
        values.darkModePrimaryColor !== darkModePrimaryColor
      ) {
        if (
          hexRegexp.test(values.darkModePrimaryColor) &&
          values.darkModePrimaryColor.length >= 7
        ) {
          setDarkModePrimaryColor(values.darkModePrimaryColor);
        } else if (
          hexRegexp.test(`#${values.darkModePrimaryColor}`) &&
          values.darkModePrimaryColor.length >= 6
        ) {
          setDarkModePrimaryColor(`#${values.darkModePrimaryColor}`);
          form.setValue(
            "darkModePrimaryColor",
            `#${values.darkModePrimaryColor}`,
          );
        }
      }
    },
  });

  async function handleDarkModeLogoChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) {
      return;
    }
    setDarkModeLogoFile(file);

    const reader = new FileReader();
    reader.onload = (e) => {
      const base64String = e.target?.result;
      if (typeof base64String !== "string") {
        return;
      }
      setDarkModeLogo(base64String as string);
    };

    reader.readAsDataURL(file);
  }

  async function handleLogoChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) {
      return;
    }
    setLogoFile(file);

    const reader = new FileReader();
    reader.onload = (e) => {
      const base64String = e.target?.result;
      if (typeof base64String !== "string") {
        return;
      }
      setLogo(base64String as string);
    };

    reader.readAsDataURL(file);
  }

  async function handleSubmit(data: z.infer<typeof schema>) {
    let logoUploadUrl, darkModeLogoUploadUrl;
    try {
      const { logoPresignedUploadUrl, darkModeLogoPresignedUploadUrl } =
        await updateProjectUISettingsMutation.mutateAsync({
          logInLayout: data.logInLayout,
          detectDarkModeEnabled: data.detectDarkModeEnabled,
          primaryColor: data.primaryColor,
          darkModePrimaryColor: data.darkModePrimaryColor,
        });

      logoUploadUrl = logoPresignedUploadUrl;
      darkModeLogoUploadUrl = darkModeLogoPresignedUploadUrl;

      // special-case local development, where the s3 that api can dial isn't
      // the same s3 that the host can dial
      if (logoPresignedUploadUrl.startsWith("http://s3:9090/")) {
        logoUploadUrl = logoPresignedUploadUrl.replace(
          "http://s3:9090/",
          "https://tesseralusercontent.example.com/",
        );
      }
      if (darkModeLogoPresignedUploadUrl.startsWith("http://s3:9090/")) {
        darkModeLogoUploadUrl = darkModeLogoPresignedUploadUrl.replace(
          "http://s3:9090/",
          "https://tesseralusercontent.example.com/",
        );
      }
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error("Failed to update vault UI settings", {
        description: message,
      });
      await refetch();
      return;
    }

    try {
      if (logoUploadUrl && logo) {
        const response = await fetch(logoUploadUrl, {
          body: logoFile,
          method: "PUT",
        });

        if (!response.ok) {
          throw new Error("Failed to upload logo");
        }
      }
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error("Failed to update logo", {
        description: message,
      });
      await refetch();
      return;
    }

    try {
      if (darkModeLogoUploadUrl && darkModeLogo) {
        const response = await fetch(darkModeLogoUploadUrl, {
          body: darkModeLogoFile,
          method: "PUT",
        });

        if (!response.ok) {
          throw new Error("Failed to upload dark mode logo");
        }
      }

      toast.success("Vault UI settings updated successfully");
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error("Failed to update dark mode logo", {
        description: message,
      });
    }

    await refetch();
  }

  useEffect(() => {
    if (getProjectUISettingsResponse) {
      form.reset({
        detectDarkModeEnabled:
          getProjectUISettingsResponse.projectUiSettings
            ?.detectDarkModeEnabled || false,
        logInLayout: (getProjectUISettingsResponse.projectUiSettings
          ?.logInLayout || "centered") as z.infer<typeof schema>["logInLayout"],
        primaryColor:
          getProjectUISettingsResponse.projectUiSettings?.primaryColor ||
          "#0f172a",
        darkModePrimaryColor:
          getProjectUISettingsResponse.projectUiSettings
            ?.darkModePrimaryColor || "#ffffff",
        logo: getProjectUISettingsResponse.projectUiSettings?.logoUrl || "",
        darkModeLogo:
          getProjectUISettingsResponse.projectUiSettings?.darkModeLogoUrl || "",
      });

      // Set logos independent from form reset, since they are not part of the form submission payload
      if (getProjectUISettingsResponse.projectUiSettings?.logoUrl) {
        setLogo(getProjectUISettingsResponse.projectUiSettings.logoUrl);
      }
      if (getProjectUISettingsResponse.projectUiSettings?.darkModeLogoUrl) {
        setDarkModeLogo(
          getProjectUISettingsResponse.projectUiSettings.darkModeLogoUrl,
        );
      }
    }
  }, [getProjectUISettingsResponse, form]);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-5 xl:grid-cols-7 2xl:grid-cols-3 gap-8 lg:items-stretch">
      <div className="col-span-1 lg:col-span-3 xl:col-span-5 2xl:col-span-2 order-2 lg:order-1">
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)}>
            <Card>
              <CardHeader>
                <CardTitle>Vault Branding Settings</CardTitle>
                <CardDescription>
                  Control the look and feel of your Vault pages.
                </CardDescription>
                <CardAction>
                  <Button
                    type="submit"
                    size="sm"
                    disabled={
                      !form.formState.isDirty ||
                      logo !==
                        getProjectUISettingsResponse?.projectUiSettings
                          ?.logoUrl ||
                      darkModeLogo !==
                        getProjectUISettingsResponse?.projectUiSettings
                          ?.darkModeLogoUrl
                    }
                  >
                    Save Changes
                  </Button>
                </CardAction>
              </CardHeader>
              <CardContent className="space-y-6">
                <FormField
                  control={form.control}
                  name="detectDarkModeEnabled"
                  render={({ field }) => (
                    <FormItem className="flex items-center justify-between gap-4">
                      <div className="space-y-2">
                        <FormLabel>Auto-detect Dark Mode</FormLabel>
                        <FormDescription>
                          Automatically switch to dark mode based on user system
                          preferences.
                        </FormDescription>
                        <FormMessage />
                      </div>
                      <FormControl>
                        <Switch
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="logInLayout"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Layout</FormLabel>
                      <FormDescription>
                        Choose the layout for your Vault Login and Signup pages.
                      </FormDescription>
                      <FormMessage />
                      <FormControl>
                        <div className="grid grid-cols-2 gap-4 mt-4">
                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <div
                                  className={cn(
                                    "p-4 border rounded-sm relative",
                                    field.value === "centered"
                                      ? "border-primary border-2 cursor-default"
                                      : "cursor-pointer",
                                  )}
                                  onClick={() => field.onChange("centered")}
                                >
                                  <div
                                    className={cn(
                                      "font-semibold text-sm",
                                      field.value === "centered"
                                        ? "text-primary"
                                        : "text-muted-foreground",
                                    )}
                                  >
                                    <SquareSquare
                                      className="inline-block mr-2"
                                      size={16}
                                    />
                                    Center card
                                    {field.value === "centered" && (
                                      <div className="h-5 w-5 text-white bg-primary rounded-full flex justify-center items-center absolute top-2 right-2">
                                        <Check size={12} />
                                      </div>
                                    )}
                                  </div>
                                  <p className="text-xs text-muted-foreground mt-2">
                                    A center-aligned card layout.
                                  </p>
                                </div>
                              </TooltipTrigger>
                              <TooltipContent className="bg-primary">
                                <div className="rounded">
                                  <img
                                    src="/images/auth-preview-centered.png"
                                    alt="Auth Preview - Centered"
                                    className="rounded max-w-[300px]"
                                  />
                                </div>
                              </TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <div
                                  className={cn(
                                    "p-4 border rounded-sm relative",
                                    field.value === "side_by_side"
                                      ? "border-primary border-2 cursor-default"
                                      : "cursor-pointer",
                                  )}
                                  onClick={() => field.onChange("side_by_side")}
                                >
                                  <div
                                    className={cn(
                                      "font-semibold text-sm",
                                      field.value === "side_by_side"
                                        ? "text-primary"
                                        : "text-muted-foreground",
                                    )}
                                  >
                                    <SquareSplitHorizontal
                                      className="inline-block mr-2"
                                      size={16}
                                    />
                                    Side by side
                                    {field.value === "side_by_side" && (
                                      <div className="h-5 w-5 text-white bg-primary rounded-full flex justify-center items-center absolute top-2 right-2">
                                        <Check size={12} />
                                      </div>
                                    )}
                                  </div>
                                  <p className="text-xs text-muted-foreground mt-2">
                                    A horizontally split layout.
                                  </p>
                                </div>
                              </TooltipTrigger>
                              <TooltipContent className="bg-primary">
                                <div className="rounded">
                                  <img
                                    src="/images/auth-preview-side-by-side.png"
                                    alt="Auth Preview - Side by Side"
                                    className="rounded max-w-[300px]"
                                  />
                                </div>
                              </TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                        </div>
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="primaryColor"
                  render={({ field }) => (
                    <FormItem className="flex items-center justify-between gap-4">
                      <div className="space-y-2">
                        <FormLabel>Primary Color</FormLabel>
                        <FormDescription>
                          The primary color for your Vault pages.
                        </FormDescription>
                        <FormMessage />
                      </div>
                      <FormControl>
                        <ColorPicker
                          className="mt-2"
                          onBlur={() => {
                            if (
                              !hexRegexp.test(field.value) ||
                              field.value.length < 7
                            ) {
                              field.onChange(primaryColor);
                            }
                          }}
                          onChange={field.onChange}
                          value={field.value}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="logo"
                  render={() => (
                    <FormItem className="flex items-center justify-between gap-4">
                      <div className="space-y-2">
                        <FormLabel>Logo</FormLabel>
                        <FormDescription>
                          Upload a logo for your Vault pages.
                        </FormDescription>
                        <FormMessage />

                        {logo && (
                          <div className="p-2 rounded-md bg-white border">
                            <img className="max-h-10 max-w-64" src={logo} />
                          </div>
                        )}
                      </div>

                      <FormControl>
                        <>
                          <Input
                            className="max-w-sm"
                            type="file"
                            accept="image/*"
                            ref={logoPickerRef}
                            onChange={handleLogoChange}
                          />
                        </>
                      </FormControl>
                    </FormItem>
                  )}
                />
                {detectDarkModeEnabled && (
                  <>
                    <div className="mt-8 mb-6 font-semibold">
                      Dark Mode Settings
                    </div>
                    <FormField
                      control={form.control}
                      name="darkModePrimaryColor"
                      render={({ field }) => (
                        <FormItem className="flex items-center justify-between gap-4">
                          <div className="spacey-y-2">
                            <FormLabel>Dark Mode Primary Color</FormLabel>
                            <FormDescription>
                              The primary color for your Vault pages in dark
                              mode.
                            </FormDescription>
                            <FormMessage />
                          </div>
                          <FormControl>
                            <ColorPicker
                              className="mt-2"
                              onBlur={() => {
                                if (
                                  !hexRegexp.test(field.value) ||
                                  field.value.length < 7
                                ) {
                                  field.onChange(darkModePrimaryColor);
                                }
                              }}
                              onChange={field.onChange}
                              value={field.value}
                            />
                          </FormControl>
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name="darkModeLogo"
                      render={() => (
                        <FormItem className="flex items-center justify-between gap-4">
                          <div className="space-y-2">
                            <FormLabel>Dark Mode Logo</FormLabel>
                            <FormDescription>
                              Upload a logo for your Vault pages in dark mode.
                            </FormDescription>
                            <FormMessage />
                            {darkModeLogo && (
                              <div className="p-2 rounded-md bg-primary border">
                                <img
                                  className="max-h-10 max-w-64"
                                  src={darkModeLogo}
                                />
                              </div>
                            )}
                          </div>
                          <FormControl>
                            <>
                              <Input
                                className="max-w-sm"
                                type="file"
                                accept="image/*"
                                ref={darkModeLogoPickerRef}
                                onChange={handleDarkModeLogoChange}
                              />
                            </>
                          </FormControl>
                        </FormItem>
                      )}
                    />
                  </>
                )}
              </CardContent>
            </Card>
          </form>
        </Form>
      </div>
      <div className="col-span-1 lg:col-span-2 2xl:col-span-1 order-1 lg:order-2">
        <Card className="bg-muted">
          <CardHeader>
            <CardTitle>Preview</CardTitle>
            <CardDescription>
              Here you can preview how your Vault pages will look with the
              current settings.
            </CardDescription>
            <CardAction className="min-w-24">
              {detectDarkModeEnabled && (
                <div className="flex items-center gap-2">
                  <div
                    className={cn(
                      "font-semibold text-xs",
                      darkMode ? "text-foreground" : "text-muted-foreground/70",
                    )}
                  >
                    Dark Mode
                  </div>
                  <div className="relative text-right">
                    <>
                      <Switch
                        checked={darkMode}
                        onCheckedChange={setDarkMode}
                      />
                      {darkMode ? (
                        <div className="absolute top-1 right-4 text-white">
                          <Moon size={12} />
                        </div>
                      ) : (
                        <div className="absolute top-1 right-1 text-muted-foreground">
                          <Sun size={12} />
                        </div>
                      )}
                    </>
                  </div>
                </div>
              )}
            </CardAction>
          </CardHeader>
          <CardContent>
            <div
              key={`${darkMode ? "dark" : "light"}-${primaryColor}-${darkModePrimaryColor}`}
              ref={previewRef}
              className={cn(
                "bg-background rounded-md py-8 px-6 border border-border/50 shadow-sm",
                darkMode ? "dark" : "",
              )}
              style={
                {
                  "--primary": darkMode ? darkModePrimaryColor : primaryColor,
                  "--primary-foreground": isColorDark(
                    darkMode ? darkModePrimaryColor : primaryColor,
                  )
                    ? "#ffffff"
                    : "#000000",
                } as React.CSSProperties
              }
            >
              <AuthCardPreview logo={darkMode ? darkModeLogo : logo} />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
