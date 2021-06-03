import React from 'react';
import './index.css';

interface DashboardProps {
  logout: () => Promise<boolean>;
  navigate: (r: string) => void;
}

interface DashboardState {
  errorMessage: string;
}

class Dashboard extends React.Component<DashboardProps, DashboardState> {
  constructor(props: DashboardProps) {
    super(props);

    this.state = {
      errorMessage: '',
    };

    this.logout = this.logout.bind(this);
  }

  public async logout(): Promise<void> {
    const { logout, navigate } = this.props;
    const status = await logout();
    // If they get 401, they weren't supposed to be here in the first place!
    if (status === 200 || status === 401) {
      navigate('/login');
    } else {
      this.setState({ errorMessage: 'Server error! Please contact [somebody] for assistance.' });
    }
  }

  render() {
    const barStyle = { width: '35%' };
    const { errorMessage } = this.state;
    const alertStyle = { 'margin-top': '1em', display: errorMessage ? 'block' : 'none' };
    return (
      <div>
        <header className="top-nav">
          <h1>
            <i className="material-icons">supervised_user_circle</i>
            User Management Dashboard
          </h1>
          <button className="button is-border" type="button" onClick={this.logout}>Logout</button>
        </header>

        <div className="alert is-error" style={alertStyle}>{errorMessage}</div>

        <div className="plan">
          <header>Startup Plan - $100/Month</header>

          <div className="plan-content">
            <div className="progress-bar">
              <div style={barStyle} className="progress-bar-usage" />
            </div>

            <h3>Users: 35/100</h3>
          </div>

          <footer>
            <button className="button is-success" type="button">Upgrade to Enterprise Plan</button>
          </footer>
        </div>
      </div>
    );
  }
}

export default Dashboard;
