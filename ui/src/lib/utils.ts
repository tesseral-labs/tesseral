import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
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
