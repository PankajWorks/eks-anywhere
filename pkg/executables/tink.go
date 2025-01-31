package executables

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tinkerbell/tink/protos/hardware"
)

const (
	tinkPath                   = "tink"
	TinkerbellCertUrlKey       = "TINKERBELL_CERT_URL"
	TinkerbellGrpcAuthorityKey = "TINKERBELL_GRPC_AUTHORITY"
)

type Tink struct {
	Executable
	tinkerbellCertUrl       string
	tinkerbellGrpcAuthority string
	envMap                  map[string]string
}

func NewTink(executable Executable, tinkerbellCertUrl, tinkerbellGrpcAuthority string) *Tink {
	return &Tink{
		Executable:              executable,
		tinkerbellCertUrl:       tinkerbellCertUrl,
		tinkerbellGrpcAuthority: tinkerbellGrpcAuthority,
		envMap: map[string]string{
			TinkerbellCertUrlKey:       tinkerbellCertUrl,
			TinkerbellGrpcAuthorityKey: tinkerbellGrpcAuthority,
		},
	}
}

func (t *Tink) PushHardware(ctx context.Context, hardware []byte) error {
	params := []string{"hardware", "push"}
	if _, err := t.Command(ctx, params...).WithStdIn(hardware).WithEnvVars(t.envMap).Run(); err != nil {
		return fmt.Errorf("error pushing hardware: %v", err)
	}
	return nil
}

func (t *Tink) GetHardware(ctx context.Context) ([]*hardware.Hardware, error) {
	params := []string{"hardware", "get", "--format", "json"}
	data, err := t.Command(ctx, params...).WithEnvVars(t.envMap).Run()
	if err != nil {
		return nil, fmt.Errorf("error getting hardware list: %v", err)
	}
	var hardwareList []*hardware.Hardware
	hardwareString := data.String()

	if len(hardwareString) > 0 {
		hardwareListData := map[string][]*hardware.Hardware{}

		if err = json.Unmarshal([]byte(hardwareString), &hardwareListData); err != nil {
			return nil, fmt.Errorf("error unmarshling hardware json: %v", err)
		}
		if len(hardwareListData["data"]) > 0 {
			hardwareList = append(hardwareList, hardwareListData["data"]...)
		}
	}

	return hardwareList, nil
}
