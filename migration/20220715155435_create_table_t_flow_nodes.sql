-- +goose Up
-- +goose StatementBegin
create table t_flow_nodes
(
    id          int auto_increment comment 'Primary Key'
        primary key,
    flow_name   varchar(63)  null comment 'Flow type name',
    node_name   varchar(63)  null comment 'Flow node name',
    node_path   varchar(255) null,
    node_detail varchar(255) null comment 'description of node',
    create_time datetime     null comment 'Create Time',
    update_time datetime     null comment 'Update Time',
    constraint flow_name_node_name
        unique (flow_name, node_name)
)
    comment 'flow node info' ;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table t_flow_nodes;
-- +goose StatementEnd
