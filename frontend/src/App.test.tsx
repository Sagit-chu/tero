// frontend/src/App.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import App from './App';
import '@testing-library/jest-dom';

describe('App Component', () => {
  beforeEach(() => {
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  test('renders loading state initially', () => {
    // Mock fetch to return a promise that never resolves
    (global.fetch as any).mockReturnValue(new Promise(() => {}));
    render(<App />);
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  test('renders success state', async () => {
    (global.fetch as any).mockResolvedValue({
      json: () => Promise.resolve({ status: 'ok' }),
    });
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('ok')).toBeInTheDocument();
    });
    expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
  });

  test('renders error state on fetch failure', async () => {
    (global.fetch as any).mockRejectedValue(new Error('Network Error'));
    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Error')).toBeInTheDocument();
    });
  });
});
