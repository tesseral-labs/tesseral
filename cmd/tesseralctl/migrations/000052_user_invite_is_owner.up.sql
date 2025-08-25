alter table user_invites
    add column is_owner boolean not null default false;
