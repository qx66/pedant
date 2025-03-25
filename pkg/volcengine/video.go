package volcengine

const (
	// 查询视频生成任务信息 API
	getTaskApi = "https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks/{id}"
	
	// 查询内容生成任务列表 API
	getTasksApi = "https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks?page_num={page_num}&page_size={page_size}&filter.status={filter.status}&filter.task_ids={filter.task_ids}&filter.model={filter.model}"
	
	// 取消或删除内容生成任务 (mothed: DELETE)
	deleteTaskApi = "https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks/{id}"
)
