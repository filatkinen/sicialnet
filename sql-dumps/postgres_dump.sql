--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3 (Ubuntu 15.3-0ubuntu0.23.04.1)
-- Dumped by pg_dump version 15.3 (Ubuntu 15.3-0ubuntu0.23.04.1)

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

ALTER TABLE IF EXISTS ONLY public.user_credentials DROP CONSTRAINT IF EXISTS user_credentials_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.token DROP CONSTRAINT IF EXISTS tokens_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.user_credentials DROP CONSTRAINT IF EXISTS user_credentials_pkey;
ALTER TABLE IF EXISTS ONLY public.token DROP CONSTRAINT IF EXISTS tokens_pkey;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.user_credentials;
DROP TABLE IF EXISTS public.token;
DROP TYPE IF EXISTS public.sex_status;
--
-- Name: sex_status; Type: TYPE; Schema: public; Owner: socialnet
--

CREATE TYPE public.sex_status AS ENUM (
    'male',
    'female'
);


ALTER TYPE public.sex_status OWNER TO socialnet;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: token; Type: TABLE; Schema: public; Owner: socialnet
--

CREATE TABLE public.token (
    hash character(64) NOT NULL,
    user_id character(36) NOT NULL,
    expires timestamp without time zone NOT NULL
);


ALTER TABLE public.token OWNER TO socialnet;

--
-- Name: user_credentials; Type: TABLE; Schema: public; Owner: socialnet
--

CREATE TABLE public.user_credentials (
    user_id character(36) NOT NULL,
    password character(60)
);


ALTER TABLE public.user_credentials OWNER TO socialnet;

--
-- Name: users; Type: TABLE; Schema: public; Owner: socialnet
--

CREATE TABLE public.users (
    user_id character(36) NOT NULL,
    first_name character varying(64) NOT NULL,
    second_name character varying(64),
    sex public.sex_status,
    biography text,
    city character varying(64),
    birthdate timestamp without time zone
);


ALTER TABLE public.users OWNER TO socialnet;

--
-- Data for Name: token; Type: TABLE DATA; Schema: public; Owner: socialnet
--

INSERT INTO public.token (hash, user_id, expires) VALUES ('ba96e5f66d29570dfef7dbf54d921ec4f3cf7d8abd49a138c25c31425df7f776', '37b48b26-03b1-9d29-57ad-b306b41edbda', '2023-07-20 17:06:00');
INSERT INTO public.token (hash, user_id, expires) VALUES ('762af679f6a8fbb5340e8b0e71f92c09af635261f45b9e150e106489ba9c6a83', '63314972-3de1-ee61-7f26-0ff16cdb46e3', '2023-07-20 17:06:00');


--
-- Data for Name: user_credentials; Type: TABLE DATA; Schema: public; Owner: socialnet
--

INSERT INTO public.user_credentials (user_id, password) VALUES ('37b48b26-03b1-9d29-57ad-b306b41edbda', '$2a$04$kenc8nkCiWk/.QQOMUDLeOunD6EnAReZcqq6KSf4sI05m7GMHMtNi');
INSERT INTO public.user_credentials (user_id, password) VALUES ('63314972-3de1-ee61-7f26-0ff16cdb46e3', '$2a$04$fzeG7EiczYA4yofcK9SaeeGIh0CRzus/RoN/m61VRgN1QAXgDsFkO');


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: socialnet
--

INSERT INTO public.users (user_id, first_name, second_name, sex, biography, city, birthdate) VALUES ('37b48b26-03b1-9d29-57ad-b306b41edbda', 'Ivan1', 'Frolov', NULL, 'Hokkey', 'Moskva', '2002-02-11 00:00:00');
INSERT INTO public.users (user_id, first_name, second_name, sex, biography, city, birthdate) VALUES ('63314972-3de1-ee61-7f26-0ff16cdb46e3', 'Masha1', 'Frolova', NULL, 'Dance', 'Piter', '2003-02-11 00:00:00');


--
-- Name: token tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: socialnet
--

ALTER TABLE ONLY public.token
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (hash);


--
-- Name: user_credentials user_credentials_pkey; Type: CONSTRAINT; Schema: public; Owner: socialnet
--

ALTER TABLE ONLY public.user_credentials
    ADD CONSTRAINT user_credentials_pkey PRIMARY KEY (user_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: socialnet
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: token tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: socialnet
--

ALTER TABLE ONLY public.token
    ADD CONSTRAINT tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- Name: user_credentials user_credentials_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: socialnet
--

ALTER TABLE ONLY public.user_credentials
    ADD CONSTRAINT user_credentials_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

