package lazy_save

import (
	"sync"
	"time"

	"github.com/cr-mao/loric/log"
)

// 延迟保存记录字典,
// key = GetLsoId(), val = lazySaveRecord
var lazySaveRecordMap = &sync.Map{}

//
// 模块初始化时就开始保存
func init() {
	startSave()
}

//
// SaveOrUpdate 保存或更新
func SaveOrUpdate(lazySaveObj LazySaveObj) {
	if nil == lazySaveObj {
		return
	}

	log.Infof(
		"记录延迟保存数据, lsoId = %s",
		lazySaveObj.GetLsoId(),
	)

	currRecord, _ := lazySaveRecordMap.Load(lazySaveObj.GetLsoId())

	if nil != currRecord {
		// 修改最后更新时间
		currRecord.(*lazySaveRecord).lastUpdateTime = time.Now().UnixMilli()
		currRecord.(*lazySaveRecord).UpdateNums += 1
		return
	}

	newRecord := &lazySaveRecord{
		objRef:         lazySaveObj,
		lastUpdateTime: time.Now().UnixMilli(),
		UpdateNums:     1,
	}

	lazySaveRecordMap.Store(lazySaveObj.GetLsoId(), newRecord)
}

// 开始保存
func startSave() {
	go func() {
		for {
			// 先休息 1 秒
			time.Sleep(time.Second)

			nowTime := time.Now().UnixMilli()
			deleteLsoIdArray := make([]string, 0, 64)

			lazySaveRecordMap.Range(func(_, val interface{}) bool {
				if nil == val {
					return true
				}

				currRecord := val.(*lazySaveRecord)

				if nowTime-currRecord.lastUpdateTime < 10000 && currRecord.UpdateNums < 10 {
					// 如果时间差小于 20 秒
					// 不进行保存
					return true
				}

				log.Infof(
					"执行延迟保存, lsoId = %s",
					currRecord.objRef.GetLsoId(),
				)

				currRecord.objRef.SaveOrUpdate(nil)
				deleteLsoIdArray = append(deleteLsoIdArray, currRecord.objRef.GetLsoId())
				return true
			})

			for _, lsoId := range deleteLsoIdArray {
				lazySaveRecordMap.Delete(lsoId)
			}
		}
	}()
}
