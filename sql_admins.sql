CREATE TABLE public.administrators (
	id varchar NOT NULL
);
CREATE UNIQUE INDEX administrators_id_idx ON public.administrators (id);


ALTER TABLE public.administrators ALTER COLUMN id TYPE int8 USING id::int8;
