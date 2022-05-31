package mysql

import (
	"time"

	"gorm.io/gorm"

	"github.com/ethan256/books/pkg/log"
)

const (
	callBackBeforeName = "core:before"
	callBackAfterName  = "core:after"
	startTime          = "_start_time"

	spanKey             = "wosai-hera-go-span"
	instrumentationName = "git.wosai-inc.com/middleware/hera-go/instrumentation/otelgorm"
)

type TracePlugin struct{}

var _ gorm.Plugin = &TracePlugin{}

func NewTracePlugin() *TracePlugin {
	return &TracePlugin{}
}

func (op *TracePlugin) Name() string {
	return "gorm:tracePlugin"
}

func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:before_query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:before_update").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:before_row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:before_raw").Register(callBackBeforeName, before)

	// 结束后
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:after_row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:after_raw").Register(callBackAfterName, after)
	return
}

var _ gorm.Plugin = &TracePlugin{}

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
	return
}

func after(db *gorm.DB) {
	if db == nil || db.Statement == nil {
		return
	}

	_ts, _ := db.InstanceGet(startTime)
	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars)
	rows := db.Statement.RowsAffected
	costSeconds := time.Since(ts).Seconds()

	log.Logger.Debug().Str("sqb", sql).Int64("rows", rows).Float64("costSeconds", costSeconds)
}
