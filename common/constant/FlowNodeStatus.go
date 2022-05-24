/*
 * @Author: Hugo
 * @Date: 2022-05-18 04:28:15
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 03:44:34
 */
package constant

const (
	FlowStart    = "flow_start"
	FlowRunning  = "flow_running"
	FlowFailed   = "flow_failed"
	FlowFinished = "flow_finish"
)

type FlowNodeStatus string

const (
	FlowNodeStart   FlowNodeStatus = "node_start"
	FlowNodeSkip    FlowNodeStatus = "node_skip"
	FlowNodeRunning FlowNodeStatus = "node_running"
	FlowNodeFailed  FlowNodeStatus = "node_failed"
	FlowNodeFinish  FlowNodeStatus = "node_finish"
)
