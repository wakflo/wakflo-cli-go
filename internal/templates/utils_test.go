package templates

import (
	"os"
	"testing"
)

func TestWriteTemplateToFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		tmpl      string
		meta      any
		wantError bool
	}{
		{
			name:      "valid template and filePath",
			filePath:  "testfile.txt",
			tmpl:      "Hello, {{.Name}}!",
			meta:      map[string]string{"Name": "World"},
			wantError: false,
		},
		//{
		//	name:      "invalid filePath",
		//	filePath:  "",
		//	tmpl:      "Hello, {{.Name}}!",
		//	meta:      map[string]string{"Name": "World"},
		//	wantError: true,
		//},
		//{
		//	name:      "invalid template structure",
		//	filePath:  "testfile.txt",
		//	tmpl:      "Hello, {{.NonExistent}}!",
		//	meta:      map[string]string{"Name": "World"},
		//	wantError: true,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteTemplateToFile(tt.filePath, tt.tmpl, tt.meta)

			if tt.wantError && err == nil {
				t.Errorf("expected an error, but got none")
			}

			if !tt.wantError && err != nil {
				t.Errorf("did not expect an error, but got: %v", err)
			}

			if !tt.wantError && tt.filePath != "" {
				_ = os.Remove(tt.filePath)
			}
		})
	}
}

func TestIsIntegrationProject(t *testing.T) {
	tests := []struct {
		name      string
		fileSetup func() // optional setup function
		want      bool
	}{
		{
			name: "file exists",
			fileSetup: func() {
				file, _ := os.Create("integration.toml")
				file.Close()
			},
			want: true,
		},
		{
			name:      "file does not exist",
			fileSetup: func() {},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fileSetup != nil {
				tt.fileSetup()
			}

			if got := IsIntegrationProject(); got != tt.want {
				t.Errorf("IsIntegrationProject() = %v, want %v", got, tt.want)
			}

			_ = os.Remove("integration.toml")
		})
	}
}
