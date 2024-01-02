create database `pedant` default character set 'utf8mb4';

create user `pedant` IDENTIFIED BY 'pdkasxAHC#3v3lk';
grant all on pedant.* to 'pedant'@'%';
flush privileges;

use pedant;

drop table if exists session;
create table if not exists session
(
    uuid        varchar(50) not null primary key,
    user_uuid   varchar(50) not null comment '用户Uuid',
    name        text,
    create_time bigint
) comment 'session表';


drop table if exists session_context;
create table if not exists session_context
(
    uuid              varchar(50) not null primary key,
    session_uuid      varchar(50) not null comment 'sessionUuid',
    user_content      text comment '用户提交内容',
    assistant_content text comment '助理返回内容',
    prompt_tokens     int default 0 comment '问题tokens数',
    completion_tokens int default 0 comment '回答tokens数',
    total_tokens      int default 0 comment 'tokens总数',
    llm               varchar(100) comment '大模型语言',
    create_time       bigint
) comment 'session上下文表';


drop table if exists multi_modal;
create table if not exists multi_modal
(
    uuid              varchar(50) not null primary key,
    user_uuid         varchar(50) not null comment '用户id',
    user_content      text comment '用户提交内容',
    images            longtext comment '图片base64',
    assistant_content text comment '助理返回内容',
    llm               varchar(100) comment '大模型语言',
    create_time       bigint
) comment '多模态';

drop table if exists image;
create table if not exists image
    (
        uuid              varchar(50) not null primary key,
        user_uuid         varchar(50) not null comment '用户id',
        prompt text comment 'prompt',
        negative_prompt text comment '反向prompt',
        images longtext,
        prompt_tokens int,
        total_tokens int,
        create_time bigint
) comment '图片';

