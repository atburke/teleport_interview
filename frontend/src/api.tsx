import axios from 'axios';

const windowData = window as any;

export async function login(username: string, password: string): Promise<boolean> {
  try {
    const response = await axios({
      url: '/login',
      method: 'post',
      baseURL: 'https://localhost:8080/api/',
      headers: {
        CSRF: windowData.csrfToken,
      },
      auth: {
        username,
        password,
      },
    });
    return response.status === 200;
  } catch (error) {
    return false;
  }
}

export async function logout(): Promise<boolean> {
  try {
    const response = await axios({
      url: '/logout',
      method: 'post',
      baseURL: 'https://localhost:8080/api/',
      headers: {
        CSRF: windowData.csrfToken,
      },
    });
    return response.status === 200;
  } catch (error) {
    return false;
  }
}
