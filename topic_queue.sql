CREATE TABLE topic_update_job (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  topic_id INTEGER NOT NULL,
  priority INTEGER NOT NULL DEFAULT '10',
  run_after TIMESTAMP  NOT NULL,
  retry_waits TEXT[] NOT NULL,
  ran_at TIMESTAMP,
  error TEXT NOT NULL DEFAULT '',
  locked boolean default '0'
);

CREATE UNIQUE INDEX topic_update_job_topic_id_idx ON public.topic_update_job (topic_id,run_after);
