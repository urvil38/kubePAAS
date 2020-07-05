package util

type VersionedConfig interface {
	GetVersion() string
	Upgrade() (VersionedConfig, error)
}
