const intermediateSessionTokenKey = "intermediate_session";

export function getIntermediateSessionToken(): string | null {
  return localStorage.getItem(intermediateSessionTokenKey);
}

export function setIntermediateSessionToken(s: string) {
  localStorage.setItem(intermediateSessionTokenKey, s);
}

export function setAccessToken(s: string) {
  localStorage.setItem("access_token", s);
}

export function getAccessToken(): string | null {
  return localStorage.getItem("access_token");
}

export function setRefreshToken(s: string) {
  localStorage.setItem("refresh_token", s);
}

export function getRefreshToken(): string | null {
  return localStorage.getItem("refresh_token");
}
