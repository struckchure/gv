import { FormEvent, useState } from "react";
import {
  LoaderFunctionArgs,
  useLoaderData,
  useRevalidator,
} from "react-router";

interface Todo {
  id: number;
  title: string;
  completed: boolean;
}

export async function loader(args: LoaderFunctionArgs) {
  const todos = (await (await fetch("/api/todos/")).json()) as Todo[];

  return { todos: [] };
}

export default function Page() {
  const { revalidate } = useRevalidator();
  const { todos } = useLoaderData<typeof loader>();

  const [title, setTitle] = useState("");

  async function createTodo(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();

    await fetch("/api/todos/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title }),
    });
    await revalidate();
    setTitle("");
  }

  async function updateTodo(id: number, completed: boolean) {
    await fetch(`/api/todos/${id}/`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ completed }),
    });
    await revalidate();
  }

  async function deleteTodo(id: number) {
    await fetch(`/api/todos/${id}/`, {
      method: "DELETE",
    });
    await revalidate();
  }

  return (
    <section className="todo-page">
      <form onSubmit={createTodo}>
        <input
          id="new-todo"
          type="text"
          placeholder="uninstall nodejs,bun,deno ..."
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          autoComplete="off"
        />
        <button>Create</button>
      </form>

      <div className="todo-list">
        {todos.map(
          (
            todo: any,
            idx: any // TODO: fix typing issue for react-router types
          ) => (
            <div className="todo-item" key={idx}>
              <input
                id={idx}
                type="checkbox"
                checked={todo.completed}
                onChange={async () =>
                  await updateTodo(todo.id, !todo.completed)
                }
              />
              <label
                htmlFor={idx}
                style={todo.completed ? { textDecoration: "line-through" } : {}}
              >
                {todo.title}
              </label>

              <button onClick={async () => await deleteTodo(todo.id)}>
                Delete
              </button>
            </div>
          )
        )}
      </div>
    </section>
  );
}
