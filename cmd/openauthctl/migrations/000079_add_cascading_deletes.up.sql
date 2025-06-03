-- organizations.project_id → projects.id
ALTER TABLE organizations DROP CONSTRAINT organizations_project_id_fkey;
ALTER TABLE organizations ADD CONSTRAINT organizations_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- users.organization_id → organizations.id
ALTER TABLE users DROP CONSTRAINT users_organization_id_fkey;
ALTER TABLE users ADD CONSTRAINT users_organization_id_fkey
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- projects.organization_id → organizations.id
ALTER TABLE projects DROP CONSTRAINT projects_organization_id_fkey;
ALTER TABLE projects ADD CONSTRAINT projects_organization_id_fkey
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- intermediate_sessions.project_id → projects.id
ALTER TABLE intermediate_sessions DROP CONSTRAINT intermediate_sessions_project_id_fkey;
ALTER TABLE intermediate_sessions ADD CONSTRAINT intermediate_sessions_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- sessions.user_id → users.id
ALTER TABLE sessions DROP CONSTRAINT sessions_user_id_fkey;
ALTER TABLE sessions ADD CONSTRAINT sessions_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- session_signing_keys.project_id → projects.id
ALTER TABLE session_signing_keys DROP CONSTRAINT session_signing_keys_project_id_fkey;
ALTER TABLE session_signing_keys ADD CONSTRAINT session_signing_keys_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- backend_api_keys.project_id → projects.id
ALTER TABLE backend_api_keys DROP CONSTRAINT project_api_keys_project_id_fkey;
ALTER TABLE backend_api_keys ADD CONSTRAINT project_api_keys_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- oauth_verified_emails.project_id → projects.id
ALTER TABLE oauth_verified_emails DROP CONSTRAINT verified_emails_project_id_fkey;
ALTER TABLE oauth_verified_emails ADD CONSTRAINT verified_emails_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- saml_connections.organization_id → organizations.id
ALTER TABLE saml_connections DROP CONSTRAINT saml_connections_organization_id_fkey;
ALTER TABLE saml_connections ADD CONSTRAINT saml_connections_organization_id_fkey
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- organization_domains.organization_id → organizations.id
ALTER TABLE organization_domains DROP CONSTRAINT organization_domains_organization_id_fkey;
ALTER TABLE organization_domains ADD CONSTRAINT organization_domains_organization_id_fkey
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- scim_api_keys.organization_id → organizations.id
ALTER TABLE scim_api_keys DROP CONSTRAINT scim_api_keys_organization_id_fkey;
ALTER TABLE scim_api_keys ADD CONSTRAINT scim_api_keys_organization_id_fkey
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- intermediate_sessions.organization_id → organizations.id
ALTER TABLE intermediate_sessions DROP CONSTRAINT intermediate_sessions_organization_id_fkey;
ALTER TABLE intermediate_sessions ADD CONSTRAINT intermediate_sessions_organization_id_fkey
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- project_ui_settings.project_id → projects.id
ALTER TABLE project_ui_settings DROP CONSTRAINT project_ui_settings_project_id_fkey;
ALTER TABLE project_ui_settings ADD CONSTRAINT project_ui_settings_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- user_impersonation_tokens.impersonator_id → users.id
ALTER TABLE user_impersonation_tokens DROP CONSTRAINT user_impersonation_tokens_impersonator_id_fkey;
ALTER TABLE user_impersonation_tokens ADD CONSTRAINT user_impersonation_tokens_impersonator_id_fkey
FOREIGN KEY (impersonator_id) REFERENCES users(id) ON DELETE CASCADE;

-- user_impersonation_tokens.impersonated_id → users.id
ALTER TABLE user_impersonation_tokens DROP CONSTRAINT user_impersonation_tokens_impersonated_id_fkey;
ALTER TABLE user_impersonation_tokens ADD CONSTRAINT user_impersonation_tokens_impersonated_id_fkey
FOREIGN KEY (impersonated_id) REFERENCES users(id) ON DELETE CASCADE;

-- sessions.impersonator_user_id → users.id
ALTER TABLE sessions DROP CONSTRAINT sessions_impersonator_user_id_fkey;
ALTER TABLE sessions ADD CONSTRAINT sessions_impersonator_user_id_fkey
FOREIGN KEY (impersonator_user_id) REFERENCES users(id) ON DELETE CASCADE;

-- passkeys.user_id → users.id
ALTER TABLE passkeys DROP CONSTRAINT passkeys_user_id_fkey;
ALTER TABLE passkeys ADD CONSTRAINT passkeys_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- user_invites.organization_id → organizations.id
ALTER TABLE user_invites DROP CONSTRAINT user_invites_organization_id_fkey;
ALTER TABLE user_invites ADD CONSTRAINT user_invites_organization_id_fkey
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- project_trusted_domains.project_id → projects.id
ALTER TABLE project_trusted_domains DROP CONSTRAINT project_trusted_domains_project_id_fkey;
ALTER TABLE project_trusted_domains ADD CONSTRAINT project_trusted_domains_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- publishable_keys.project_id → projects.id
ALTER TABLE publishable_keys DROP CONSTRAINT publishable_keys_project_id_fkey;
ALTER TABLE publishable_keys ADD CONSTRAINT publishable_keys_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- vault_domain_settings.project_id → projects.id
ALTER TABLE vault_domain_settings DROP CONSTRAINT vault_domain_settings_project_id_fkey;
ALTER TABLE vault_domain_settings ADD CONSTRAINT vault_domain_settings_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- project_email_quota_daily_usage.project_id → projects.id
ALTER TABLE project_email_quota_daily_usage DROP CONSTRAINT project_email_quota_daily_usage_project_id_fkey;
ALTER TABLE project_email_quota_daily_usage ADD CONSTRAINT project_email_quota_daily_usage_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- roles.project_id → projects.id
ALTER TABLE roles DROP CONSTRAINT roles_project_id_fkey;
ALTER TABLE roles ADD CONSTRAINT roles_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- roles.organization_id → organizations.id
ALTER TABLE roles DROP CONSTRAINT roles_organization_id_fkey;
ALTER TABLE roles ADD CONSTRAINT roles_organization_id_fkey
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- actions.project_id → projects.id
ALTER TABLE actions DROP CONSTRAINT actions_project_id_fkey;
ALTER TABLE actions ADD CONSTRAINT actions_project_id_fkey
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- user_invites.role_id → roles.id
ALTER TABLE user_invites DROP CONSTRAINT user_invites_role_id_fkey;
ALTER TABLE user_invites ADD CONSTRAINT user_invites_role_id_fkey
FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;
