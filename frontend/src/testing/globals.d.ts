declare function describe(description: string, specDefinitions: () => void): void;
declare function it(description: string, expectation: () => void): void;
declare function beforeEach(action: () => void): void;
declare function afterEach(action: () => void): void;
declare function expect<T>(actual: T): any;
declare function spyOn<T>(object: T, method: keyof T): any;
declare function fail(message?: string): never;
