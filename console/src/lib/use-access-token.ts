import { DOGFOOD_PROJECT_ID } from '@/config';
import { useMemo } from 'react';

const ACCESS_TOKEN_NAME = `tesseral_${DOGFOOD_PROJECT_ID}_access_token`;

export function useAccessToken() {
  return useMemo(() => {
    return document.cookie.split(';').find((row) => row.startsWith(`${ACCESS_TOKEN_NAME}=`))?.split('=')[1];
  }, [document.cookie])
}
