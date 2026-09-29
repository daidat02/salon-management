import type { Metadata } from "next";
import CategoryHeading from "./components/CategoryHeading";
import CategoryToolbar from "./components/CategoryToolbar";
import CategoryGrid from "./components/CategoryGrid";

export const metadata: Metadata = { title: "Danh mục | SALON ADMIN" };

export default function DanhMucPage() {
  return (
    <div className="space-y-4">
      <CategoryHeading />
      <CategoryToolbar />
      <CategoryGrid />
    </div>
  );
}
