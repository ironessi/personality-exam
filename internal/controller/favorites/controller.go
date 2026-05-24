package favorites

import (
	"PersonalityTest/internal/dao"
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

type GetFavoritesReq struct {
	g.Meta     `path:"/" method:"get" tags:"收藏" summary:"获取收藏列表"`
	Page       int    `json:"page" dc:"页码" d:"1"`
	PageSize   int    `json:"pageSize" dc:"每页数量" d:"10"`
	TestTypeId int64  `json:"testTypeId" dc:"测试类型ID（筛选）"`
	Keyword    string `json:"keyword" dc:"搜索关键词"`
	SortBy     string `json:"sortBy" dc:"排序字段" d:"createdTime"`
	SortOrder  string `json:"sortOrder" dc:"排序方向" d:"desc"`
}

type FavoriteListItem struct {
	Id                  int64       `json:"id"`
	TestTypeId          int64       `json:"testTypeId"`
	TestTypeName        string      `json:"testTypeName"`
	TestTypeDescription string      `json:"testTypeDescription"`
	TestTypeIcon        string      `json:"testTypeIcon"`
	ResultId            *int64      `json:"resultId"`
	ResultType          string      `json:"resultType"`
	ResultTitle         string      `json:"resultTitle"`
	UserCount           int         `json:"userCount"`
	QuestionCount       int         `json:"questionCount"`
	EstimatedTime       int         `json:"estimatedTime"`
	FavoritedAt         *gtime.Time `json:"favoritedAt"`
	CompletedAt         *gtime.Time `json:"completedAt"`
}

type GetFavoritesRes struct {
	List       []FavoriteListItem `json:"list"`
	Pagination struct {
		Page      int  `json:"page"`
		PageSize  int  `json:"pageSize"`
		Total     int  `json:"total"`
		TotalPage int  `json:"totalPages"`
		HasNext   bool `json:"hasNext"`
		HasPrev   bool `json:"hasPrev"`
	} `json:"pagination"`
}

type AddFavoriteReq struct {
	g.Meta     `path:"/" method:"post" tags:"收藏" summary:"添加收藏"`
	TestTypeId int64 `json:"testTypeId" v:"required|min:1#测试类型ID不能为空|测试类型ID无效"`
	ResultId   int64 `json:"resultId"`
}

type AddFavoriteRes struct {
	FavoriteId   int64       `json:"favoriteId"`
	TestTypeId   int64       `json:"testTypeId"`
	TestTypeName string      `json:"testTypeName"`
	FavoritedAt  *gtime.Time `json:"favoritedAt"`
}

type RemoveFavoriteReq struct {
	g.Meta     `path:"/{favorite_id}" method:"delete" tags:"收藏" summary:"取消收藏"`
	FavoriteId int64 `json:"favorite_id" v:"required|min:1#收藏ID不能为空|收藏ID无效"`
}

type RemoveFavoriteRes struct {
	FavoriteId   int64  `json:"favoriteId"`
	TestTypeId   int64  `json:"testTypeId"`
	TestTypeName string `json:"testTypeName"`
}

type CheckFavoriteReq struct {
	g.Meta     `path:"/check/{test_type_id}" method:"get" tags:"收藏" summary:"检查收藏状态"`
	TestTypeId int64 `json:"test_type_id" v:"required|min:1#测试类型ID不能为空|测试类型ID无效"`
}

type CheckFavoriteRes struct {
	IsFavorited bool        `json:"isFavorited"`
	FavoriteId  *int64      `json:"favoriteId"`
	FavoritedAt *gtime.Time `json:"favoritedAt"`
}

type GetFavoriteStatisticsReq struct {
	g.Meta `path:"/statistics" method:"get" tags:"收藏" summary:"获取收藏统计"`
}

type FavoriteCategoryStat struct {
	CategoryId   int64  `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	Count        int    `json:"count"`
}

type MostFavoritedTest struct {
	TestTypeId   int64  `json:"testTypeId"`
	TestTypeName string `json:"testTypeName"`
	Count        int    `json:"count"`
}

type GetFavoriteStatisticsRes struct {
	TotalFavorites      int                    `json:"totalFavorites"`
	FavoritesByCategory []FavoriteCategoryStat `json:"favoritesByCategory"`
	RecentFavorites     int                    `json:"recentFavorites"`
	MostFavoritedTest   *MostFavoritedTest     `json:"mostFavoritedTest"`
}

func (c *Controller) GetFavorites(ctx context.Context, req *GetFavoritesReq) (res *GetFavoritesRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	sortColumn := "f.created_time"
	switch req.SortBy {
	case "updatedTime":
		sortColumn = "f.updated_time"
	case "testTypeId":
		sortColumn = "f.test_type_id"
	}
	sortOrder := "DESC"
	if strings.EqualFold(req.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	base := dao.Favorites.Ctx(ctx).
		As("f").
		LeftJoin("test_types tt", "tt.id=f.test_type_id").
		LeftJoin("user_test_results utr", "utr.id=f.result_id").
		Where("f.user_id", userID).
		WhereNull("f.deleted_time")

	if req.TestTypeId > 0 {
		base = base.Where("f.test_type_id", req.TestTypeId)
	}
	if req.Keyword != "" {
		base = base.WhereLike("tt.name", "%"+req.Keyword+"%").WhereOrLike("tt.description", "%"+req.Keyword+"%")
	}

	total, err := base.Clone().Count()
	if err != nil {
		return nil, gerror.New("查询收藏失败")
	}

	records, err := base.Fields(`
		f.id, f.test_type_id, f.result_id, f.created_time, 
		tt.name AS test_type_name, tt.description AS test_type_description, tt.img AS test_type_icon, tt.test_time AS estimated_time,
		utr.result_data, utr.test_date AS completed_at,
		(SELECT COUNT(1) FROM questions q WHERE q.test_type_id = f.test_type_id) AS question_count,
		(SELECT COUNT(1) FROM user_test_results ur2 WHERE ur2.test_type_id = f.test_type_id) AS user_count
	`).
		Page(req.Page, req.PageSize).
		Order(sortColumn + " " + sortOrder).
		All()
	if err != nil {
		return nil, gerror.New("查询收藏失败")
	}

	res = &GetFavoritesRes{List: make([]FavoriteListItem, 0, len(records))}
	for _, r := range records {
		var resultID *int64
		if !r["result_id"].IsNil() {
			v := r["result_id"].Int64()
			resultID = &v
		}
		resultType, resultTitle := extractResultMeta(r["result_data"].String())
		item := FavoriteListItem{
			Id:                  r["id"].Int64(),
			TestTypeId:          r["test_type_id"].Int64(),
			TestTypeName:        r["test_type_name"].String(),
			TestTypeDescription: r["test_type_description"].String(),
			TestTypeIcon:        r["test_type_icon"].String(),
			ResultId:            resultID,
			ResultType:          resultType,
			ResultTitle:         resultTitle,
			UserCount:           r["user_count"].Int(),
			QuestionCount:       r["question_count"].Int(),
			EstimatedTime:       r["estimated_time"].Int(),
			FavoritedAt:         r["created_time"].GTime(),
			CompletedAt:         r["completed_at"].GTime(),
		}
		res.List = append(res.List, item)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}
	res.Pagination.Page = req.Page
	res.Pagination.PageSize = req.PageSize
	res.Pagination.Total = total
	res.Pagination.TotalPage = totalPages
	res.Pagination.HasNext = req.Page < totalPages
	res.Pagination.HasPrev = req.Page > 1
	return res, nil
}

func (c *Controller) AddFavorite(ctx context.Context, req *AddFavoriteReq) (res *AddFavoriteRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	exists, err := dao.TestTypes.Ctx(ctx).Where("id", req.TestTypeId).Count()
	if err != nil {
		return nil, gerror.New("查询测试类型失败")
	}
	if exists == 0 {
		return nil, gerror.New("测试不存在")
	}

	if req.ResultId > 0 {
		resultCount, e := dao.UserTestResults.Ctx(ctx).
			Where("id", req.ResultId).
			Where("user_id", userID).
			Where("test_type_id", req.TestTypeId).
			Count()
		if e != nil {
			return nil, gerror.New("查询测试结果失败")
		}
		if resultCount == 0 {
			ownedResultCount, _ := dao.UserTestResults.Ctx(ctx).Where("id", req.ResultId).Count()
			if ownedResultCount == 0 {
				return nil, gerror.New("测试不存在")
			}
			return nil, gerror.New("无权收藏该测试")
		}
	}

	var existing gdb.Record
	existing, err = dao.Favorites.Ctx(ctx).
		Where("user_id", userID).
		Where("test_type_id", req.TestTypeId).
		One()
	if err != nil {
		return nil, gerror.New("查询收藏失败")
	}

	var favoriteID int64
	favoritedAt := gtime.Now()
	if !existing.IsEmpty() {
		if existing["deleted_time"].IsNil() {
			return nil, gerror.New("您已经收藏过该测试")
		}
		favoriteID = existing["id"].Int64()
		_, err = dao.Favorites.Ctx(ctx).
			Where("id", favoriteID).
			Data(g.Map{
				"result_id":    nil,
				"deleted_time": nil,
				"created_time": favoritedAt,
			}).
			Update()
		if req.ResultId > 0 {
			_, err = dao.Favorites.Ctx(ctx).
				Where("id", favoriteID).
				Data(g.Map{"result_id": req.ResultId}).
				Update()
		}
		if err != nil {
			return nil, gerror.New("收藏失败")
		}
	} else {
		data := g.Map{
			"user_id":      userID,
			"test_type_id": req.TestTypeId,
			"created_time": favoritedAt,
		}
		if req.ResultId > 0 {
			data["result_id"] = req.ResultId
		}
		result, e := dao.Favorites.Ctx(ctx).Data(data).Insert()
		if e != nil {
			return nil, gerror.New("收藏失败")
		}
		favoriteID, _ = result.LastInsertId()
	}

	var testTypeName string
	if v, e := dao.TestTypes.Ctx(ctx).Where("id", req.TestTypeId).Fields("name").Value(); e == nil && v != nil {
		testTypeName = v.String()
	}
	return &AddFavoriteRes{
		FavoriteId:   favoriteID,
		TestTypeId:   req.TestTypeId,
		TestTypeName: testTypeName,
		FavoritedAt:  favoritedAt,
	}, nil
}

func (c *Controller) RemoveFavorite(ctx context.Context, req *RemoveFavoriteReq) (res *RemoveFavoriteRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	record, err := dao.Favorites.Ctx(ctx).Where("id", req.FavoriteId).One()
	if err != nil {
		return nil, gerror.New("查询收藏失败")
	}
	if record.IsEmpty() || !record["deleted_time"].IsNil() {
		return nil, gerror.New("收藏记录不存在")
	}
	if record["user_id"].Int64() != userID {
		return nil, gerror.New("无权操作该收藏记录")
	}

	_, err = dao.Favorites.Ctx(ctx).
		Where("id", req.FavoriteId).
		Data(g.Map{"deleted_time": gtime.Now()}).
		Update()
	if err != nil {
		return nil, gerror.New("取消收藏失败")
	}

	var testTypeName string
	if v, e := dao.TestTypes.Ctx(ctx).Where("id", record["test_type_id"].Int64()).Fields("name").Value(); e == nil && v != nil {
		testTypeName = v.String()
	}
	return &RemoveFavoriteRes{
		FavoriteId:   req.FavoriteId,
		TestTypeId:   record["test_type_id"].Int64(),
		TestTypeName: testTypeName,
	}, nil
}

func (c *Controller) CheckFavorite(ctx context.Context, req *CheckFavoriteReq) (res *CheckFavoriteRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	record, err := dao.Favorites.Ctx(ctx).
		Where("user_id", userID).
		Where("test_type_id", req.TestTypeId).
		WhereNull("deleted_time").
		One()
	if err != nil {
		return nil, gerror.New("查询收藏状态失败")
	}

	if record.IsEmpty() {
		return &CheckFavoriteRes{IsFavorited: false}, nil
	}
	id := record["id"].Int64()
	return &CheckFavoriteRes{
		IsFavorited: true,
		FavoriteId:  &id,
		FavoritedAt: record["created_time"].GTime(),
	}, nil
}

func (c *Controller) GetFavoriteStatistics(ctx context.Context, req *GetFavoriteStatisticsReq) (res *GetFavoriteStatisticsRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	total, err := dao.Favorites.Ctx(ctx).
		Where("user_id", userID).
		WhereNull("deleted_time").
		Count()
	if err != nil {
		return nil, gerror.New("查询收藏统计失败")
	}

	categoryRows, err := dao.Favorites.Ctx(ctx).
		As("f").
		InnerJoin("test_types tt", "tt.id=f.test_type_id").
		LeftJoin("test_categories tc", "tc.id=tt.category_id").
		Where("f.user_id", userID).
		WhereNull("f.deleted_time").
		Fields("tc.id AS category_id, tc.name AS category_name, COUNT(1) AS count").
		Group("tc.id, tc.name").
		Order("count DESC, category_id ASC").
		All()
	if err != nil {
		return nil, gerror.New("查询收藏统计失败")
	}

	recent, err := dao.Favorites.Ctx(ctx).
		Where("user_id", userID).
		WhereNull("deleted_time").
		WhereGTE("created_time", gtime.Now().AddDate(0, 0, -7)).
		Count()
	if err != nil {
		return nil, gerror.New("查询收藏统计失败")
	}

	mostRow, err := dao.Favorites.Ctx(ctx).
		As("f1").
		InnerJoin("favorites f2", "f2.test_type_id=f1.test_type_id AND f2.deleted_time IS NULL").
		InnerJoin("test_types tt", "tt.id=f1.test_type_id").
		Where("f1.user_id", userID).
		WhereNull("f1.deleted_time").
		Fields("tt.id AS test_type_id, tt.name AS test_type_name, COUNT(f2.id) AS count").
		Group("tt.id, tt.name").
		Order("count DESC, test_type_id ASC").
		Limit(1).
		One()
	if err != nil {
		return nil, gerror.New("查询收藏统计失败")
	}

	res = &GetFavoriteStatisticsRes{
		TotalFavorites:      total,
		FavoritesByCategory: make([]FavoriteCategoryStat, 0, len(categoryRows)),
		RecentFavorites:     recent,
	}

	for _, row := range categoryRows {
		res.FavoritesByCategory = append(res.FavoritesByCategory, FavoriteCategoryStat{
			CategoryId:   row["category_id"].Int64(),
			CategoryName: row["category_name"].String(),
			Count:        row["count"].Int(),
		})
	}
	if !mostRow.IsEmpty() {
		res.MostFavoritedTest = &MostFavoritedTest{
			TestTypeId:   mostRow["test_type_id"].Int64(),
			TestTypeName: mostRow["test_type_name"].String(),
			Count:        mostRow["count"].Int(),
		}
	}

	return res, nil
}

func extractResultMeta(resultData string) (resultType, resultTitle string) {
	if resultData == "" {
		return "", ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(resultData), &m); err != nil {
		return "", ""
	}
	if v, ok := m["level"]; ok {
		resultType = gconv.String(v)
	}
	if v, ok := m["title"]; ok {
		resultTitle = gconv.String(v)
	}
	if resultTitle == "" {
		if v, ok := m["result_title"]; ok {
			resultTitle = gconv.String(v)
		}
	}
	if resultType == "" {
		if v, ok := m["result_type"]; ok {
			resultType = gconv.String(v)
		}
	}
	return
}
