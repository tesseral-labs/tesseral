import { base64urlDecode } from "./utils";

interface AccessTokenClaims {
  exp: number;
}

export function parseAccessToken(accessToken: string): AccessTokenClaims {
  const claimsPart = accessToken.split(".")[1];
  const decodedClaims = base64urlDecode(claimsPart);
  return JSON.parse(decodedClaims) as AccessTokenClaims;
}
