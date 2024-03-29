--
-- PostgreSQL database dump
--

-- Dumped from database version 14.6 (Homebrew)
-- Dumped by pg_dump version 15.1

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
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: blogger_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.blogger_status AS ENUM (
    'new',
    'info_saved',
    'medias_found',
    'all_medias_parsed',
    'done',
    'invalid'
);


--
-- Name: cursor_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.cursor_type AS ENUM (
    'followers'
);


--
-- Name: dataset_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.dataset_type AS ENUM (
    'followers',
    'phone_numbers',
    'likes_and_comments'
);


--
-- Name: pgqueue_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.pgqueue_status AS ENUM (
    'new',
    'must_retry',
    'no_attempts_left',
    'cancelled',
    'succeeded'
);


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
    is_correct boolean DEFAULT false NOT NULL,
    status public.blogger_status DEFAULT 'new'::public.blogger_status NOT NULL,
    is_private boolean DEFAULT false NOT NULL,
    is_verified boolean DEFAULT false NOT NULL
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
    commented_per_post integer DEFAULT 0 NOT NULL,
    type public.dataset_type DEFAULT 'likes_and_comments'::public.dataset_type NOT NULL,
    followers_count integer DEFAULT 0 NOT NULL
);


--
-- Name: fetch_cursors; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.fetch_cursors (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    inst_pk bigint NOT NULL,
    dataset_id uuid NOT NULL,
    locked_until timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);


--
-- Name: full_targets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_targets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    dataset_id uuid NOT NULL,
    parsed_at timestamp without time zone DEFAULT now() NOT NULL,
    username text NOT NULL,
    inst_pk bigint NOT NULL,
    full_name text NOT NULL,
    is_private boolean NOT NULL,
    is_verified boolean NOT NULL,
    is_business boolean NOT NULL,
    is_potential_business boolean NOT NULL,
    has_anonymous_profile_picture boolean NOT NULL,
    biography text NOT NULL,
    external_url text NOT NULL,
    media_count integer NOT NULL,
    follower_count integer NOT NULL,
    following_count integer NOT NULL,
    category text NOT NULL,
    city_name text NOT NULL,
    contact_phone_number text NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    public_email text NOT NULL,
    public_phone_country_code text NOT NULL,
    public_phone_number text NOT NULL,
    whatsapp_number text NOT NULL,
    bio_links text DEFAULT ''::text NOT NULL
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
-- Name: medias; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.medias (
    pk bigint NOT NULL,
    id text NOT NULL,
    dataset_id uuid NOT NULL,
    media_type integer NOT NULL,
    code text NOT NULL,
    has_more_comments boolean NOT NULL,
    caption text NOT NULL,
    width integer NOT NULL,
    height integer NOT NULL,
    like_count integer NOT NULL,
    taken_at integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: pgqueue; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pgqueue (
    id bigint NOT NULL,
    kind smallint NOT NULL,
    payload bytea NOT NULL,
    external_key text,
    status public.pgqueue_status DEFAULT 'new'::public.pgqueue_status NOT NULL,
    messages text[] DEFAULT ARRAY[]::text[] NOT NULL,
    attempts_left smallint NOT NULL,
    attempts_elapsed smallint DEFAULT 0 NOT NULL,
    delayed_till timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
)
WITH (fillfactor='80');


--
-- Name: pgqueue_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.pgqueue_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pgqueue_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.pgqueue_id_seq OWNED BY public.pgqueue.id;


--
-- Name: targets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.targets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL,
    user_id bigint NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    is_private boolean DEFAULT false NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    full_name text DEFAULT ''::text NOT NULL,
    media_pk bigint NOT NULL,
    dataset_id uuid NOT NULL
);


--
-- Name: pgqueue id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pgqueue ALTER COLUMN id SET DEFAULT nextval('public.pgqueue_id_seq'::regclass);


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
-- Name: fetch_cursors fetch_cursors_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fetch_cursors
    ADD CONSTRAINT fetch_cursors_pkey PRIMARY KEY (id);


--
-- Name: full_targets full_targets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_targets
    ADD CONSTRAINT full_targets_pkey PRIMARY KEY (id);


--
-- Name: medias medias_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.medias
    ADD CONSTRAINT medias_pkey PRIMARY KEY (pk, dataset_id);


--
-- Name: pgqueue pgqueue_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pgqueue
    ADD CONSTRAINT pgqueue_pkey PRIMARY KEY (id);


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
-- Name: fetch_cursors_inst_pk_dataset_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX fetch_cursors_inst_pk_dataset_id_idx ON public.fetch_cursors USING btree (inst_pk, dataset_id);


--
-- Name: full_targets_uniq_user_per_dataset; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX full_targets_uniq_user_per_dataset ON public.full_targets USING btree (inst_pk, dataset_id);


--
-- Name: pgqueue_broken_tasks_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pgqueue_broken_tasks_idx ON public.pgqueue USING btree (kind, created_at) WHERE (status = 'no_attempts_left'::public.pgqueue_status);


--
-- Name: pgqueue_idempotency_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX pgqueue_idempotency_idx ON public.pgqueue USING btree (kind, external_key);


--
-- Name: pgqueue_open_tasks_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pgqueue_open_tasks_idx ON public.pgqueue USING btree (kind, delayed_till) WHERE (status = ANY (ARRAY['new'::public.pgqueue_status, 'must_retry'::public.pgqueue_status]));


--
-- Name: pgqueue_terminal_tasks_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pgqueue_terminal_tasks_idx ON public.pgqueue USING btree (kind, updated_at) WHERE (status = ANY (ARRAY['cancelled'::public.pgqueue_status, 'succeeded'::public.pgqueue_status]));


--
-- Name: targets_uniq_user_per_dataset; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX targets_uniq_user_per_dataset ON public.targets USING btree (user_id, dataset_id);


--
-- Name: uniq_bloggers_per_dataset; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uniq_bloggers_per_dataset ON public.bloggers USING btree (username, dataset_id);


--
-- Name: bloggers bloggers_dataset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bloggers
    ADD CONSTRAINT bloggers_dataset_id_fkey FOREIGN KEY (dataset_id) REFERENCES public.datasets(id);


--
-- Name: fetch_cursors fetch_cursors_dataset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fetch_cursors
    ADD CONSTRAINT fetch_cursors_dataset_id_fkey FOREIGN KEY (dataset_id) REFERENCES public.datasets(id);


--
-- Name: full_targets full_targets_dataset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_targets
    ADD CONSTRAINT full_targets_dataset_id_fkey FOREIGN KEY (dataset_id) REFERENCES public.datasets(id);


--
-- Name: medias medias_dataset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.medias
    ADD CONSTRAINT medias_dataset_id_fkey FOREIGN KEY (dataset_id) REFERENCES public.datasets(id);


--
-- Name: targets medias_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.targets
    ADD CONSTRAINT medias_fk FOREIGN KEY (media_pk, dataset_id) REFERENCES public.medias(pk, dataset_id);


--
-- PostgreSQL database dump complete
--

