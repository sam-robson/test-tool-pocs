import { readFileSync } from 'fs';
import { join } from 'path';
import { findUp } from 'find-up';
import { Config } from './types.js';

export async function findRepoRoot(): Promise<string> {
  const configPath = await findUp('config.json');
  if (!configPath) {
    throw new Error('config.json not found');
  }
  return configPath.replace('/config.json', '');
}

export function loadConfig(rootDir: string): Config {
  const data = readFileSync(join(rootDir, 'config.json'), 'utf8');
  return JSON.parse(data);
}