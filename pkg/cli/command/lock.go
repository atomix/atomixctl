package command

import (
	"errors"
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/lock"
	"github.com/spf13/cobra"
)

func newLockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock {create,lock,get,unlock,delete}",
		Short: "Manage the state of a distributed lock",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the lock name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_locks"},
	}
	cmd.MarkFlagRequired("name")
	cmd.AddCommand(newLockCreateCommand())
	cmd.AddCommand(newLockLockCommand())
	cmd.AddCommand(newLockGetCommand())
	cmd.AddCommand(newLockUnlockCommand())
	cmd.AddCommand(newLockDeleteCommand())
	return cmd
}

func newLockFromName(cmd *cobra.Command) lock.Lock {
	name, _ := cmd.PersistentFlags().GetString("name")
	if name == "" {
		ExitWithError(ExitBadArgs, errors.New("--name is a required flag"))
	}

	group := newGroupFromName(cmd, name)
	m, err := group.GetLock(newTimeoutContext(cmd), getPrimitiveName(name))
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newLockCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		Run:  runLockCreateCommand,
	}
}

func runLockCreateCommand(cmd *cobra.Command, _ []string) {
	lock := newLockFromName(cmd)
	lock.Close()
	ExitWithOutput(fmt.Sprintf("Created %s", lock.Name().String()))
}

func newLockDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		Run:  runLockDeleteCommand,
	}
}

func runLockDeleteCommand(cmd *cobra.Command, _ []string) {
	lock := newLockFromName(cmd)
	err := lock.Delete()
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(fmt.Sprintf("Deleted %s", lock.Name().String()))
	}
}

func newLockLockCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "lock",
		Args: cobra.NoArgs,
		Run:  runLockLockCommand,
	}
}

func runLockLockCommand(cmd *cobra.Command, _ []string) {
	lock := newLockFromName(cmd)
	version, err := lock.Lock(newTimeoutContext(cmd))
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(version)
	}
}

func newLockGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runLockGetCommand,
	}
	cmd.Flags().Uint64P("version", "v", 0, "the lock version")
	return cmd
}

func runLockGetCommand(cmd *cobra.Command, _ []string) {
	l := newLockFromName(cmd)
	version, _ := cmd.Flags().GetUint64("version")
	if version == 0 {
		locked, err := l.IsLocked(newTimeoutContext(cmd))
		if err != nil {
			ExitWithError(ExitError, err)
		} else {
			ExitWithOutput(locked)
		}
	} else {
		locked, err := l.IsLocked(newTimeoutContext(cmd), lock.WithIsVersion(version))
		if err != nil {
			ExitWithError(ExitError, err)
		} else {
			ExitWithOutput(locked)
		}
	}
}

func newLockUnlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unlock",
		Args: cobra.NoArgs,
		Run:  runLockUnlockCommand,
	}
	cmd.Flags().Uint64P("version", "v", 0, "the lock version")
	return cmd
}

func runLockUnlockCommand(cmd *cobra.Command, _ []string) {
	l := newLockFromName(cmd)
	version, _ := cmd.Flags().GetUint64("version")
	if version == 0 {
		unlocked, err := l.Unlock(newTimeoutContext(cmd))
		if err != nil {
			ExitWithError(ExitError, err)
		} else {
			ExitWithOutput(unlocked)
		}
	} else {
		unlocked, err := l.Unlock(newTimeoutContext(cmd), lock.WithVersion(version))
		if err != nil {
			ExitWithError(ExitError, err)
		} else {
			ExitWithOutput(unlocked)
		}
	}
}
