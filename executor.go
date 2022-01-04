package wgorm

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/hanzeingithub/wgorm/build_sqls"
)

type Executor interface {
	Create(ctx context.Context, value interface{}) error
	BatchCreate(ctx context.Context, value interface{}) error

	MustGet(ctx context.Context, result interface{}, id interface{}, opts ...Option) error
	MGet(ctx context.Context, result interface{}, ids []interface{}, opts ...Option) error
	Search(ctx context.Context, result interface{}, where interface{}, opts ...Option) error

	Update(ctx context.Context, value interface{}, update interface{}, where interface{}) (int64, error)
	UpdateByID(ctx context.Context, value interface{}, update interface{}, id interface{}) (int64, error)
	Save(ctx context.Context, value interface{}) (int64, error)

	Delete(ctx context.Context, value interface{}, where interface{}) (int64, error)
	DeleteByID(ctx context.Context, value interface{}, id interface{}) (int64, error)

	// transaction
	Begin(ctx context.Context) (context.Context, *gorm.DB, error)
	CommitORRollback(ctx context.Context, err error, panicRecover interface{}) (context.Context, error)
}

type execution struct {
	db                    *gorm.DB
	whereBuilder          *build_sqls.WhereBuilder
	updateBuilder         *build_sqls.UpdateBuilder
	transactionKey        *string
	reentryTransactionKey *string
}

func NewExecutor(dsn string) (Executor, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	client := execution{
		db:                    db,
		whereBuilder:          &build_sqls.WhereBuilder{},
		updateBuilder:         &build_sqls.UpdateBuilder{},
		transactionKey:        thrift.StringPtr("wgorm_transaction_key"),
		reentryTransactionKey: thrift.StringPtr("wgorm_reentry_transaction_key"),
	}
	return &client, err
}
