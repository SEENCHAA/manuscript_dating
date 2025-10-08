--
-- PostgreSQL database dump
--

-- Dumped from database version 12.22 (Debian 12.22-1.pgdg120+1)
-- Dumped by pg_dump version 12.22 (Debian 12.22-1.pgdg120+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: letters; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.letters (
    id bigint NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    period_start bigint NOT NULL,
    period_end bigint NOT NULL,
    details text NOT NULL,
    image_url text,
    is_active boolean DEFAULT true
);


ALTER TABLE public.letters OWNER TO postgres;

--
-- Name: letters_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.letters_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.letters_id_seq OWNER TO postgres;

--
-- Name: letters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.letters_id_seq OWNED BY public.letters.id;


--
-- Name: manuscript_letters; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.manuscript_letters (
    manuscript_id bigint NOT NULL,
    letter_id bigint NOT NULL,
    quantity bigint DEFAULT 1,
    is_active boolean DEFAULT true
);


ALTER TABLE public.manuscript_letters OWNER TO postgres;

--
-- Name: manuscripts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.manuscripts (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    status character varying(20) DEFAULT 'draft'::character varying NOT NULL,
    created_at timestamp with time zone,
    submitted_at timestamp with time zone,
    finished_at timestamp with time zone,
    moderator_id bigint
);


ALTER TABLE public.manuscripts OWNER TO postgres;

--
-- Name: manuscripts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.manuscripts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.manuscripts_id_seq OWNER TO postgres;

--
-- Name: manuscripts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.manuscripts_id_seq OWNED BY public.manuscripts.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    is_moderator boolean DEFAULT false
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: letters id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.letters ALTER COLUMN id SET DEFAULT nextval('public.letters_id_seq'::regclass);


--
-- Name: manuscripts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.manuscripts ALTER COLUMN id SET DEFAULT nextval('public.manuscripts_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: letters; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.letters (id, name, description, period_start, period_end, details, image_url, is_active) FROM stdin;
1	в•ЁРЎв•¤Р“в•Ёв•‘в•Ёв–“в•Ёв–‘ в•¤Рі (в•¤Рџв•¤Р’в•¤Рњ)	в•ЁР в•Ёв•‘в•¤Р’в•Ёв••в•Ёв–“в•Ёв•њв•Ёв•› в•Ёв••в•¤Р‘в•Ёв”ђв•Ёв•›в•Ёв•—в•¤Рњв•Ёв•–в•¤Р“в•Ёв•Ўв•¤Р’в•¤Р‘в•¤Рџ в•¤Р‘ XI в•Ёв”ђв•Ёв•› XVIII в•Ёв–“в•Ёв–“.	1000	1750	в•ЁРЎв•¤Р“в•Ёв•‘в•Ёв–“в•Ёв–‘ в•Ёв”ђв•¤Рђв•Ёв••в•¤Рв•Ёв•—в•Ёв–‘ в•Ёв••в•Ёв•– в•¤Р‘в•¤Р’в•Ёв–‘в•¤Рђв•Ёв•›в•¤Р‘в•Ёв•—в•Ёв–‘в•Ёв–“в•¤Рџв•Ёв•њв•¤Р‘в•Ёв•‘в•Ёв•›в•Ёв•Ј в•Ёв”ђв•Ёв••в•¤Р‘в•¤Рњв•Ёв•ќв•Ёв•Ўв•Ёв•њв•Ёв•њв•Ёв•›в•¤Р‘в•¤Р’в•Ёв••. в•ЁРў в•Ёв•ўв•Ёв••в•Ёв–“в•Ёв•›в•Ёв•Ј в•¤Рђв•Ёв•Ўв•¤Р—в•Ёв•• в•Ёв•‘ XVII в•Ёв–“в•Ёв•Ўв•Ёв•‘в•¤Р“ в•¤Рђв•Ёв–‘в•Ёв•–в•Ёв•—в•Ёв••в•¤Р—в•Ёв••в•Ёв•Ў в•Ёв•ќв•Ёв•Ўв•Ёв•ўв•Ёв”¤в•¤Р“ в•¤Рі в•Ёв•• в•Ёв•Ў в•Ёв••в•¤Р‘в•¤Р—в•Ёв•Ўв•Ёв•–в•Ёв•—в•Ёв•›, в•Ёв•њв•Ёв•› в•Ёв–“ в•Ёв”ђв•Ёв••в•¤Р‘в•¤Рњв•Ёв•ќв•Ёв•Ўв•Ёв•њв•Ёв•њв•Ёв•›в•¤Р‘в•¤Р’в•Ёв•• в•Ёв–’в•¤Р“в•Ёв•‘в•Ёв–“в•Ёв–‘ в•¤Р‘в•Ёв•›в•¤Р•в•¤Рђв•Ёв–‘в•Ёв•њв•¤Рџв•Ёв•—в•Ёв–‘в•¤Р‘в•¤Рњ в•Ёв”¤в•Ёв•› в•¤Рђв•Ёв•Ўв•¤Р”в•Ёв•›в•¤Рђв•Ёв•ќв•¤Р› 1918 в•Ёв”‚в•Ёв•›в•Ёв”¤в•Ёв–‘. в•ЁРЇв•Ёв•›в•¤Р‘в•Ёв•—в•Ёв•Ў в•¤Рђв•Ёв•Ўв•¤Р”в•Ёв•›в•¤Рђв•Ёв•ќв•¤Р› в•Ёв•Ўв•¤РЎ в•Ёв”ђв•Ёв•›в•Ёв•—в•Ёв•њв•Ёв•›в•¤Р‘в•¤Р’в•¤Рњв•¤Рћ в•Ёв•–в•Ёв–‘в•Ёв•ќв•Ёв•Ўв•Ёв•њв•Ёв••в•Ёв•—в•Ёв•• в•Ёв•њв•Ёв–‘ в•Ёв•Ў.	http://127.0.0.1:9000/manuscripts/yat.jpg	t
2	в•ЁРЎв•¤Р“в•Ёв•‘в•Ёв–“в•Ёв–‘ в•¤в•Ў (в•Ёв••в•Ёв•ўв•Ёв••в•¤Р–в•Ёв–‘)	в•ЁР®в•Ёв•‘в•Ёв•›в•Ёв•њв•¤Р—в•Ёв–‘в•¤Р’в•Ёв•Ўв•Ёв•—в•¤Рњв•Ёв•њв•Ёв•› в•Ёв••в•¤Р‘в•¤Р—в•Ёв•Ўв•Ёв•–в•Ёв–‘в•Ёв•Ўв•¤Р’ в•Ёв–“ XVIII в•Ёв–“в•Ёв•Ўв•Ёв•‘в•Ёв•Ў.	1000	1700	в•ЁРЁв•Ёв•ўв•Ёв••в•¤Р–в•Ёв–‘ в•Ёв••в•¤Р‘в•Ёв”ђв•Ёв•›в•Ёв•—в•¤Рњв•Ёв•–в•Ёв•›в•Ёв–“в•Ёв–‘в•Ёв•—в•Ёв–‘в•¤Р‘в•¤Рњ в•Ёв”¤в•Ёв•—в•¤Рџ в•Ёв”ђв•Ёв•Ўв•¤Рђв•Ёв•Ўв•Ёв”¤в•Ёв–‘в•¤Р—в•Ёв•• в•Ёв”‚в•¤Рђв•Ёв•Ўв•¤Р—в•Ёв•Ўв•¤Р‘в•Ёв•‘в•Ёв•›в•Ёв•Ј в•Ёв–’в•¤Р“в•Ёв•‘в•Ёв–“в•¤Р› 'в•§Р•'. в•ЁРів•Ёв•ўв•Ёв•Ў в•Ёв–“ XVII в•Ёв–“в•Ёв•Ўв•Ёв•‘в•Ёв•Ў в•Ёв••в•Ёв•– в•¤Р“в•Ёв”ђв•Ёв•›в•¤Р’в•¤Рђв•Ёв•Ўв•Ёв–’в•Ёв•—в•Ёв•Ўв•Ёв•њв•Ёв••в•¤Рџ в•Ёв”ђв•Ёв•›в•¤Р—в•¤Р’в•Ёв•• в•Ёв••в•¤Р‘в•¤Р—в•Ёв•Ўв•Ёв•–в•Ёв•—в•Ёв–‘, в•Ёв–“в•¤Р‘в•¤Р’в•¤Рђв•Ёв•Ўв•¤Р—в•Ёв–‘в•Ёв•—в•Ёв–‘в•¤Р‘в•¤Рњ в•¤Р’в•Ёв•›в•Ёв•—в•¤Рњв•Ёв•‘в•Ёв•› в•Ёв–“ в•¤Р–в•Ёв•Ўв•¤Рђв•Ёв•‘в•Ёв•›в•Ёв–“в•Ёв•њв•¤Р›в•¤Р• в•Ёв•‘в•Ёв•њв•Ёв••в•Ёв”‚в•Ёв–‘в•¤Р•. в•ЁРў 1735 в•Ёв”‚в•Ёв•›в•Ёв”¤в•¤Р“ в•Ёв••в•Ёв•ўв•Ёв••в•¤Р–в•Ёв–‘ в•Ёв–’в•¤Р›в•Ёв•—в•Ёв–‘ в•Ёв••в•¤Р‘в•Ёв•‘в•Ёв•—в•¤Рћв•¤Р—в•Ёв•Ўв•Ёв•њв•Ёв–‘ в•Ёв••в•Ёв•– в•Ёв”‚в•¤Рђв•Ёв–‘в•Ёв•ўв•Ёв”¤в•Ёв–‘в•Ёв•њв•¤Р‘в•Ёв•‘в•Ёв•›в•Ёв”‚в•Ёв•› в•¤Рв•¤Рђв•Ёв••в•¤Р”в•¤Р’в•Ёв–‘ в•¤Р“в•Ёв•‘в•Ёв–‘в•Ёв•–в•Ёв•›в•Ёв•ќ в•ЁР в•Ёв•‘в•Ёв–‘в•Ёв”¤в•Ёв•Ўв•Ёв•ќв•Ёв••в•Ёв•• в•Ёв•њв•Ёв–‘в•¤Р“в•Ёв•‘.	http://127.0.0.1:9000/manuscripts/izhitsa.jpg	t
3	в•ЁРЎв•¤Р“в•Ёв•‘в•Ёв–“в•Ёв–‘ в•¤в–“ (в•¤Р”в•Ёв••в•¤Р’в•Ёв–‘)	в•ЁРЁв•¤Р‘в•¤Р—в•Ёв•Ўв•Ёв•–в•Ёв–‘в•Ёв•Ўв•¤Р’ в•Ёв–“ в•Ёв•њв•Ёв–‘в•¤Р—в•Ёв–‘в•Ёв•—в•Ёв•Ў XX в•Ёв–“в•Ёв•Ўв•Ёв•‘в•Ёв–‘ в•Ёв”ђв•Ёв•›в•¤Р‘в•Ёв•—в•Ёв•Ў в•¤Рђв•Ёв•Ўв•¤Р”в•Ёв•›в•¤Рђв•Ёв•ќв•¤Р›.	1000	1918	в•ЁРґв•Ёв••в•¤Р’в•Ёв–‘ в•Ёв••в•¤Р‘в•Ёв”ђв•Ёв•›в•Ёв•—в•¤Рњв•Ёв•–в•Ёв•›в•Ёв–“в•Ёв–‘в•Ёв•—в•Ёв–‘в•¤Р‘в•¤Рњ в•Ёв–“ в•Ёв•–в•Ёв–‘в•Ёв••в•Ёв•ќв•¤Р‘в•¤Р’в•Ёв–“в•Ёв•›в•Ёв–“в•Ёв–‘в•Ёв•њв•Ёв•њв•¤Р›в•¤Р• в•Ёв”‚в•¤Рђв•Ёв•Ўв•¤Р—в•Ёв•Ўв•¤Р‘в•Ёв•‘в•Ёв••в•¤Р• в•¤Р‘в•Ёв•—в•Ёв•›в•Ёв–“в•Ёв–‘в•¤Р• в•Ёв”¤в•Ёв•—в•¤Рџ в•Ёв”ђв•Ёв•Ўв•¤Рђв•Ёв•Ўв•Ёв”¤в•Ёв–‘в•¤Р—в•Ёв•• в•Ёв•–в•Ёв–“в•¤Р“в•Ёв•‘в•Ёв–‘ 'в•¤Р”'. в•ЁРў в•¤Рђв•¤Р“в•¤Р‘в•¤Р‘в•Ёв•‘в•Ёв•›в•Ёв•ќ в•¤Рџв•Ёв•–в•¤Р›в•Ёв•‘в•Ёв•Ў в•Ёв•–в•Ёв–“в•¤Р“в•Ёв•‘в•Ёв–‘ [в•¬в••] в•Ёв•њв•Ёв••в•Ёв•‘в•Ёв•›в•Ёв”‚в•Ёв”¤в•Ёв–‘ в•Ёв•њв•Ёв•Ў в•¤Р‘в•¤Р“в•¤Р™в•Ёв•Ўв•¤Р‘в•¤Р’в•Ёв–“в•Ёв•›в•Ёв–“в•Ёв–‘в•Ёв•—в•Ёв•›, в•Ёв”ђв•Ёв•›в•¤Рќв•¤Р’в•Ёв•›в•Ёв•ќв•¤Р“ в•¤Р”в•Ёв••в•¤Р’в•Ёв–‘ в•Ёв•–в•Ёв–“в•¤Р“в•¤Р—в•Ёв–‘в•Ёв•—в•Ёв–‘ в•¤Р’в•Ёв–‘в•Ёв•‘ в•Ёв•ўв•Ёв•Ў, в•Ёв•‘в•Ёв–‘в•Ёв•‘ в•Ёв•›в•Ёв–’в•¤Р›в•¤Р—в•Ёв•њв•Ёв–‘в•¤Рџ в•¤Р”. в•ЁРў в•¤Рђв•Ёв•Ўв•Ёв•–в•¤Р“в•Ёв•—в•¤Рњв•¤Р’в•Ёв–‘в•¤Р’в•Ёв•Ў в•Ёв•›в•Ёв•њв•Ёв–‘ в•¤Р‘в•¤Р’в•Ёв–‘в•Ёв•—в•Ёв–‘ в•Ёв••в•Ёв•–в•Ёв–’в•¤Р›в•¤Р’в•Ёв•›в•¤Р—в•Ёв•њв•Ёв•›в•Ёв•Ј в•Ёв•• в•Ёв–’в•¤Р›в•Ёв•—в•Ёв–‘ в•Ёв••в•¤Р‘в•Ёв•‘в•Ёв•—в•¤Рћв•¤Р—в•Ёв•Ўв•Ёв•њв•Ёв–‘ в•Ёв–“ 1918 в•Ёв”‚в•Ёв•›в•Ёв”¤в•¤Р“, в•Ёв•–в•Ёв–‘в•Ёв•ќв•Ёв•Ўв•Ёв•њв•Ёв••в•Ёв–“в•¤Рв•Ёв••в•¤Р‘в•¤Рњ в•Ёв•њв•Ёв–‘ в•¤Р”.	http://127.0.0.1:9000/manuscripts/fita.jpg	t
4	в•ЁРІв•Ёв–“в•¤РЎв•¤Рђв•Ёв”¤в•¤Р›в•Ёв•Ј в•Ёв•–в•Ёв•њв•Ёв–‘в•Ёв•‘ в•Ёв•њв•Ёв–‘ в•Ёв•‘в•Ёв•›в•Ёв•њв•¤Р–в•Ёв•Ў в•¤Р‘в•Ёв•—в•Ёв•›в•Ёв–“	в•ЁР®в•Ёв–’в•¤Рџв•Ёв•–в•Ёв–‘в•¤Р’в•Ёв•Ўв•Ёв•—в•Ёв•Ўв•Ёв•њ в•Ёв”¤в•Ёв•› в•¤Рђв•Ёв•Ўв•¤Р”в•Ёв•›в•¤Рђв•Ёв•ќв•¤Р› 1918 в•Ёв”‚в•Ёв•›в•Ёв”¤в•Ёв–‘.	1000	1918	в•ЁР¤в•Ёв•› в•¤Рђв•Ёв•Ўв•¤Р”в•Ёв•›в•¤Рђв•Ёв•ќв•¤Р› 1918 в•Ёв”‚в•Ёв•›в•Ёв”¤в•Ёв–‘ в•¤Р’в•Ёв–“в•¤РЎв•¤Рђв•Ёв”¤в•¤Р›в•Ёв•Ј в•Ёв•–в•Ёв•њв•Ёв–‘в•Ёв•‘ в•¤Р‘в•¤Р’в•Ёв–‘в•Ёв–“в•Ёв••в•Ёв•—в•¤Р‘в•¤Рџ в•Ёв•њв•Ёв–‘ в•Ёв•‘в•Ёв•›в•Ёв•њв•¤Р–в•Ёв•Ў в•Ёв•‘в•Ёв–‘в•Ёв•ўв•Ёв”¤в•Ёв•›в•Ёв”‚в•Ёв•› в•¤Р‘в•Ёв•—в•Ёв•›в•Ёв–“в•Ёв–‘ в•Ёв”ђв•Ёв•›в•¤Р‘в•Ёв•—в•Ёв•Ў в•¤Р‘в•Ёв•›в•Ёв”‚в•Ёв•—в•Ёв–‘в•¤Р‘в•Ёв•њв•Ёв•›в•Ёв•Ј. в•ЁРЅв•¤Р’в•Ёв–‘ в•Ёв•њв•Ёв•›в•¤Рђв•Ёв•ќв•Ёв–‘ в•¤Р‘в•Ёв••в•Ёв•—в•¤Рњв•Ёв•њв•Ёв•› в•¤Р“в•¤Р’в•¤Рџв•Ёв•ўв•Ёв•Ўв•Ёв•—в•¤Рџв•Ёв•—в•Ёв–‘ в•¤Р’в•Ёв•Ўв•Ёв•‘в•¤Р‘в•¤Р’в•¤Р›: в•¤Р‘в•¤Р—в•Ёв••в•¤Р’в•Ёв–‘в•Ёв•Ўв•¤Р’в•¤Р‘в•¤Рџ, в•¤Р—в•¤Р’в•Ёв•› в•Ёв•›в•Ёв•‘в•Ёв•›в•Ёв•—в•Ёв•› 4С‚РђРЈ5% в•Ёв”ђв•Ёв•Ўв•¤Р—в•Ёв–‘в•¤Р’в•Ёв•њв•Ёв•›в•Ёв•Ј в•Ёв”ђв•Ёв•—в•Ёв•›в•¤Р™в•Ёв–‘в•Ёв”¤в•Ёв•• в•Ёв–“ в•Ёв•‘в•Ёв•њв•Ёв••в•Ёв”‚в•Ёв–‘в•¤Р• в•Ёв•–в•Ёв–‘в•Ёв•њв•Ёв••в•Ёв•ќв•Ёв–‘в•Ёв•—в•Ёв•• в•Ёв•њв•Ёв•Ўв•Ёв•њв•¤Р“в•Ёв•ўв•Ёв•њв•¤Р›в•Ёв•Ў в•¤Р’в•Ёв–“в•¤РЎв•¤Рђв•Ёв”¤в•¤Р›в•Ёв•Ў в•Ёв•–в•Ёв•њв•Ёв–‘в•Ёв•‘в•Ёв••. в•ЁРў 1918 в•Ёв”‚в•Ёв•›в•Ёв”¤в•¤Р“ в•Ёв•›в•Ёв–’в•¤Рџв•Ёв•–в•Ёв–‘в•¤Р’в•Ёв•Ўв•Ёв•—в•¤Рњв•Ёв•њв•Ёв•›в•Ёв•Ў в•Ёв•њв•Ёв–‘в•Ёв”ђв•Ёв••в•¤Р‘в•Ёв–‘в•Ёв•њв•Ёв••в•Ёв•Ў в•Ёв•њв•Ёв–‘ в•Ёв•‘в•Ёв•›в•Ёв•њв•¤Р–в•Ёв•Ў в•¤Р‘в•Ёв•—в•Ёв•›в•Ёв–“ в•Ёв•›в•¤Р’в•Ёв•ќв•Ёв•Ўв•Ёв•њв•Ёв••в•Ёв•—в•Ёв••.	http://127.0.0.1:9000/manuscripts/hardsign.jpg	t
5	в•ЁРЎв•¤Р“в•Ёв•‘в•Ёв–“в•Ёв–‘ i	в•ЁРЁв•¤Р‘в•Ёв”ђв•Ёв•›в•Ёв•—в•¤Рњв•Ёв•–в•¤Р“в•Ёв•Ўв•¤Р’в•¤Р‘в•¤Рџ в•Ёв”¤в•Ёв•› в•¤Рђв•Ёв•Ўв•¤Р”в•Ёв•›в•¤Рђв•Ёв•ќв•¤Р› 1918 в•Ёв”‚в•Ёв•›в•Ёв”¤в•Ёв–‘.	1000	1918	в•ЁРЅв•¤Р’в•Ёв–‘ в•Ёв–’в•¤Р“в•Ёв•‘в•Ёв–“в•Ёв–‘ в•Ёв–’в•¤Р›в•Ёв•—в•Ёв–‘ в•Ёв–“в•Ёв–“в•Ёв•Ўв•Ёв”¤в•Ёв•Ўв•Ёв•њв•Ёв–‘ в•Ёв”¤в•Ёв•—в•¤Рџ в•Ёв–’в•Ёв•›в•Ёв•—в•Ёв•Ўв•Ёв•Ў в•¤Р’в•Ёв•›в•¤Р—в•Ёв•њв•Ёв•›в•Ёв”‚в•Ёв•› в•¤Р‘в•Ёв•›в•Ёв•›в•¤Р’в•Ёв–“в•Ёв•Ўв•¤Р’в•¤Р‘в•¤Р’в•Ёв–“в•Ёв••в•¤Рџ в•Ёв”‚в•¤Рђв•Ёв•Ўв•¤Р—в•Ёв•Ўв•¤Р‘в•Ёв•‘в•Ёв•›в•Ёв•ќв•¤Р“ в•Ёв”ђв•Ёв••в•¤Р‘в•¤Рњв•Ёв•ќв•¤Р“ в•Ёв•• в•Ёв”¤в•Ёв•—в•¤Рџ в•¤Рђв•Ёв–‘в•Ёв•–в•Ёв•—в•Ёв••в•¤Р—в•Ёв•Ўв•Ёв•њв•Ёв••в•¤Рџ в•Ёв•›в•Ёв•ќв•Ёв•›в•Ёв•њв•Ёв••в•Ёв•ќв•Ёв•›в•Ёв–“. в•ЁР­в•Ёв•› в•Ёв•‘ XIX в•Ёв–“в•Ёв•Ўв•Ёв•‘в•¤Р“ в•¤Рђв•Ёв–‘в•Ёв•–в•Ёв•—в•Ёв••в•¤Р—в•Ёв••в•Ёв•Ў в•¤Р‘в•¤Р’в•Ёв–‘в•Ёв•—в•Ёв•› в•¤Р”в•Ёв•›в•¤Рђв•Ёв•ќв•Ёв–‘в•Ёв•—в•¤Рњв•Ёв•њв•¤Р›в•Ёв•ќ в•Ёв•• в•Ёв•ќв•Ёв•Ўв•¤Рв•Ёв–‘в•Ёв•—в•Ёв•› в•Ёв•›в•Ёв–’в•¤Р“в•¤Р—в•Ёв•Ўв•Ёв•њв•Ёв••в•¤Рћ. в•ЁРў 1918 в•Ёв”‚в•Ёв•›в•Ёв”¤в•¤Р“ в•Ёв•Ўв•¤РЎ в•Ёв•›в•¤Р’в•Ёв•ќв•Ёв•Ўв•Ёв•њв•Ёв••в•Ёв•—в•Ёв••, в•Ёв•• в•Ёв–“в•¤Р‘в•Ёв•Ў в•¤Р‘в•Ёв•—в•¤Р“в•¤Р—в•Ёв–‘в•Ёв•• в•Ёв•њв•Ёв–‘в•Ёв”ђв•Ёв••в•¤Р‘в•Ёв–‘в•Ёв•њв•Ёв••в•¤Рџ в•Ёв•–в•Ёв–‘в•Ёв•ќв•Ёв•Ўв•Ёв•њв•Ёв••в•Ёв•—в•Ёв•• в•Ёв•њв•Ёв–‘ в•Ёв••.	http://127.0.0.1:9000/manuscripts/decimalI.jpg	t
\.


--
-- Data for Name: manuscript_letters; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.manuscript_letters (manuscript_id, letter_id, quantity, is_active) FROM stdin;
\.


--
-- Data for Name: manuscripts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.manuscripts (id, user_id, status, created_at, submitted_at, finished_at, moderator_id) FROM stdin;
1	1	draft	2025-10-02 20:54:26.315124+00	\N	\N	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, password, is_moderator) FROM stdin;
1	scribe1	hashed_password1	f
2	moderator1	hashed_password2	t
\.


--
-- Name: letters_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.letters_id_seq', 5, true);


--
-- Name: manuscripts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.manuscripts_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 7, true);


--
-- Name: letters letters_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.letters
    ADD CONSTRAINT letters_pkey PRIMARY KEY (id);


--
-- Name: manuscript_letters manuscript_letters_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.manuscript_letters
    ADD CONSTRAINT manuscript_letters_pkey PRIMARY KEY (manuscript_id, letter_id);


--
-- Name: manuscripts manuscripts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.manuscripts
    ADD CONSTRAINT manuscripts_pkey PRIMARY KEY (id);


--
-- Name: users uni_users_username; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_username UNIQUE (username);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: manuscript_letters fk_manuscript_letters_letter; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.manuscript_letters
    ADD CONSTRAINT fk_manuscript_letters_letter FOREIGN KEY (letter_id) REFERENCES public.letters(id);


--
-- Name: manuscript_letters fk_manuscripts_letters; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.manuscript_letters
    ADD CONSTRAINT fk_manuscripts_letters FOREIGN KEY (manuscript_id) REFERENCES public.manuscripts(id);


--
-- Name: manuscripts fk_manuscripts_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.manuscripts
    ADD CONSTRAINT fk_manuscripts_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

