import { ErrorDetail } from '@/gen/tesseral/common/v1/common_pb';

interface ApiError {
  code: string;
  details: { debug: ErrorDetail }[];
  message: string;
}

/* eslint-disable */
export const parseErrorMessage = (error: any): string => {
  let message = !!error.message ? error.message : error;

  if (!!error.details && error.details.length > 0) {
    const err = error as ApiError;
    message = `${err.details[0].debug.description}`;
  }

  return message;
};
/* eslint-enable */
