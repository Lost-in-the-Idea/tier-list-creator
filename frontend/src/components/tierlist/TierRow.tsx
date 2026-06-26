import { useDroppable } from '@dnd-kit/core'
import React, { useState, FC } from 'react'

export interface TierRowProps {
  children: React.ReactNode
}

const TierRow: FC<TierRowProps> = ({ children }) => {

  const {
    active,
    over,
    isOver,
  } = useDroppable({
    id: 'row',
  })


  return (
    <div>

    </div>
  )
}

export default TierRow