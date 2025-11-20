-- auto-generated definition
create table boss_user
(
    id         uuid    default gen_random_uuid() not null
        primary key,
    is_delete  boolean default false             not null,
    created_by uuid,
    updated_by uuid,
    created_at timestamp with time zone          not null,
    updated_at timestamp with time zone          not null,
    username   varchar                           not null,
    password   varchar                           not null,
    age        integer                           not null
);

comment on column boss_user.id is '主键';

comment on column boss_user.is_delete is '是否删除:false表示不删除';

comment on column boss_user.created_by is '创建人';

comment on column boss_user.updated_by is '更新人';

comment on column boss_user.created_at is '创建时间';

comment on column boss_user.updated_at is '更新时间';

comment on column boss_user.username is '用户名';

comment on column boss_user.password is '密码';

comment on column boss_user.age is '年龄';

alter table boss_user
    owner to postgres;

create index user_is_delete
    on boss_user (is_delete);

