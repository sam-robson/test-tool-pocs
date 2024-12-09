import { Config, TestItem } from './types.js';
import { readFileSync } from 'fs';
import { join } from 'path';

export const loadConfig = (rootDir: string): Config => {
  const configPath = join(rootDir, 'config.json');
  const configFile = readFileSync(configPath, 'utf-8');
  const configObject = JSON.parse(configFile);
  const tests: TestItem[] = Object.entries(configObject.tests).map(([path, details]: [string, any]) => ({
    path,
    description: details.description,
    command: details.command,
    type: details.type
  }));

  return { tests };
};