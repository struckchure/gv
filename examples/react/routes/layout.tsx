import { Link, Outlet } from "react-router";

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
