import { redirect } from "next/navigation";

// Ngoài cùng không có page — vào `/` tự đá về landing.
export default function RootRedirectPage() {
  redirect("/landing");
}
