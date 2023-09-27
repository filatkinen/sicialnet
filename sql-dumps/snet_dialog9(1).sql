--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4 (Ubuntu 15.4-0ubuntu0.23.04.1)
-- Dumped by pg_dump version 15.4 (Ubuntu 15.4-0ubuntu0.23.04.1)

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

DROP DATABASE IF EXISTS snet_dialog;
--
-- Name: snet_dialog; Type: DATABASE; Schema: -; Owner: socialnet
--

CREATE DATABASE snet_dialog WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.UTF-8';


ALTER DATABASE snet_dialog OWNER TO socialnet;

\connect snet_dialog

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
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: dialogs; Type: TABLE; Schema: public; Owner: socialnet
--

CREATE TABLE public.dialogs (
    dialog_id character(36) NOT NULL,
    user_id character(36),
    friend_id character(36),
    message text
);


ALTER TABLE public.dialogs OWNER TO socialnet;

--
-- Data for Name: dialogs; Type: TABLE DATA; Schema: public; Owner: socialnet
--

INSERT INTO public.dialogs VALUES ('G35T6I43R2ZAQDVUZZHXE2RC44          ', '628ad2a9-1676-5aba-f27c-dd38363128b8', 'aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b', 'message 0101');
INSERT INTO public.dialogs VALUES ('B35I3LJ7NKF2F3VA45KGKW2FKE          ', '628ad2a9-1676-5aba-f27c-dd38363128b8', 'aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b', 'message 0102');
INSERT INTO public.dialogs VALUES ('NFSSGNNXCQHBCUTVOMWHIAWMJM          ', 'aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b', '628ad2a9-1676-5aba-f27c-dd38363128b8', 'message 0201');
INSERT INTO public.dialogs VALUES ('QQMG3B3R6JPRDEIUXWPUERDRM4          ', 'aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b', '628ad2a9-1676-5aba-f27c-dd38363128b8', 'message 0202');
INSERT INTO public.dialogs VALUES ('LUA4AEHPSQOHA7BRAQNKCGWYFY          ', '628ad2a9-1676-5aba-f27c-dd38363128b8', '902c7d85-ace2-319d-6852-4224b2866067', 'message1 From Katya to Anna');
INSERT INTO public.dialogs VALUES ('OKBS3V55TVLXLCARD6HLH3LPSM          ', '628ad2a9-1676-5aba-f27c-dd38363128b8', '902c7d85-ace2-319d-6852-4224b2866067', 'message1 From Katya to Anna');
INSERT INTO public.dialogs VALUES ('CV5YKJXOZ2RVX7NV4YOCEPOQQM          ', '628ad2a9-1676-5aba-f27c-dd38363128b8', '902c7d85-ace2-319d-6852-4224b2866067', 'message2 From Katya to Anna');
INSERT INTO public.dialogs VALUES ('6YATVFPW5SRDYFOHHW2BAZKCVQ          ', '902c7d85-ace2-319d-6852-4224b2866067', '628ad2a9-1676-5aba-f27c-dd38363128b8', 'message1 From Anna to Katya');
INSERT INTO public.dialogs VALUES ('PO3MGBAQAHJ7PHSZCH26WMEDVI          ', '902c7d85-ace2-319d-6852-4224b2866067', '628ad2a9-1676-5aba-f27c-dd38363128b8', 'message2 From Anna to Katya');
INSERT INTO public.dialogs VALUES ('VW6NV7OBZBXXQGLPEEQGGPDDJ4          ', '628ad2a9-1676-5aba-f27c-dd38363128b8', '902c7d85-ace2-319d-6852-4224b2866067', 'message1 From Katya to Anna');
INSERT INTO public.dialogs VALUES ('MCR3XWT2MDGDZ2KCOO5HNC5U7U          ', '628ad2a9-1676-5aba-f27c-dd38363128b8', '902c7d85-ace2-319d-6852-4224b2866067', 'message1 From Katya to Anna');


--
-- PostgreSQL database dump complete
--

