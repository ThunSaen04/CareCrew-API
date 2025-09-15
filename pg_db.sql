--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5 (Ubuntu 17.5-1.pgdg24.10+1)
-- Dumped by pg_dump version 17.5

-- Started on 2025-09-16 00:47:20

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 238 (class 1259 OID 24735)
-- Name: FCM_Tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."FCM_Tokens" (
    id integer NOT NULL,
    personnel_id integer NOT NULL,
    token text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public."FCM_Tokens" OWNER TO postgres;

--
-- TOC entry 239 (class 1259 OID 24758)
-- Name: FCM_Tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."FCM_Tokens" ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."FCM_Tokens_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 235 (class 1259 OID 16524)
-- Name: Reports; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Reports" (
    report_id integer NOT NULL,
    detail text,
    location text,
    created_at timestamp with time zone DEFAULT now(),
    file text,
    title text,
    personnel_id integer NOT NULL
);


ALTER TABLE public."Reports" OWNER TO postgres;

--
-- TOC entry 234 (class 1259 OID 16523)
-- Name: Guest_Reports_report_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Reports" ALTER COLUMN report_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Guest_Reports_report_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 217 (class 1259 OID 16389)
-- Name: Personnels; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Personnels" (
    personnel_id integer NOT NULL,
    first_name text,
    last_name text,
    phone text,
    role_type_id integer
);


ALTER TABLE public."Personnels" OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 16394)
-- Name: Personnels_personnel_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Personnels" ALTER COLUMN personnel_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Personnels_personnel_id_seq"
    START WITH 2056000
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 227 (class 1259 OID 16439)
-- Name: Priority_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Priority_types" (
    priority_type_id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public."Priority_types" OWNER TO postgres;

--
-- TOC entry 226 (class 1259 OID 16438)
-- Name: Priority_types_priority_type_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Priority_types" ALTER COLUMN priority_type_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Priority_types_priority_type_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 219 (class 1259 OID 16395)
-- Name: Role_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Role_types" (
    role_type_id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public."Role_types" OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 16400)
-- Name: Role_types_role_type_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Role_types" ALTER COLUMN role_type_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Role_types_role_type_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 225 (class 1259 OID 16431)
-- Name: Status_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Status_types" (
    status_type_id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public."Status_types" OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 16430)
-- Name: Status_types_status_type_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Status_types" ALTER COLUMN status_type_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Status_types_status_type_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 221 (class 1259 OID 16401)
-- Name: Task_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Task_types" (
    task_type_id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public."Task_types" OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 16406)
-- Name: Task_types_task_type_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Task_types" ALTER COLUMN task_type_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Task_types_task_type_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 229 (class 1259 OID 16447)
-- Name: Tasks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Tasks" (
    task_id integer NOT NULL,
    task_type_id integer DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    status_type_id integer DEFAULT 3,
    priority_type_id integer DEFAULT 3,
    completed boolean DEFAULT false NOT NULL,
    completed_at timestamp with time zone
);


ALTER TABLE public."Tasks" OWNER TO postgres;

--
-- TOC entry 237 (class 1259 OID 16538)
-- Name: Tasks_assignment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Tasks_assignment" (
    assignment_id integer NOT NULL,
    personnel_id integer NOT NULL,
    task_id integer NOT NULL,
    accecp_at timestamp with time zone DEFAULT now(),
    submit boolean DEFAULT false NOT NULL,
    submit_at timestamp with time zone
);


ALTER TABLE public."Tasks_assignment" OWNER TO postgres;

--
-- TOC entry 236 (class 1259 OID 16537)
-- Name: Tasks_assignment_Assignment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Tasks_assignment" ALTER COLUMN assignment_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Tasks_assignment_Assignment_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 231 (class 1259 OID 16492)
-- Name: Tasks_attachments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Tasks_attachments" (
    attachment_id integer NOT NULL,
    assignment_id integer NOT NULL,
    file text,
    uploaded_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public."Tasks_attachments" OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 16491)
-- Name: Tasks_attachments_attachment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Tasks_attachments" ALTER COLUMN attachment_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Tasks_attachments_attachment_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 233 (class 1259 OID 16505)
-- Name: Tasks_detail; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Tasks_detail" (
    task_detail_id integer NOT NULL,
    task_id integer NOT NULL,
    detail text,
    location text,
    people_needed integer DEFAULT 1 NOT NULL,
    assigned_by integer,
    title text
);


ALTER TABLE public."Tasks_detail" OWNER TO postgres;

--
-- TOC entry 232 (class 1259 OID 16504)
-- Name: Tasks_detail_task_detail_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Tasks_detail" ALTER COLUMN task_detail_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Tasks_detail_task_detail_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 228 (class 1259 OID 16446)
-- Name: Tasks_task_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public."Tasks" ALTER COLUMN task_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public."Tasks_task_id_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- TOC entry 223 (class 1259 OID 16407)
-- Name: Users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Users" (
    personnel_id integer NOT NULL,
    password text NOT NULL,
    last_login timestamp with time zone
);


ALTER TABLE public."Users" OWNER TO postgres;

--
-- TOC entry 3347 (class 2606 OID 24750)
-- Name: FCM_Tokens FCM_Tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."FCM_Tokens"
    ADD CONSTRAINT "FCM_Tokens_pkey" PRIMARY KEY (id);


--
-- TOC entry 3343 (class 2606 OID 16531)
-- Name: Reports Guest_Reports_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Reports"
    ADD CONSTRAINT "Guest_Reports_pkey" PRIMARY KEY (report_id);


--
-- TOC entry 3324 (class 2606 OID 16413)
-- Name: Personnels Personnels_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Personnels"
    ADD CONSTRAINT "Personnels_pkey" PRIMARY KEY (personnel_id);


--
-- TOC entry 3334 (class 2606 OID 16445)
-- Name: Priority_types Priority_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Priority_types"
    ADD CONSTRAINT "Priority_types_pkey" PRIMARY KEY (priority_type_id);


--
-- TOC entry 3326 (class 2606 OID 16415)
-- Name: Role_types Role_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Role_types"
    ADD CONSTRAINT "Role_types_pkey" PRIMARY KEY (role_type_id);


--
-- TOC entry 3332 (class 2606 OID 16437)
-- Name: Status_types Status_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Status_types"
    ADD CONSTRAINT "Status_types_pkey" PRIMARY KEY (status_type_id);


--
-- TOC entry 3328 (class 2606 OID 16417)
-- Name: Task_types Task_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Task_types"
    ADD CONSTRAINT "Task_types_pkey" PRIMARY KEY (task_type_id);


--
-- TOC entry 3338 (class 2606 OID 16498)
-- Name: Tasks_attachments Tasks_attachments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks_attachments"
    ADD CONSTRAINT "Tasks_attachments_pkey" PRIMARY KEY (attachment_id);


--
-- TOC entry 3340 (class 2606 OID 16512)
-- Name: Tasks_detail Tasks_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks_detail"
    ADD CONSTRAINT "Tasks_detail_pkey" PRIMARY KEY (task_detail_id);


--
-- TOC entry 3336 (class 2606 OID 16454)
-- Name: Tasks Tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks"
    ADD CONSTRAINT "Tasks_pkey" PRIMARY KEY (task_id);


--
-- TOC entry 3330 (class 2606 OID 16419)
-- Name: Users Users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "Users_pkey" PRIMARY KEY (personnel_id);


--
-- TOC entry 3345 (class 2606 OID 16542)
-- Name: Tasks_assignment assignment_id_PK; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks_assignment"
    ADD CONSTRAINT "assignment_id_PK" PRIMARY KEY (assignment_id) INCLUDE (assignment_id);


--
-- TOC entry 3349 (class 2606 OID 24752)
-- Name: FCM_Tokens unique_token; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."FCM_Tokens"
    ADD CONSTRAINT unique_token UNIQUE (token);


--
-- TOC entry 3341 (class 1259 OID 16558)
-- Name: fki_assigned_by_FK; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "fki_assigned_by_FK" ON public."Tasks_detail" USING btree (assigned_by);


--
-- TOC entry 3356 (class 2606 OID 16553)
-- Name: Tasks_detail assigned_by_FK; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks_detail"
    ADD CONSTRAINT "assigned_by_FK" FOREIGN KEY (assigned_by) REFERENCES public."Personnels"(personnel_id) ON DELETE CASCADE;


--
-- TOC entry 3355 (class 2606 OID 24730)
-- Name: Tasks_attachments assignment_id_FK; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks_attachments"
    ADD CONSTRAINT "assignment_id_FK" FOREIGN KEY (assignment_id) REFERENCES public."Tasks_assignment"(assignment_id) ON DELETE CASCADE;


--
-- TOC entry 3351 (class 2606 OID 16420)
-- Name: Users personnel_id_FK; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Users"
    ADD CONSTRAINT "personnel_id_FK" FOREIGN KEY (personnel_id) REFERENCES public."Personnels"(personnel_id) ON DELETE CASCADE;


--
-- TOC entry 3361 (class 2606 OID 24753)
-- Name: FCM_Tokens personnel_id_FK; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."FCM_Tokens"
    ADD CONSTRAINT "personnel_id_FK" FOREIGN KEY (personnel_id) REFERENCES public."Personnels"(personnel_id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3358 (class 2606 OID 32956)
-- Name: Reports personnel_id_FK; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Reports"
    ADD CONSTRAINT "personnel_id_FK" FOREIGN KEY (personnel_id) REFERENCES public."Personnels"(personnel_id) ON DELETE CASCADE;


--
-- TOC entry 3359 (class 2606 OID 16543)
-- Name: Tasks_assignment personnel_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks_assignment"
    ADD CONSTRAINT personnel_id_fk FOREIGN KEY (personnel_id) REFERENCES public."Personnels"(personnel_id) ON DELETE CASCADE;


--
-- TOC entry 3352 (class 2606 OID 16465)
-- Name: Tasks priority_type_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks"
    ADD CONSTRAINT priority_type_id FOREIGN KEY (priority_type_id) REFERENCES public."Priority_types"(priority_type_id) ON DELETE SET NULL;


--
-- TOC entry 3350 (class 2606 OID 16425)
-- Name: Personnels role_type_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Personnels"
    ADD CONSTRAINT role_type_id_fk FOREIGN KEY (role_type_id) REFERENCES public."Role_types"(role_type_id);


--
-- TOC entry 3353 (class 2606 OID 16460)
-- Name: Tasks status_type_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks"
    ADD CONSTRAINT status_type_id FOREIGN KEY (status_type_id) REFERENCES public."Status_types"(status_type_id) ON DELETE SET NULL;


--
-- TOC entry 3357 (class 2606 OID 16513)
-- Name: Tasks_detail task_id_FK; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks_detail"
    ADD CONSTRAINT "task_id_FK" FOREIGN KEY (task_id) REFERENCES public."Tasks"(task_id) ON DELETE CASCADE;


--
-- TOC entry 3360 (class 2606 OID 16548)
-- Name: Tasks_assignment task_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks_assignment"
    ADD CONSTRAINT task_id_fk FOREIGN KEY (task_id) REFERENCES public."Tasks"(task_id) ON DELETE CASCADE;


--
-- TOC entry 3354 (class 2606 OID 16455)
-- Name: Tasks task_type_id_FK; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Tasks"
    ADD CONSTRAINT "task_type_id_FK" FOREIGN KEY (task_type_id) REFERENCES public."Task_types"(task_type_id) ON DELETE SET NULL;


-- Completed on 2025-09-16 00:47:21

--
-- PostgreSQL database dump complete
--

