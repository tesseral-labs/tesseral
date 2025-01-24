CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create the Dogfood Project
INSERT INTO projects (id, display_name, log_in_with_google_enabled, log_in_with_microsoft_enabled, log_in_with_password_enabled, auth_domain, custom_auth_domain)
	VALUES ('56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid, 'Tesseral Local', true, true, true, 'auth.app.tesseral.example.com', 'auth.app.tesseral.example.com');

-- Create the Dogfood Organization
INSERT INTO organizations (id, display_name, project_id, override_log_in_methods, saml_enabled, scim_enabled)
  VALUES ('7a76decb-6d79-49ce-9449-34fcc53151df'::uuid, 'Tesseral Local', (SELECT id FROM projects LIMIT 1), false, true, false);

-- Update the Dogfood Project
UPDATE projects
  SET organization_id = '7a76decb-6d79-49ce-9449-34fcc53151df'::uuid;

-- Create the Dogfood User
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
  VALUES (gen_random_uuid(), 'root@app.tesseral.example.com', crypt('testpassword', gen_salt('bf', 14)), (SELECT id FROM organizations LIMIT 1), true);

-- Create Project UI Settings
INSERT INTO project_ui_settings (id, project_id)
  VALUES (gen_random_uuid(), '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid);

-- Create a Dogfood Project API Key
INSERT INTO project_api_keys (id, project_id, secret_token_sha256, display_name)
  VALUES (gen_random_uuid(), (SELECT id FROM projects LIMIT 1), digest(uuid_send('F938657E-65FC-4C43-B2F1-CE875A0B64D6'::uuid), 'sha256'), 'localhost');

-- Create a Session Signing Key for the Dogfood Project
INSERT INTO session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time) 
  VALUES (
    gen_random_uuid(), 
    '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid, 
    decode('3059301306072a8648ce3d020106082a8648ce3d03010703420004a82072a20d2217055f0c5f9f9283e128d5bc26334b19024c93f6ad50619bbe83bc565a2fbdc05e02dc3f1452ff273d7ec2534e2cbe7fe395443d887b128dd7b8', 'hex'), 
    decode('a1931242e0770f54e2e8365053ff4b72dc72faba0830cff2099655d78aa188f750b9b1557e70566f00449fed97a5b8a94a113e8049a6ea71436a08e135f35a7b86863f47f36e3e0b62dad8da491f28aba812a93e7a2a44913c6b2377c7ea4d89991eba682d9cfb17d5bcfa3f608e973dd61aa9910453e8d48058ea80ccbd0d5961de3fd25dcfe893dbdd84a43112d1533b4ebae65e35b0e8eca25b1af53eec97304899cb542ac850e59a6c5521ecbee5549329a451c8c948d82f1d6858a6d2680d987e72945ad5b4166c3529b70ce1106573874fb68847ed823567a9edfeac712d464ac5b339f80365be985ab69703d7100c65c872765b04a9ee575002edadef', 'hex'), 
    (SELECT NOW() + INTERVAL '1 year')
  );