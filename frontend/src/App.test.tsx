import { render, screen } from '@testing-library/react';
import App from './App';
import '@testing-library/jest-dom';

test('renders dashboard header', () => {
  render(<App />);
  const linkElement = screen.getByText(/Flvx Monitor Dashboard/i);
  expect(linkElement).toBeInTheDocument();
});