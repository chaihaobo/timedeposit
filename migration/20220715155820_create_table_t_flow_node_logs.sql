-- +goose Up
-- +goose StatementBegin
create table t_flow_node_logs
(
    id          int unsigned auto_increment
        primary key,
    account_id  varchar(50)  null,
    flow_id     varchar(50)  null,
    flow_name   varchar(255) null,
    node_name   varchar(255) null,
    node_result varchar(255) null,
    node_msg    varchar(255) null,
    create_time datetime     null,
    update_time datetime     null
);
-- +goose StatementEnd

-- +goose StatementBegin
create index t_flow_node_logs_flow_id_index
    on t_flow_node_logs (flow_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table t_flow_node_logs;
-- +goose StatementEnd
