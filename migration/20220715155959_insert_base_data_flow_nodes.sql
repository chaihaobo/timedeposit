-- +goose Up
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (1, 'eod_flow', 'start_node', 'StartNode', 'start_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (2, 'eod_flow', 'undo_maturity_node', 'UndoMaturityNode', 'undo_maturity_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (3, 'eod_flow', 'start_new_maturity_node', 'StartNewMaturityNode', 'start_new_maturity_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (4, 'eod_flow', 'apply_profit_node', 'ApplyProfitNode', 'apply_profit_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (5, 'eod_flow', 'withdraw_netprofit_node', 'WithdrawNetprofitNode', 'withdraw_netprofit_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (6, 'eod_flow', 'deposit_netprofit_node', 'DepositNetprofitNode', 'deposit_netprofit_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (7, 'eod_flow', 'withdraw_balance_node', 'WithdrawBalanceNode', 'withdraw_balance_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (12, 'eod_flow', 'patch_account_node', 'PatchAccountNode', 'patch_account_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (13, 'eod_flow', 'close_account_node', 'CloseAccountNode', 'close_account_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (14, 'eod_flow', 'end_node', 'EndNode', 'end_node', '2022-05-30 10:43:44', '2022-05-30 10:43:47');
INSERT INTO t_flow_nodes (id, flow_name, node_name, node_path, node_detail, create_time, update_time) VALUES (16, 'eod_flow', 'additional_profit_node', 'AdditionalProfitNode', 'additional_profit_node', '2022-06-09 02:30:06', '2022-06-09 02:30:06');


-- +goose Down
SELECT 'down SQL query';
