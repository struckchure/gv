import { useState } from "react";
import { Link } from "react-router";

type PageProps = import("./page.types").PageProps;

export default function Page({ data }: PageProps) {
  const [count, setCount] = useState(0);
  return (
    <>
      <button onClick={() => setCount(count + 1)}>Count: {count}</button>
      <Link to="/">Home</Link>
      {data.names.map((name, idx) => (
        <span key={idx}>{name}</span>
      ))}
    </>
  );
}
