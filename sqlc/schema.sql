--
-- PostgreSQL database dump
--

-- Dumped from database version 15.8 (Debian 15.8-1.pgdg120+1)
-- Dumped by pg_dump version 17.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: auth_method; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.auth_method AS ENUM (
    'email',
    'google',
    'microsoft'
);


ALTER TYPE public.auth_method OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: intermediate_sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.intermediate_sessions (
    id uuid NOT NULL,
    project_id uuid NOT NULL,
    unverified_email character varying,
    verified_email character varying,
    created_time timestamp with time zone DEFAULT now() NOT NULL,
    expire_time timestamp with time zone NOT NULL,
    token character varying NOT NULL,
    token_sha256 bytea,
    revoked boolean DEFAULT false NOT NULL
);


ALTER TABLE public.intermediate_sessions OWNER TO postgres;

--
-- Name: method_verification_challenges; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.method_verification_challenges (
    id uuid NOT NULL,
    project_id uuid NOT NULL,
    complete_time timestamp with time zone,
    email character varying NOT NULL,
    auth_method public.auth_method NOT NULL,
    expire_time timestamp with time zone NOT NULL,
    secret_token_sha256 bytea NOT NULL
);


ALTER TABLE public.method_verification_challenges OWNER TO postgres;

--
-- Name: organizations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.organizations (
    id uuid NOT NULL,
    project_id uuid NOT NULL,
    display_name character varying NOT NULL,
    override_log_in_with_password_enabled boolean,
    override_log_in_with_google_enabled boolean,
    override_log_in_with_microsoft_enabled boolean,
    google_hosted_domain character varying,
    microsoft_tenant_id character varying
);


ALTER TABLE public.organizations OWNER TO postgres;

--
-- Name: projects; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.projects (
    id uuid NOT NULL,
    organization_id uuid,
    log_in_with_password_enabled boolean NOT NULL,
    log_in_with_google_enabled boolean NOT NULL,
    log_in_with_microsoft_enabled boolean NOT NULL,
    google_oauth_client_id character varying,
    google_oauth_client_secret character varying,
    microsoft_oauth_client_id character varying,
    microsoft_oauth_client_secret character varying
);


ALTER TABLE public.projects OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_time timestamp with time zone DEFAULT now() NOT NULL,
    expire_time timestamp with time zone NOT NULL,
    token character varying NOT NULL,
    token_sha256 bytea,
    revoked boolean DEFAULT false NOT NULL
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    organization_id uuid NOT NULL,
    unverified_email character varying,
    verified_email character varying,
    password_bcrypt character varying,
    google_user_id character varying,
    microsoft_user_id character varying
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: intermediate_sessions intermediate_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intermediate_sessions
    ADD CONSTRAINT intermediate_sessions_pkey PRIMARY KEY (id);


--
-- Name: organizations organizations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_pkey PRIMARY KEY (id);


--
-- Name: organizations organizations_project_id_google_hosted_domain_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_project_id_google_hosted_domain_key UNIQUE (project_id, google_hosted_domain);


--
-- Name: organizations organizations_project_id_microsoft_tenant_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_project_id_microsoft_tenant_id_key UNIQUE (project_id, microsoft_tenant_id);


--
-- Name: projects projects_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT projects_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: users users_organization_id_google_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_organization_id_google_user_id_key UNIQUE (organization_id, google_user_id);


--
-- Name: users users_organization_id_microsoft_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_organization_id_microsoft_user_id_key UNIQUE (organization_id, microsoft_user_id);


--
-- Name: users users_organization_id_unverified_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_organization_id_unverified_email_key UNIQUE (organization_id, unverified_email);


--
-- Name: users users_organization_id_verified_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_organization_id_verified_email_key UNIQUE (organization_id, verified_email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: intermediate_sessions intermediate_sessions_project_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intermediate_sessions
    ADD CONSTRAINT intermediate_sessions_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.projects(id);


--
-- Name: method_verification_challenges method_verification_challenges_project_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.method_verification_challenges
    ADD CONSTRAINT method_verification_challenges_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.projects(id);


--
-- Name: organizations organizations_project_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.projects(id);


--
-- Name: projects projects_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT projects_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users users_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- PostgreSQL database dump complete
--

