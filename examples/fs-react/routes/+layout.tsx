type LayoutProps = import("./layout.types").LayoutProps;

export default function Layout(props: LayoutProps) {
  return (
    <html>
      <head>
        <title>React + GV</title>
      </head>

      <body>{props.children}</body>
    </html>
  );
}
