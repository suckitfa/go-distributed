package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredServices []ServiceName // 服务依赖项
	ServiceUpdateURL string        // 服务更新地址
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
)

// 每一条服务更新
type pathEntry struct {
	Name ServiceName
	URL  string
}

// 服务更新
type patch struct {
	Added   []pathEntry
	Removed []pathEntry
}
