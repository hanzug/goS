package dao

//
//import (
//	"context"
//	"go.uber.org/zap"
//	"sync"
//	"time"
//
//	"github.com/hanzug/goS/consts"
//	logs "github.com/hanzug/goS/pkg/logger"
//	"github.com/hanzug/goS/repository/mysql/model"
//	"github.com/hanzug/goS/types"
//)
//
//type MySqlDirectUpload struct {
//	ctx     context.Context
//	doneCtx context.Context
//
//	data   []*model.InputData // 数据
//	upData []*model.InputData // 上传的数据
//	wLock  *sync.Mutex
//	upLock *sync.RWMutex
//	task   *types.Task
//
//	done func()
//}
//
//var GobalMysqlDirectUpload *MySqlDirectUpload
//
//func InitMysqlDirectUpload(ctx context.Context) {
//	zap.S().Info(logs.RunFuncName())
//	task := &types.Task{
//		Columns:    []string{"doc_id", "title", "body", "url"},
//		BiTable:    "data",
//		SourceType: consts.DataSourceCSV,
//	}
//	up := NewMySqlDirectUpload(ctx, task)
//	GobalMysqlDirectUpload = up
//}
//
//// NewMySqlDirectUpload 新建一个上传的对象
//func NewMySqlDirectUpload(ctx context.Context, task *types.Task) *MySqlDirectUpload {
//
//	zap.S().Info(logs.RunFuncName())
//
//	ctx, done := context.WithCancel(ctx)
//
//	directUpload := &MySqlDirectUpload{
//		ctx:    ctx,
//		data:   make([]*model.InputData, 0, 1e5),
//		upData: make([]*model.InputData, 0),
//		wLock:  &sync.Mutex{},
//		upLock: &sync.RWMutex{},
//		task:   task,
//		done:   done,
//	}
//
//	go func() {
//		defer func() {
//			if err := recover(); err != nil {
//				zap.S().Errorf("NewMySqlDirectUpload-消费出现错误 :%+v", err)
//			}
//		}()
//		directUpload.consume()
//	}()
//
//	return directUpload
//}
//
//func (d *MySqlDirectUpload) consume() {
//	zap.S().Info(logs.RunFuncName())
//	gapTime := 5 * time.Second
//	for {
//		select {
//		case <-time.After(gapTime):
//			zap.S().Infof("direct upload")
//			_, err := d.StreamUpload()
//			if err != nil {
//				zap.S().Error("err", err)
//			}
//		case <-d.doneCtx.Done(): // when the program end, upload the data what in memory into database
//			_, err := d.StreamUpload()
//			if err != nil {
//				zap.S().Error("err", err)
//			}
//		}
//	}
//}
//
//func (d *MySqlDirectUpload) StreamUpload() (count int, err error) {
//
//	zap.S().Info(logs.RunFuncName())
//	// 写数据库
//	d.wLock.Lock()
//	if len(d.data) == 0 {
//		d.upData = d.data
//	} else {
//		d.upData = append(d.upData, d.data...)
//	}
//	d.data = make([]*model.InputData, 0)
//	count = len(d.upData)
//	d.wLock.Unlock()
//
//	// 开始上报数据
//	d.upLock.Lock()
//	defer d.upLock.Unlock()
//
//	err = NewInputDataDao(d.ctx).BatchCreateInputData(d.upData)
//	if err != nil {
//		zap.S().Error("BatchCreateInputData ", err)
//	}
//
//	// 重制 updata
//	d.wLock.Lock()
//	d.upData = make([]*model.InputData, 0)
//	d.wLock.Unlock()
//
//	return
//}
//
//func (d *MySqlDirectUpload) Finish() {
//	zap.S().Info(logs.RunFuncName())
//	d.done()
//}
//
//func (d *MySqlDirectUpload) Push(records *model.InputData) int {
//	zap.S().Info(logs.RunFuncName())
//	d.wLock.Lock()
//	defer d.wLock.Unlock()
//	d.data = append(d.data, records)
//	zap.S().Infof("direct_upload push bi_table:%s", d.task.BiTable)
//
//	return len(d.data)
//}
