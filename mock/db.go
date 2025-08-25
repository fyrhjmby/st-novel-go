package mock

import "sync"

var (
	once sync.Once

	UserStore    UserDataStore
	NovelStore   NovelDataStore
	AIStore      AIDataStore
	SettingStore SettingDataStore
)

// InitDatabase 初始化所有模块的内存数据库
// 使用 sync.Once 确保在程序生命周期中只执行一次
func InitDatabase() {
	once.Do(func() {
		initUserData()
		initNovelData()
		initAIData()
		initSettingData()
	})
}
