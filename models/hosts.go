package models

// HostGroup 表示一个hosts配置组
type HostGroup struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	IsActive    bool     `json:"isActive"`
	IsRemote    bool     `json:"isRemote"`
	RemoteURL   string   `json:"remoteUrl,omitempty"`
	LastUpdated string   `json:"lastUpdated"` // 使用字符串表示时间
	CreatedAt   string   `json:"createdAt"`   // 使用字符串表示时间
	Tags        []string `json:"tags,omitempty"`
	Parent      string   `json:"parent,omitempty"` // TODO 父组ID，用于实现分组功能
}

// SystemConfig 应用配置
type SystemConfig struct {
	Theme        string `json:"theme"`        // 主题 light/dark
	AutoBackup   bool   `json:"autoBackup"`   // 自动备份
	BackupFolder string `json:"backupFolder"` // 备份文件夹
	UpdateFreq   int    `json:"updateFreq"`   // 远程hosts更新频率(分钟)
	HostsMode    string `json:"hostsMode"`    // hosts模式 append/exclusive
}

// HostsFile 系统hosts文件信息
type HostsFile struct {
	Path    string `json:"path"`    // hosts文件路径
	Content string `json:"content"` // 当前hosts文件内容
}
