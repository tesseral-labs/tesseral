import { expect, test } from 'vitest';
import { base64UrlDecode } from './utils';

test('decodes access token with missing padding', () => {
    const accessToken = 'eyJpc3MiOiJodHRwczovL3Byb2plY3QtNzlsZHd3d3p5Ym42NmR4YTkxdWRpN21uMy50ZXNzZXJhbC5hcHAiLCJzdWIiOiJ1c2VyXzB0bGk4dWlwZWtpcHJ0dTNxamRzNnVqOWUiLCJhdWQiOiJodHRwczovL3Byb2plY3QtNzlsZHd3d3p5Ym42NmR4YTkxdWRpN21uMy50ZXNzZXJhbC5hcHAiLCJleHAiOjE3NDgzMDMwNTEsIm5iZiI6MTc0ODMwMjc1MSwiaWF0IjoxNzQ4MzAyNzUxLCJzZXNzaW9uIjp7ImlkIjoic2Vzc2lvbl8wM2UwMGE1anM1dGpoZXNrczZyaW54bzZ2In0sInVzZXIiOnsiaWQiOiJ1c2VyXzB0bGk4dWlwZWtpcHJ0dTNxamRzNnVqOWUiLCJlbWFpbCI6ImRpbGxvbkBjZWxlc3QuZGV2IiwiZGlzcGxheU5hbWUiOiJEaWxsb24gTnlzIiwicHJvZmlsZVBpY3R1cmVVcmwiOiJodHRwczovL2F2YXRhcnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tL3UvMjQ3NDA4NjM_dj00In0sIm9yZ2FuaXphdGlvbiI6eyJpZCI6Im9yZ181d3d3b3VpemRxc3RlMWE0NWptd2k4MDVzIiwiZGlzcGxheU5hbWUiOiJBQ01FIENvcnAifX0';
    const claims = '{"iss":"https://project-79ldwwwzybn66dxa91udi7mn3.tesseral.app","sub":"user_0tli8uipekiprtu3qjds6uj9e","aud":"https://project-79ldwwwzybn66dxa91udi7mn3.tesseral.app","exp":1748303051,"nbf":1748302751,"iat":1748302751,"session":{"id":"session_03e00a5js5tjhesks6rinxo6v"},"user":{"id":"user_0tli8uipekiprtu3qjds6uj9e","email":"dillon@celest.dev","displayName":"Dillon Nys","profilePictureUrl":"https://avatars.githubusercontent.com/u/24740863?v=4"},"organization":{"id":"org_5wwwouizdqste1a45jmwi805s","displayName":"ACME Corp"}}';
    expect(base64UrlDecode(accessToken)).toStrictEqual(claims);
});
