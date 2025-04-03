package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	AppName    string `env:"APP_NAME"`
	AppVersion string `env:"APP_VERSION"`
}

func TestLoad(t *testing.T) { //nolint:tparallel
	t.Parallel()

	tt := []struct {
		name        string
		setup       func()
		cleanup     func()
		expectError bool
		err         error
	}{
		{
			name: "successful load from env file",
			setup: func() {
				_ = os.WriteFile(envPath, []byte("APP_NAME=test_app\nAPP_VERSION=dev"), 0o644)
			},
			cleanup: func() {
				_ = os.Remove(envPath)
			},
			expectError: false,
			err:         nil,
		},
		{
			name: "stat error on env file",
			setup: func() {
				_ = os.WriteFile(envPath, []byte("APP_NAME=test_app\nAPP_VERSION=dev"), 0o644)
				_ = os.Chmod(envPath, 0o000)
			},
			cleanup: func() {
				_ = os.Remove(envPath)
			},
			expectError: true,
			err:         errLoad,
		},
		{
			name:        "no env file",
			setup:       func() {},
			cleanup:     func() {},
			expectError: true,
			err:         errStat,
		},
	}

	for i := range tt { //nolint:paralleltest
		tc := &tt[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			cfg, err := Load[TestConfig]()

			if tc.expectError {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)
			}

			tc.cleanup()
		})
	}
}
