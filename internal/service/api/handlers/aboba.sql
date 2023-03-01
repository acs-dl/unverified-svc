SELECT t.username, t.module, t.created_at, m.name, m.phone, m.email
FROM (
         SELECT username, MAX(created_at) as created_at, array_agg(module) as module
             FROM users
             GROUP BY username
     ) t JOIN (SELECT DISTINCT ON (username) username, name, phone, email FROM users) m ON m.username = t.username
WHERE (t.username ILIKE '%mhr%' OR m.phone ILIKE '%mhr%' OR m.email ILIKE '%mhr%' OR m.name ILIKE '%mhr%')
ORDER BY t.created_at desc;