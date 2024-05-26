package cmd

import (
	"log"
	"os"

	"github.com/antham/wo/cmd/internal/completion"
	"github.com/antham/wo/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "wo",
		Short: "Manage workspace in shell",
	}
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	w, err := newWorkspaceManager()
	if err != nil {
		log.Fatal(err)
	}

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

	rootCmd.AddCommand(newCreateCmd(w, wksCompMgr))
	rootCmd.AddCommand(newEditCmd(w, envCompMgr))
	rootCmd.AddCommand(newListCmd(w))
	rootCmd.AddCommand(newLoadCmd(w, wksCompMgr))
	rootCmd.AddCommand(newRemoveCmd(w, wksCompMgr))
	rootCmd.AddCommand(newRunCmd(w, funcCompMgr))
	rootCmd.AddCommand(newShowCmd(w, wksCompMgr))
	rootCmd.AddCommand(configCmd)
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
