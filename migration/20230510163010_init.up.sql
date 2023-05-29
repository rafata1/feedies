create table source
(
    id         int auto_increment primary key,
    url        varchar(500) not null unique,
    logo_url   varchar(3000),
    status     int           not null default 2,
    created_at datetime      not null default current_timestamp,
    updated_at datetime      not null default current_timestamp on update current_timestamp
);

create table feed_item
(
    id          int auto_increment primary key,
    source_id   int           not null,
    title       varchar(1000) not null,
    description mediumtext null,
    content     mediumtext null,
    url         varchar(500) not null unique,
    image_url   varchar(3000) null,
    categories  varchar(500) null,
    status      int                    default 2 null,
    created_at  datetime      not null default current_timestamp,
    updated_at  datetime      not null default current_timestamp on update current_timestamp
);
