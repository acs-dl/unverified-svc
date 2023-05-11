-- +migrate Up

create table if not exists users (
    id bigserial primary key,
    module_id text not null,
    username text,
    phone text,
    email text,
    name text,
    submodule text not null,
    module text not null,
    created_at timestamp without time zone default current_timestamp,

    unique (module_id, module, submodule)
);

create index if not exists users_username_idx on users(username);
create index if not exists users_phone_idx on users(phone);
create index if not exists users_email_idx on users(email);
create index if not exists users_module_idx on users(module);
create index if not exists users_moduleid_idx on users(module_id);
create index if not exists users_submodule_idx on users(submodule);


-- +migrate Down

drop index if exists users_submodule_idx;
drop index if exists users_moduleid_idx;
drop index if exists users_module_idx;
drop index if exists users_email_idx;
drop index if exists users_phone_idx;
drop index if exists users_username_idx;


drop table if exists users;