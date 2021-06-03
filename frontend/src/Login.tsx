import React from 'react';
import './index.css';

interface LoginProps {
  login: (user: string, pass: string) => Promise<boolean>;
  navigate: (r: string) => void;
}

interface LoginState {
  username: string;
  password: string;
  errorMessage: string;
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
    const success = await login(username, password);
    if (success === 200) {
      navigate('/dashboard');
    } else if (success === 401) {
      this.setState({ errorMessage: 'Invalid email/password.' });

    // Most likely 500, but a user won't need to know differently if it isn't.
    } else {
      this.setState({ errorMessage: 'Server error! Please contact [somebody] for assistance.' });
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
