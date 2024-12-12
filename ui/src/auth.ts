const intermediateSessionTokenKey = 'intermediate_session'

export function getIntermediateSessionToken(): string | null {
  return localStorage.getItem(intermediateSessionTokenKey)
}

export function setIntermediateSessionToken(s: string) {
  localStorage.setItem(intermediateSessionTokenKey, s)
}
