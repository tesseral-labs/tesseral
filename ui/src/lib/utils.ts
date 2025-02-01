import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

const base32Alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567'

export const base32Encode = (buffer: Uint8Array): string => {
  let binaryString = ''
  for (let byte of buffer) {
    binaryString += byte.toString(2).padStart(8, '0')
  }

  let base32String = ''
  for (let i = 0; i < binaryString.length; i += 5) {
    const segment = binaryString.substring(i, i + 5).padEnd(5, '0')
    const index = parseInt(segment, 2)
    base32String += base32Alphabet[index]
  }

  return base32String
}

export const base64Decode = (s: string): string => {
  const binaryString = atob(s)

  const bytes = new Uint8Array(binaryString.length)
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }

  return new TextDecoder().decode(bytes)
}

export const base64urlEncode = (buffer: ArrayBuffer): string => {
  let binary = ''
  let bytes = new Uint8Array(buffer)
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return btoa(binary)
    .replace(/\+/g, '-') // Replace '+' with '-'
    .replace(/\//g, '_') // Replace '/' with '_'
    .replace(/=+$/, '') // Remove padding '='
}

export const cn = (...inputs: ClassValue[]) => {
  return twMerge(clsx(inputs))
}

export const isColorDark = (hex: string) => {
  // Ensure hex is valid
  if (!/^#([0-9A-F]{3}|[0-9A-F]{6})$/i.test(hex)) {
    throw new Error('Invalid hex color')
  }

  // Normalize shorthand hex (e.g., #abc -> #aabbcc)
  if (hex.length === 4) {
    hex = '#' + [...hex.slice(1)].map((char) => char + char).join('')
  }

  // Convert hex to RGB
  const r = parseInt(hex.slice(1, 3), 16) / 255
  const g = parseInt(hex.slice(3, 5), 16) / 255
  const b = parseInt(hex.slice(5, 7), 16) / 255

  // Linearize RGB values
  const linearize = (value: number) =>
    value <= 0.03928 ? value / 12.92 : Math.pow((value + 0.055) / 1.055, 2.4)
  const rLin = linearize(r)
  const gLin = linearize(g)
  const bLin = linearize(b)

  // Calculate luminance
  const luminance = 0.2126 * rLin + 0.7152 * gLin + 0.0722 * bLin

  // Return true if dark (luminance below 0.5)
  return luminance < 0.5
}
