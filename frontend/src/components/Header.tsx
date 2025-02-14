import { useEffect, useState } from "react";
import { whoAmI } from "../api/api-calls";

export const Header = () => {
  const [user, setUser] = useState({ username: "", discord_id: "" });

  useEffect(() => {
    whoAmI()
      .then((response) => {
        const username = response.data.user.username;
        const discord_id = response.data.user.discord_id;
        setUser({ username: username, discord_id: discord_id });
      })
      .catch((error) => {
        console.log(error);
      });
  }, []);

  return (
    <div className="flex">
      <h1 className="mr-4">Username: {user.username}</h1>
      <h1>Discord ID: {user.discord_id}</h1>
    </div>
  );
};
