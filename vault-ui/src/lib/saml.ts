export interface SAMLMetadata {
  idpEntityId: string;
  idpRedirectUrl: string;
  idpX509Certificate: string;
}

export function parseSamlMetadata(xmlString: string): SAMLMetadata {
  const parser = new DOMParser();
  const xml = parser.parseFromString(xmlString, "application/xml");

  if (xml.getElementsByTagName("parsererror").length > 0) {
    throw new Error("Invalid XML");
  }

  const mdNS = "urn:oasis:names:tc:SAML:2.0:metadata";
  const dsNS = "http://www.w3.org/2000/09/xmldsig#";

  const entityDescriptor = xml.documentElement;
  const idpEntityId = entityDescriptor.getAttribute("entityID");
  if (!idpEntityId) throw new Error("Missing entityID");

  const idpNode = entityDescriptor.getElementsByTagNameNS(
    mdNS,
    "IDPSSODescriptor",
  )[0];
  if (!idpNode) {
    throw new Error("Missing IDPSSODescriptor");
  }

  // Parse <SingleSignOnService> entries
  const ssoBindings: { binding: string; location: string }[] = [];
  const ssoNodes = idpNode.getElementsByTagNameNS(mdNS, "SingleSignOnService");
  for (const el of Array.from(ssoNodes)) {
    const binding = el.getAttribute("Binding");
    const location = el.getAttribute("Location");
    if (binding && location) {
      ssoBindings.push({ binding, location });
    }
  }

  // Parse <KeyDescriptor><ds:KeyInfo><ds:X509Data><ds:X509Certificate>
  const certs: string[] = [];
  const keyDescriptorNodes = idpNode.getElementsByTagNameNS(
    mdNS,
    "KeyDescriptor",
  );
  for (const keyDescriptor of Array.from(keyDescriptorNodes)) {
    const certNodes = keyDescriptor.getElementsByTagNameNS(
      dsNS,
      "X509Certificate",
    );
    for (const certNode of Array.from(certNodes)) {
      const cert = certNode.textContent?.trim();
      if (cert) certs.push(cert);
    }
  }

  // Parse <NameIDFormat>
  const nameIDFormats: string[] = [];
  const nameIDNodes = idpNode.getElementsByTagNameNS(mdNS, "NameIDFormat");
  for (const el of Array.from(nameIDNodes)) {
    const val = el.textContent?.trim();
    if (val) nameIDFormats.push(val);
  }

  // Get redirect URL
  if (ssoBindings.length === 0) {
    throw new Error("No Single Sign-On service found in metadata");
  }
  const firstBinding = ssoBindings[0];
  const idpRedirectUrl = firstBinding.location;

  // Get the first certificate or throw if none found
  if (certs.length === 0) {
    throw new Error("No X.509 certificate found in metadata");
  }
  const firstCert = certs[0];
  const idpX509Certificate = `-----BEGIN CERTIFICATE-----\n${firstCert}\n-----END CERTIFICATE-----`;

  return {
    idpEntityId,
    idpRedirectUrl,
    idpX509Certificate,
  };
}
