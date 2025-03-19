alter table publishable_keys
    add column support_relayed_sessions boolean not null default false;
