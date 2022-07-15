-- +goose Up
-- +goose StatementBegin
create table t_mambu_request_logs
(
    id            bigint auto_increment
        primary key,
    flow_id       varchar(100)                       null,
    account_id    varchar(100)                       null,
    node_name     varchar(100)                       null,
    type          varchar(50)                        null,
    request_url   varchar(255)                       null,
    request_body  text                               null,
    response_code int                                null,
    response_body text                               null,
    error         text                               null,
    create_time   datetime default CURRENT_TIMESTAMP null,
    update_time   datetime default CURRENT_TIMESTAMP null
);
-- +goose StatementEnd

-- +goose StatementBegin
create index t_mambu_request_log_flow_id_index
    on t_mambu_request_logs (flow_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table t_mambu_request_logs;
-- +goose StatementEnd
