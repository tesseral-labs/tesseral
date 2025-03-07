import { Organization } from '@/gen/tesseral/intermediate/v1/intermediate_pb';

export enum PrimaryLoginFactor {
  Email = 'email',
  GoogleOAuth = 'google_oauth',
  MicrosoftOAuth = 'microsoft_oauth',
}

const primaryLoginFactorToOrganizationSettingMap: Record<
  PrimaryLoginFactor,
  string
> = {
  [PrimaryLoginFactor.Email]: 'logInWithEmail',
  [PrimaryLoginFactor.GoogleOAuth]: 'logInWithGoogle',
  [PrimaryLoginFactor.MicrosoftOAuth]: 'logInWithMicrosoft',
};

export const isValidPrimaryLoginFactor = (
  primaryLoginFactor: PrimaryLoginFactor,
  organization: Organization,
) => {
  const organizationSetting =
    primaryLoginFactorToOrganizationSettingMap[primaryLoginFactor];

  console.log('organizationSetting', organizationSetting);

  return !!organization[organizationSetting as keyof Organization];
};
