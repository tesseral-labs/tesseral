CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create the Dogfood Project
INSERT INTO projects (id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, auth_domain, custom_auth_domain)
	VALUES ('56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid, 'Tesseral Local Development', true, true, true, true, true, true, true, 'auth.console.tesseral.example.com', 'auth.console.tesseral.example.com');

-- Create the Dogfood Project's backing organization
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
  VALUES ('7a76decb-6d79-49ce-9449-34fcc53151df'::uuid, 'project_54vwf0clhh0caqe20eujxgpeq Backing Organization', '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2', true, false, true, true, true, true, true, true);

UPDATE projects SET organization_id = '7a76decb-6d79-49ce-9449-34fcc53151df'::uuid where id = '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid;

-- have dogfood project support passkeys from auth.console.tesseral.example.com (vault) and console.tesseral.example.com (self-built login flow)
insert into
    project_passkey_rp_ids (project_id, rp_id)
values
    ('56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2', 'auth.console.tesseral.example.com'),
    ('56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2', 'console.tesseral.example.com');

-- Create a user in the dogfood project
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
  VALUES (gen_random_uuid(), 'root@app.tesseral.example.com', crypt('testpassword', gen_salt('bf', 14)), '7a76decb-6d79-49ce-9449-34fcc53151df', true);

-- Create project UI settings for the dogfood project
INSERT INTO project_ui_settings (id, project_id)
  VALUES (gen_random_uuid(), '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid);

-- Create a project API key in the dogfood project
INSERT INTO project_api_keys (id, project_id, secret_token_sha256, display_name)
  VALUES (gen_random_uuid(), (SELECT id FROM projects LIMIT 1), digest(uuid_send('F938657E-65FC-4C43-B2F1-CE875A0B64D6'::uuid), 'sha256'), 'localhost');

-- Create a session signing key in the dogfood project
INSERT INTO session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time) 
  VALUES (
    gen_random_uuid(), 
    '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2'::uuid, 
    decode('3059301306072a8648ce3d020106082a8648ce3d03010703420004a82072a20d2217055f0c5f9f9283e128d5bc26334b19024c93f6ad50619bbe83bc565a2fbdc05e02dc3f1452ff273d7ec2534e2cbe7fe395443d887b128dd7b8', 'hex'), 
    decode('a1931242e0770f54e2e8365053ff4b72dc72faba0830cff2099655d78aa188f750b9b1557e70566f00449fed97a5b8a94a113e8049a6ea71436a08e135f35a7b86863f47f36e3e0b62dad8da491f28aba812a93e7a2a44913c6b2377c7ea4d89991eba682d9cfb17d5bcfa3f608e973dd61aa9910453e8d48058ea80ccbd0d5961de3fd25dcfe893dbdd84a43112d1533b4ebae65e35b0e8eca25b1af53eec97304899cb542ac850e59a6c5521ecbee5549329a451c8c948d82f1d6858a6d2680d987e72945ad5b4166c3529b70ce1106573874fb68847ed823567a9edfeac712d464ac5b339f80365be985ab69703d7100c65c872765b04a9ee575002edadef', 'hex'), 
    (SELECT NOW() + INTERVAL '1 year')
  );

-- Create customer1's project
insert into projects (id, log_in_with_password, log_in_with_google, log_in_with_microsoft, display_name, custom_auth_domain, auth_domain)
    values ('7abd6d2e-c314-456e-b9c5-bdbb62f0345f'::uuid, true, false, false, 'Customer One', 'auth.customer1.example.com', 'auth.customer1.example.com');

-- Create customer1's project's backing organization
INSERT INTO organizations (id, display_name, project_id, log_in_with_saml, scim_enabled, log_in_with_password)
VALUES ('8648d50b-baa1-4929-be0f-bc7238f685ab'::uuid, 'project_79ldwwwzybn66dxa91udi7mn3 Backing Organization', '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2', false, false, true);

update projects set organization_id = '8648d50b-baa1-4929-be0f-bc7238f685ab'::uuid where id = '7abd6d2e-c314-456e-b9c5-bdbb62f0345f'::uuid;

insert into project_passkey_rp_ids (project_id, rp_id) values ('7abd6d2e-c314-456e-b9c5-bdbb62f0345f', 'auth.customer1.example.com');

-- Create a user in customer1
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
VALUES (gen_random_uuid(), 'user1@company1.example.com', crypt('testpassword', gen_salt('bf', 14)), '8648d50b-baa1-4929-be0f-bc7238f685ab', true);

-- create customer1's session signing keys
insert into session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time)
values (
           gen_random_uuid(),
           '7abd6d2e-c314-456e-b9c5-bdbb62f0345f'::uuid,
           decode('3059301306072a8648ce3d020106082a8648ce3d0301070342000473bbd17732bc07085a24ad9385edb16eb6e882deb60efb140dc32790f0a37f8dfd9631f2f60f345c84611ecf1a055748c4b786d84e28f1b91a4b1dfe34742aec', 'hex'),
           decode('16b4dfd43beccde193bba4a02392fc2ac18ad45521caf94e55ee61e3957ba5d4e060c9cf2493597b2aa5d61642007f6d190b64fc3cfef43ec7aa8e3735276912424a6e6795a53a3516e2527f16a938f733346ab96db4aa1f8312026c666e5cb34e80803a09cee1ed52da411451b6d19230105ae0ef6bc9c2cc8ed02c30ae3d59abf67e4b33949353ceb35572dde287d4a197b63c69dbce3cb19177111fccb4e36de68fb1b9f4c60dfe9661026bca72c932f47e05b2dff6767eb3a38d62398ae62d56432e1079f621adc819ee5d93c526ce6ff1484c288103f6a2136c8892a43b33f2c4b386d17a1da81cb6f0d2476867d9d7829818ef0535afc4910eb53541f3', 'hex'),
           (select now() + interval '1 year')
       );

-- Create customer2's project
insert into projects (id, log_in_with_password, log_in_with_google, log_in_with_microsoft, display_name, custom_auth_domain, auth_domain)
values ('24ba0dd5-e178-460e-8f7a-f3f72cf6a1e7'::uuid, true, false, false, 'Customer Two', 'auth.customer2.example.com', 'auth.customer2.example.com');

-- Create customer1's project's backing organization
INSERT INTO organizations (id, display_name, project_id, log_in_with_saml, scim_enabled, log_in_with_password)
VALUES ('8b5972b6-c878-4c6c-a351-9e01da20f776'::uuid, 'project_269wse1l6u0jnvs8afpq44f6v Backing Organization', '56bfa2b3-4f5a-4c68-8fc5-db3bf20731a2', false, false, true);

update projects set organization_id = '8b5972b6-c878-4c6c-a351-9e01da20f776'::uuid where id = '24ba0dd5-e178-460e-8f7a-f3f72cf6a1e7'::uuid;

insert into project_passkey_rp_ids (project_id, rp_id) values ('24ba0dd5-e178-460e-8f7a-f3f72cf6a1e7', 'auth.customer2.example.com');

-- Create a user in customer2
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
VALUES (gen_random_uuid(), 'user1@company2.example.com', crypt('testpassword', gen_salt('bf', 14)), '8b5972b6-c878-4c6c-a351-9e01da20f776', true);

-- create customer2's session signing keys
insert into session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time)
values (
           gen_random_uuid(),
           '24ba0dd5-e178-460e-8f7a-f3f72cf6a1e7'::uuid,
           decode('3059301306072a8648ce3d020106082a8648ce3d0301070342000473bbd17732bc07085a24ad9385edb16eb6e882deb60efb140dc32790f0a37f8dfd9631f2f60f345c84611ecf1a055748c4b786d84e28f1b91a4b1dfe34742aec', 'hex'),
           decode('16b4dfd43beccde193bba4a02392fc2ac18ad45521caf94e55ee61e3957ba5d4e060c9cf2493597b2aa5d61642007f6d190b64fc3cfef43ec7aa8e3735276912424a6e6795a53a3516e2527f16a938f733346ab96db4aa1f8312026c666e5cb34e80803a09cee1ed52da411451b6d19230105ae0ef6bc9c2cc8ed02c30ae3d59abf67e4b33949353ceb35572dde287d4a197b63c69dbce3cb19177111fccb4e36de68fb1b9f4c60dfe9661026bca72c932f47e05b2dff6767eb3a38d62398ae62d56432e1079f621adc819ee5d93c526ce6ff1484c288103f6a2136c8892a43b33f2c4b386d17a1da81cb6f0d2476867d9d7829818ef0535afc4910eb53541f3', 'hex'),
           (select now() + interval '1 year')
       );
