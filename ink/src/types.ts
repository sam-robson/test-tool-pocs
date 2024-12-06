export interface Test {
    description: string;
    command: string;
    type: string;
  }
  
  export interface Config {
    tests: {
      [key: string]: Test;
    };
  }
  
  export interface TestItem {
    path: string;
    test: Test;
    selected: boolean;
    isAvailable: boolean;
  }