import React from 'react';
import {
  render, fireEvent, screen, waitFor,
} from '@testing-library/react';
import Login from './Login';

test('navigates on successful login', async () => {
  const login = async () => 200;
  const navigate = jest.fn(() => {});

  render(<Login login={login} navigate={navigate} />);

  const username = await screen.findByLabelText('Email Address');
  const password = await screen.findByLabelText('Password');
  const submit = await screen.findByText('Login to my Dashboard');

  fireEvent.change(username, { target: { value: 'me@example.com' } });
  fireEvent.change(password, { target: { value: 'topsneaky' } });
  fireEvent.click(submit);

  await waitFor(() => expect(navigate).toHaveBeenCalledWith('/dashboard'));
  expect(() => screen.getByText('Invalid email/password.')).toThrow();
});

test('shows error message on failed login', async () => {
  const login = async () => 401;
  const navigate = jest.fn(() => {});

  render(<Login login={login} navigate={navigate} />);

  const username = await screen.findByLabelText('Email Address');
  const password = await screen.findByLabelText('Password');
  const submit = await screen.findByText('Login to my Dashboard');

  fireEvent.change(username, { target: { value: 'me@example.com' } });
  fireEvent.change(password, { target: { value: 'topsneaky' } });
  fireEvent.click(submit);

  await waitFor(() => expect(screen.getByText('Invalid email/password.')).toBeTruthy());
  expect(navigate).not.toHaveBeenCalledWith('/dashboard');
});

test('shows different error message on server error', async () => {
  const login = async () => 500;
  const navigate = jest.fn(() => {});

  render(<Login login={login} navigate={navigate} />);

  const username = await screen.findByLabelText('Email Address');
  const password = await screen.findByLabelText('Password');
  const submit = await screen.findByText('Login to my Dashboard');

  fireEvent.change(username, { target: { value: 'me@example.com' } });
  fireEvent.change(password, { target: { value: 'topsneaky' } });
  fireEvent.click(submit);

  await waitFor(() => expect(screen.getByText('Server error! Please contact [somebody] for assistance.')).toBeTruthy());
  expect(navigate).not.toHaveBeenCalledWith('/dashboard');
});
