import "./App.css";
import { Frontpage } from "./components/Frontpage";
import { Header } from "./components/Header";

function App() {
  return (
    <div className="min-h-screen min-w-[365px] flex items-center flex-col">
      <div className="max-w-[1280px] p-4 w-full">
        <Header />
        <Frontpage />
      </div>
    </div>
  );
}

export default App;
