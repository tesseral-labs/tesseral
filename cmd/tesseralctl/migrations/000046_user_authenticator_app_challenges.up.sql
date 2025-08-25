create table user_authenticator_app_challenges
(
    user_id                                 uuid not null references users (id) on delete cascade,
    authenticator_app_secret_ciphertext     bytea
);
