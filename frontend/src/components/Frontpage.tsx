import { useEffect, useState } from "react";
import { getAllTierlists } from "../api/api-calls";

export const Frontpage = () => {
  const [tierlists, setTierlists] = useState([]);

  useEffect(() => {
    getAllTierlists()
      .then((response) => {
        console.log(response);
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  return (
    <div className="min-h-screen min-w-[365px] flex items-center flex-col">
      <div className="max-w-[1280px] p-4 w-full">
        <h1 className="text-3xl font-bold underline text-center">
          Hello World
        </h1>
      </div>
    </div>
  );
};
