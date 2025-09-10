import '@testing-library/jest-dom';

// Minimal mock for scrollTo to avoid errors in JSDOM
Object.defineProperty(window, 'scrollTo', { value: () => {}, writable: true });

