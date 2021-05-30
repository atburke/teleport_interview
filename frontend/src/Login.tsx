import React from 'react';
import './index.css';

interface LoginProps {
  login: (string, string) => Promise<boolean>;
  navigate: (string) => void;
}

interface LoginState {
  username: string;
  password: string;
  errorMessage: string;
}

class Login extends React.Component<LoginProps, LoginState> {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      errorMessage: '',
    };

    this.setUsername = this.setUsername.bind(this);
    this.setPassword = this.setPassword.bind(this);
    this.login = this.login.bind(this);
  }

  public setUsername(event): void {
    this.setState({ username: event.target.value });
  }

  public setPassword(event): void {
    this.setState({ password: event.target.value });
  }

  public async login() {
    const { username, password } = this.state;
    const { login, navigate } = this.props;
    const success = await login(username, password);
    if (success) {
      navigate('/dashboard');
    } else {
      this.setState({ errorMessage: 'Invalid email/password.' });
    }
  }

  render() {
    const { errorMessage } = this.state;
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
          onClick={this.login}
        />
        <div className="alert is-error">{errorMessage}</div>
      </form>
    );
  }
}

export default Login;
