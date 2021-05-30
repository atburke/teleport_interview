import React from 'react';
import Login from './Login';
import Dashboard from './Dashboard';
// import logo from './logo.svg';
// import './App.css';

interface AppProps {
  login: (user: string, pass: string) => Promise<boolean>;
  logout: () => Promise<boolean>;
}

interface AppState {
  route: string;
}

class App extends React.Component<AppProps, AppState> {
  constructor(props: AppProps) {
    super(props);
    this.state = {

      // Since we really only have 2 views, I'm doing this instead of proper
      // routing.
      route: '/login',
    };
  }

  public navTo(route: string): void {
    this.setState({ route });
  }

  render() {
    const { route } = this.state;
    const { login, logout } = this.props;
    if (route === '/login') {
      return (
        <Login
          login={(uname, pass) => login(uname, pass)}
          navigate={(r) => this.navTo(r)}
        />
      );
    } if (route === '/dashboard') {
      return (
        <Dashboard
          logout={() => logout()}
          navigate={(r) => this.navTo(r)}
        />
      );
    }

    return (
      <div>
        Error: bad route
        {route}
      </div>
    );
  }
}

// function App() {
//   return (
//     <div className="App">
//       <header className="App-header">
//         <img src={logo} className="App-logo" alt="logo" />
//         <p>
//           Edit
//           {' '}
//           <code>src/App.tsx</code>
//           {' '}
//           and save to reload.
//         </p>
//         <a
//           className="App-link"
//           href="https://reactjs.org"
//           target="_blank"
//           rel="noopener noreferrer"
//         >
//           Learn React
//         </a>
//       </header>
//     </div>
//   );
// }

export default App;
