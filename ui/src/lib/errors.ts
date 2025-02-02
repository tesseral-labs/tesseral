import { ErrorDetail } from '@/gen/openauth/common/v1/common_pb'

interface ApiError {
  code: string
  details: { debug: ErrorDetail }[]
  message: string
}

export const parseErrorMessage = (error: any): string => {
  let message = error

  if (error.details) {
    const err = error as ApiError
    message = `${err.details[0].debug.description}`
  }

  return message
}
