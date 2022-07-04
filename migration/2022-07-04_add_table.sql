-- auto-generated definition
create table t_flow_account_maturity_day
(
    id            bigint auto_increment
        primary key,
    flow_id       varchar(100)                       not null,
    account_id    varchar(100)                       not null,
    maturity_date datetime                           null,
    create_time   datetime default CURRENT_TIMESTAMP null,
    update_time   datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP
);

create index t_account_maturity_day_account_id_index
    on t_flow_account_maturity_day (account_id);

create index t_account_maturity_day_flow_id_index
    on t_flow_account_maturity_day (flow_id);

