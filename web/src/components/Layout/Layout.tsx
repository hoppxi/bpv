import { ReactNode } from "react";
import Header from "./Header";
import Sidebar from "./Sidebar";
import "@/styles/Layout/Layout.css";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  return (
    <div className="layout">
      <Header />
      <div className="layout-content">
        <Sidebar />
        <main className="main-content">{children}</main>
      </div>
    </div>
  );
}
