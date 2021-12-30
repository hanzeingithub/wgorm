package wgorm

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"gouqinggan/wgorm/log"
)

// Begin 的确把事务存在context里面要好很多，至少来说不用频繁调用delete来删除缓存中的东西，但是context会一直往下透传，也就是说这个事务也会一直透传，导致context过大
// 这种写法就不能像gorm那样支持嵌套事务部分commit了，那么只能在重复begin的时候返回相同的事务
func (e *execution) Begin(ctx context.Context) (context.Context, *gorm.DB, error) {
	transaction := ctx.Value(e.transactionKey)
	var tx *gorm.DB

	// 防止嵌套事务，如果嵌套则返回原事务
	if transaction != nil {
		ctx = context.WithValue(ctx, e.reentryTransactionKey, true)
		return ctx, transaction.(*gorm.DB), nil
	} else {
		tx = e.db.Begin()
		if err := tx.Error; err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, e.transactionKey, tx)
	}

	return ctx, tx, nil
}

func (e *execution) CommitORRollback(ctx context.Context, err error, panicRecover interface{}) (context.Context, error) {
	defer func() {
		if panicRecover != nil {
			panic(panicRecover)
		}
	}()

	if isReentry := ctx.Value(e.reentryTransactionKey); isReentry != nil {
		log.Logger.WarnOf(ctx, "reentry transaction, can't commit")
		return ctx, nil
	}

	tx := ctx.Value(e.transactionKey)
	if tx == nil {
		return ctx, fmt.Errorf("commit or rollback transaction error, invaild transaction")
	}

	transaction := tx.(*gorm.DB)
	if err != nil || panicRecover != nil {
		if err = transaction.Rollback().Error; err != nil {
			log.Logger.ErrorOf(ctx, "rollback error", err)
			return ctx, err
		}
	}

	if err = transaction.Commit().Error; err != nil {
		log.Logger.ErrorOf(ctx, "commit error", err)
		return ctx, err
	}

	return nil, err
}
