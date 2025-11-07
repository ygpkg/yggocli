package fetch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ForestProcessor 实现了 Processor 接口
type ForestProcessor struct{}

func (p *ForestProcessor) Route() string {
	return "/v3/forest.GetOriginResource"
}

func (p *ForestProcessor) PreProcessData(f *FlagSet) RequestBody {
	return RequestBody{
		Request: struct {
			ResourceIDs  []int  `json:"resource_ids"`
			ResourceType string `json:"resource_type"`
		}{
			ResourceIDs:  f.ResourceIDs,
			ResourceType: f.ResourceType,
		},
	}
}

func (p *ForestProcessor) Process(data []*GenericResource, baseDir string) error {
	type ForestMeta struct {
		Forest   KnownowForest `json:"forest"`
		FileList []*FileItem   `json:"file_list"`
	}

	for _, resource := range data {
		metaBytes, err := json.Marshal(resource.Meta)
		if err != nil {
			fmt.Printf("Error processing meta (marshal): %v\n", err)
			continue
		}

		var forestMeta ForestMeta
		if err := json.Unmarshal(metaBytes, &forestMeta); err != nil {
			fmt.Printf("Error processing meta (unmarshal): %v\n", err)
			continue
		}

		forest := forestMeta.Forest
		fileList := forestMeta.FileList

		forestDirName := fmt.Sprintf("%d_%s", forest.ID, forest.Name)
		fullForestDirPath := filepath.Join(baseDir, forestDirName)

		if err := os.MkdirAll(fullForestDirPath, 0755); err != nil {
			return fmt.Errorf("error creating forest directory %s: %w", fullForestDirPath, err)
		}
		fmt.Printf("📁 Processing Forest: %s\n", fullForestDirPath)

		for _, fileItem := range fileList {
			targetFileName := fmt.Sprintf("%d_%s", fileItem.ID, fileItem.Name)
			fullFilePath := filepath.Join(fullForestDirPath, targetFileName)

			if err := downloadFile(fileItem.PublicUrl, fullFilePath); err != nil {
				fmt.Printf("❌ Download failed for %s: %v\n", fileItem.PublicUrl, err)
			}
		}
	}
	return nil
}
