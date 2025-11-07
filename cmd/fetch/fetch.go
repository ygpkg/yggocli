package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:   "fetch",
	Short: "A CLI tool to fetch and download resources (e.g., forest files) from the server.",
	Long: `fetch is a command-line interface (CLI) tool designed to retrieve 
and download specific resources, such as 'forest' data, from a remote server.

It provides flexibility by allowing configuration via command-line flags or 
a dedicated configuration file (e.g., config.yaml).

Features:
  - Fetches resources by type ('forest', etc.) and specific IDs.
  - Downloads associated files to a specified local directory.
  - Supports authorization via an authentication key.

Usage:

1. Using command-line flags (all parameters provided inline):
   go run . -a "your_auth_key" -o "http://api.server.com" -r "forest" -i 930 -i 931

2. Using a configuration file:
   (First, create config.yaml with all parameters defined)
   go run . -f config.yaml

   If the -f (--config-file) flag is provided, all other flags 
   (like -a, -o, -r, -i, -t) will be ignored, and settings will be read 
   exclusively from the file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initConfig(flags.ConfigFile); err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		if err := flags.Validate(); err != nil {
			fmt.Println(err)
			return
		}

		processor, err := createProcessor(flags.ResourceType)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := fetchResources(processor); err != nil {
			fmt.Println("Error fetching resources:", err)
			return
		}
	},
}

var flags = &FlagSet{}

func init() {

	Cmd.Flags().StringVarP(&flags.ConfigFile, "config-file", "f", "", "Configuration file path (e.g., config.yaml). If provided, all other flags are ignored.")

	// 标准参数
	Cmd.Flags().StringVarP(&flags.Target, "target", "t", fmt.Sprintf("./fetch_download_dir_%v", time.Now().Unix()), "Target directory for downloaded files")
	Cmd.Flags().StringVarP(&flags.Origin, "origin", "o", "http://localhost:8080", "API URL (e.g., https://api.example.com)")
	Cmd.Flags().StringVarP(&flags.ResourceType, "resource_type", "r", "forest", "Resource type (e.g., forest)")
	Cmd.Flags().IntSliceVarP(&flags.ResourceIDs, "resource_ids", "i", []int{}, "Resource IDs (e.g., 1,2,3)")
	Cmd.Flags().StringVarP(&flags.AuthKey, "auth-key", "a", "", "Authorization key for API access")
}

// FlagSet 结构体 (Viper 将映射到这里)
type FlagSet struct {
	ConfigFile   string
	Target       string `mapstructure:"target"`
	Origin       string `mapstructure:"origin"`
	ResourceType string `mapstructure:"resource_type"`
	ResourceIDs  []int  `mapstructure:"resource_ids"`
	AuthKey      string `mapstructure:"auth_key"`
}

// Validate 验证 flags (在 Viper 加载后运行)
func (f *FlagSet) Validate() error {
	switch f.ResourceType {
	case "forest":
	default:
		return fmt.Errorf("invalid resource type: %s, options: [forest]", f.ResourceType)
	}
	if len(f.AuthKey) <= 0 {
		return fmt.Errorf("invalid auth key: %s", f.AuthKey)
	}
	if len(f.ResourceIDs) <= 0 {
		return fmt.Errorf("resource_ids are required")
	}
	return nil
}

// initConfig 使用 Viper 加载配置
func initConfig(configFile string) error {
	// 如果 -f (configFile) 未提供，则 Cobra 已自动填充 flags，无需操作
	if configFile == "" {
		fmt.Println("Using command-line flags...")
		return nil
	}

	fmt.Printf("Loading configuration from %s...\n", configFile)
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(flags); err != nil {
		return fmt.Errorf("error unmarshaling config to flags: %w", err)
	}
	return nil
}

func createProcessor(resourceType string) (Processor, error) {
	switch resourceType {
	case "forest":
		return &ForestProcessor{}, nil
	// case "file":
	// 	return &FileProcessor{}, nil
	default:
		return nil, fmt.Errorf("no processor found for resource type: %s", resourceType)
	}
}

// fetchResources 负责发送请求并将响应交给 Processor 处理
func fetchResources(processor Processor) error {
	requestData := processor.PreProcessData(flags)
	jsonBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}
	bodyReader := bytes.NewBuffer(jsonBody)

	url := flags.Origin + processor.Route()
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", flags.AuthKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("🚀 Sending POST request to %s with body: %s\n", url, jsonBody)

	resp, err := (&http.Client{
		Timeout: 10 * time.Second,
	}).Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code for response: %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	fmt.Println("--- Response Body ---")
	fmt.Println(string(responseBody))
	fmt.Println("---------------------")

	var apiResponse GenericApiResponse
	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return fmt.Errorf("error unmarshalling generic response: %w", err)
	}
	fmt.Printf("Response Code: %d\n", apiResponse.Code)

	if err = processor.Process(apiResponse.Response.Data, flags.Target); err != nil {
		return fmt.Errorf("error processing response: %w", err)
	}
	return nil
}

func downloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching URL %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code for URL %s: %s", url, resp.Status)
	}
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying data to file %s: %w", err)
	}
	fmt.Printf("✅ Successfully downloaded: %s\n", filePath)
	return nil
}
