drop table if exists models;
drop table if exists api_keys;
drop table if exists fine_grained_keys;

create table models
(
    id   integer primary key autoincrement,
    name text -- available model names
);

create table api_keys
(
    id  integer primary key autoincrement,
    key text -- openai api key, e.g. sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
);

create table fine_grained_keys
(
    id           integer primary key autoincrement,
    hash         text,    -- sha2566(fine_grained_key)
    parent_id    integer, -- api_keys.id
    desc         text,    -- part of the fine-grained key, e.g. sk-xx...xxxx
    type         text,    -- whitelist / blacklist
    list         text,    -- json array save model name
    expire       integer, -- unix timestamp
    remain_calls integer  -- the number of remaining calls
);
