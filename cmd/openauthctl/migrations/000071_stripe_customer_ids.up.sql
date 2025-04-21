alter table projects
    add column stripe_customer_id            varchar,
    add column entitled_custom_vault_domains bool not null default false,
    add column entitled_backend_api_keys     bool not null default false;

create index on projects (stripe_customer_id);
