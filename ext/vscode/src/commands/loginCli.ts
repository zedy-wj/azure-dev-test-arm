// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

import { IActionContext } from '@microsoft/vscode-azext-utils';
import { createAzureDevCli } from '../utils/azureDevCli';
import { executeAsTask } from '../utils/executeAsTask';
import { getAzDevTerminalTitle } from './cmdUtil';

export async function loginCli(context: IActionContext, shouldPrompt: boolean = true): Promise<void> {
    const azureCli = await createAzureDevCli(context);
    const command = azureCli.commandBuilder.withArg('login');
    await executeAsTask(command.build(), getAzDevTerminalTitle(), { alwaysRunNew: true, focus: true });
}
