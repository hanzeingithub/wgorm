package build_sqls

import (
	"fmt"
	"testing"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
)

type whereCondition struct {
	MiniProjectId                *int64     `sql_field:"mini_project_id"`                             // 主键id
	MiniProjectUsername          *string    `sql_field:"mini_project_username"`                       // 用户名字
	MiniProjectState             *int8      `sql_field:"mini_project_state"`                          // 任务状态
	MiniProjectBeginTimeAfter    *time.Time `sql_field:"mini_project_begin_time" sql_operator:">="`   // 任务开始时间
	MiniProjectBeginTimeBefore   *time.Time `sql_field:"mini_project_begin_time" sql_operator:"<="`   // 任务开始时间
	MiniProjectExpiredTimeAfter  *time.Time `sql_field:"mini_project_expired_time" sql_operator:">="` // 任务过期时间
	MiniProjectExpiredTimeBefore *time.Time `sql_field:"mini_project_expired_time" sql_operator:"<="` // 任务过期时间
	MiniProjectDelete            *int16     `sql_field:"mini_project_delete"`                         // 逻辑删除
}

func TestBuildSql(t *testing.T) {
	sqlBuilder := WhereBuilder{}
	where := &whereCondition{
		MiniProjectId:       thrift.Int64Ptr(10086),
		MiniProjectUsername: thrift.StringPtr("test"),
	}
	query, _, err := sqlBuilder.BuildSQL(where)
	fmt.Println(query)
	fmt.Println(err)
}
