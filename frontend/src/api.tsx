import axios from 'axios';

const windowData = window as any;

export interface ApiResponse {
  ok: boolean;
  error: string;
}

export async function login(username: string, password: string): Promise<ApiResponse> {
  try {
    await axios({
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
    return { ok: true, error: '' };
  } catch (error) {
    if (error.response) {
      if (error.response.status === 401) {
        return { ok: false, error: 'auth' };
      }

      // Treat other errors the same, since user probably can't do anything
      // about it
      return { ok: false, error: 'server' };
    }

    return { ok: false, error: 'network' };
  }
}

export async function logout(): Promise<ApiResponse> {
  try {
    await axios({
      url: '/logout',
      method: 'post',
      baseURL: 'https://localhost:8080/api/',
      headers: {
        CSRF: windowData.csrfToken,
      },
    });
    return { ok: true, error: '' };
  } catch (error) {
    if (error.response) {
      if (error.response.status === 401) {
        return { ok: false, error: 'auth' };
      }

      return { ok: false, error: 'server' };
    }

    return { ok: false, error: 'network' };
  }
}
