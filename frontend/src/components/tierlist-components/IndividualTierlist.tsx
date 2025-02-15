interface TierlistProps {
  name: string;
  description: string;
  tiers: Array<{
    id: number;
    tierlist_id: number;
    text: string;
    colour: string;
  }>;
  items: Array<{
    id: number;
    tierlist_id: number;
    text: string;
    image: string;
    tier_text: string;
  }>;
}

export const IndividualTierlist: React.FC<TierlistProps> = ({
  name,
  description,
  tiers,
  items,
}) => {
  return (
    <div className="border rounded-lg p-4 mb-4 w-full">
      <h1>{name}</h1>
      <h2>{description}</h2>
      <h3>Tiers</h3>
      <ul>
        {tiers.map((tier) => (
          <li key={tier.id}>{tier.text}</li>
        ))}
      </ul>
      <h3>Items</h3>
      <ul>
        {items.map((item) => (
          <li key={item.id}>
            <p>
              {item.text} {item.tier_text}
            </p>
          </li>
        ))}
      </ul>
    </div>
  );
};
