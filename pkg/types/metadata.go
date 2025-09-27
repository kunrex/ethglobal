package types

type VersionMetaData struct {
	Version    uint32 `json:"version"`
	CommitHash string `json:"commit_hash"`
}
