-- +goose Up
-- +goose StatementBegin
create table t_flow_task_infos
(
    id            int unsigned auto_increment
        primary key,
    flow_id       varchar(50)          null,
    account_id    varchar(50)          null,
    flow_name     varchar(255)         not null,
    flow_status   varchar(255)         null,
    cur_node_name varchar(255)         null,
    cur_status    varchar(255)         null,
    end_status    varchar(255)         null,
    start_time    datetime             null,
    end_time      datetime             null,
    enable        tinyint(1) default 1 null,
    create_time   datetime             null,
    update_time   datetime             null
);
-- +goose StatementEnd


-- +goose StatementBegin
create index t_flow_task_infos_account_id_index
    on t_flow_task_infos (account_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table t_flow_task_infos;
-- +goose StatementEnd
