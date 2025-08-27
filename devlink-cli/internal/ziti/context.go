package ziti

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openziti/sdk-golang/ziti"
	"github.com/spf13/cobra"
)

type AppContext struct {
	ZitiContext ziti.Context
}

type appContextKey struct{}


func AttachAppContext(cmd *cobra.Command) error {
	identityPath, _ := cmd.Flags().GetString("identity")

	if identityPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %v", err)
		}
		identityPath = filepath.Join(home, ".devlink", "identity.json")

		if _, err := os.Stat(identityPath); os.IsNotExist(err) {
			return errors.New("identity file not found. please run 'devlink init <invitation-token>' to configure your environment")
		}
	}

	
	zitiConfig, err := ziti.NewConfigFromFile(identityPath)
	if err != nil {
		return fmt.Errorf("error loading ziti config from file %s: %v", identityPath, err)
	}

	zitiContext, err := ziti.NewContext(zitiConfig)
	if err != nil {
		return fmt.Errorf("error creating ziti context: %v", err)
	}

	AppContext := &AppContext{
		ZitiContext: zitiContext,
	}

	ctx := context.WithValue(cmd.Context(), appContextKey{}, AppContext)

	cmd.SetContext(ctx)

	return nil
}

func AppContextFrom(cmd *cobra.Command) (*AppContext, bool) {
	appCtx, ok := cmd.Context().Value(appContextKey{}).(*AppContext)
	return appCtx, ok
}