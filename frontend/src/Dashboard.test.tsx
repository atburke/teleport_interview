import React from 'react';
import {
  render, fireEvent, screen, waitFor,
} from '@testing-library/react';
import Dashboard from './Dashboard';

test('navigates on successful logout', async () => {
  const logout = async () => ({ ok: true, error: '' });
  const navigate = jest.fn(() => {});

  render(<Dashboard logout={logout} navigate={navigate} />);

  const submit = await screen.findByText('Logout');
  fireEvent.click(submit);

  await waitFor(() => expect(navigate).toHaveBeenCalledWith('/login'));
  expect(() => screen.getByText('Unexpected server error. Please contact [somebody] for assistance.')).toThrow();
});

test('show error message on server error', async () => {
  const logout = async () => ({ ok: false, error: 'server' });
  const navigate = jest.fn(() => {});

  render(<Dashboard logout={logout} navigate={navigate} />);

  const submit = await screen.findByText('Logout');
  fireEvent.click(submit);

  await waitFor(() => expect(screen.getByText('Unexpected server error. Please contact [somebody] for assistance.')).toBeTruthy());
  expect(navigate).not.toHaveBeenCalledWith('/login');
});
