package cmd

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/antham/wo/internal/cmd/internal/completion"
	"github.com/antham/wo/internal/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var env string

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "wo",
		Short: "Manage workspaces in shell",
	}
	rootCmd.SetOut(os.Stdout)

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
	rootCmd.AddCommand(newAliasCmd(w))
	rootCmd.AddCommand(newCreateCmd(w, dirCompMgr))
	rootCmd.AddCommand(newEditCmd(w, wksCompMgr))
	rootCmd.AddCommand(newListCmd(w))
	rootCmd.AddCommand(newRemoveCmd(w, wksCompMgr))
	rootCmd.AddCommand(newRunCmd(w, funcCompMgr))
	rootCmd.AddCommand(newShowCmd(w, wksCompMgr))
	rootCmd.AddCommand(newVersionCmd())
	return rootCmd
}

func newWorkspaceManager() (workspaceManager, error) {
	editor, hasEditor := os.LookupEnv("EDITOR")
	visual, hasVisual := os.LookupEnv("VISUAL")
	shell, hasShell := os.LookupEnv("SHELL")
	configPath, hasConfigPath := os.LookupEnv("WO_CONFIG_PATH")
	if !hasEditor && !hasVisual {
		return nil, errors.New("missing EDITOR or VISUAL environment variable")
	}
	if !hasShell {
		return nil, errors.New("missing SHELL environment variable")
	}
	options := []func(*workspace.WorkspaceManager){
		workspace.WithEditor(editor, visual),
		workspace.WithShellPath(shell),
	}
	if hasConfigPath {
		options = append(options, workspace.WithConfigPath(configPath))
	}
	return workspace.NewWorkspaceManager(options...)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := newRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
