import React from "https://esm.sh/react@19.1.0";

export default function Page() {
  const [count, setCount] = React.useState(0);

  return (
    <>
      <title>Root</title>

      <p>Root</p>
      <button onClick={() => setCount(count + 1)}>Count {count}</button>
    </>
  );
}
