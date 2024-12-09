import { findUp } from 'find-up';

export async function findRepoRoot(): Promise<string> {
  const configPath = await findUp('config.json');
  if (!configPath) {
    throw new Error('config.json not found');
  }
  return configPath.replace('/config.json', '');
}