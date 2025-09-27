package types

type UploadResponse struct {
	Success bool   `json:"success"`
	CID     string `json:"cid"`
	Message string `json:"message"`
	Size    int64  `json:"size"`
}

type DownloadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Size    int64  `json:"size"`
}

type FileInfo struct {
	CID      string `json:"cid"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}
