// frontend/src/App.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { describe, test, expect, beforeEach, afterEach, vi } from 'vitest';
import App from './App';
import '@testing-library/jest-dom';

describe('App Component', () => {
  beforeEach(() => {
    globalThis.fetch = vi.fn();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  test('renders loading state initially', () => {
    // Mock fetch to return a promise that never resolves
    (globalThis.fetch as any).mockReturnValue(new Promise(() => {}));
    render(<App />);
    expect(screen.getAllByText('Loading...').length).toBeGreaterThan(0);
  });

  test('renders success state', async () => {
    (globalThis.fetch as any).mockImplementation((url: string) => {
      if (url === '/api/status') {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ status: 'ok', node_status: 'Alive' }),
        });
      }
      if (url === '/api/nodes') {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([]),
        });
      }
      if (url === '/api/config') {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({}),
        });
      }
      return Promise.reject(new Error('not found'));
    });
    
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Online')).toBeInTheDocument();
      expect(screen.getByText('Alive')).toBeInTheDocument();
    });
    expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
  });

  test('renders error state on fetch failure', async () => {
    (globalThis.fetch as any).mockRejectedValue(new Error('Network Error'));
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getAllByText('Error').length).toBeGreaterThan(0);
    });
  });
});
