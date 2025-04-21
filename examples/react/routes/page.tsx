import { useState } from "react";

export default function Page() {
  const [count, setCount] = useState(0);

  return (
    <>
      <title>Root</title>

      <p>Root</p>
      <button onClick={() => setCount(count + 2)}>Count {count}</button>
    </>
  );
}
