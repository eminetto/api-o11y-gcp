create table user (id varchar(50) not null primary key,email varchar(255),password varchar(255),first_name varchar(100), last_name varchar(100), created_at datetime, updated_at datetime);
insert into user values ('adb8101e-cfe6-4a71-8594-ebc80af3a86d','eminetto@email.com','2672275fe0c456fb671e4f417fb2f9892c7573ba', 'Elton', 'Minetto', datetime(), null);
create table feedback (id varchar(50) not null primary key,email varchar(255),title varchar(255),body text,  created_at datetime, updated_at datetime);
create table vote (id varchar(50)  not null primary key,email varchar(255),talk_name varchar(255), score int,  created_at datetime, updated_at datetime);