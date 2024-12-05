alter table projects
    drop column google_oauth_client_secret,
    drop column microsoft_oauth_client_secret,
    add column google_oauth_client_secret_ciphertext    varchar,
    add column microsoft_oauth_client_secret_ciphertext varchar;
