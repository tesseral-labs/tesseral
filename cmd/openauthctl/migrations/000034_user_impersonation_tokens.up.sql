create table user_impersonation_tokens (
    id uuid not null primary key,
    impersonator_id uuid not null references users(id),
    impersonated_id uuid not null references users(id),

    create_time timestamp with time zone not null default now(),
    expire_time timestamp with time zone not null,
    redeem_time timestamp with time zone,

    secret_token_sha256 bytea
);
