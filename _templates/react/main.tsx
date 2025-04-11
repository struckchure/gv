import { useState } from "react";
import { createRoot } from "react-dom/client";

function App() {
  const [count, setCount] = useState(0);

  return (
    <>
      <title>GV + React</title>

      <main>
        <h4>GV + React</h4>
        <button onClick={() => setCount(count + 1)}>Count : {count}</button>
      </main>
    </>
  );
}

createRoot(document.querySelector("#root")!).render(<App />);
