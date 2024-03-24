package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var dbsecret = "{\"password\":\"Tig#fD[XED,)S:G;<.ruAm9\",\"dbname\":\"exitus\",\"engine\":\"postgres\",\"port\":5432,\"host\":\"abc123abc123abc123.abc123abc123.ap-southeast-2.rds.amazonaws.com\",\"username\":\"exitus\"}"

func TestConfig_DBSecret(t *testing.T) {
	assert := require.New(t)

	cfg := &Config{}

	cfg.DbSecrets = dbsecret

	err := cfg.parseDbSecrets()
	assert.Nil(err)

	assert.Equal("postgres://exitus@abc123abc123abc123.abc123abc123.ap-southeast-2.rds.amazonaws.com:5432/exitus?password=Tig%23fD%5BXED%2C%29S%3AG%3B%3C.ruAm9", cfg.PGDatasource)
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		ListenAddr     string
		DefaultTimeout int
		Debug          bool
		Stage          string
		Branch         string
		DbSecrets      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all required values should validate",
			fields: fields{
				Stage:  "dev",
				Branch: "master",
			},
			wantErr: false,
		},
		{
			name: "missing stage should not validate",
			fields: fields{
				Branch: "master",
			},
			wantErr: true,
		},
		{
			name: "missing branch should not validate",
			fields: fields{
				Stage: "dev",
			},
			wantErr: true,
		},
		{
			name: "missing branch should not validate",
			fields: fields{
				Branch:    "master",
				Stage:     "dev",
				DbSecrets: "{\"password\":\"XXXXXXXX\",\"dbname\":\"exitus\",\"engine\":\"postgres\",\"port\":5432,\"host\":\"abc123abc123abc123.abc123abc123.ap-southeast-2.rds.amazonaws.com\",\"username\":\"exitus\"}",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Addr:   tt.fields.ListenAddr,
				Debug:  tt.fields.Debug,
				Stage:  tt.fields.Stage,
				Branch: tt.fields.Branch,
			}
			if err := cfg.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
