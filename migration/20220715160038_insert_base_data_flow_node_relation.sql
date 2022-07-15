-- +goose Up
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (1, 'eod_flow', 'start_node', 'success', 'undo_maturity_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (2, 'eod_flow', 'undo_maturity_node', 'success', 'start_new_maturity_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (3, 'eod_flow', 'start_new_maturity_node', 'success', 'apply_profit_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (4, 'eod_flow', 'apply_profit_node', 'success', 'withdraw_netprofit_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (5, 'eod_flow', 'withdraw_netprofit_node', 'success', 'deposit_netprofit_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (6, 'eod_flow', 'deposit_netprofit_node', 'success', 'withdraw_balance_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (7, 'eod_flow', 'withdraw_balance_node', 'success', 'additional_profit_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (12, 'eod_flow', 'patch_account_node', 'success', 'close_account_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (13, 'eod_flow', 'close_account_node', 'success', 'end_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (14, 'eod_flow', 'undo_maturity_node', 'skip', 'start_new_maturity_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (15, 'eod_flow', 'start_new_maturity_node', 'skip', 'apply_profit_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (16, 'eod_flow', 'apply_profit_node', 'skip', 'withdraw_netprofit_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (17, 'eod_flow', 'withdraw_netprofit_node', 'skip', 'deposit_netprofit_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (18, 'eod_flow', 'deposit_netprofit_node', 'skip', 'withdraw_balance_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (19, 'eod_flow', 'withdraw_balance_node', 'skip', 'additional_profit_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (23, 'eod_flow', 'patch_account_node', 'skip', 'close_account_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (24, 'eod_flow', 'close_account_node', 'skip', 'end_node', '2022-05-30 10:50:11', '2022-05-30 10:50:12');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (25, 'eod_flow', 'additional_profit_node', 'success', 'patch_account_node', '2022-06-09 10:36:17', '2022-06-09 10:36:19');
INSERT INTO t_flow_node_relations (id, flow_name, node_name, result_code, next_node, create_time, update_time) VALUES (26, 'eod_flow', 'additional_profit_node', 'skip', 'patch_account_node', '2022-06-09 10:36:17', '2022-06-09 10:36:19');


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
