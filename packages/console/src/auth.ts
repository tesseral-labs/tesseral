const intermediateSessionTokenKey = 'intermediate_session';

export const getIntermediateSessionToken = (): string | null => {
  return localStorage.getItem(intermediateSessionTokenKey);
};

export const setIntermediateSessionToken = (s: string) => {
  localStorage.setItem(intermediateSessionTokenKey, s);
};

export const setAccessToken = (s: string) => {
  localStorage.setItem('access_token', s);
};

export const getAccessToken = (): string | null => {
  return localStorage.getItem('access_token');
};

export const setRefreshToken = (s: string) => {
  localStorage.setItem('refresh_token', s);
};

export const getRefreshToken = (): string | null => {
  return localStorage.getItem('refresh_token');
};
