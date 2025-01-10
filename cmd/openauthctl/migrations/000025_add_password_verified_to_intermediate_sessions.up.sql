alter table intermediate_sessions 
add column password_verified boolean,
add column organization_id uuid references organizations(id);