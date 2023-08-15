# Репликация в Postgres
sudo setfacl -m u:$(id -u):rwx -R pgslave


1. Создаем сеть, запоминаем адрес

docker network create socialnet
docker network inspect socialnet | grep Subnet 

172.28.0.0/16

2. Определяем в docker-compose.yaml прометеус, графану и cadviser для мониторинга постгреса 

2. Поднимаем мастер

cd ~/dev/highload/socialnet/labs/lab03/postgres/p01
создаем docker-compose.yaml, в котором маппим папку с дампом в папку docker-entrypoint-initdb.d для postgres
Также в нем определяем postgres-exporter 
docker-compose up

<!-- docker run -dit -v $PWD/pgmaster/:/var/lib/postgresql/data -e POSTGRES_PASSWORD=pass -p 5432:5432 --restart=unless-stopped --network=pgnet --name=pgmaster postgres -->



3. Меняем postgresql.conf на мастере

nano  postgres01/postgresql.conf

ssl = off
wal_level = replica
max_wal_senders = 4 # expected slave num

4.1 Устанавливаем репликацию в режим асинхронный - режим local(локальная запись WAL)
synchronous_commit = local		# synchronization level;

4. Подключаемся к мастеру и создаем пользователя для репликации

<!-- docker exec -it pgmaster su - postgres -c psql -->
docker exec -it postgres01 psql -U postgres

create role replicator with login replication password 'pass';

6. Добавляем запись в pg_hba.conf с ip с первого шага

host    replication  replicator  172.28.0.0/16  md5

7. Перезапустим мастера

docker restart postgres01

8.  Сделаем бэкап для реплик

docker exec -it postgres01 bash

mkdir /pgslave

pg_basebackup -h pgmaster -D /pgslave -U replicator -v -P --wal-method=stream

9. Копируем директорию себе

docker cp postgres01:/pgslave postgres02

10. Создадим файл, чтобы реплика узнала, что она реплика

touch postgres02/standby.signal

11. Меняем postgresql.conf на реплике

primary_conninfo = 'host=postgres01 port=5432 user=replicator password=pass application_name=postgres02'

12. Запускаем реплику

docker run -dit -v $PWD/pgslave/:/var/lib/postgresql/data -e POSTGRES_PASSWORD=pass -p 15432:5432 --network=pgnet --restart=unless-stopped --name=pgslave postgres


cd ~/dev/highload/socialnet/labs/lab03/postgres/p02
docker-compose up



13. Запустим вторую реплику

docker cp postgres01:/pgslave postgres03

primary_conninfo = 'host=postgres01 port=5432 user=replicator password=pass application_name=postgres03'

touch postgres03/standby.signal





13.1 Уточним статус реплик:
 docker exec -it postgres01 psql -U postgres
 postgres=# select application_name, sync_state from pg_stat_replication;
 application_name | sync_state 
------------------+------------
 postgres02       | async
 postgres03       | async
(2 rows)



14. Включаем синхронную репликацию на мастере:

synchronous_commit = on
synchronous_standby_names = 'FIRST 1 (postgres01, postgres02)'

docker exec -it postgres01 psql -U postgres
select pg_reload_conf();

postgres=# select application_name, sync_state from pg_stat_replication;
 application_name | sync_state 
------------------+------------
 postgres02       | sync
 postgres03       | async
(2 rows)


15. Создадим тестовую таблицу на мастере

docker exec -it postgres01 psql -U postgres


CREATE TABLE IF NOT EXISTS test
(
    id  SERIAL
        CONSTRAINT test_pk
            PRIMARY KEY,
    uid CHAR(36)
);

ALTER TABLE test
    OWNER TO socialnet;
ALTER TABLE

15.1 запускаем нагрузку на запись

lab03/load : docker-compose up

15.2 Убиваем мастер:

docker exec -it postgres01 /bin/bash
killall -9 postgres

В логах контейнера с прикладом загрузки:
pq: the database system is in recovery mode
Last insert ID=125535, RowsCount=108121
Time to take: 5m43.43817654s


15.3 Проверяем что записалось в синх реплику и асинкреплику

docker exec -it postgres02 psql -U postgres
snet=# select id from test
          ORDER BY id DESC
          limit 1;
   id   
--------
 125536
(1 row)


docker exec -it postgres03 psql -U postgres
snet=# select id from test
          ORDER BY id DESC
          limit 1;
   id   
--------
 125536
(1 row)


И синхронная и асинхронная реплики успели записать изменения, хотя приклад вернул последнее значение на 1 меньше. 
Но наиболее вероятное объяснение  следущее: так как  запрос, который пишет в базу идет с retuning_id, то вероятно запрос записал и на полпути в момент возвращения записанного значения постгрес получил kill -9 и не успел вернуть.




16. Запромоутим реплику pgslave

docker stop postgres01

docker exec -it postgres02 psql -U postgres

select * from pg_promote();

synchronous_commit = on
synchronous_standby_names = 'ANY 1 (pgmaster, pgasyncslave)'

17. Подключим вторую реплику к новому мастеру

primary_conninfo = 'host=postgres02 port=5432 user=replicator password=pass application_name=postgres03'


18. Смотрим статус репликации на новом мастере

 select application_name, sync_state from pg_stat_replication;
 application_name | sync_state 
------------------+------------
 postgres03       | quorum
(1 row)

19. Вставляем запись на мастере в таблицу
snet=# INSERT INTO test (uid) VALUES ('10');
INSERT 0 1

select * from test where uid='10';
   id   |                 uid                  
--------+--------------------------------------
 125560 | 10                                  
(1 row)


20. Проверяяем что запись появлась на слейве:
select * from test where uid='10';
   id   |                 uid                  
--------+--------------------------------------
 125560 | 10                                  
(1 row)










18. Восстановим мастер в качестве реплики

touch pgmaster/standby.signal

primary_conninfo = 'host=pgslave port=5432 user=replicator password=pass application_name=pgmaster'


19. Настроим логическую репликацию с текущего мастера (pgslave) на новый сервер

wal_level = logical

docker restart pgslave

20. Создадим публикацию

GRANT CONNECT ON DATABASE postgres TO replicator;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;
create publication pg_pub for table test;

21. Создадим новый сервер для логической репликации

docker run -dit -v $PWD/pgstandalone/:/var/lib/postgresql/data -e POSTGRES_PASSWORD=pass -p 35432:5432 --restart=unless-stopped --network=pgnet --name=pgstandalone postgres

22. Копируем файлы
docker exec -it pgslave su - postgres

pg_dumpall -U postgres -r -h pgslave -f /var/lib/postgresql/roles.dmp
pg_dump -U postgres -Fc -h pgslave -f /var/lib/postgresql/schema.dmp -s postgres


docker cp pgslave:/var/lib/postgresql/roles.dmp .
docker cp roles.dmp pgstandalone:/var/lib/postgresql/roles.dmp
docker cp pgslave:/var/lib/postgresql/schema.dmp .
docker cp schema.dmp pgstandalone:/var/lib/postgresql/schema.dmp


docker exec -it pgstandalone su - postgres
psql -f roles.dmp
pg_restore -d postgres -C schema.dmp

23. Создаем подписку

CREATE SUBSCRIPTION pg_sub CONNECTION 'host=pgslave port=5432 user=replicator password=pass dbname=postgres' PUBLICATION pg_pub;

24. Сделаем конфликт в данных

На sub:
insert into test values(9);

На pub:
insert into test values(9);

В логах видим:
2023-03-27 16:15:02.753 UTC [258] ERROR:  duplicate key value violates unique constraint "test_pkey"
2023-03-27 16:15:02.753 UTC [258] DETAIL:  Key (id)=(9) already exists.
2023-03-28 18:30:42.893 UTC [108] CONTEXT:  processing remote data for replication origin "pg_16395" during message type "INSERT" for replication target relation "public.test" in transaction 739, finished at 0/3026450

25. Исправляем конфликт

select * from pg_subscription;
SELECT pg_replication_origin_advance('pg_16395', '0/3026C28'::pg_lsn); <- message from log + 1
ALTER SUBSCRIPTION pg_sub ENABLE;
