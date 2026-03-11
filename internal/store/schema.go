package store

import "time"

// ServerProfile stores the configuration for a GameVault server connection.
type ServerProfile struct {
	ID          uint      `gorm:"primarykey"`
	DisplayName string    `gorm:"not null"`
	ServerURL   string    `gorm:"not null;uniqueIndex"`
	Username    string    `gorm:"not null"`
	Active      bool      `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GameCache caches game metadata retrieved from a server to reduce API calls.
type GameCache struct {
	ID           uint      `gorm:"primarykey"`
	ServerID     uint      `gorm:"not null;index"`
	GameID       int       `gorm:"not null"`
	Title        string
	MetadataJSON string    `gorm:"type:text"`
	CachedAt     time.Time
}

// Download tracks the state of a game download operation.
type Download struct {
	ID              uint      `gorm:"primarykey"`
	ServerID        uint      `gorm:"not null;index"`
	GameID          int       `gorm:"not null"`
	Status          string    `gorm:"not null;default:'queued'"` // queued|downloading|paused|complete|failed
	BytesDownloaded int64     `gorm:"default:0"`
	TotalBytes      int64     `gorm:"default:0"`
	InstallPath     string
	PartPath        string
	UpdatedAt       time.Time
}

// SavePath maps a game on a server to its local cloud save directory.
type SavePath struct {
	ID           uint      `gorm:"primarykey"`
	ServerID     uint      `gorm:"not null;index"`
	GameID       int       `gorm:"not null"`
	LocalPath    string    `gorm:"not null"`
	LastSyncedAt time.Time
	LastChecksum string
}

// AppSetting stores application-level key/value configuration.
type AppSetting struct {
	Key       string    `gorm:"primarykey"`
	Value     string    `gorm:"type:text"`
	UpdatedAt time.Time
}
