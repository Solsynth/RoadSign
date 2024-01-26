import Navbar from "./shared/Navbar";

export default function RootLayout(props: any) {
  return (
    <div>
      <Navbar />

      <main class="h-[calc(100vh-64px)]">{props.children}</main>
    </div>
  );
}
