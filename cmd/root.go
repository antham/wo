package cmd

import (
	"log"
	"log/slog"
	"os"

	"github.com/antham/wo/cmd/internal/completion"
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "wo",
		Short: "Manage workspaces in shell",
	}

	err := viper.BindEnv("WO_DEBUG")
	if err != nil {
		log.Fatal(err)
	}
	if viper.GetBool("WO_DEBUG") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	w, err := newWorkspaceManager()
	if err != nil {
		log.Fatal(err)
	}

	dirCompMgr := completion.New(
		w, []completion.Decorator{
			completion.NoOp,
			completion.FindDirs,
		},
	)
	wksCompMgr := completion.New(
		w, []completion.Decorator{
			completion.FindWorkspaces,
		},
	)
	funcCompMgr := completion.New(
		w, []completion.Decorator{
			completion.FindWorkspaces,
			completion.FindFunctions,
		},
	)
	envCompMgr := completion.New(
		w, []completion.Decorator{
			completion.FindWorkspaces,
			completion.FindEnvs,
		},
	)

	configCmd := newConfigCmd()
	configCmd.AddCommand(newSetCmd(w, wksCompMgr))

	envCmd := newEnvCmd()
	envCmd.AddCommand(newCreateEnvCmd(w, wksCompMgr))
	envCmd.AddCommand(newEditEnvCmd(w, envCompMgr))

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(newCreateCmd(w, dirCompMgr))
	rootCmd.AddCommand(newEditCmd(w, wksCompMgr))
	rootCmd.AddCommand(newListCmd(w))
	rootCmd.AddCommand(newLoadCmd(w, envCompMgr))
	rootCmd.AddCommand(newRemoveCmd(w, wksCompMgr))
	rootCmd.AddCommand(newRunCmd(w, funcCompMgr))
	rootCmd.AddCommand(newShowCmd(w, wksCompMgr))
	rootCmd.AddCommand(newCdCmd(w, wksCompMgr))
	rootCmd.AddCommand(newVersionCmd())
	return rootCmd
}

func newWorkspaceManager() (workspaceManager, error) {
	viper.AutomaticEnv()
	options := []func(*workspace.WorkspaceManager){
		workspace.WithEditor(viper.GetString("EDITOR"), viper.GetString("VISUAL")),
		workspace.WithShellPath(viper.GetString("SHELL")),
	}
	viper.SetEnvPrefix("WO")
	if viper.IsSet("CONFIG_PATH") {
		options = append(options, workspace.WithConfigPath(viper.GetString("CONFIG_PATH")))
	}
	wksManager, err := workspace.NewWorkspaceManager(options...)
	if err != nil {
		return nil, err
	}
	return wksManager, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := newRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
