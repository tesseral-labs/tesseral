import React, { FC } from 'react'

interface TextDividerProps {
  text: string
}

const TextDivider: FC<TextDividerProps> = ({ text }) => {
  return (
    <div className="relative w-full mx-auto my-8 flex flex-row justify-center">
      <div className="w-full border-t absolute mt-5 z-0" />
      <div className="px-8 py-2 inline-block relative z-1 m-auto flex-grow-0 bg-card text-border uppercase font-semibold">
        {text}
      </div>
    </div>
  )
}

export default TextDivider
