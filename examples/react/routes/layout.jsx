import { Link, Outlet } from "https://esm.sh/react-router@7.5.0";
import React from "https://esm.sh/react@19.1.0";

export default function Layout() {
  return (
    <main>
      <nav style={{ display: "flex", gap: "1rem" }}>
        <Link to="/">Home</Link>
        <Link to="/login">Login</Link>
        <Link to="/register">Register</Link>
      </nav>

      <Outlet />
    </main>
  );
}
