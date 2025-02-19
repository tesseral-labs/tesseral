// https://ui.shadcn.com/docs/installation/manual#add-a-cn-helper

import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export const base64Decode = (s: string): string => {
  const binaryString = atob(s);

  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }

  return new TextDecoder().decode(bytes);
};

export const base64urlEncode = (buffer: ArrayBuffer): string => {
  let binary = '';
  const bytes = new Uint8Array(buffer);
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary)
    .replace(/\+/g, '-') // Replace '+' with '-'
    .replace(/\//g, '_') // Replace '/' with '_'
    .replace(/=+$/, ''); // Remove padding '='
};

export const cn = (...inputs: ClassValue[]) => {
  return twMerge(clsx(inputs));
};
