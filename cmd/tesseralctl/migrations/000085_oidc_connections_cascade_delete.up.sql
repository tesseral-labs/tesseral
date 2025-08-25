ALTER TABLE oidc_connections DROP CONSTRAINT oidc_connections_organization_id_fkey;
ALTER TABLE oidc_connections ADD CONSTRAINT oidc_connections_organization_id_fkey
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

ALTER TABLE oidc_intermediate_sessions DROP CONSTRAINT oidc_intermediate_sessions_oidc_connection_id_fkey;
ALTER TABLE oidc_intermediate_sessions ADD CONSTRAINT oidc_intermediate_sessions_oidc_connection_id_fkey
    FOREIGN KEY (oidc_connection_id) REFERENCES oidc_connections(id) ON DELETE CASCADE;
