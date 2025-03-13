import {
  Organization,
  PrimaryAuthFactor,
} from '@/gen/tesseral/intermediate/v1/intermediate_pb';

export const isValidPrimaryAuthFactor = (
  primaryAuthFactor: PrimaryAuthFactor,
  organization: Organization,
) => {
  switch (primaryAuthFactor) {
    case PrimaryAuthFactor.EMAIL:
      return organization.logInWithEmail === true;
    case PrimaryAuthFactor.GOOGLE:
      return organization.logInWithGoogle === true;
    case PrimaryAuthFactor.MICROSOFT:
      return organization.logInWithMicrosoft === true;
  }
};
