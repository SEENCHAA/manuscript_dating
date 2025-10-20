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
    moderator_id bigint,
    calculated_period character varying(50)
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
1	Буква ѣ (ять)	Активно использовалась с XI по XVIII вв.	1000	1750	Буква пришла из старославянской письменности. В живой речи к XVII веку различие между ѣ и е исчезло, но в письменности буква сохранялась до реформы 1918 года. После реформы её полностью заменили на е.	http://127.0.0.1:9000/manuscripts/yat.jpg	t
2	Буква ѵ (ижица)	Окончательно исчезает в XVIII веке.	1000	1700	Ижица использовалась для передачи греческой буквы "υ". Уже в XVII веке из употребления почти исчезла, встречалась только в церковных книгах. В 1735 году ижица была исключена из гражданского шрифта указом Академии наук.	http://127.0.0.1:9000/manuscripts/izhitsa.jpg	t
3	Буква Ѳ (фита)	Исчезает в начале XX века после реформы.	1000	1918	Фита использовалась в заимствованных греческих словах для передачи звука "ф". В русском языке звука [θ] никогда не существовало, поэтому фита звучала так же, как обычная ф. В результате она стала избыточной и была исключена в 1918 году, заменившись на ф.	http://127.0.0.1:9000/manuscripts/fita.jpg	t
4	Твёрдый знак	Употреблялся до реформы 1918 года.	1000	1918	До реформы 1918 года твёрдый знак использовался на конце слов после согласных. После реформы его употребление было ограничено только определёнными случаями.	http://127.0.0.1:9000/manuscripts/hardsign.jpg	t
5	Буква i	Отменена реформой 1918 года.	1000	1918	Буква i использовалась наряду с и. После реформы 1918 года была упразднена, а все случаи её употребления заменены на и.	http://127.0.0.1:9000/manuscripts/decimalI.jpg	t
6	Буква ѯ (кси)	Использовалась в древнерусских текстах с XI по XVIII век\r\n	1000	1700	Буква ѯ (кси) пришла из греческого алфавита, где обозначала звук 'кс'. В старославянской письменности употреблялась преимущественно в заимствованных и церковных словах греческого происхождения, таких как 'ѯристос' (Христос). В повседневных текстах встречалась редко. После реформы гражданского шрифта при Петре I постепенно вышла из употребления и полностью исчезла к XVIII веку.	http://127.0.0.1:9000/manuscripts/ksi.png	t
\.


--
-- Data for Name: manuscript_letters; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.manuscript_letters (manuscript_id, letter_id, quantity, is_active) FROM stdin;
1	2	5	t
1	3	1	t
2	1	1	t
2	2	1	t
3	1	1	t
4	1	1	t
4	2	1	t
5	2	1	t
6	1	232	t
\.


--
-- Data for Name: manuscripts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.manuscripts (id, user_id, status, created_at, submitted_at, finished_at, moderator_id, calculated_period) FROM stdin;
1	1	deleted	2025-10-02 20:54:26.315124+00	\N	\N	\N	\N
2	1	deleted	2025-10-08 10:03:46.859673+00	\N	\N	\N	\N
3	1	deleted	2025-10-08 10:15:04.051241+00	\N	\N	\N	\N
4	1	deleted	2025-10-08 10:59:13.6621+00	\N	\N	\N	\N
5	1	deleted	2025-10-08 11:20:11.391584+00	\N	\N	\N	\N
6	1	deleted	2025-10-08 14:41:44.812765+00	\N	\N	\N	\N
7	1	deleted	2025-10-16 11:53:44.028419+00	\N	\N	\N	\N
8	1	deleted	2025-10-16 11:53:51.460736+00	\N	\N	\N	\N
9	1	deleted	2025-10-16 11:53:58.740095+00	\N	\N	\N	\N
10	1	draft	2025-10-16 11:54:17.731994+00	\N	\N	\N	\N
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

SELECT pg_catalog.setval('public.letters_id_seq', 1, true);


--
-- Name: manuscripts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.manuscripts_id_seq', 10, true);


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
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

