export enum LoginPage {
  ChooseAdditionalFactor = 'choose_additional_factor',
  ChooseProject = 'choose_project',
  CreateProject = 'create_project',
  RegisterAuthenticatorApp = 'register_authenticator_app',
  RegisterPasskey = 'register_passkey',
  RegisterPassword = 'register_password',
  StartLogin = 'start_login',
  VerifyAuthenticatorApp = 'verify_authenticator_app',
  VerifyEmail = 'verify_email',
  VerifyPasskey = 'verify_passkey',
  VerifyPassword = 'verify_password',
}

export const LoginPageMap = {
  [LoginPage.ChooseAdditionalFactor]: '/choose-additional-factor',
  [LoginPage.ChooseProject]: '/choose-project',
  [LoginPage.CreateProject]: '/create-project',
  [LoginPage.RegisterAuthenticatorApp]: '/register-authenticator-app',
  [LoginPage.RegisterPasskey]: '/register-passkey',
  [LoginPage.RegisterPassword]: '/register-password',
  [LoginPage.StartLogin]: '/start-login',
  [LoginPage.VerifyAuthenticatorApp]: '/verify-authenticator-app',
  [LoginPage.VerifyEmail]: '/verify-email',
  [LoginPage.VerifyPasskey]: '/verify-passkey',
  [LoginPage.VerifyPassword]: '/verify-password',
};
