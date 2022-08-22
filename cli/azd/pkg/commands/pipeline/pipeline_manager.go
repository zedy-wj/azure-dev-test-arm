// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/azure/azure-dev/cli/azd/pkg/application_context"
	"github.com/azure/azure-dev/cli/azd/pkg/azd_context"
	"github.com/azure/azure-dev/cli/azd/pkg/commands/global_command_options"
	"github.com/azure/azure-dev/cli/azd/pkg/environment"
	"github.com/azure/azure-dev/cli/azd/pkg/input"
	"github.com/azure/azure-dev/cli/azd/pkg/project"
	"github.com/azure/azure-dev/cli/azd/pkg/tools"
)

// pipelineManager takes care of setting up the scm and pipeline.
// The manager allows to use and test scm providers without a cobra command.
type pipelineManager struct {
	scmProvider
	ciProvider
	askOne                       input.Asker
	azdCtx                       *azd_context.AzdContext
	rootOptions                  *global_command_options.GlobalCommandOptions
	pipelineServicePrincipalName string
	pipelineRemoteName           string
	pipelineRoleName             string
}

func (i *pipelineManager) requiredTools() []tools.ExternalTool {
	reqTools := i.scmProvider.requiredTools()
	reqTools = append(reqTools, i.ciProvider.requiredTools()...)
	return reqTools
}

// validateDependencyInjection panic if the manager did not received all the
// mandatory dependencies to work
func validateDependencyInjection(manager *pipelineManager) {
	if manager.azdCtx == nil {
		log.Panic("missing azd context for pipeline manager")
	}
	if manager.scmProvider == nil {
		log.Panic("missing scm provider for pipeline manager")
	}
	if manager.ciProvider == nil {
		log.Panic("missing CI provider for pipeline manager")
	}
}

func (manager *pipelineManager) configure(ctx context.Context) error {

	// check that scm and ci providers are set
	validateDependencyInjection(manager)

	if err := project.EnsureProject(manager.azdCtx.ProjectPath()); err != nil {
		return err
	}

	// check all required tools are installed
	azCli := application_context.GetAzCliFromContext(ctx)
	requiredTools := manager.requiredTools()
	requiredTools = append(requiredTools, azCli)
	if err := tools.EnsureInstalled(ctx, requiredTools...); err != nil {
		return err
	}

	// Read or init env
	_, err := environment.LoadOrInitEnvironment(ctx, &manager.rootOptions.EnvironmentName, manager.azdCtx, manager.askOne)
	if err != nil {
		return fmt.Errorf("loading environment: %w", err)
	}

	return nil
}
