package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/assert"

	"github.com/hanzeingithub/wgorm"
	"github.com/hanzeingithub/wgorm/log"
)

func InitClient() wgorm.Executor {
	dsn := fmt.Sprintf("root:@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "127.0.0.1:3306", "my_learning")
	client, err := wgorm.NewExecutor(dsn)
	if err != nil {
		panic(err)
	}
	return client
}

func TestFind(t *testing.T) {
	as := assert.New(t)
	client := InitClient()
	ctx := context.Background()
	ctx, _, _ = client.Begin(ctx)
	defer func() {
		log.Logger.InfoOf(ctx, "rollback")
		ctx, err := client.CommitORRollback(ctx, nil, recover())
		if err != nil {
			log.Logger.WarnOf(ctx, "commit error", err)
		}
	}()
	whereInfo := &WhereFileInfo{
		FileName: thrift.StringPtr("kafka_2.12-2.7.0.tgz"),
	}
	result := make([]*FileInfo, 0)
	err := client.Search(ctx, &result, whereInfo)
	fmt.Println(len(result))
	as.Nil(err)

	err = client.MGet(ctx, &result, []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(len(result))

	theOne := &FileInfo{}
	err = client.MustGet(ctx, theOne, 30)
	fmt.Println(theOne)

	update := &UpdateFileInfo{FileName: thrift.StringPtr("kafka2.tgz")}
	//_, err = client.Update(ctx,&FileInfo{},update,whereInfo)
	//as.Nil(err)

	_, err = client.UpdateByID(ctx, &FileInfo{}, update, 10)
	as.Nil(err)
}

func TestTx(t *testing.T) {
	client := InitClient()
	ctx := context.Background()
	ctx, _, _ = client.Begin(ctx)
	defer func() {
		ctx, err := client.CommitORRollback(ctx, nil, recover())
		if err != nil {
			log.Logger.WarnOf(ctx, "commit error", err)
		}
	}()
	something(ctx, client)
}

func something(ctx context.Context, client wgorm.Executor) {
	ctx, _, _ = client.Begin(ctx)
	defer func() {
		ctx, err := client.CommitORRollback(ctx, nil, recover())
		if err != nil {
			log.Logger.WarnOf(ctx, "commit error", err)
		}
	}()

}
