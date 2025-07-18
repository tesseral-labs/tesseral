// https://ui.shadcn.com/docs/installation/manual#add-a-cn-helper
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function base64Decode(s: string): string {
  const binaryString = atob(s);

  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }

  return new TextDecoder().decode(bytes);
}

export function base64urlEncode(buffer: ArrayBuffer): string {
  let binary = "";
  const bytes = new Uint8Array(buffer);
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary)
    .replace(/\+/g, "-") // Replace '+' with '-'
    .replace(/\//g, "_") // Replace '/' with '_'
    .replace(/=+$/, ""); // Remove padding '='
}

export function base64urlDecode(encoded: string): string {
  // Normalize to standard base64
  encoded = encoded.replace(/-/g, "+").replace(/_/g, "/");

  // Add padding if necessary
  const padding = encoded.length % 4;
  if (padding) {
    encoded += "=".repeat(4 - padding);
  }

  encoded = atob(encoded);

  const bytes = new Uint8Array(encoded.length);
  for (let i = 0; i < encoded.length; i++) {
    bytes[i] = encoded.charCodeAt(i);
  }
  return new TextDecoder().decode(bytes);
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function hexToHSL(hex: string): string {
  // Remove the "#" if present
  hex = hex.replace(/^#/, "");

  // Convert to RGB
  const r = parseInt(hex.substring(0, 2), 16) / 255;
  const g = parseInt(hex.substring(2, 4), 16) / 255;
  const b = parseInt(hex.substring(4, 6), 16) / 255;

  // Get max and min values
  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);

  let h: number = 0;
  let s: number = 0;
  const l: number = (max + min) / 2;

  if (max === min) {
    h = s = 0; // Achromatic (gray)
  } else {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }

    h *= 60;
  }

  return `${Math.round(h)} ${Math.round(s * 100)}% ${Math.round(l * 100)}%`;
}

export function isColorDark(hex: string) {
  // Ensure hex is valid
  if (!/^#([0-9A-F]{3}|[0-9A-F]{6})$/i.test(hex)) {
    throw new Error("Invalid hex color");
  }

  // Normalize shorthand hex (e.g., #abc -> #aabbcc)
  if (hex.length === 4) {
    hex = "#" + [...hex.slice(1)].map((char) => char + char).join("");
  }

  // Convert hex to RGB
  const r = parseInt(hex.slice(1, 3), 16) / 255;
  const g = parseInt(hex.slice(3, 5), 16) / 255;
  const b = parseInt(hex.slice(5, 7), 16) / 255;

  // Linearize RGB values
  function linearize(value: number) {
    return value <= 0.03928
      ? value / 12.92
      : Math.pow((value + 0.055) / 1.055, 2.4);
  }

  const rLin = linearize(r);
  const gLin = linearize(g);
  const bLin = linearize(b);

  // Calculate luminance
  const luminance = 0.2126 * rLin + 0.7152 * gLin + 0.0722 * bLin;

  // Return true if dark (luminance below 0.5)
  return luminance < 0.5;
}

export function titleCaseSlug(str: string, lastId: boolean = false) {
  if (str.includes("_")) {
    if (!lastId) {
      return str.substring(0, str.lastIndexOf("_") + 6) + "...";
    }

    return str;
  }

  // Convert to lowercase and split by hyphen
  const words = str.toLowerCase().split("-");

  // Capitalize the first letter of each word
  const capitalizedWords = words.map((word) => {
    if (word.length === 0) return "";
    return word.charAt(0).toUpperCase() + word.slice(1);
  });

  // Join the words back together with spaces
  return capitalizedWords
    .join(" ")
    .replace(/(api)/gi, "API")
    .replace(/(oidc)/gi, "OIDC")
    .replace(/(saml)/gi, "SAML")
    .replace(/(scim)/gi, "SCIM")
    .replace(/(ui)/gi, "UI")
    .replace(/(rbac)/gi, "RBAC");
}

export function toTitleCase(str: string): string {
  return str
    .toLowerCase()
    .split(" ")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}
