-- 创建数据库 yarns
CREATE DATABASE IF NOT EXISTS yarns DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;

-- 使用数据库 yarns
USE yarns;

-- 创建 user表
create table user(
    openid varchar(64) primary key,
    nick_name varchar(64),
    avatar_url varchar(128)
);


-- 创建 image表
create table image(
    id int primary key auto_increment,
    filename varchar(128),
    openid varchar(64),
    foreign key(openid) references user(openid)
);

-- 创建 counter表
create table counter(
    id int primary key auto_increment,
    name varchar(128),
    openid varchar(64),
    count int,
    foreign key(openid) references user(openid)
);

-- 创建 my_course 表
drop table if exists my_course;
create table my_course(
    id int primary key auto_increment,
    openid varchar(64),
    course_id varchar(64),
    foreign key(openid) references user(openid)
);
-- 添加index course_id
create index course_id_index on my_course(course_id);

-- 创建 finished 表
create table finished(
    id int primary key auto_increment,
    openid varchar(64),
    object_id varchar(128),
    foreign key(openid) references user(openid)
);
-- 添加index object_id
create index unique_id_index on finished(object_id);


-- 创建 reward 表
create table reward(
    openid varchar(64) primary key,
    crochet_count int,
    knitting_count int,
    lv1_count int,
    lv2_count int,
    lv3_count int,
    share_count int
);

