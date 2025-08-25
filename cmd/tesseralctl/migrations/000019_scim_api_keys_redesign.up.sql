alter table scim_api_keys
    drop column create_time,
    drop column revoke_time,
    add column display_name varchar not null;
