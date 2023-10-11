package lazy_save

// 延迟保存记录
type lazySaveRecord struct {
	// 延迟保存对象
	objRef LazySaveObj
	// 最后修改时间   	// 20秒内有更新则不进行保存
	lastUpdateTime int64
	// 更新次数 执行update次数   大于一定次数则进行保存
	UpdateNums int32
}
