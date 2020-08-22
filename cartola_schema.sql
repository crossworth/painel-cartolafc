--
-- PostgreSQL database dump
--

-- Dumped from database version 11.8 (Debian 11.8-1.pgdg100+1)
-- Dumped by pg_dump version 12.3

-- Started on 2020-08-21 21:52:32

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
-- TOC entry 2 (class 3079 OID 192549)
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- TOC entry 3067 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


--
-- TOC entry 3 (class 3079 OID 136270)
-- Name: unaccent; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS unaccent WITH SCHEMA public;


--
-- TOC entry 3068 (class 0 OID 0)
-- Dependencies: 3
-- Name: EXTENSION unaccent; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION unaccent IS 'text search dictionary that removes accents';


--
-- TOC entry 230 (class 1255 OID 136277)
-- Name: f_unaccent(text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.f_unaccent(text) RETURNS text
    LANGUAGE sql IMMUTABLE STRICT PARALLEL SAFE
    AS $_$
SELECT public.unaccent('public.unaccent', $1)  -- schema-qualify function and dictionary
$_$;


ALTER FUNCTION public.f_unaccent(text) OWNER TO postgres;

--
-- TOC entry 247 (class 1255 OID 152576)
-- Name: insert_full_text_search(); Type: FUNCTION; Schema: public; Owner: cartola
--

CREATE FUNCTION public.insert_full_text_search() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO full_text_search_topic
    (topic_id, tsv, date)  VALUES (NEW.id, setweight(to_tsvector(NEW.title), 'A'), NEW.created_at)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.insert_full_text_search() OWNER TO cartola;

--
-- TOC entry 263 (class 1255 OID 159325)
-- Name: insert_full_text_search_comment(); Type: FUNCTION; Schema: public; Owner: cartola
--

CREATE FUNCTION public.insert_full_text_search_comment() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO full_text_search_comment
    (topic_id, comment_id, tsv, date)  VALUES (NEW.topic_id, NEW.id, setweight(to_tsvector(noquote(NEW.text)), 'A'), NEW.date)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.insert_full_text_search_comment() OWNER TO cartola;

--
-- TOC entry 231 (class 1255 OID 19500)
-- Name: insert_profile(); Type: FUNCTION; Schema: public; Owner: cartola
--

CREATE FUNCTION public.insert_profile() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
  BEGIN
    INSERT INTO profile_names 
     (profile_id, first_name, last_name, screen_name, photo, date) 
    VALUES
      (NEW.id, NEW.first_name, NEW.last_name, NEW.screen_name, NEW.photo, extract(epoch from now()))
    ON CONFLICT DO NOTHING;
    RETURN NEW;
  END;
$$;


ALTER FUNCTION public.insert_profile() OWNER TO cartola;

--
-- TOC entry 262 (class 1255 OID 193837)
-- Name: noquote(text); Type: FUNCTION; Schema: public; Owner: cartola
--

CREATE FUNCTION public.noquote(txt text) RETURNS text
    LANGUAGE plpgsql IMMUTABLE
    AS $$
	BEGIN
		return regexp_replace(txt, '\[(id|club)\d+:?(bp-\d+_\d+)?\|([@\d\w\s_\-A-Za-z*\xAA\xB5\xBA\xC0-\xD6\xD8-\xF6\xF8-\u02C1\u02C6-\u02D1\u02E0-\u02E4\u02EC\u02EE\u0370-\u0374\u0376\u0377\u037A-\u037D\u0386\u0388-\u038A\u038C\u038E-\u03A1\u03A3-\u03F5\u03F7-\u0481\u048A-\u0527\u0531-\u0556\u0559\u0561-\u0587\u05D0-\u05EA\u05F0-\u05F2\u0620-\u064A\u066E\u066F\u0671-\u06D3\u06D5\u06E5\u06E6\u06EE\u06EF\u06FA-\u06FC\u06FF\u0710\u0712-\u072F\u074D-\u07A5\u07B1\u07CA-\u07EA\u07F4\u07F5\u07FA\u0800-\u0815\u081A\u0824\u0828\u0840-\u0858\u08A0\u08A2-\u08AC\u0904-\u0939\u093D\u0950\u0958-\u0961\u0971-\u0977\u0979-\u097F\u0985-\u098C\u098F\u0990\u0993-\u09A8\u09AA-\u09B0\u09B2\u09B6-\u09B9\u09BD\u09CE\u09DC\u09DD\u09DF-\u09E1\u09F0\u09F1\u0A05-\u0A0A\u0A0F\u0A10\u0A13-\u0A28\u0A2A-\u0A30\u0A32\u0A33\u0A35\u0A36\u0A38\u0A39\u0A59-\u0A5C\u0A5E\u0A72-\u0A74\u0A85-\u0A8D\u0A8F-\u0A91\u0A93-\u0AA8\u0AAA-\u0AB0\u0AB2\u0AB3\u0AB5-\u0AB9\u0ABD\u0AD0\u0AE0\u0AE1\u0B05-\u0B0C\u0B0F\u0B10\u0B13-\u0B28\u0B2A-\u0B30\u0B32\u0B33\u0B35-\u0B39\u0B3D\u0B5C\u0B5D\u0B5F-\u0B61\u0B71\u0B83\u0B85-\u0B8A\u0B8E-\u0B90\u0B92-\u0B95\u0B99\u0B9A\u0B9C\u0B9E\u0B9F\u0BA3\u0BA4\u0BA8-\u0BAA\u0BAE-\u0BB9\u0BD0\u0C05-\u0C0C\u0C0E-\u0C10\u0C12-\u0C28\u0C2A-\u0C33\u0C35-\u0C39\u0C3D\u0C58\u0C59\u0C60\u0C61\u0C85-\u0C8C\u0C8E-\u0C90\u0C92-\u0CA8\u0CAA-\u0CB3\u0CB5-\u0CB9\u0CBD\u0CDE\u0CE0\u0CE1\u0CF1\u0CF2\u0D05-\u0D0C\u0D0E-\u0D10\u0D12-\u0D3A\u0D3D\u0D4E\u0D60\u0D61\u0D7A-\u0D7F\u0D85-\u0D96\u0D9A-\u0DB1\u0DB3-\u0DBB\u0DBD\u0DC0-\u0DC6\u0E01-\u0E30\u0E32\u0E33\u0E40-\u0E46\u0E81\u0E82\u0E84\u0E87\u0E88\u0E8A\u0E8D\u0E94-\u0E97\u0E99-\u0E9F\u0EA1-\u0EA3\u0EA5\u0EA7\u0EAA\u0EAB\u0EAD-\u0EB0\u0EB2\u0EB3\u0EBD\u0EC0-\u0EC4\u0EC6\u0EDC-\u0EDF\u0F00\u0F40-\u0F47\u0F49-\u0F6C\u0F88-\u0F8C\u1000-\u102A\u103F\u1050-\u1055\u105A-\u105D\u1061\u1065\u1066\u106E-\u1070\u1075-\u1081\u108E\u10A0-\u10C5\u10C7\u10CD\u10D0-\u10FA\u10FC-\u1248\u124A-\u124D\u1250-\u1256\u1258\u125A-\u125D\u1260-\u1288\u128A-\u128D\u1290-\u12B0\u12B2-\u12B5\u12B8-\u12BE\u12C0\u12C2-\u12C5\u12C8-\u12D6\u12D8-\u1310\u1312-\u1315\u1318-\u135A\u1380-\u138F\u13A0-\u13F4\u1401-\u166C\u166F-\u167F\u1681-\u169A\u16A0-\u16EA\u1700-\u170C\u170E-\u1711\u1720-\u1731\u1740-\u1751\u1760-\u176C\u176E-\u1770\u1780-\u17B3\u17D7\u17DC\u1820-\u1877\u1880-\u18A8\u18AA\u18B0-\u18F5\u1900-\u191C\u1950-\u196D\u1970-\u1974\u1980-\u19AB\u19C1-\u19C7\u1A00-\u1A16\u1A20-\u1A54\u1AA7\u1B05-\u1B33\u1B45-\u1B4B\u1B83-\u1BA0\u1BAE\u1BAF\u1BBA-\u1BE5\u1C00-\u1C23\u1C4D-\u1C4F\u1C5A-\u1C7D\u1CE9-\u1CEC\u1CEE-\u1CF1\u1CF5\u1CF6\u1D00-\u1DBF\u1E00-\u1F15\u1F18-\u1F1D\u1F20-\u1F45\u1F48-\u1F4D\u1F50-\u1F57\u1F59\u1F5B\u1F5D\u1F5F-\u1F7D\u1F80-\u1FB4\u1FB6-\u1FBC\u1FBE\u1FC2-\u1FC4\u1FC6-\u1FCC\u1FD0-\u1FD3\u1FD6-\u1FDB\u1FE0-\u1FEC\u1FF2-\u1FF4\u1FF6-\u1FFC\u2071\u207F\u2090-\u209C\u2102\u2107\u210A-\u2113\u2115\u2119-\u211D\u2124\u2126\u2128\u212A-\u212D\u212F-\u2139\u213C-\u213F\u2145-\u2149\u214E\u2183\u2184\u2C00-\u2C2E\u2C30-\u2C5E\u2C60-\u2CE4\u2CEB-\u2CEE\u2CF2\u2CF3\u2D00-\u2D25\u2D27\u2D2D\u2D30-\u2D67\u2D6F\u2D80-\u2D96\u2DA0-\u2DA6\u2DA8-\u2DAE\u2DB0-\u2DB6\u2DB8-\u2DBE\u2DC0-\u2DC6\u2DC8-\u2DCE\u2DD0-\u2DD6\u2DD8-\u2DDE\u2E2F\u3005\u3006\u3031-\u3035\u303B\u303C\u3041-\u3096\u309D-\u309F\u30A1-\u30FA\u30FC-\u30FF\u3105-\u312D\u3131-\u318E\u31A0-\u31BA\u31F0-\u31FF\u3400-\u4DB5\u4E00-\u9FCC\uA000-\uA48C\uA4D0-\uA4FD\uA500-\uA60C\uA610-\uA61F\uA62A\uA62B\uA640-\uA66E\uA67F-\uA697\uA6A0-\uA6E5\uA717-\uA71F\uA722-\uA788\uA78B-\uA78E\uA790-\uA793\uA7A0-\uA7AA\uA7F8-\uA801\uA803-\uA805\uA807-\uA80A\uA80C-\uA822\uA840-\uA873\uA882-\uA8B3\uA8F2-\uA8F7\uA8FB\uA90A-\uA925\uA930-\uA946\uA960-\uA97C\uA984-\uA9B2\uA9CF\uAA00-\uAA28\uAA40-\uAA42\uAA44-\uAA4B\uAA60-\uAA76\uAA7A\uAA80-\uAAAF\uAAB1\uAAB5\uAAB6\uAAB9-\uAABD\uAAC0\uAAC2\uAADB-\uAADD\uAAE0-\uAAEA\uAAF2-\uAAF4\uAB01-\uAB06\uAB09-\uAB0E\uAB11-\uAB16\uAB20-\uAB26\uAB28-\uAB2E\uABC0-\uABE2\uAC00-\uD7A3\uD7B0-\uD7C6\uD7CB-\uD7FB\uF900-\uFA6D\uFA70-\uFAD9\uFB00-\uFB06\uFB13-\uFB17\uFB1D\uFB1F-\uFB28\uFB2A-\uFB36\uFB38-\uFB3C\uFB3E\uFB40\uFB41\uFB43\uFB44\uFB46-\uFBB1\uFBD3-\uFD3D\uFD50-\uFD8F\uFD92-\uFDC7\uFDF0-\uFDFB\uFE70-\uFE74\uFE76-\uFEFC\uFF21-\uFF3A\uFF41-\uFF5A\uFF66-\uFFBE\uFFC2-\uFFC7\uFFCA-\uFFCF\uFFD2-\uFFD7\uFFDA-\uFFDC]+)\]','', 'g');
	END;
$$;


ALTER FUNCTION public.noquote(txt text) OWNER TO cartola;

--
-- TOC entry 248 (class 1255 OID 152733)
-- Name: update_full_text_search(); Type: FUNCTION; Schema: public; Owner: cartola
--

CREATE FUNCTION public.update_full_text_search() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    if NEW.deleted = true then
        DELETE FROM full_text_search_topic WHERE topic_id = NEW.id;
        DELETE FROM full_text_search_comment WHERE topic_id = NEW.id;
    else
        UPDATE full_text_search_topic SET
                                          tsv = setweight(to_tsvector(NEW.title), 'A'),
                                          date = NEW.created_at
        WHERE topic_id = NEW.id;
    end if;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_full_text_search() OWNER TO cartola;

--
-- TOC entry 264 (class 1255 OID 159328)
-- Name: update_full_text_search_comment(); Type: FUNCTION; Schema: public; Owner: cartola
--

CREATE FUNCTION public.update_full_text_search_comment() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE full_text_search_comment SET
        tsv = setweight(to_tsvector(noquote(NEW.text), 'A')),
        date = NEW.date
    WHERE topic_id = NEW.topic_id AND comment_id = NEW.id;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_full_text_search_comment() OWNER TO cartola;

--
-- TOC entry 220 (class 1255 OID 19502)
-- Name: update_profile(); Type: FUNCTION; Schema: public; Owner: cartola
--

CREATE FUNCTION public.update_profile() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
  BEGIN
    INSERT INTO profile_names 
     (profile_id, first_name, last_name, screen_name, photo, date) 
    VALUES
      (NEW.id, NEW.first_name, NEW.last_name, NEW.screen_name, NEW.photo, extract(epoch from now()))
    ON CONFLICT DO NOTHING;
    RETURN NEW;
  END;
$$;


ALTER FUNCTION public.update_profile() OWNER TO cartola;

SET default_tablespace = '';

--
-- TOC entry 205 (class 1259 OID 137829)
-- Name: administrators; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.administrators (
    id bigint NOT NULL
);


ALTER TABLE public.administrators OWNER TO cartola;

--
-- TOC entry 198 (class 1259 OID 16386)
-- Name: attachments; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.attachments (
    content text NOT NULL,
    comment_id bigint NOT NULL
);


ALTER TABLE public.attachments OWNER TO cartola;

--
-- TOC entry 199 (class 1259 OID 16392)
-- Name: comments; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.comments (
    id bigint NOT NULL,
    from_id bigint NOT NULL,
    date bigint NOT NULL,
    text text NOT NULL,
    likes integer NOT NULL,
    reply_to_uid bigint NOT NULL,
    reply_to_cid bigint NOT NULL,
    topic_id bigint NOT NULL,
    profile_id bigint NOT NULL
);


ALTER TABLE public.comments OWNER TO cartola;

--
-- TOC entry 210 (class 1259 OID 137865)
-- Name: full_text_search_comment; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.full_text_search_comment (
    topic_id bigint NOT NULL,
    comment_id bigint NOT NULL,
    tsv tsvector NOT NULL,
    date integer NOT NULL
);


ALTER TABLE public.full_text_search_comment OWNER TO cartola;

--
-- TOC entry 209 (class 1259 OID 137856)
-- Name: full_text_search_topic; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.full_text_search_topic (
    topic_id bigint NOT NULL,
    tsv tsvector NOT NULL,
    date bigint NOT NULL
);


ALTER TABLE public.full_text_search_topic OWNER TO cartola;

--
-- TOC entry 200 (class 1259 OID 16398)
-- Name: poll_answers; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.poll_answers (
    id bigint NOT NULL,
    text character varying(500) NOT NULL,
    votes integer NOT NULL,
    rate real NOT NULL,
    poll_id bigint NOT NULL
);


ALTER TABLE public.poll_answers OWNER TO cartola;

--
-- TOC entry 201 (class 1259 OID 16404)
-- Name: polls; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.polls (
    id bigint NOT NULL,
    question character varying(500) NOT NULL,
    votes integer NOT NULL,
    multiple boolean NOT NULL,
    end_date bigint NOT NULL,
    closed boolean NOT NULL,
    topic_id bigint NOT NULL
);


ALTER TABLE public.polls OWNER TO cartola;

--
-- TOC entry 204 (class 1259 OID 19484)
-- Name: profile_names; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.profile_names (
    profile_id bigint NOT NULL,
    first_name character varying(250) NOT NULL,
    last_name character varying(250) NOT NULL,
    screen_name character varying(250) NOT NULL,
    photo character varying(250) NOT NULL,
    date bigint NOT NULL
);


ALTER TABLE public.profile_names OWNER TO cartola;

--
-- TOC entry 202 (class 1259 OID 16410)
-- Name: profiles; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.profiles (
    id bigint NOT NULL,
    first_name character varying(250) NOT NULL,
    last_name character varying(250) NOT NULL,
    screen_name character varying(250) NOT NULL,
    photo character varying(250) NOT NULL
);


ALTER TABLE public.profiles OWNER TO cartola;

--
-- TOC entry 208 (class 1259 OID 137849)
-- Name: settings; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.settings (
    name character varying NOT NULL,
    value text NOT NULL
);


ALTER TABLE public.settings OWNER TO cartola;

--
-- TOC entry 207 (class 1259 OID 137835)
-- Name: topic_update_job; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.topic_update_job (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    topic_id integer NOT NULL,
    priority integer DEFAULT 10 NOT NULL,
    run_after timestamp without time zone NOT NULL,
    retry_waits text[] NOT NULL,
    ran_at timestamp without time zone,
    error text DEFAULT ''::text NOT NULL,
    locked boolean DEFAULT false
);


ALTER TABLE public.topic_update_job OWNER TO cartola;

--
-- TOC entry 206 (class 1259 OID 137833)
-- Name: topic_update_job_id_seq; Type: SEQUENCE; Schema: public; Owner: cartola
--

CREATE SEQUENCE public.topic_update_job_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.topic_update_job_id_seq OWNER TO cartola;

--
-- TOC entry 3069 (class 0 OID 0)
-- Dependencies: 206
-- Name: topic_update_job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: cartola
--

ALTER SEQUENCE public.topic_update_job_id_seq OWNED BY public.topic_update_job.id;


--
-- TOC entry 203 (class 1259 OID 16416)
-- Name: topics; Type: TABLE; Schema: public; Owner: cartola
--

CREATE TABLE public.topics (
    id bigint NOT NULL,
    title character varying(500) NOT NULL,
    is_closed boolean NOT NULL,
    is_fixed boolean NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    created_by bigint NOT NULL,
    updated_by bigint NOT NULL,
    deleted boolean NOT NULL
);


ALTER TABLE public.topics OWNER TO cartola;

--
-- TOC entry 2893 (class 2604 OID 137838)
-- Name: topic_update_job id; Type: DEFAULT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.topic_update_job ALTER COLUMN id SET DEFAULT nextval('public.topic_update_job_id_seq'::regclass);


--
-- TOC entry 2901 (class 2606 OID 19464)
-- Name: comments PK_comments; Type: CONSTRAINT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT "PK_comments" PRIMARY KEY (id);


--
-- TOC entry 2911 (class 2606 OID 19466)
-- Name: polls PK_poll; Type: CONSTRAINT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.polls
    ADD CONSTRAINT "PK_poll" PRIMARY KEY (id);


--
-- TOC entry 2908 (class 2606 OID 19468)
-- Name: poll_answers PK_poll_answers; Type: CONSTRAINT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.poll_answers
    ADD CONSTRAINT "PK_poll_answers" PRIMARY KEY (id);


--
-- TOC entry 2914 (class 2606 OID 19470)
-- Name: profiles PK_profiles; Type: CONSTRAINT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT "PK_profiles" PRIMARY KEY (id);


--
-- TOC entry 2917 (class 2606 OID 19472)
-- Name: topics PK_topics; Type: CONSTRAINT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT "PK_topics" PRIMARY KEY (id);


--
-- TOC entry 2899 (class 2606 OID 19474)
-- Name: attachments attachments_pkey; Type: CONSTRAINT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.attachments
    ADD CONSTRAINT attachments_pkey PRIMARY KEY (comment_id, content);


--
-- TOC entry 2923 (class 2606 OID 19494)
-- Name: profile_names profile_names_pk; Type: CONSTRAINT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.profile_names
    ADD CONSTRAINT profile_names_pk PRIMARY KEY (profile_id, first_name, last_name, screen_name, photo);


--
-- TOC entry 2926 (class 2606 OID 137847)
-- Name: topic_update_job topic_update_job_pkey; Type: CONSTRAINT; Schema: public; Owner: cartola
--

ALTER TABLE ONLY public.topic_update_job
    ADD CONSTRAINT topic_update_job_pkey PRIMARY KEY (id);


--
-- TOC entry 2924 (class 1259 OID 137832)
-- Name: administrators_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE UNIQUE INDEX administrators_id_idx ON public.administrators USING btree (id);


--
-- TOC entry 2902 (class 1259 OID 106341)
-- Name: comments_date_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX comments_date_idx ON public.comments USING btree (date);


--
-- TOC entry 2903 (class 1259 OID 56784)
-- Name: comments_from_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX comments_from_id_idx ON public.comments USING btree (from_id, date);


--
-- TOC entry 2904 (class 1259 OID 106340)
-- Name: comments_likes_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX comments_likes_idx ON public.comments USING btree (likes, from_id);


--
-- TOC entry 2905 (class 1259 OID 106337)
-- Name: comments_topic_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX comments_topic_id_idx ON public.comments USING btree (topic_id);


--
-- TOC entry 2932 (class 1259 OID 137871)
-- Name: full_text_search_comment_date_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX full_text_search_comment_date_idx ON public.full_text_search_comment USING btree (date);


--
-- TOC entry 2933 (class 1259 OID 137872)
-- Name: full_text_search_comment_topic_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE UNIQUE INDEX full_text_search_comment_topic_id_idx ON public.full_text_search_comment USING btree (topic_id, comment_id);


--
-- TOC entry 2929 (class 1259 OID 137862)
-- Name: full_text_search_topic_date_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX full_text_search_topic_date_idx ON public.full_text_search_topic USING btree (date);


--
-- TOC entry 2930 (class 1259 OID 137863)
-- Name: full_text_search_topic_topic_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE UNIQUE INDEX full_text_search_topic_topic_id_idx ON public.full_text_search_topic USING btree (topic_id);


--
-- TOC entry 2934 (class 1259 OID 137873)
-- Name: ids_full_text_c; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX ids_full_text_c ON public.full_text_search_comment USING gin (tsv);


--
-- TOC entry 2931 (class 1259 OID 137864)
-- Name: ids_full_text_t; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX ids_full_text_t ON public.full_text_search_topic USING gin (tsv);


--
-- TOC entry 2909 (class 1259 OID 106344)
-- Name: poll_answers_poll_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX poll_answers_poll_id_idx ON public.poll_answers USING btree (poll_id);


--
-- TOC entry 2912 (class 1259 OID 106343)
-- Name: polls_topic_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX polls_topic_id_idx ON public.polls USING btree (topic_id);


--
-- TOC entry 2915 (class 1259 OID 106342)
-- Name: profiles_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX profiles_id_idx ON public.profiles USING btree (id);


--
-- TOC entry 2928 (class 1259 OID 137855)
-- Name: settings_name_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE UNIQUE INDEX settings_name_idx ON public.settings USING btree (name);


--
-- TOC entry 2906 (class 1259 OID 193850)
-- Name: text_trgm_gin; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX text_trgm_gin ON public.comments USING gin (public.noquote(text) public.gin_trgm_ops);


--
-- TOC entry 2918 (class 1259 OID 193859)
-- Name: title_trgm_gin; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX title_trgm_gin ON public.topics USING gin (title public.gin_trgm_ops);


--
-- TOC entry 2927 (class 1259 OID 137848)
-- Name: topic_update_job_topic_id_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE UNIQUE INDEX topic_update_job_topic_id_idx ON public.topic_update_job USING btree (topic_id, run_after);


--
-- TOC entry 2919 (class 1259 OID 56557)
-- Name: topics_created_at_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX topics_created_at_idx ON public.topics USING btree (created_at);


--
-- TOC entry 2920 (class 1259 OID 56556)
-- Name: topics_created_by_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX topics_created_by_idx ON public.topics USING btree (created_by);


--
-- TOC entry 2921 (class 1259 OID 133831)
-- Name: topics_updated_at_idx; Type: INDEX; Schema: public; Owner: cartola
--

CREATE INDEX topics_updated_at_idx ON public.topics USING btree (updated_at);


--
-- TOC entry 2935 (class 2620 OID 159327)
-- Name: comments insert_comment_trigger; Type: TRIGGER; Schema: public; Owner: cartola
--

CREATE TRIGGER insert_comment_trigger AFTER INSERT ON public.comments FOR EACH ROW EXECUTE PROCEDURE public.insert_full_text_search_comment();


--
-- TOC entry 2937 (class 2620 OID 134083)
-- Name: profiles insert_profile_trigger; Type: TRIGGER; Schema: public; Owner: cartola
--

CREATE TRIGGER insert_profile_trigger AFTER INSERT ON public.profiles FOR EACH ROW EXECUTE PROCEDURE public.insert_profile();


--
-- TOC entry 2939 (class 2620 OID 152654)
-- Name: topics insert_topic_trigger; Type: TRIGGER; Schema: public; Owner: cartola
--

CREATE TRIGGER insert_topic_trigger AFTER INSERT ON public.topics FOR EACH ROW EXECUTE PROCEDURE public.insert_full_text_search();


--
-- TOC entry 2936 (class 2620 OID 159330)
-- Name: comments update_comment_trigger; Type: TRIGGER; Schema: public; Owner: cartola
--

CREATE TRIGGER update_comment_trigger AFTER UPDATE ON public.comments FOR EACH ROW EXECUTE PROCEDURE public.update_full_text_search_comment();


--
-- TOC entry 2938 (class 2620 OID 134084)
-- Name: profiles update_profile_trigger; Type: TRIGGER; Schema: public; Owner: cartola
--

CREATE TRIGGER update_profile_trigger AFTER UPDATE ON public.profiles FOR EACH ROW EXECUTE PROCEDURE public.update_profile();


--
-- TOC entry 2940 (class 2620 OID 152785)
-- Name: topics update_topic_trigger; Type: TRIGGER; Schema: public; Owner: cartola
--

CREATE TRIGGER update_topic_trigger AFTER UPDATE ON public.topics FOR EACH ROW EXECUTE PROCEDURE public.update_full_text_search();


-- Completed on 2020-08-21 21:53:10

--
-- PostgreSQL database dump complete
--

