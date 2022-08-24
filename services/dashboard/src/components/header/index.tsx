import { Component } from "solid-js";
import { Link } from "@solidjs/router";

const Header: Component = () => {
  return (
    <nav class="bg-slate-800 border-gray-200 px-2 sm:px-4 py-2.5  dark:bg-gray-900">
      <div class="container flex flex-wrap justify-between items-center mx-auto">
        <Link href="/">
          <a class="text-gray-300 font-semibold text-lg hover:text-gray-600">
            Home
          </a>
        </Link>
        <Link href="/suppliers">
          <a class="text-gray-300 font-semibold text-lg hover:text-gray-600">
            Suppliers
          </a>
        </Link>
        <Link href="/customers">
          <a class="text-gray-300 font-semibold text-lg hover:text-gray-600">
            Customers
          </a>
        </Link>
      </div>
    </nav>
  );
};
export default Header;
