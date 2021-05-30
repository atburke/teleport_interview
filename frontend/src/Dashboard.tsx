import React from 'react';
import './index.css';

interface DashboardProps {
  logout: () => Promise<boolean>;
  navigate: (r: string) => void;
}

class Dashboard extends React.Component<DashboardProps, any> {
  constructor(props: DashboardProps) {
    super(props);

    this.logout = this.logout.bind(this);
  }

  public async logout(): Promise<void> {
    const { logout, navigate } = this.props;
    const success = await logout();
    if (success) {
      navigate('/login');
    }
  }

  render() {
    const barStyle = { width: '35%' };
    return (
      <div>
        <header className="top-nav">
          <h1>
            <i className="material-icons">supervised_user_circle</i>
            User Management Dashboard
          </h1>
          <button className="button is-border" type="button" onClick={this.logout}>Logout</button>
        </header>

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
