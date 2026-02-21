package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/zamedic/labradoc-cli/internal/cli"

	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "File operations via the API",
}

var (
	filesStatus     []string
	filesPageSize   int
	filesPageNumber int
)

var filesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List files",
	RunE: func(cmd *cobra.Command, _ []string) error {
		opts, err := resolveAPIConfig()
		if err != nil {
			return err
		}
		if opts.APIKey == "" && opts.Token == "" {
			return fmt.Errorf("missing api token (use --api-token, --token, api_token, or --use-auth-token)")
		}

		query := url.Values{}
		for _, s := range filesStatus {
			if strings.TrimSpace(s) != "" {
				query.Add("status", strings.TrimSpace(s))
			}
		}
		if filesPageSize > 0 {
			query.Set("pageSize", fmt.Sprintf("%d", filesPageSize))
		}
		if filesPageNumber > 0 {
			query.Set("pageNumber", fmt.Sprintf("%d", filesPageNumber))
		}
		path := "/api/user/files"
		if qs := query.Encode(); qs != "" {
			path = path + "?" + qs
		}

		resp, err := cli.DoRequest(cmd.Context(), "GET", path, nil, opts)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("request failed: %s", resp.Status)
		}
		return nil
	},
}

var (
	uploadFilePath string
)

var filesUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if uploadFilePath == "" {
			return fmt.Errorf("missing --file")
		}
		opts, err := resolveAPIConfig()
		if err != nil {
			return err
		}
		if opts.APIKey == "" && opts.Token == "" {
			return fmt.Errorf("missing api token (use --api-token, --token, api_token, or --use-auth-token)")
		}

		f, err := os.Open(uploadFilePath)
		if err != nil {
			return err
		}
		defer f.Close()

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", filepath.Base(uploadFilePath))
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, f); err != nil {
			return err
		}
		if err := writer.Close(); err != nil {
			return err
		}

		opts.Headers = map[string]string{
			"Content-Type": writer.FormDataContentType(),
		}
		resp, err := cli.DoRequest(cmd.Context(), "PUT", "/api/user/files", &body, opts)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("request failed: %s", resp.Status)
		}
		return nil
	},
}

var (
	fileID       string
	filesOutPath string
)

var filesGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get file metadata",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if fileID == "" {
			return fmt.Errorf("missing --id")
		}
		return simpleGet(cmd, fmt.Sprintf("/api/user/files/%s", fileID), "")
	},
}

var filesContentCmd = &cobra.Command{
	Use:   "content",
	Short: "Get file content",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if fileID == "" {
			return fmt.Errorf("missing --id")
		}
		return simpleGet(cmd, fmt.Sprintf("/api/user/files/%s/content", fileID), filesOutPath)
	},
}

var filesOcrCmd = &cobra.Command{
	Use:   "ocr",
	Short: "Get OCR text for a file",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if fileID == "" {
			return fmt.Errorf("missing --id")
		}
		return simpleGet(cmd, fmt.Sprintf("/api/user/files/%s/ocr", fileID), filesOutPath)
	},
}

var filesDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download original file",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if fileID == "" {
			return fmt.Errorf("missing --id")
		}
		out := filesOutPath
		if out == "" {
			out = fileID + ".pdf"
		}
		return simpleGet(cmd, fmt.Sprintf("/api/user/files/%s/download", fileID), out)
	},
}

var (
	questionText string
	questionFile string
)

var filesQuestionCmd = &cobra.Command{
	Use:   "question",
	Short: "Ask a question about a file",
	RunE: func(cmd *cobra.Command, _ []string) error {
		if fileID == "" {
			return fmt.Errorf("missing --id")
		}
		body, err := readQuestionBody()
		if err != nil {
			return err
		}
		if body == nil {
			return fmt.Errorf("missing question body (--question, --body, or --body-file)")
		}
		return simplePost(cmd, fmt.Sprintf("/api/user/files/%s/question", fileID), body, "text/plain", filesOutPath)
	},
}

var filesSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search across files using the agent",
	RunE: func(cmd *cobra.Command, _ []string) error {
		body, err := readQuestionBody()
		if err != nil {
			return err
		}
		if body == nil {
			return fmt.Errorf("missing question body (--question, --body, or --body-file)")
		}
		return simplePost(cmd, "/api/user/files", body, "text/plain", filesOutPath)
	},
}

func init() {
	filesCmd.AddCommand(filesListCmd)
	filesCmd.AddCommand(filesUploadCmd)
	filesCmd.AddCommand(filesGetCmd)
	filesCmd.AddCommand(filesContentCmd)
	filesCmd.AddCommand(filesOcrCmd)
	filesCmd.AddCommand(filesDownloadCmd)
	filesCmd.AddCommand(filesQuestionCmd)
	filesCmd.AddCommand(filesSearchCmd)

	filesListCmd.Flags().StringSliceVar(&filesStatus, "status", nil, "Filter by status (repeatable)")
	filesListCmd.Flags().IntVar(&filesPageSize, "page-size", 0, "Page size")
	filesListCmd.Flags().IntVar(&filesPageNumber, "page-number", 0, "Page number")

	filesUploadCmd.Flags().StringVar(&uploadFilePath, "file", "", "Path to the file to upload")

	filesGetCmd.Flags().StringVar(&fileID, "id", "", "File ID")
	filesContentCmd.Flags().StringVar(&fileID, "id", "", "File ID")
	filesOcrCmd.Flags().StringVar(&fileID, "id", "", "File ID")
	filesDownloadCmd.Flags().StringVar(&fileID, "id", "", "File ID")
	filesQuestionCmd.Flags().StringVar(&fileID, "id", "", "File ID")

	filesContentCmd.Flags().StringVar(&filesOutPath, "out", "", "Write response to file instead of stdout")
	filesOcrCmd.Flags().StringVar(&filesOutPath, "out", "", "Write response to file instead of stdout")
	filesDownloadCmd.Flags().StringVar(&filesOutPath, "out", "", "Write response to file instead of stdout")
	filesQuestionCmd.Flags().StringVar(&filesOutPath, "out", "", "Write response to file instead of stdout")
	filesSearchCmd.Flags().StringVar(&filesOutPath, "out", "", "Write response to file instead of stdout")

	filesQuestionCmd.Flags().StringVar(&questionText, "question", "", "Question text (plain)")
	filesQuestionCmd.Flags().StringVar(&questionText, "body", "", "Question text (plain)")
	filesQuestionCmd.Flags().StringVar(&questionFile, "body-file", "", "Question text file ('-' for stdin)")
	filesSearchCmd.Flags().StringVar(&questionText, "question", "", "Question text (plain)")
	filesSearchCmd.Flags().StringVar(&questionText, "body", "", "Question text (plain)")
	filesSearchCmd.Flags().StringVar(&questionFile, "body-file", "", "Question text file ('-' for stdin)")
}

func readQuestionBody() (io.Reader, error) {
	if questionFile != "" {
		if questionFile == "-" {
			return os.Stdin, nil
		}
		b, err := os.ReadFile(questionFile)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(b), nil
	}
	if questionText != "" {
		return strings.NewReader(questionText), nil
	}
	return nil, nil
}

func simpleGet(cmd *cobra.Command, path string, outPath string) error {
	opts, err := resolveAPIConfig()
	if err != nil {
		return err
	}
	if opts.APIKey == "" && opts.Token == "" {
		return fmt.Errorf("missing api token (use --api-token, --token, api_token, or --use-auth-token)")
	}
	resp, err := cli.DoRequest(cmd.Context(), "GET", path, nil, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return writeResponse(resp, outPath)
}

func simplePost(cmd *cobra.Command, path string, body io.Reader, contentType string, outPath string) error {
	opts, err := resolveAPIConfig()
	if err != nil {
		return err
	}
	if opts.APIKey == "" && opts.Token == "" {
		return fmt.Errorf("missing api token (use --api-token, --token, api_token, or --use-auth-token)")
	}
	opts.Headers = map[string]string{
		"Content-Type": contentType,
	}
	resp, err := cli.DoRequest(cmd.Context(), "POST", path, body, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return writeResponse(resp, outPath)
}

func writeResponse(resp *http.Response, outPath string) error {
	var out io.Writer = os.Stdout
	if outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed: %s", resp.Status)
	}
	return nil
}
