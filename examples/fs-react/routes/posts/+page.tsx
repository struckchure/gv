type PageProps = import("./page.types").PageProps;

export default function Page({ data }: PageProps) {
  return (
    <div>
      <label>{data.title}</label>
      <p>{data.description}</p>
    </div>
  );
}
