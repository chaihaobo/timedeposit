-- +goose Up
-- +goose StatementBegin
-- auto-generated definition
create table t_flow_node_relations
(
    id          int auto_increment comment 'Primary Key'
        primary key,
    flow_name   varchar(63) null comment 'Flow type name',
    node_name   varchar(63) null comment 'Flow node name',
    result_code varchar(63) null comment 'description of node',
    next_node   varchar(63) null comment 'the flow node name of next step',
    create_time datetime    null comment 'Create Time',
    update_time datetime    null comment 'Update Time'
)
    comment 'flow node info';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table t_flow_node_relations;
-- +goose StatementEnd
