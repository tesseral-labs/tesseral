-- Ensures that the organization_id column in the projects table is nullable, 
-- so we can create projects before the managing Organization has been created.

alter table projects
alter column organization_id drop not null;