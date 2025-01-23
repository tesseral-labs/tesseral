CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create the Dogfood Project
INSERT INTO projects (id, display_name, log_in_with_google_enabled, log_in_with_microsoft_enabled, log_in_with_password_enabled, auth_domain, custom_auth_domain)
	VALUES ('56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid, 'localhost', true, true, true, 'auth.app.tesseral.com', 'auth.app.tesseral.com');

-- Create the Dogfood Organization
INSERT INTO organizations (id, display_name, project_id, override_log_in_methods, saml_enabled, scim_enabled)
  VALUES ('7a76decb-6d79-49ce-9449-34fcc53151df'::uuid, 'localhost', (SELECT id FROM projects LIMIT 1), false, true, false);

-- Update the Dogfood Project
UPDATE projects
  SET organization_id = '7a76decb-6d79-49ce-9449-34fcc53151df'::uuid;

-- Create the Dogfood User
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
  VALUES (gen_random_uuid(), 'root@app.tesseral.example.com', crypt('this_is_a_very_sensitive_password_' || encode(gen_random_bytes(16), 'hex'), gen_salt('bf', 14)), (SELECT id FROM organizations LIMIT 1), true);

-- Create a Dogfood Project API Key
INSERT INTO project_api_keys (id, project_id, secret_token_sha256, display_name)
  VALUES (gen_random_uuid(), (SELECT id FROM projects LIMIT 1), digest(uuid_send('F938657E-65FC-4C43-B2F1-CE875A0B64D6'::uuid), 'sha256'), 'localhost');

-- Create a Session Signing Key for the Dogfood Project
INSERT INTO session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time) 
  VALUES (gen_random_uuid(), '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid, pg_read_binary_file('/tmp/local/session-signing-key.encrypted'), pg_read_binary_file('/tmp/local/session-signing-public-key.pem'), (SELECT NOW() + INTERVAL '1 year'));