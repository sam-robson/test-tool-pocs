import React, { useState, useEffect } from 'react';
import { Box, Text, useInput, useApp } from 'ink';
import { relative } from 'path';
import { execa } from 'execa';
import { TestItem } from './types.js';
import { findRepoRoot, loadConfig } from './config.js';

const TestRunner: React.FC = () => {
  const { exit } = useApp();
  const [items, setItems] = useState<TestItem[]>([]);
  const [cursor, setCursor] = useState(0);
  const [rootDir, setRootDir] = useState('');
  const [currentDir, setCurrentDir] = useState('');

  useEffect(() => {
    const init = async () => {
      try {
        const root = await findRepoRoot();
        const config = loadConfig(root);
        const cwd = process.cwd();
        const relDir = relative(root, cwd);

        const testItems = Object.entries(config.tests).map(([path, test]) => ({
          path,
          test: test,
          selected: false,
          isAvailable: relDir === '.' || path.startsWith(relDir),
        }));

        setRootDir(root);
        setCurrentDir(relDir || '.');
        setItems(testItems);
      } catch (error) {
        console.error('Error initializing:', error);
        exit();
      }
    };

    init();
  }, []);

  useInput((input, key) => {
    if (input === 'q') {
      exit();
    }

    if (key.upArrow && cursor > 0) {
      setCursor(cursor - 1);
    }

    if (key.downArrow && cursor < items.length - 1) {
      setCursor(cursor + 1);
    }

    if (input === ' ') {
      const newItems = [...items];
      if (newItems[cursor].isAvailable) {
        newItems[cursor].selected = !newItems[cursor].selected;
        setItems(newItems);
      }
    }

    if (key.return) {
      runTests();
    }
  });

  const runTests = async () => {
    const selectedTests = items.filter(item => item.selected);
    
    for (const item of selectedTests) {
      console.log(`\n=== Running ${item.path} ===\n`);
      
      const [cmd, ...args] = item.test.command.split(' ');
      try {
        await execa(cmd, [...args, item.path], {
          cwd: rootDir,
          stdio: 'inherit',
        });
      } catch (error) {
        console.error(`\nError running test ${item.path}:`, error);
      }
    }

    // Reset selections
    setItems(items.map(item => ({ ...item, selected: false })));
  };

  return (
    <Box flexDirection="column">
      <Box marginBottom={1}>
        <Text bold color="blue">Test Runner</Text>
      </Box>
      <Box marginBottom={1}>
        <Text color="gray">Directory: {currentDir}</Text>
      </Box>
      {items.map((item, i) => (
        <Box key={item.path}>
          <Text>
            {cursor === i ? '❯ ' : '  '}
            {item.selected ? '[✓]' : '[ ]'} 
            <Text color={!item.isAvailable ? "gray" : undefined}>
              {item.path.padEnd(40)}  
              [{item.test.type.padEnd(13)}]  
              {item.test.description}
            </Text>
          </Text>
        </Box>
      ))}
      <Box marginTop={1}>
        <Text color="gray">space: select • enter: run • q: quit</Text>
      </Box>
    </Box>
  );
};

export default TestRunner;