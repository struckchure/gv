import { Outlet } from "react-router";

export default function Layout() {
  return (
    <>
      <title>GV + React's Todo App</title>
      <main>
        <h4>My Todos</h4>

        <Outlet />
      </main>
    </>
  );
}
