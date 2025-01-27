// https://ui.shadcn.com/docs/installation/manual#add-a-cn-helper

import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export const base64Decode = (s: string): string => {
  const binaryString = atob(s)

  const bytes = new Uint8Array(binaryString.length)
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }

  return new TextDecoder().decode(bytes)
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
