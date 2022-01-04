package wgorm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/apache/thrift/lib/go/thrift"
	"gorm.io/gorm"

	"github.com/hanzeingithub/wgorm/log"
)

type whereID struct {
	Int64ID *int64  `gorm:"column:id" sql_field:"id"`
	StrID   *string `gorm:"column:id" sql_field:"id"`
}

func parseID(id interface{}) *whereID {
	idType := reflect.TypeOf(id)
	if idType.Kind() == reflect.Ptr {
		idType = idType.Elem()
	}
	switch idType.Kind() {
	case reflect.String:
		return &whereID{StrID: thrift.StringPtr(id.(string))}
	case reflect.Int:
		return &whereID{Int64ID: thrift.Int64Ptr(int64(id.(int)))}
	case reflect.Int8:
		return &whereID{Int64ID: thrift.Int64Ptr(int64(id.(int8)))}
	case reflect.Int16:
		return &whereID{Int64ID: thrift.Int64Ptr(int64(id.(int16)))}
	case reflect.Int32:
		return &whereID{Int64ID: thrift.Int64Ptr(int64(id.(int32)))}
	case reflect.Int64:
		return &whereID{Int64ID: thrift.Int64Ptr(id.(int64))}
	}
	return nil
}

func (e *execution) Create(ctx context.Context, value interface{}) error {
	db := e.DB(ctx)

	if db.Error != nil {
		return db.Error
	}

	if err := db.Create(value).Error; err != nil {
		log.Logger.ErrorOf(ctx, "[wgorm] create err %w", err)
		return err
	}
	return nil
}

func (e *execution) BatchCreate(ctx context.Context, value interface{}) error {
	return e.Create(ctx, value)
}

func (e *execution) MustGet(ctx context.Context, result interface{}, id interface{}, opts ...Option) (err error) {
	if err = judgeReturnSingle(result); err != nil {
		log.Logger.ErrorOf(ctx, "[wgorm] target type error:%w", err)
		return err
	}

	db := e.DB(ctx)
	if db.Error != nil {
		return db.Error
	}

	if err = db.First(result, id).Error; err != nil {
		log.Logger.ErrorOf(ctx, "[wgorm] get error:%w", err)
		return err
	}

	return nil

}

func (e *execution) MGet(ctx context.Context, result interface{}, ids []interface{}, opts ...Option) (err error) {
	if err = judgeReturnMulti(result); err != nil {
		log.Logger.ErrorOf(ctx, "[wgorm] target type error:%w", err)
		return err
	}

	db := e.DB(ctx)
	if db.Error != nil {
		return db.Error
	}

	if err = db.Find(result, ids).Error; err != nil {
		log.Logger.ErrorOf(ctx, "[wgorm] get error:%w", err)
		return err
	}
	return nil
}

func (e *execution) Search(ctx context.Context, result interface{}, where interface{}, opts ...Option) (err error) {
	if err = judgeReturnMulti(result); err != nil {
		log.Logger.ErrorOf(ctx, "[wgorm] target type error:%w", err)
		return err
	}

	query, args, err := e.whereBuilder.BuildSQL(where)
	if err != nil {
		return err
	}

	db := e.DB(ctx)
	if db.Error != nil {
		return db.Error
	}

	err = db.Where(query, args...).Find(result).Error // ignore_security_alert
	if err != nil {
		log.Logger.ErrorOf(ctx, "[wgorm] get error:%w", err)
		return err
	}

	return nil
}

func (e *execution) Update(ctx context.Context, value interface{}, update interface{}, where interface{}) (int64, error) {
	db := e.DB(ctx)
	if db.Error != nil {
		return 0, db.Error
	}

	query, args, err := e.whereBuilder.BuildSQL(where)
	if err != nil {
		return 0, err
	}
	if len(args) == 0 {
		return 0, fmt.Errorf("[wgorm] update error: where is empty")
	}

	updatesMap, err := e.updateBuilder.BuildSQL(update)
	if err != nil {
		return 0, err
	}
	res := db.Model(value).Where(query, args...).Updates(updatesMap) // ignore_security_alert
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (e *execution) UpdateByID(ctx context.Context, value interface{}, update interface{}, id interface{}) (int64, error) {
	return e.Update(ctx, value, update, parseID(id))
}

func (e *execution) Save(ctx context.Context, value interface{}) (int64, error) {
	if err := judgeReturnSingle(value); err != nil {
		log.Logger.ErrorOf(ctx, "[wgorm] save target type error:%w", err)
		return 0, err
	}

	db := e.DB(ctx)
	if db.Error != nil {
		return 0, db.Error
	}

	res := db.Save(value)
	return res.RowsAffected, res.Error
}

func (e *execution) Delete(ctx context.Context, value interface{}, where interface{}) (int64, error) {
	db := e.DB(ctx)
	if db.Error != nil {
		return 0, db.Error
	}

	query, args, err := e.whereBuilder.BuildSQL(where)
	if err != nil {
		return 0, err
	}
	if len(args) == 0 {
		return 0, fmt.Errorf("[wgorm] delete error: where is empty")
	}

	res := db.Where(query, args...).Delete(value) // ignore_security_alert
	return res.RowsAffected, res.Error
}
func (e *execution) DeleteByID(ctx context.Context, value interface{}, id interface{}) (int64, error) {
	return e.Delete(ctx, value, parseID(id))
}

func (e *execution) DB(ctx context.Context, options ...Option) *gorm.DB {
	transaction := ctx.Value(e.transactionKey)
	if transaction == nil {
		return e.db.WithContext(ctx)
	}
	return transaction.(*gorm.DB).WithContext(ctx)
}

func judgeReturnSingle(result interface{}) error {
	if result == nil {
		return fmt.Errorf("[wgorm] target can't be nil")
	}
	rt := reflect.TypeOf(result)
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("[wgorm] result must be a pointer, but got %v", rt.Kind())
	}
	rt = rt.Elem()
	if rt.Kind() != reflect.Struct {
		return fmt.Errorf("[wgorm] result must be a pointer to struct, but got %v", rt.Kind())
	}
	return nil
}

func judgeReturnMulti(result interface{}) error {
	if result == nil {
		return fmt.Errorf("[wgorm] target can't be nil")
	}
	rt := reflect.TypeOf(result)
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("[wgorm] result must be a pointer, but got %v", rt.Kind())
	}
	rt = rt.Elem()
	if rt.Kind() != reflect.Slice {
		return fmt.Errorf("[wgorm] result must be a slice, but got %v", rt.Kind())
	}
	return nil
}
