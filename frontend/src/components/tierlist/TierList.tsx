import { closestCenter, DndContext, DragEndEvent, KeyboardSensor, PointerSensor, UniqueIdentifier, useSensor, useSensors } from '@dnd-kit/core'
import { arrayMove, horizontalListSortingStrategy, SortableContext, sortableKeyboardCoordinates, verticalListSortingStrategy } from '@dnd-kit/sortable'
import TierItem from './TierItem'
import { useState } from 'react';
import TierRow from './TierRow';
import { createRange } from '../../utils/helpers';

type Items = Record<UniqueIdentifier, UniqueIdentifier[]>;

interface Props {
  initialItems?: Items;
}

const TierList = ({
  initialItems,
}: Props) => {

  const itemCount = 4;

  const [items, setItems] = useState<Items>(
    () =>
      initialItems ?? {
        A: createRange(itemCount, (index) => `A${index + 1}`),
        B: createRange(itemCount, (index) => `B${index + 1}`),
        C: createRange(itemCount, (index) => `C${index + 1}`),
        D: createRange(itemCount, (index) => `D${index + 1}`),
      }
  );

  const [rows, setRows] = useState(
    Object.keys(items) as UniqueIdentifier[]
  )

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );


  console.log(rows)

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
    >

    </DndContext>
  )

  // function handleDragEnd(event: DragEndEvent) {
  //   const { active, over } = event;

  //   if (active.id !== over?.id && over) {
  //     if (items.includes(active.id)) {
  //       setItems((items) => {
  //         const oldIndex = items.indexOf(active.id);
  //         const newIndex = items.indexOf(over.id);
  
  //         return arrayMove(items, oldIndex, newIndex);
  //       });
  //     }

  //     if (items2.includes(active.id)) {
  //       setItems2((items) => {
  //         const oldIndex = items.indexOf(active.id);
  //         const newIndex = items.indexOf(over.id);
  
  //         return arrayMove(items, oldIndex, newIndex);
  //       });
  //     }
  //   }
  // }

}

export default TierList