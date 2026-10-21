package user_test

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/logic/evaluation"
	"PersonalityTest/internal/logic/rag"
	"PersonalityTest/internal/logic/result"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	// 导入 API 定义，用于 AnswerRecord

	v1 "PersonalityTest/api/test_type/v1"
)

type Controller struct {
	result    *result.Result
	evaluator *evaluation.Evaluation
}

func New() *Controller {
	return &Controller{
		result:    result.New(),
		evaluator: evaluation.New(),
	}
}

// ==================== 获取用户测试记录列表 ====================

type GetUserTestsReq struct {
	g.Meta   `path:"/" method:"get" tags:"用户测试" summary:"获取用户测试记录列表"`
	Page     int `json:"page" v:"min:1#页码最小为1" dc:"页码" d:"1"`
	PageSize int `json:"limit" v:"min:1|max:100#每页数量最小为1|每页数量最多100" dc:"每页数量" d:"10"`
}

type TestRecordItem struct {
	Id         int64       `json:"id"`
	TestId     int64       `json:"test_id"`
	TestName   string      `json:"test_name"`
	TotalScore int         `json:"total_score"`
	Level      string      `json:"level"`
	Completed  bool        `json:"completed"`
	CreatedAt  *gtime.Time `json:"created_at"`
}

type GetUserTestsRes struct {
	Tests      []*TestRecordItem `json:"tests"`
	Pagination struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	} `json:"pagination"`
}

func (c *Controller) GetUserTests(ctx context.Context, req *GetUserTestsReq) (res *GetUserTestsRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	// 查询用户测试结果
	model := dao.UserTestResults.Ctx(ctx).
		Where("user_id", userID).
		OrderDesc("test_date")

	// 获取总数
	total, err := model.Count()
	if err != nil {
		return nil, gerror.New("查询失败")
	}

	// 分页查询
	var records []entity.UserTestResults
	err = model.Page(req.Page, req.PageSize).Scan(&records)
	if err != nil {
		return nil, gerror.New("查询失败")
	}

	// 构造响应
	res = &GetUserTestsRes{
		Tests: make([]*TestRecordItem, 0),
	}

	for _, record := range records {
		// 解析结果数据
		var resultData map[string]interface{}
		if err := json.Unmarshal([]byte(record.ResultData), &resultData); err != nil {
			continue
		}

		// 获取测试类型名称
		var testType entity.TestTypes
		dao.TestTypes.Ctx(ctx).Where("id", record.TestTypeId).Scan(&testType)

		item := &TestRecordItem{
			Id:         record.Id,
			TestId:     record.TestTypeId,
			TestName:   testType.Name,
			TotalScore: gconv.Int(resultData["total_score"]),
			Level:      gconv.String(resultData["level"]),
			Completed:  true,
			CreatedAt:  record.TestDate,
		}
		res.Tests = append(res.Tests, item)
	}

	res.Pagination.Page = req.Page
	res.Pagination.Limit = req.PageSize
	res.Pagination.Total = total

	return res, nil
}

// ==================== 获取测试详情 ====================

type GetTestDetailReq struct {
	g.Meta   `path:"/{record_id}" method:"get" tags:"用户测试" summary:"获取测试详情"`
	RecordID int64 `json:"record_id" v:"required|min:1#记录ID必填|记录ID无效"`
}

type GetTestDetailRes struct {
	Id        int64                    `json:"id"`
	TestId    int64                    `json:"test_id"`
	TestName  string                   `json:"test_name"`
	Answers   []map[string]interface{} `json:"answers"`
	Result    map[string]interface{}   `json:"result"`
	Completed bool                     `json:"completed"`
	CreatedAt *gtime.Time              `json:"created_at"`
}

func (c *Controller) GetTestDetail(ctx context.Context, req *GetTestDetailReq) (res *GetTestDetailRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	// 查询记录
	var record entity.UserTestResults
	err = dao.UserTestResults.Ctx(ctx).
		Where("id", req.RecordID).
		Where("user_id", userID).
		Scan(&record)
	if err != nil || record.Id == 0 {
		return nil, gerror.New("记录不存在")
	}

	// 获取测试类型名称
	var testType entity.TestTypes
	dao.TestTypes.Ctx(ctx).Where("id", record.TestTypeId).Scan(&testType)

	// 解析结果数据
	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(record.ResultData), &resultData); err != nil {
		return nil, gerror.New("解析结果失败")
	}

	// 构造响应
	res = &GetTestDetailRes{
		Id:        record.Id,
		TestId:    record.TestTypeId,
		TestName:  testType.Name,
		Result:    resultData,
		Completed: true,
		CreatedAt: record.TestDate,
	}

	// 如果有用户答案数据，也返回
	if answers, ok := resultData["user_answers"]; ok {
		if answerList, ok := answers.([]interface{}); ok {
			res.Answers = make([]map[string]interface{}, len(answerList))
			for i, answer := range answerList {
				if answerMap, ok := answer.(map[string]interface{}); ok {
					res.Answers[i] = map[string]interface{}{
						"question_id": answerMap["question_id"],
						"answer":      answerMap["selected_tag"],
						"score":       answerMap["score"],
					}
				}
			}
		}
	}

	return res, nil
}

// ==================== 保存测试进度 ====================

type SaveTestProgressReq struct {
	g.Meta   `path:"/save" method:"post" tags:"用户测试" summary:"保存测试进度"`
	TestId   int64                    `json:"test_id" v:"required|min:1#测试ID必填|测试ID无效"`
	TestName string                   `json:"test_name"`
	Answers  []map[string]interface{} `json:"answers"`
}

type SaveTestProgressRes struct {
	Message  string `json:"message"`
	RecordId int64  `json:"record_id"`
}

func (c *Controller) SaveTestProgress(ctx context.Context, req *SaveTestProgressReq) (res *SaveTestProgressRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	// 检查测试类型是否存在
	var testType entity.TestTypes
	err = dao.TestTypes.Ctx(ctx).Where("id", req.TestId).Scan(&testType)
	if err != nil || testType.Id == 0 {
		return nil, gerror.New("测试类型不存在")
	}

	// ============ 1. 处理测试记录表 (test_records) ============
	// 查找用户最近的一条进行中的测试记录
	var testRecord entity.TestRecords
	err = dao.TestRecords.Ctx(ctx).
		Where("user_id", userID).
		Where("test_type_id", req.TestId).
		Where("status", "进行中").
		OrderDesc("start_time").
		Limit(1).
		Scan(&testRecord)

	var recordID int64
	if testRecord.Id > 0 {
		// 如果存在，更新更新时间
		recordID = testRecord.Id
		_, err = dao.TestRecords.Ctx(ctx).
			Where("id", recordID).
			Data(g.Map{"updated_time": gtime.Now()}).
			Update()
	} else {
		// 如果不存在，插入新记录
		result, err := dao.TestRecords.Ctx(ctx).Data(do.TestRecords{
			UserId:     userID,
			TestTypeId: req.TestId,
			Status:     "进行中",
			StartTime:  gtime.Now(),
		}).Insert()
		if err != nil {
			return nil, gerror.New("创建进度记录失败")
		}
		recordID, _ = result.LastInsertId()
	}

	// ============ 2. 保存用户答案 (user_answers) ============
	// 遍历答案列表
	for _, answer := range req.Answers {
		// 提取字段
		questionID := gconv.Int64(answer["question_id"])
		if questionID == 0 {
			continue
		}

		// 处理 answer
		answerStr := gconv.String(answer["answer"])
		// 处理 option_id (如果前端传了)
		var optionID interface{}
		if optID, ok := answer["option_id"]; ok {
			if v := gconv.Int64(optID); v > 0 {
				optionID = v
			}
		}

		// 先删除该用户针对该题目的旧答案，避免重复
		_, err := dao.UserAnswers.Ctx(ctx).
			Where("user_id", userID).
			Where("test_type_id", req.TestId).
			Where("question_id", questionID).
			Delete()
		if err != nil {
			g.Log().Warningf(ctx, "清理旧答案失败: %v", err)
			continue
		}

		// 插入新答案
		_, err = dao.UserAnswers.Ctx(ctx).Data(do.UserAnswers{
			UserId:     userID,
			TestTypeId: req.TestId,
			QuestionId: questionID,
			OptionId:   optionID, // 如果没有传，默认为0 or NULL
			Answer:     answerStr,
			// CreatedTime: gtime.Now(), // 自动生成
		}).Insert()
		if err != nil {
			g.Log().Errorf(ctx, "保存答案失败: %v", err)
			// 注意：这里是否报错返回取决于业务要求，一般不阻断整体流程
		}
	}

	res = &SaveTestProgressRes{
		Message:  "进度保存成功",
		RecordId: recordID,
	}
	return res, nil
}

// ==================== 提交测试结果 ====================

type SubmitTestResultReq struct {
	g.Meta   `path:"/submit" method:"post" tags:"用户测试" summary:"提交测试结果"`
	TestId   int64                    `json:"test_id" v:"required|min:1#测试ID必填|测试ID无效"`
	TestName string                   `json:"test_name"`
	Answers  []map[string]interface{} `json:"answers" v:"required#答案必填"`
}

type SubmitTestResultRes struct {
	Message  string                 `json:"message"`
	RecordId int64                  `json:"record_id"`
	Result   map[string]interface{} `json:"result"`
}

func (c *Controller) SubmitTestResult(ctx context.Context, req *SubmitTestResultReq) (res *SubmitTestResultRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	// 检查测试类型是否存在
	var testType entity.TestTypes
	err = dao.TestTypes.Ctx(ctx).Where("id", req.TestId).Scan(&testType)
	if err != nil || testType.Id == 0 {
		return nil, gerror.New("测试类型不存在")
	}

	// ============ 【新增】查找并更新 test_records 表中的进度记录 ============
	// 查找用户最近的一条进行中的测试记录
	var testRecord entity.TestRecords // 第351行：这里声明了 testRecord
	err = dao.TestRecords.Ctx(ctx).
		Where("user_id", userID).
		Where("test_type_id", req.TestId).
		Where("status", "进行中").
		OrderDesc("start_time").
		Limit(1).
		Scan(&testRecord)

	// 如果找到进行中的记录，则更新为已完成
	if testRecord.Id > 0 {
		_, err = dao.TestRecords.Ctx(ctx).
			Where("id", testRecord.Id).
			Data(g.Map{
				"status":       "已完成",
				"end_time":     gtime.Now(),
				"updated_time": gtime.Now(),
			}).
			Update()
		if err != nil {
			g.Log().Warningf(ctx, "更新测试记录状态失败: %v", err)
		} else {
			g.Log().Infof(ctx, "更新测试记录状态成功: 记录ID=%d", testRecord.Id)
		}
	}
	// ============ 查找并更新结束 ============

	// ============ 【优化】提交时同步保存用户答案 (user_answers) ============
	// 无论用户之前保存过没有，这里都做一次保存，确保答案明细完整。
	// 遍历答案列表
	for _, answer := range req.Answers {
		// 提取字段
		questionID := gconv.Int64(answer["question_id"])
		if questionID == 0 {
			continue
		}

		// 处理 answer
		answerStr := gconv.String(answer["answer"])
		// 处理 option_id (如果前端传了)
		var optionID interface{}
		if optID, ok := answer["option_id"]; ok {
			if v := gconv.Int64(optID); v > 0 {
				optionID = v
			}
		}

		// 先删除该用户针对该题目的旧答案，避免重复
		_, err := dao.UserAnswers.Ctx(ctx).
			Where("user_id", userID).
			Where("test_type_id", req.TestId).
			Where("question_id", questionID).
			Delete()
		if err != nil {
			g.Log().Warningf(ctx, "清理旧答案失败: %v", err)
			continue
		}

		// 插入新答案
		_, err = dao.UserAnswers.Ctx(ctx).Data(do.UserAnswers{
			UserId:     userID,
			TestTypeId: req.TestId,
			QuestionId: questionID,
			OptionId:   optionID, // 如果没有传，默认为0 or NULL
			Answer:     answerStr,
			// CreatedTime: gtime.Now(), // 自动生成
		}).Insert()
		if err != nil {
			g.Log().Errorf(ctx, "保存答案失败: %v", err)
			// 注意：这里是否报错返回取决于业务要求，一般不阻断整体流程
		}
	}
	// ============ 保存用户答案结束 ============

	// ============ 将 answers 转换为 v1.AnswerRecord 格式 ============
	var answerRecords []v1.AnswerRecord
	for _, answer := range req.Answers {
		var questionID int
		var answerStr string
		var score int64

		// 处理 question_id
		switch v := answer["question_id"].(type) {
		case float64:
			questionID = int(v)
		case int:
			questionID = v
		case int64:
			questionID = int(v)
		default:
			// 尝试使用 gconv 转换
			questionID = gconv.Int(answer["question_id"])
		}

		// 处理 answer
		if ans, ok := answer["answer"].(string); ok {
			answerStr = ans
		} else {
			// 如果 answer 不是 string，尝试转换
			answerStr = gconv.String(answer["answer"])
		}

		// 处理 score
		switch v := answer["score"].(type) {
		case float64:
			score = int64(v)
		case int:
			score = int64(v)
		case int64:
			score = v
		default:
			// 尝试使用 gconv 转换
			score = gconv.Int64(answer["score"])
		}

		answerRecords = append(answerRecords, v1.AnswerRecord{
			QuestionID: questionID,
			Answer:     answerStr,
			Score:      score,
		})
	}

	// ============ 调用 evaluation 包的评估方法 ============
	resultData, err := c.evaluator.Evaluate(int(req.TestId), answerRecords)
	if err != nil {
		return nil, gerror.Newf("评估失败: %v", err)
	}

	// ============ 处理评估结果 ============
	// 确保结果是 map[string]interface{} 类型
	var resultMap map[string]interface{}

	// 类型断言处理
	switch v := resultData.(type) {
	case map[string]interface{}:
		resultMap = v
	default:
		// 如果不是 map，使用 json 序列化/反序列化来转换
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, gerror.New("结果转换失败")
		}
		if err := json.Unmarshal(jsonData, &resultMap); err != nil {
			// 如果还是失败，创建默认结果
			resultMap = make(map[string]interface{})
			resultMap["raw_result"] = v
		}
	}

	// 添加额外信息到结果中
	resultMap["test_name"] = req.TestName
	resultMap["user_answers"] = req.Answers
	resultMap["test_type_id"] = req.TestId

	// 确保结果包含必要的字段
	ensureResultFields(resultMap, answerRecords)

	// ============ RAG 增强（旁路，失败绝不影响主流程）============
	// ADR-0007：同步最多等 3 秒，超时转异步由 worker 回填。
	// 这里**不能**返回错误——增强是加分项，不是测评结果的一部分。
	pendingRAG := rag.Prepare(ctx, int(req.TestId), req.TestName, resultMap)
	ragAttached := false
	if pendingRAG != nil {
		if ragMap := pendingRAG.TryInline(ctx); ragMap != nil {
			resultMap["rag"] = ragMap
			// 顶层镜像一份 rag_status：前端只看它就知道该不该轮询（§6.3）。
			if status, ok := ragMap["status"].(string); ok {
				resultMap["rag_status"] = status
			}
			ragAttached = true
		} else {
			// 没拿到同步结果 → 下面会入队。先落 pending，用户才知道"稍后会来"，
			// 而不是以为这个测评根本不支持增强。
			resultMap["rag_status"] = "pending"
		}
	}

	// 将结果序列化为 JSON
	resultJSON, err := json.Marshal(resultMap)
	if err != nil {
		return nil, gerror.New("结果序列化失败")
	}

	// 保存到数据库
	dbResult, err := dao.UserTestResults.Ctx(ctx).Data(do.UserTestResults{
		UserId:     userID,
		TestTypeId: req.TestId,
		ResultData: string(resultJSON),
		TestDate:   gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, gerror.New("提交结果失败")
	}

	recordID, _ := dbResult.LastInsertId()

	// 同步路径没拿到结果时入队，交给 worker 用完整 30 秒超时跑完再回填。
	// 入队失败只记日志：用户已经拿到测评结果，队列不是他的问题。
	if pendingRAG != nil && !ragAttached && recordID > 0 {
		if enqueueErr := pendingRAG.Enqueue(ctx, recordID); enqueueErr != nil {
			g.Log().Warningf(ctx, "RAG: 增强任务入队失败（不影响测评结果）: %v", enqueueErr)
			// 队列不可达时结果永远不会来，把 pending 改成 failed，
			// 否则前端会一直轮询一个不存在的任务。
			if markErr := rag.MarkFailed(ctx, recordID); markErr != nil {
				g.Log().Warningf(ctx, "RAG: 标记增强失败状态也失败: %v", markErr)
			}
		}
	}

	// 返回响应
	res = &SubmitTestResultRes{
		Message:  "测试提交成功",
		RecordId: recordID,
		Result:   resultMap,
	}

	g.Log().Infof(ctx, "用户提交测试结果成功: 用户ID=%d, 测试类型=%d, 记录ID=%d", userID, req.TestId, recordID)
	return res, nil
}

// ensureResultFields 确保结果包含必要的字段
func ensureResultFields(result map[string]interface{}, answers []v1.AnswerRecord) {
	// 计算总分
	totalScore := 0
	for _, answer := range answers {
		totalScore += int(answer.Score)
	}

	// 确保必要的字段存在
	if _, ok := result["total_score"]; !ok {
		result["total_score"] = totalScore
	}

	if _, ok := result["level"]; !ok {
		result["level"] = getDefaultLevel(totalScore, len(answers))
	}

	if _, ok := result["title"]; !ok {
		result["title"] = "测试结果"
	}

	if _, ok := result["analysis"]; !ok {
		result["analysis"] = "测试分析完成"
	}

	if _, ok := result["suggestions"]; !ok {
		result["suggestions"] = []string{
			"保持积极的生活态度",
			"定期进行心理健康评估",
			"寻求专业心理咨询帮助",
		}
	}
}

// getDefaultLevel 获取默认等级
func getDefaultLevel(score int, questionCount int) string {
	if questionCount == 0 {
		return "未知"
	}

	average := float64(score) / float64(questionCount)
	if average >= 4 {
		return "优秀"
	} else if average >= 3 {
		return "良好"
	} else if average >= 2 {
		return "一般"
	} else {
		return "需要关注"
	}
}

// ==================== 删除测试记录 ====================

type DeleteTestRecordReq struct {
	g.Meta   `path:"/{record_id}" method:"delete" tags:"用户测试" summary:"删除测试记录"`
	RecordID int64 `json:"record_id" v:"required|min:1#记录ID必填|记录ID无效"`
}

type DeleteTestRecordRes struct {
	Message string `json:"message"`
}

func (c *Controller) DeleteTestRecord(ctx context.Context, req *DeleteTestRecordReq) (res *DeleteTestRecordRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	// 删除记录
	_, err = dao.UserTestResults.Ctx(ctx).
		Where("id", req.RecordID).
		Where("user_id", userID).
		Delete()
	if err != nil {
		return nil, gerror.New("删除失败")
	}

	res = &DeleteTestRecordRes{
		Message: "记录删除成功",
	}
	return res, nil
}
