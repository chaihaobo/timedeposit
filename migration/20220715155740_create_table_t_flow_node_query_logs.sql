-- +goose Up
-- +goose StatementBegin
create table t_flow_node_query_logs
(
    id          bigint auto_increment
        primary key,
    flow_id     varchar(50)                        not null comment 'flow id',
    node_name   varchar(100)                       not null,
    query_type  varchar(100)                       null,
    data        longtext                           null,
    create_time datetime default CURRENT_TIMESTAMP null,
    update_time datetime                           null
);

-- +goose StatementEnd

-- +goose StatementBegin
create index t_flow_node_query_log_flow_id_index
    on t_flow_node_query_logs (flow_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table t_flow_node_query_logs;
-- +goose StatementEnd
