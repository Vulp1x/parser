--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Homebrew)
-- Dumped by pg_dump version 14.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_table_access_method = heap;

--
-- Name: bloggers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bloggers (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    dataset_id uuid NOT NULL,
    username text NOT NULL,
    user_id bigint NOT NULL,
    followers_count bigint DEFAULT '-1'::integer NOT NULL,
    is_initial boolean NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    parsed_at timestamp without time zone,
    updated_at timestamp without time zone,
    parsed boolean DEFAULT false NOT NULL,
    is_private boolean DEFAULT false NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    is_business boolean DEFAULT false NOT NULL,
    followings_count integer DEFAULT '-1'::integer NOT NULL,
    contact_phone_number text,
    public_phone_number text,
    public_phone_country_code text,
    city_name text,
    public_email text
);


--
-- Name: bots; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bots (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL,
    session_id text NOT NULL,
    proxy jsonb NOT NULL,
    is_blocked boolean NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    locked_until timestamp without time zone
);


--
-- Name: datasets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.datasets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    phone_code integer,
    status smallint NOT NULL,
    title text NOT NULL,
    manager_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    started_at timestamp with time zone,
    stopped_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    posts_per_blogger integer DEFAULT 0 NOT NULL,
    liked_per_post integer DEFAULT 0 NOT NULL,
    commented_per_post integer DEFAULT 0 NOT NULL
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: targets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.targets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    dataset_id uuid NOT NULL,
    username text NOT NULL,
    user_id bigint NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    is_private boolean DEFAULT false NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    is_business boolean DEFAULT false NOT NULL,
    followers_count integer DEFAULT '-1'::integer NOT NULL,
    followings_count integer DEFAULT '-1'::integer NOT NULL,
    contact_phone_number text,
    public_phone_number text,
    public_phone_country_code text,
    city_name text,
    public_email text
);


--
-- Name: bloggers bloggers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bloggers
    ADD CONSTRAINT bloggers_pkey PRIMARY KEY (id);


--
-- Name: bots bots_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bots
    ADD CONSTRAINT bots_pkey PRIMARY KEY (id);


--
-- Name: bots bots_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bots
    ADD CONSTRAINT bots_username_key UNIQUE (username);


--
-- Name: datasets datasets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.datasets
    ADD CONSTRAINT datasets_pkey PRIMARY KEY (id);


--
-- Name: targets targets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.targets
    ADD CONSTRAINT targets_pkey PRIMARY KEY (id);


--
-- Name: bots_session_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX bots_session_id_idx ON public.bots USING btree (session_id);


--
-- Name: target_users_uniq_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX target_users_uniq_idx ON public.targets USING btree (dataset_id, user_id);


--
-- Name: bloggers bloggers_dataset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bloggers
    ADD CONSTRAINT bloggers_dataset_id_fkey FOREIGN KEY (dataset_id) REFERENCES public.datasets(id);


--
-- Name: targets targets_dataset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.targets
    ADD CONSTRAINT targets_dataset_id_fkey FOREIGN KEY (dataset_id) REFERENCES public.datasets(id);


--
-- PostgreSQL database dump complete
--

