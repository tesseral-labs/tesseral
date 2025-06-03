import { parseAccessToken } from "./parse-access-token";

test("should parse valid access token", () => {
  const token =
    "eyJraWQiOiJzZXNzaW9uX3NpZ25pbmdfa2V5XzU3Zmt1NWFpd2dpdWxpYjlvMHN4emF6YmQiLCJhbGciOiJFUzI1NiJ9.eyJpc3MiOiJodHRwczovL3Byb2plY3QtNzlsZHd3d3p5Ym42NmR4YTkxdWRpN21uMy50ZXNzZXJhbC5hcHAiLCJzdWIiOiJ1c2VyXzdjcnR3Y3VvdWZmdDJzYnY2dDN6azJreGQiLCJhdWQiOiJodHRwczovL3Byb2plY3QtNzlsZHd3d3p5Ym42NmR4YTkxdWRpN21uMy50ZXNzZXJhbC5hcHAiLCJleHAiOjE3NDg0NTEzMzMsIm5iZiI6MTc0ODQ1MTAzMywiaWF0IjoxNzQ4NDUxMDMzLCJzZXNzaW9uIjp7ImlkIjoic2Vzc2lvbl8wM2UwZHhxNmRtZjBncmwwZHh1ZGhyamE0In0sInVzZXIiOnsiaWQiOiJ1c2VyXzdjcnR3Y3VvdWZmdDJzYnY2dDN6azJreGQiLCJlbWFpbCI6ImRpbGxvbi5hbmRyZS5ueXNAZ21haWwuY29tIiwiZGlzcGxheU5hbWUiOiJEaWxsb24gTnlzIiwicHJvZmlsZVBpY3R1cmVVcmwiOiJodHRwczovL2F2YXRhcnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tL3UvMjQ3NDA4NjM_dj00In0sIm9yZ2FuaXphdGlvbiI6eyJpZCI6Im9yZ19hemI3YmF1eXEyeHU1c3M0Njh5Nmh3N3gxIiwiZGlzcGxheU5hbWUiOiJBQ01FIENvcnAifX0.f8IjcCw9qrK_cEuiyZ_4iR6mgMoc0aaOOhBKumdzpS3uZsj3szP7mprgUc_yTnivWtMwp5V5GvRYGLF8__LPbQ";
  const expected = {
    iss: "https://project-79ldwwwzybn66dxa91udi7mn3.tesseral.app",
    sub: "user_7crtwcuoufft2sbv6t3zk2kxd",
    aud: "https://project-79ldwwwzybn66dxa91udi7mn3.tesseral.app",
    exp: 1748451333,
    nbf: 1748451033,
    iat: 1748451033,
    session: {
      id: "session_03e0dxq6dmf0grl0dxudhrja4",
    },
    user: {
      id: "user_7crtwcuoufft2sbv6t3zk2kxd",
      email: "dillon.andre.nys@gmail.com",
      displayName: "Dillon Nys",
      profilePictureUrl: "https://avatars.githubusercontent.com/u/24740863?v=4",
    },
    organization: {
      id: "org_azb7bauyq2xu5ss468y6hw7x1",
      displayName: "ACME Corp",
    },
  };
  expect(parseAccessToken(token)).toEqual(expected);
});

test("should parse UTF-8 values in access token", () => {
  const claims = { sub: "😃" };
  const encodedClaims = Buffer.from(JSON.stringify(claims)).toString(
    "base64url",
  );
  const token = `testHeader.${encodedClaims}.testSignature`;
  expect(parseAccessToken(token)).toEqual(claims);
});

test("should throw error for invalid access token", () => {
  const token = "invalid.token.string";
  expect(() => parseAccessToken(token)).toThrow();
});
