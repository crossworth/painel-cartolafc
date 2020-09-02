-- topics by day (last year)
select d.date,
       count(t.id)
from (
         select to_char(date_trunc('day', (current_date - offs)), 'YYYY-MM-DD') as date
         from generate_series(0, 365, 1) as offs) d
         left outer join topics t on
        d.date = to_char(date_trunc('day', to_timestamp(t.created_at)), 'YYYY-MM-DD')
group by d.date;

-- topics by day (last year) from user
select d.date,
       count(t.id)
from (
         select to_char(date_trunc('day', (current_date - offs)), 'YYYY-MM-DD') as date
         from generate_series(0, 365, 1) as offs) d
         left outer join topics t on
            d.date = to_char(date_trunc('day', to_timestamp(t.created_at)), 'YYYY-MM-DD') and
            t.created_by = 259592548
group by d.date
order by date;
