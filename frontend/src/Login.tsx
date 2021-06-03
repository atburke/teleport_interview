import React from 'react';
import './index.css';
import { ApiResponse } from './api';

interface LoginProps {
  login: (user: string, pass: string) => Promise<ApiResponse>;
  navigate: (r: string) => void;
}

interface LoginState {
  username: string;
  password: string;
  errorMessage: string;
}

// should check to see if airbnb style is ok with using objects as dicts
function interpretError(error: string): string {
  if (error === 'auth') {
    return 'Invalid email/password.';

  // server or network error
  }
  return `Unexpected ${error} error. Please contact [somebody] for assistance.`;
}

class Login extends React.Component<LoginProps, LoginState> {
  constructor(props: LoginProps) {
    super(props);
    this.state = {
      username: '',
      password: '',
      errorMessage: '',
    };

    this.setUsername = this.setUsername.bind(this);
    this.setPassword = this.setPassword.bind(this);
    this.onSubmit = this.onSubmit.bind(this);
    this.login = this.login.bind(this);
  }

  public onSubmit(event: any): void {
    event.preventDefault();
    this.login();
  }

  public setUsername(event: any): void {
    this.setState({ username: event.target.value });
  }

  public setPassword(event: any): void {
    this.setState({ password: event.target.value });
  }

  public async login() {
    const { username, password } = this.state;
    const { login, navigate } = this.props;
    const { ok, error } = await login(username, password);
    if (ok) {
      navigate('/dashboard');
    } else {
      this.setState({ errorMessage: interpretError(error) });
    }
  }

  render() {
    const { errorMessage } = this.state;
    const alertStyle = { 'margin-top': '1em', display: errorMessage ? 'block' : 'none' };
    return (
      <form className="login-form">
        <h1>Sign Into Your Account</h1>
        <div>
          <label htmlFor="email">
            Email Address
            <input type="email" id="email" className="field" onChange={this.setUsername} />
          </label>
        </div>
        <div>
          <label htmlFor="password">
            Password
            <input type="password" id="password" className="field" onChange={this.setPassword} />
          </label>
        </div>
        <input
          type="submit"
          value="Login to my Dashboard"
          className="button block"
          onClick={this.onSubmit}
        />
        <div className="alert is-error" style={alertStyle}>{errorMessage}</div>
      </form>
    );
  }
}

export default Login;
