export interface TestItem {
    path: string;
    description: string;
    command: string;
    type: string;
  }
  
  export interface Config {
    tests: TestItem[];
  }