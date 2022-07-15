-- +goose Up
-- +goose StatementBegin
create table t_flow_transactions
(
    id                   int unsigned auto_increment
        primary key,
    trans_id             varchar(255)   null,
    flow_id              varchar(100)   null,
    mambu_trans_id       varchar(100)   null,
    terminal_rrn         varchar(255)   null,
    source_account_no    varchar(50)    null,
    source_account_name  varchar(255)   null,
    benefit_account_no   varchar(50)    null,
    benefit_account_name varchar(255)   null,
    amount               double         null,
    channel              varchar(255)   null,
    transaction_type     varchar(50)    null,
    result               int default -1 null comment '1:succeed, 0:failed',
    encoded_key          varchar(255)   null,
    create_time          datetime       null,
    update_time          datetime       null,
    error_msg            varchar(1000)  null
);
-- +goose StatementEnd

-- +goose StatementBegin
create index t_flow_transactions_flow_id_index
    on t_flow_transactions (flow_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table t_flow_transactions;
-- +goose StatementEnd
