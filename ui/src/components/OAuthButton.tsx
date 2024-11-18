import React, { FC } from 'react'
import { cn } from '@/lib/utils'
import { Button, ButtonProps } from '@/components/ui/button'

export enum OAuthMethods {
  google = 'Google',
  microsoft = 'Microsoft',
}

interface OAuthButtonProps extends ButtonProps {
  method: OAuthMethods
}

const OAuthButton: FC<OAuthButtonProps> = ({ method, ...props }) => {
  return (
    <Button {...props}>
      <div className="mr-3">
        {method === OAuthMethods.google ? (
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M22.56 12.25C22.56 11.47 22.49 10.72 22.36 10H12V14.26H17.92C17.66 15.63 16.88 16.79 15.71 17.57V20.34H19.28C21.36 18.42 22.56 15.6 22.56 12.25Z"
              fill="#4285F4"
            />
            <path
              d="M12 23C14.97 23 17.46 22.02 19.28 20.34L15.71 17.57C14.73 18.23 13.48 18.63 12 18.63C9.13999 18.63 6.70999 16.7 5.83999 14.1H2.17999V16.94C3.98999 20.53 7.69999 23 12 23Z"
              fill="#34A853"
            />
            <path
              d="M5.84 14.09C5.62 13.43 5.49 12.73 5.49 12C5.49 11.27 5.62 10.57 5.84 9.91V7.07H2.18C1.43 8.55 1 10.22 1 12C1 13.78 1.43 15.45 2.18 16.93L5.03 14.71L5.84 14.09Z"
              fill="#FBBC05"
            />
            <path
              d="M12 5.38C13.62 5.38 15.06 5.94 16.21 7.02L19.36 3.87C17.45 2.09 14.97 1 12 1C7.69999 1 3.98999 3.47 2.17999 7.07L5.83999 9.91C6.70999 7.31 9.13999 5.38 12 5.38Z"
              fill="#EA4335"
            />
          </svg>
        ) : null}
        {method === OAuthMethods.microsoft ? (
          <svg
            width="24"
            height="23"
            viewBox="0 0 24 23"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <g clip-path="url(#clip0_65_191)">
              <path d="M0 0H24V23H0V0Z" fill="#F3F3F3" />
              <path d="M1.04347 1H11.4783V11H1.04347V1Z" fill="#F35325" />
              <path d="M12.5217 1H22.9565V11H12.5217V1Z" fill="#81BC06" />
              <path d="M1.04347 12H11.4783V22H1.04347V12Z" fill="#05A6F0" />
              <path d="M12.5217 12H22.9565V22H12.5217V12Z" fill="#FFBA08" />
            </g>
            <defs>
              <clipPath id="clip0_65_191">
                <rect width="24" height="23" fill="white" />
              </clipPath>
            </defs>
          </svg>
        ) : null}
      </div>
      Continue with {method}
    </Button>
  )
}

export default OAuthButton
