import { useEffect, useState } from "react";
import { getAllTierlists } from "../api/api-calls";
import { IndividualTierlist } from "./tierlist-components/IndividualTierlist";

interface Tier {
  id: number;
  tierlist_id: number;
  text: string;
  colour: string;
}

interface Item {
  id: number;
  tierlist_id: number;
  text: string;
  image: string;
  tier_text: string;
}

interface Tierlist {
  id: number;
  name: string;
  description: string;
  creator_id: number;
  tiers: Tier[];
  items: Item[];
  version: number;
}

export const Frontpage = () => {
  const [tierlists, setTierlists] = useState<Tierlist[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setIsLoading(true);
    getAllTierlists()
      .then((response) => {
        setTierlists(response.data);
        setIsLoading(false);
      })
      .catch((error) => {
        setError(error.message);
        console.error(error);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  if (isLoading) return <div>Loading tierlists...</div>;
  if (error) return <div>Error: {error}</div>;
  if (tierlists.length === 0) return <div>No tierlists found</div>;

  return (
    <div className="min-h-screen min-w-[365px] flex items-center flex-col">
      <div className="max-w-[1280px] pt-4 w-full">
        {tierlists.map((tierlist) => (
          <IndividualTierlist
            key={tierlist.id}
            name={tierlist.name}
            description={tierlist.description}
            tiers={tierlist.tiers}
            items={tierlist.items}
          />
        ))}
      </div>
    </div>
  );
};
