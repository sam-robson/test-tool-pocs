#!/usr/bin/env node
import { Command } from 'commander';
import { execa } from 'execa';
import { loadConfig } from './config.js';
import { TestItem } from './types.js';
import { findRepoRoot } from './utils.js';
import { resolve, relative, isAbsolute, dirname } from 'path';
import { cwd } from 'process';

const program = new Command();

program
  .version('1.0.0')
  .description('Test Runner CLI')
  .option('-t, --type <type>', 'Type of tests to run')
  .action(async (options) => {
    const rootDir = await findRepoRoot();
    const currentDir = cwd();
    const config = loadConfig(rootDir);
    const testsToRun = config.tests.filter((test: TestItem) => {
      if (test.type !== options.type) {
        return false;
      }
      const absolutePath = resolve(rootDir, test.path);
      const relativePath = relative(currentDir, absolutePath);
      return !relativePath.startsWith('..') && !isAbsolute(relativePath);
    });
    console.log(`Running ${testsToRun.length} ${options.type} tests from ${currentDir}`);
    for (const test of testsToRun) {
      const absolutePath = resolve(rootDir, test.path);
      const testDir = dirname(absolutePath);
      console.log(`Running test: ${test.path}`);
      await execa(test.command, { cwd: testDir, stdio: 'inherit', shell: true });
    }
    
  });

program.parse(process.argv);