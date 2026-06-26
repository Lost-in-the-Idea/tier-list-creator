import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities';

interface TierItemProps {
  id: number;
}

const TierItem = ({ id }: TierItemProps) => {

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
  } = useSortable({id: id});
  
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };
  
  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      Test
    </div>
  );
}

export default TierItem;
