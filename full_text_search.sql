INSERT INTO full_text_search_topic (topic_id, tsv)
SELECT id, setweight(to_tsvector(title), 'A')
FROM topics;

INSERT INTO full_text_search_comment (topic_id, comment_id, tsv)
SELECT topic_id, id, setweight(to_tsvector(text), 'A')
FROM "comments";

-- falta unique index comments, topics
-- falta triggers
-- index CREATE INDEX ix_scenes_tsv ON scenes USING GIN(tsv);

select
    *,
    ts_headline((select title from topics t2 where t2.id = t1.topic_id ), q),
    ts_rank(tsv, q) as rank
from
    full_text_search_topic t1,
    plainto_tsquery('messi mito') q
where
        tsv @@ q;

-- https://stackoverflow.com/questions/41484577/full-text-search-doesnt-find-anything-when-query-has-accents
