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
  console.log("Tiers:", tiers);
  console.log("Items:", items);
  console.log(
    "Items with their tiers:",
    items.map((item) => `${item.text} - ${item.tier_text}`)
  );
  return (
    <div className="border rounded-lg p-4 mb-4 w-full">
      <h1 className="text-2xl font-bold mb-2">{name}</h1>
      <h2 className="text-lg mb-4">{description}</h2>

      <div className="space-y-2">
        {tiers.map((tier) => {
          const tierItems = items.filter(
            (item) => item.tier_text === tier.text
          );

          return (
            <div key={tier.id} className="flex border rounded">
              <div className="w-24 p-4 flex items-center justify-center border-r">
                {tier.text}
              </div>
              <div className="flex flex-wrap p-2 gap-2">
                {tierItems.map((item) => (
                  <div key={item.id} className="p-2 bg-white rounded shadow">
                    {item.text}
                  </div>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
