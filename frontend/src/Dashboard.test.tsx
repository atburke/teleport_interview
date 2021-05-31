import React from 'react';
import {
  render, fireEvent, screen, waitFor,
} from '@testing-library/react';
import Dashboard from './Dashboard';

test('navigates on successful logout', async () => {
  const logout = async () => true;
  const navigate = jest.fn(() => {});

  render(<Dashboard logout={logout} navigate={navigate} />);

  const submit = await screen.findByText('Logout');
  fireEvent.click(submit);

  await waitFor(() => expect(navigate).toHaveBeenCalledWith('/login'));
});
