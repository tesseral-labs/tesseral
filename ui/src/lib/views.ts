import { createContext } from 'react'

export enum LoginLayouts {
  Centered = 'centered',
  SideBySide = 'side_by_side',
}

export enum LoginViews {
  ChooseAdditionalFactor = 'choose-additional-factor',
  ChooseOrganization = 'choose-organization',
  ChooseOrganizationPrimaryLoginFactor = 'choose-organization-primary-login-factor',
  CreateOrganization = 'create-organization',
  Login = 'login',
  RegisterPassword = 'register-password',
  RegisterAuthenticatorApp = 'register-totp',
  RegisterPasskey = 'register-webauthn',
  VerifyEmail = 'verify-email',
  VerifyPassword = 'verify-password',
  VerifyAuthenticatorApp = 'verify-totp',
  VerifyPasskey = 'verify-webauthn',
}
