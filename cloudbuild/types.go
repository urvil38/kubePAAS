package cloudbuild

import (
	"time"
)

type customBuild struct {
	Type  string `json:"@type"`
	Build struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Source struct {
			StorageSource struct {
				Bucket string `json:"bucket"`
				Object string `json:"object"`
			} `json:"storageSource"`
		} `json:"source"`
		CreateTime time.Time `json:"createTime"`
		Steps      []struct {
			Name string   `json:"name"`
			Args []string `json:"args"`
		} `json:"steps"`
		Timeout          string   `json:"timeout"`
		Images           []string `json:"images"`
		ProjectID        string   `json:"projectId"`
		LogsBucket       string   `json:"logsBucket"`
		SourceProvenance struct {
			ResolvedStorageSource struct {
				Bucket     string `json:"bucket"`
				Object     string `json:"object"`
				Generation string `json:"generation"`
			} `json:"resolvedStorageSource"`
		} `json:"sourceProvenance"`
		Options struct {
			LogStreamingOption string `json:"logStreamingOption"`
			Logging            string `json:"logging"`
		} `json:"options"`
		LogURL    string `json:"logUrl"`
		Artifacts struct {
			Images []string `json:"images"`
		} `json:"artifacts"`
	} `json:"build"`
}